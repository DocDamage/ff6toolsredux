package plugins

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"sync"
	"time"
)

// PluginSignature represents a digital signature for a plugin
type PluginSignature struct {
	PluginID    string    // Unique plugin identifier
	Hash        string    // SHA256 hash of plugin code
	Signature   string    // Base64-encoded digital signature
	Certificate string    // PEM-encoded public key certificate
	SignedAt    time.Time // Signature creation time
	ExpiresAt   time.Time // Signature expiration time
	SignedBy    string    // Signer identity/name
	Revoked     bool      // Whether signature is revoked
	RevokedAt   time.Time // Revocation time
}

// SecurityManager handles plugin security operations
type SecurityManager struct {
	privateKey     *rsa.PrivateKey
	publicKey      *rsa.PublicKey
	signatures     map[string]*PluginSignature // pluginID -> signature
	trustedKeys    map[string]string           // signer -> public key PEM
	auditLog       []*SecurityEvent
	mu             sync.RWMutex
	maxHistorySize int
}

// SecurityEvent logs security-related operations
type SecurityEvent struct {
	EventID   string
	PluginID  string
	EventType string // "sign", "verify", "revoke", "trust", "untrust"
	Status    string // "success", "failed", "denied"
	Error     string
	Timestamp time.Time
	Details   map[string]string
}

// NewSecurityManager creates a new security manager
func NewSecurityManager() *SecurityManager {
	return &SecurityManager{
		signatures:     make(map[string]*PluginSignature),
		trustedKeys:    make(map[string]string),
		auditLog:       make([]*SecurityEvent, 0, 10000),
		maxHistorySize: 10000,
	}
}

// GenerateKeyPair generates RSA key pair for signing (2048 bits)
func (sm *SecurityManager) GenerateKeyPair() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	sm.privateKey = privateKey
	sm.publicKey = &privateKey.PublicKey

	sm.logSecurityEvent("system", "generate_keys", "success", "", map[string]string{
		"key_size": "2048",
	})

	return nil
}

