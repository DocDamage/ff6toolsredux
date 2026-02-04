package selections

import (
	"fmt"

	"ffvi_editor/global"
	"ffvi_editor/ui/forms/editors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type (
	Editor struct {
		widget.BaseWidget
	}
)

func NewEditor() *Editor {
	s := &Editor{}
	s.ExtendBaseWidget(s)
	return s
}

func (s *Editor) CreateRenderer() fyne.WidgetRenderer {
	fmt.Println("[DEBUG Editor] Creating main editor tabs...")
	global.Log("[Editor] Creating main editor tabs...")
	
	global.Log("[Editor] Creating Characters tab...")
	characters := NewCharacters()
	
	global.Log("[Editor] Creating Inventory tab...")
	inventory := NewInventory()
	
	global.Log("[Editor] Creating Skills tab...")
	skills := editors.NewSkills()
	
	global.Log("[Editor] Creating Espers tab...")
	espers := editors.NewEsper()
	
	global.Log("[Editor] Creating Party tab...")
	party := editors.NewParty()
	
	global.Log("[Editor] Creating Map tab...")
	m := editors.NewMapData()
	
	global.Log("[Editor] Creating Veldt tab...")
	veldt := editors.NewVeldt()
	
	tabs := container.NewAppTabs(
		container.NewTabItem("Characters", characters),
		container.NewTabItem("Inventory", inventory),
		container.NewTabItem("Skills", skills),
		container.NewTabItem("Espers", espers),
		container.NewTabItem("Party", party),
		container.NewTabItem("Map", m),
		container.NewTabItem("Veldt", veldt),
	)
	global.Log("[Editor] Main editor tabs created")
	fmt.Println("[DEBUG Editor] Main editor tabs created")
	return widget.NewSimpleRenderer(tabs)
}
