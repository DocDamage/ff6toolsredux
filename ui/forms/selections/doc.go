// Package selections provides selection widgets.
//
// The selections package contains:
//   - Character selection lists
//   - Item selection grids
//   - Spell selection trees
//   - Multi-selection lists
//
// These widgets provide intuitive interfaces for
// selecting game entities.
//
// Example usage:
//
//	// Create character selection
//	selector := selections.NewCharacterSelector()
//	selector.SetSelected(charID)
//	selector.OnChanged = func(id int) {
//	    // Handle selection change
//	}
//
//	// Create item grid
//	grid := selections.NewItemGrid()
//	grid.SetItems(items)
package selections
