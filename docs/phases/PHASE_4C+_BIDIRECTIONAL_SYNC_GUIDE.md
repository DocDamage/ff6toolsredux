# Bidirectional Sync Guide

**Date:** January 16, 2026  
**Status:** Phase 4C+ Complete

---

## Overview

Complete bidirectional sync system design for FF6 Save Editor marketplace, enabling seamless synchronization between local plugin state and GitHub registry.

## Sync Architecture

### Core Sync Flow

```
┌────────────────┐         ┌──────────────┐
│  Local State   │◄───────►│   Registry   │
│  - Plugins     │  Sync   │  - Plugins   │
│  - Ratings     │ Manager │  - Ratings   │
│  - Config      │         │  - Metadata  │
└────────────────┘         └──────────────┘
       ▲                            ▲
       │                            │
       └────── Conflict ────────────┘
            Resolution
```

## Sync Operations

### Pull Sync (Registry → Local)
Fetches latest changes from registry and merges locally.

**Process:**
1. Fetch remote state
2. Detect changes
3. Merge changes
4. Update local state
5. Verify consistency

### Push Sync (Local → Registry)
Sends local changes to registry.

**Process:**
1. Collect local changes
2. Validate changes
3. Send to registry
4. Verify upload
5. Update sync state

### Bidirectional Merge
Merges changes from both directions with conflict resolution.

**Strategies:**
- LATEST_WINS: Remote wins
- LOCAL_WINS: Local wins
- REMOTE_WINS: Registry wins
- MANUAL: User decides

## Conflict Resolution

### Detection
Identifies when both sides have changes to same item.

### Resolution Strategies
1. **Automatic Merge** - Non-conflicting changes
2. **Three-way Merge** - Compare with common ancestor
3. **Manual Resolution** - User chooses
4. **Rollback** - Revert to previous state

## State Management

### Tracking
- Version numbers
- Timestamps
- Change hashes
- Sync markers

### Consistency
- Pre-sync verification
- Post-sync validation
- Atomic operations
- Rollback capability

## Implementation Patterns

### Pull Operation
```go
syncer.Pull()  // Fetch & merge from registry
```

### Push Operation
```go
syncer.Push()  // Send changes to registry
```

### Merge with Strategy
```go
result := syncer.Merge(ConflictStrategy.LATEST_WINS)
```

### Dry Run
```go
preview := syncer.PreviewSync()  // See what would change
```

## Advanced Scenarios

### Offline Support
- Queue changes while offline
- Sync when connection restored
- Conflict resolution on reconnect

### Large Datasets
- Batch operations
- Incremental sync
- Progress tracking

### Error Recovery
- Retry logic with backoff
- Partial sync rollback
- State repair

---

**For complete details, see other Phase 4C+ documentation files.**
