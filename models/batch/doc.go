// Package batch provides batch operation definitions.
//
// The batch package handles:
//   - Batch operation types and configurations
//   - Operation sequencing
//   - Batch result tracking
//
// Supported operations:
//   - Max all character stats
//   - Max all items
//   - Learn all magic
//   - Custom Lua script execution
//
// Example usage:
//
//	// Create batch operation
//	op := &batch.Operation{
//	    Type: batch.OpMaxStats,
//	    Target: "all",
//	}
//
//	// Execute batch
//	result, err := batch.Execute(save, op)
//	if err != nil {
//	    return err
//	}
package batch
