// Package file provides low-level file I/O operations for save files.
//
// This package handles the reading and writing of raw save file bytes,
// including support for different save file formats and encryption.
//
// Supported Formats:
//
//	JSON      - Plain text JSON saves
//	Encrypted - Encrypted save files
//	Steam     - Steam cloud saves
//
// Functions:
//
//	LoadFile(path string, saveType SaveFileType) ([]byte, bool, error)
//	SaveFile(data []byte, path string, trimmed bool, saveType SaveFileType) error
//
// The package automatically handles format detection and conversion
// between different save file types.
package file
