package pr

import (
	"testing"
)

// TestInventoryAddNeeded tests adding items to inventory
func TestInventoryAddNeeded(t *testing.T) {
	inv := GetInventory()
	inv.Clear()

	// Test adding items to empty slots
	needed := map[int]int{
		1: 5,
		2: 10,
	}

	err := inv.AddNeeded(needed)
	if err != nil {
		t.Errorf("AddNeeded failed: %v", err)
	}

	// Verify items were added
	if row, found := inv.Get(1); !found || row.Count != 5 {
		t.Errorf("Item 1: expected count 5, got %d", row.Count)
	}
	if row, found := inv.Get(2); !found || row.Count != 10 {
		t.Errorf("Item 2: expected count 10, got %d", row.Count)
	}
}

// TestInventoryAddNeededFull tests error when inventory is full
func TestInventoryAddNeededFull(t *testing.T) {
	inv := GetInventory()
	inv.Clear()

	// Fill the inventory
	for i := 0; i < inv.Size; i++ {
		inv.Rows[i] = &Row{ItemID: i + 1, Count: 1}
	}

	// Try to add another item
	needed := map[int]int{999: 1}
	err := inv.AddNeeded(needed)
	if err == nil {
		t.Error("Expected error when inventory is full")
	}
}

// TestInventoryAddNeededExisting tests adding to existing items
func TestInventoryAddNeededExisting(t *testing.T) {
	inv := GetInventory()
	inv.Clear()

	// Set up existing item
	inv.Set(0, Row{ItemID: 1, Count: 5})

	// Add more of the same item (AddNeeded actually marks existing as found and removes from needed)
	// So this test is really about new items being added
	needed := map[int]int{2: 3}
	err := inv.AddNeeded(needed)
	if err != nil {
		t.Errorf("AddNeeded failed: %v", err)
	}

	// Verify item was added
	if row, found := inv.Get(2); !found || row.Count != 3 {
		t.Errorf("Item 2: expected count 3, got %d", row.Count)
	}
}

// TestInventoryGetRowsForPrSave tests filtering for save
func TestInventoryGetRowsForPrSave(t *testing.T) {
	inv := GetInventory()
	inv.Clear()

	// Add various items
	inv.Set(0, Row{ItemID: 1, Count: 5})    // Valid
	inv.Set(1, Row{ItemID: 0, Count: 0})    // Empty - should be filtered
	inv.Set(2, Row{ItemID: 1000, Count: 1}) // Invalid ID - should be filtered
	inv.Set(3, Row{ItemID: 50, Count: 0})   // Zero count - should be filtered
	inv.Set(4, Row{ItemID: 50, Count: 100}) // Valid but count > 99

	rows := inv.GetRowsForPrSave()

	// Should only include valid items
	if len(rows) != 1 {
		t.Errorf("Expected 1 valid row, got %d", len(rows))
	}

	if len(rows) > 0 && rows[0].ItemID != 1 {
		t.Errorf("Expected ItemID 1, got %d", rows[0].ItemID)
	}
}

// TestInventoryGetItemLookup tests the item lookup map
func TestInventoryGetItemLookup(t *testing.T) {
	inv := GetInventory()
	inv.Clear()

	inv.Set(0, Row{ItemID: 1, Count: 5})
	inv.Set(1, Row{ItemID: 2, Count: 10})

	lookup := inv.GetItemLookup()

	if lookup[1] != 5 {
		t.Errorf("Lookup[1]: expected 5, got %d", lookup[1])
	}
	if lookup[2] != 10 {
		t.Errorf("Lookup[2]: expected 10, got %d", lookup[2])
	}
}

// TestInventoryReset tests clearing all items
func TestInventoryReset(t *testing.T) {
	inv := GetInventory()
	inv.Clear()

	inv.Set(0, Row{ItemID: 1, Count: 5})
	inv.Reset()

	for i, row := range inv.Rows {
		if row.ItemID != 0 || row.Count != 0 {
			t.Errorf("Row %d: expected empty, got ItemID=%d, Count=%d", i, row.ItemID, row.Count)
		}
	}
}

// TestInventorySetExpansion tests that Set expands the slice
func TestInventorySetExpansion(t *testing.T) {
	inv := GetInventory()
	inv.Clear()

	// Set an index beyond current size
	inv.Set(300, Row{ItemID: 99, Count: 1})

	if len(inv.Rows) <= 300 {
		t.Error("Expected rows to be expanded")
	}

	if inv.Rows[300].ItemID != 99 {
		t.Errorf("Expected ItemID 99 at index 300, got %d", inv.Rows[300].ItemID)
	}
}

// TestInventoryAddNeededPartial tests AddNeeded with partial existing items
func TestInventoryAddNeededPartial(t *testing.T) {
	inv := GetInventory()
	inv.Clear()

	// Add some existing items
	inv.Set(0, Row{ItemID: 1, Count: 5})
	inv.Set(1, Row{ItemID: 2, Count: 10})

	// Add needed items - some existing, some new
	needed := map[int]int{
		1: 10, // Already have 5, but AddNeeded just marks as found
		2: 20, // Already have 10
		3: 15, // New item
	}

	err := inv.AddNeeded(needed)
	if err != nil {
		t.Errorf("AddNeeded failed: %v", err)
	}

	// Item 3 should be added
	if row, found := inv.Get(3); !found || row.Count != 15 {
		t.Errorf("Item 3: expected count 15, got %d", row.Count)
	}
}

// TestInventoryImportant tests important inventory
func TestInventoryImportant(t *testing.T) {
	inv := GetImportantInventory()
	inv.Clear()

	if inv.Size != 100 {
		t.Errorf("Important inventory size: expected 100, got %d", inv.Size)
	}

	// Test basic operations
	inv.Set(0, Row{ItemID: 1, Count: 1})
	if row, found := inv.Get(1); !found || row.Count != 1 {
		t.Errorf("Important item: expected count 1, got %d", row.Count)
	}
}

// BenchmarkInventoryAddNeeded benchmarks adding items
func BenchmarkInventoryAddNeeded(b *testing.B) {
	inv := &Inventory{Size: 255}
	inv.Clear()

	needed := map[int]int{1: 5, 2: 10, 3: 15}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inv.Clear()
		inv.AddNeeded(needed)
	}
}
