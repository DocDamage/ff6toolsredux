package selections

import (
	"fmt"

	"ffvi_editor/global"
	"ffvi_editor/models/pr"
	"ffvi_editor/ui/forms/editors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type (
	Characters struct {
		widget.BaseWidget
		top    fyne.CanvasObject
		middle *container.Scroll
	}
)

func NewCharacters() *Characters {
	s := &Characters{}
	s.ExtendBaseWidget(s)

	// Initialize middle section FIRST (before dropdown callback)
	s.middle = container.NewScroll(container.NewStack())

	// Character selector
	charDropdown := widget.NewSelect(pr.CharacterNamesHumanSelect(), func(name string) {
		global.Log("[Characters] Selected character: '%s'", name)
		fmt.Printf("[DEBUG Characters] Selected character: '%s'\n", name)
		if name == "" {
			global.Log("[Characters] Empty name, ignoring")
			fmt.Println("[DEBUG Characters] Empty name, ignoring")
			return
		}
		content := container.NewStack()
		c := pr.GetCharacter(name)
		if c == nil {
			global.Log("[Characters] ERROR: Character '%s' not found", name)
			fmt.Printf("[DEBUG Characters] ERROR: Character '%s' not found\n", name)
			content.Add(widget.NewLabel("Error: Character '" + name + "' not found"))
			s.middle.Content = content
			s.middle.Refresh()
			return
		}
		global.Log("[Characters] Found character: ID=%d, RootName=%s", c.ID, c.RootName)
		fmt.Printf("[DEBUG Characters] Found character: ID=%d, RootName=%s\n", c.ID, c.RootName)
		content.Add(container.NewAppTabs(
			container.NewTabItem("Stats", editors.NewCharacter(c)),
			container.NewTabItem("Magic", editors.NewMagic(c)),
			container.NewTabItem("Equipment", editors.NewEquipment(c)),
			container.NewTabItem("Commands", editors.NewCommands(c)),
		))
		s.middle.Content = content
		s.middle.Refresh()
		global.Log("[Characters] AppTabs added to UI")
		fmt.Println("[DEBUG Characters] AppTabs added to UI")
	})

	// Fixed top section with padding
	s.top = container.NewVBox(
		container.NewHBox(widget.NewLabel("Character:"), charDropdown),
		container.NewHBox(
			widget.NewButton("Max All Characters", func() {
				for _, name := range pr.CharacterNamesHumanSelect() {
					c := pr.GetCharacter(name)
					c.Level = 99
					c.Exp = 1848184
					c.HP.Max = 9999
					c.HP.Current = 9999
					c.MP.Max = 999
					c.MP.Current = 999
					c.Vigor = 255
					c.Stamina = 255
					c.Speed = 255
					c.Magic = 255
				}
			}),
			widget.NewButton("Heal All", func() {
				for _, name := range pr.CharacterNamesHumanSelect() {
					c := pr.GetCharacter(name)
					c.HP.Current = c.HP.Max
					c.MP.Current = c.MP.Max
				}
			}),
			widget.NewButton("Reset All", func() {
				for _, name := range pr.CharacterNamesHumanSelect() {
					c := pr.GetCharacter(name)
					c.Level = 1
					c.Exp = 0
					c.HP.Max = 0
					c.HP.Current = 0
					c.MP.Max = 0
					c.MP.Current = 0
					c.Vigor = 0
					c.Stamina = 0
					c.Speed = 0
					c.Magic = 0
				}
			}),
		),
	)

	// NOW set the initial selection (after middle is initialized)
	charDropdown.SetSelectedIndex(0)

	return s
}

func (s *Characters) CreateRenderer() fyne.WidgetRenderer {
	// Use Border layout with scrollable content area
	return widget.NewSimpleRenderer(
		container.NewBorder(s.top, nil, nil, nil, s.middle),
	)
}