// SignPlugin creates a digital signature for plugin code
func (sm *SecurityManager) SignPlugin(pluginID, codePath string) (*PluginSignature, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.privateKey == nil {
		return nil, fmt.Errorf("no private key available for signing")
	}

	// Read plugin code
	codeBytes, err := os.ReadFile(codePath)
	if err != nil {
		sm.logSecurityEvent(pluginID, "sign", "failed", err.Error(), nil)
		return nil, fmt.Errorf("failed to read plugin code: %w", err)
	}

	// Compute SHA256 hash
	hash := sha256.Sum256(codeBytes)
	hashStr := hex.EncodeToString(hash[:])

	// Sign the hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, sm.privateKey, 0, hash[:])
	if err != nil {
		sm.logSecurityEvent(pluginID, "sign", "failed", err.Error(), nil)
		return nil, fmt.Errorf("failed to sign plugin: %w", err)
	}

	// Get public key certificate
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(sm.publicKey)
	if err != nil {
		sm.logSecurityEvent(pluginID, "sign", "failed", err.Error(), nil)
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	plugSig := &PluginSignature{
		PluginID:    pluginID,
		Hash:        hashStr,
		Signature:   hex.EncodeToString(signature),
		Certificate: string(pubKeyPEM),
		SignedAt:    time.Now(),
		ExpiresAt:   time.Now().AddDate(1, 0, 0), // 1-year validity
		SignedBy:    "ff6editor-system",
		Revoked:     false,
	}

	sm.signatures[pluginID] = plugSig

	sm.logSecurityEvent(pluginID, "sign", "success", "", map[string]string{
		"hash":    hashStr,
		"expires": plugSig.ExpiresAt.Format(time.RFC3339),
	})

	return plugSig, nil
}

// VerifyPlugin validates a plugin signature
func (sm *SecurityManager) VerifyPlugin(pluginID, codePath string) (bool, error) {
	sm.mu.RLock()
	sig, exists := sm.signatures[pluginID]
	sm.mu.RUnlock()

	if !exists {
		sm.mu.Lock()
		sm.logSecurityEvent(pluginID, "verify", "failed", "no signature found", nil)
		sm.mu.Unlock()
		return false, fmt.Errorf("no signature found for plugin")
	}

	if sig.Revoked {
		sm.mu.Lock()
		sm.logSecurityEvent(pluginID, "verify", "denied", "signature revoked", nil)
		sm.mu.Unlock()
		return false, fmt.Errorf("plugin signature has been revoked")
	}

	if time.Now().After(sig.ExpiresAt) {
		sm.mu.Lock()
		sm.logSecurityEvent(pluginID, "verify", "denied", "signature expired", nil)
		sm.mu.Unlock()
		return false, fmt.Errorf("plugin signature has expired")
	}

	// Read plugin code
	codeBytes, err := os.ReadFile(codePath)
	if err != nil {
		sm.mu.Lock()
		sm.logSecurityEvent(pluginID, "verify", "failed", err.Error(), nil)
		sm.mu.Unlock()
		return false, fmt.Errorf("failed to read plugin code: %w", err)
	}

	// Compute hash
	hash := sha256.Sum256(codeBytes)
	hashStr := hex.EncodeToString(hash[:])

	// Check if hash matches
	if hashStr != sig.Hash {
		sm.mu.Lock()
		sm.logSecurityEvent(pluginID, "verify", "denied", "hash mismatch (code modified)", nil)
		sm.mu.Unlock()
		return false, fmt.Errorf("plugin code has been modified (hash mismatch)")
	}

	// Decode signature
	sigBytes, err := hex.DecodeString(sig.Signature)
	if err != nil {
		sm.mu.Lock()
		sm.logSecurityEvent(pluginID, "verify", "failed", err.Error(), nil)
		sm.mu.Unlock()
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	// Verify signature using public key
	err = rsa.VerifyPKCS1v15(sm.publicKey, 0, hash[:], sigBytes)
	if err != nil {
		sm.mu.Lock()
		sm.logSecurityEvent(pluginID, "verify", "denied", "invalid signature", nil)
		sm.mu.Unlock()
		return false, fmt.Errorf("invalid signature")
	}

	sm.mu.Lock()
	sm.logSecurityEvent(pluginID, "verify", "success", "", map[string]string{
		"hash": hashStr,
	})
	sm.mu.Unlock()

	return true, nil
}

// GetPluginHash computes SHA256 hash of plugin code
func (sm *SecurityManager) GetPluginHash(codePath string) (string, error) {
	codeBytes, err := os.ReadFile(codePath)
	if err != nil {
		return "", fmt.Errorf("failed to read plugin: %w", err)
	}

	hash := sha256.Sum256(codeBytes)
	return hex.EncodeToString(hash[:]), nil
}

// IsPluginTrusted checks if plugin is signed and verified
func (sm *SecurityManager) IsPluginTrusted(pluginID string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sig, exists := sm.signatures[pluginID]
	if !exists {
		return false
	}

	if sig.Revoked || time.Now().After(sig.ExpiresAt) {
		return false
	}

	return true
}

// RevokeSignature marks a plugin signature as revoked
func (sm *SecurityManager) RevokeSignature(pluginID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sig, exists := sm.signatures[pluginID]
	if !exists {
		return fmt.Errorf("signature not found: %s", pluginID)
	}

	sig.Revoked = true
	sig.RevokedAt = time.Now()

	sm.logSecurityEvent(pluginID, "revoke", "success", "", map[string]string{
		"revoked_at": sig.RevokedAt.Format(time.RFC3339),
	})

	return nil
}

// AddTrustedKey adds a trusted signer's public key
func (sm *SecurityManager) AddTrustedKey(signerID, publicKeyPEM string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Validate PEM format
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		sm.logSecurityEvent("system", "trust", "failed", "invalid PEM format", nil)
		return fmt.Errorf("invalid PEM format")
	}

	sm.trustedKeys[signerID] = publicKeyPEM

	sm.logSecurityEvent("system", "trust", "success", "", map[string]string{
		"signer": signerID,
	})

	return nil
}

// RemoveTrustedKey removes a trusted signer's key
func (sm *SecurityManager) RemoveTrustedKey(signerID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.trustedKeys[signerID]; !exists {
		return fmt.Errorf("trusted key not found: %s", signerID)
	}

	delete(sm.trustedKeys, signerID)

	sm.logSecurityEvent("system", "untrust", "success", "", map[string]string{
		"signer": signerID,
	})

	return nil
}

// GetSignature returns the signature for a plugin
func (sm *SecurityManager) GetSignature(pluginID string) *PluginSignature {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.signatures[pluginID]
}

// GetSignatures returns all signatures
func (sm *SecurityManager) GetSignatures() map[string]*PluginSignature {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sigsCopy := make(map[string]*PluginSignature)
	for k, v := range sm.signatures {
		sigsCopy[k] = v
	}
	return sigsCopy
}

// ExportSignatures exports all signatures to a file (JSON format)
func (sm *SecurityManager) ExportSignatures(filePath string) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	type exportData struct {
		Timestamp  time.Time
		Signatures map[string]*PluginSignature
	}

	data := exportData{
		Timestamp:  time.Now(),
		Signatures: sm.signatures,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ImportSignatures imports signatures from a file
func (sm *SecurityManager) ImportSignatures(filePath string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Verify file is readable
	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	type exportData struct {
		Timestamp  time.Time
		Signatures map[string]*PluginSignature
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var data exportData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate each signature before importing
	for pluginID, sig := range data.Signatures {
		if sig == nil {
			continue
		}
		// Ensure basic fields are populated
		if sig.PluginID == "" {
			sig.PluginID = pluginID
		}
		sm.signatures[pluginID] = sig
	}

	sm.logSecurityEvent("system", "import_signatures", "success", "", map[string]string{
		"count": fmt.Sprintf("%d", len(data.Signatures)),
	})

	return nil
}

// ClearSignatures removes all signatures
func (sm *SecurityManager) ClearSignatures() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.signatures = make(map[string]*PluginSignature)

	sm.logSecurityEvent("system", "clear_sigs", "success", "", nil)
}

// logSecurityEvent logs a security operation
func (sm *SecurityManager) logSecurityEvent(pluginID, eventType, status, errMsg string, details map[string]string) {
	event := &SecurityEvent{
		EventID:   fmt.Sprintf("sec_%d", time.Now().UnixNano()),
		PluginID:  pluginID,
		EventType: eventType,
		Status:    status,
		Error:     errMsg,
		Timestamp: time.Now(),
		Details:   details,
	}

	sm.auditLog = append(sm.auditLog, event)

	// Keep history size bounded
	if len(sm.auditLog) > sm.maxHistorySize {
		sm.auditLog = sm.auditLog[len(sm.auditLog)-sm.maxHistorySize:]
	}
}

// GetAuditLog returns the security audit log
func (sm *SecurityManager) GetAuditLog() []*SecurityEvent {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	logCopy := make([]*SecurityEvent, len(sm.auditLog))
	copy(logCopy, sm.auditLog)
	return logCopy
}

// GetAuditLogByPlugin returns audit events for a specific plugin
func (sm *SecurityManager) GetAuditLogByPlugin(pluginID string) []*SecurityEvent {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var events []*SecurityEvent
	for _, e := range sm.auditLog {
		if e.PluginID == pluginID {
			events = append(events, e)
		}
	}
	return events
}

// ClearAuditLog clears the audit log
func (sm *SecurityManager) ClearAuditLog() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.auditLog = make([]*SecurityEvent, 0, 10000)
}

// GetSecurityStats returns security statistics
func (sm *SecurityManager) GetSecurityStats() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	stats := map[string]interface{}{
		"total_signatures":   len(sm.signatures),
		"trusted_keys":       len(sm.trustedKeys),
		"audit_events":       len(sm.auditLog),
		"revoked_plugins":    0,
		"expired_signatures": 0,
		"verified_plugins":   0,
	}

	revokedCount := 0
	expiredCount := 0
	verifiedCount := 0

	for _, sig := range sm.signatures {
		if sig.Revoked {
			revokedCount++
		}
		if time.Now().After(sig.ExpiresAt) {
			expiredCount++
		}
		if !sig.Revoked && time.Now().Before(sig.ExpiresAt) {
			verifiedCount++
		}
	}

	stats["revoked_plugins"] = revokedCount
	stats["expired_signatures"] = expiredCount
	stats["verified_plugins"] = verifiedCount

	return stats
}
