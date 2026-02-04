package editors

import (
	"fmt"
	"image/color"
	"sort"
	"strings"

	"ffvi_editor/models/pr"
	"ffvi_editor/settings"
	"ffvi_editor/ui/forms/inputs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	defaultWorldWidth  = 256
	defaultWorldHeight = 256
)

type (
	// MapData is the main widget for map editing
	MapData struct {
		widget.BaseWidget
		search  *widget.Entry
		results *widget.TextGrid
		logic   *mapDataLogic
	}
)

// NewMapData creates a new MapData widget
func NewMapData() *MapData {
	e := &MapData{
		search:  widget.NewEntry(),
		results: widget.NewTextGrid(),
		logic:   newMapDataLogic(),
	}

	e.search.OnChanged = func(s string) {
		e.updateSearchResults(s)
	}
	e.ExtendBaseWidget(e)
	return e
}

// updateSearchResults updates the search results text grid
func (e *MapData) updateSearchResults(s string) {
	s = strings.ToLower(s)
	if len(s) < 2 {
		e.results.SetText("")
		return
	}

	var sb strings.Builder
	if len(s) > 2 {
		var found []string
		for k, v := range mapLookup {
			if strings.Contains(strings.ToLower(v), s) {
				found = append(found, fmt.Sprintf("%3d - %s\n", k, v))
			}
		}
		sort.Strings(found)
		for _, v := range found {
			sb.WriteString(v)
		}
	}
	e.results.SetText(sb.String())
}

// CreateRenderer creates the renderer for the MapData widget
func (e *MapData) CreateRenderer() fyne.WidgetRenderer {
	data := pr.GetMapData()
	transport := pr.Transportations

	// --- Interactive Map Section ---
	worldSelect := widget.NewSelect([]string{"World of Balance", "World of Ruin"}, nil)
	worldSelect.SetSelected("World of Balance")

	// --- Zoomable Map Section ---
	bg := canvas.NewRectangle(&color.NRGBA{R: 220, G: 220, B: 255, A: 255})
	bg.SetMinSize(fyne.NewSize(1, 1))

	coordLabel := widget.NewLabel("")
	snapToLandmark := widget.NewCheck("Snap to nearest landmark", nil)
	snapToLandmark.SetChecked(false)

	locationSearch := widget.NewEntry()
	locationSearch.SetPlaceHolder("Filter locations...")
	selectedLocationName := ""
	calibratingName := ""
	calibLabel := widget.NewLabel("")

	mapWidget := NewMapClickWidget(getMapImagePath(1), nil, nil)
	mapWidget.SetMinSize(fyne.NewSize(500, 500)) // Fixed minimum size to prevent layout jumps
	mapContainer := container.NewMax(bg, mapWidget)

	getWorld := func() int {
		if worldSelect.Selected == "World of Ruin" {
			return 2
		}
		return 1
	}

	getGameSize := func() (int, int) {
		gameWidth, gameHeight := data.Gps.Width, data.Gps.Height
		if gameWidth <= 0 {
			gameWidth = defaultWorldWidth
		}
		if gameHeight <= 0 {
			gameHeight = defaultWorldHeight
		}
		return gameWidth, gameHeight
	}

	toImageMarker := func(gameX, gameY float64) (float64, float64, bool) {
		gameWidth, gameHeight := getGameSize()
		if gameWidth <= 1 || gameHeight <= 1 || mapWidget.imgWidth <= 1 || mapWidget.imgHeight <= 1 {
			fmt.Printf("[DEBUG toImageMarker] Invalid dimensions: game=%dx%d, img=%dx%d\n", gameWidth, gameHeight, mapWidget.imgWidth, mapWidget.imgHeight)
			return -1, -1, false
		}
		mx := float64(mapWidget.imgWidth-1) * gameX / float64(gameWidth-1)
		my := float64(mapWidget.imgHeight-1) * gameY / float64(gameHeight-1)
		return mx, my, true
	}

	updateMarker := func() {
		gameWidth, gameHeight := getGameSize()
		world := getWorld()

		mapWidget.showMarker = false
		mapWidget.markerX, mapWidget.markerY = -1, -1
		if mapWidget.imgWidth <= 0 || mapWidget.imgHeight <= 0 {
			mapWidget.Refresh()
			return
		}
		validPlayer := data.Player.X >= 0 && data.Player.Y >= 0 && data.Player.X <= float64(gameWidth-1) && data.Player.Y <= float64(gameHeight-1)
		if !validPlayer {
			mapWidget.Refresh()
			return
		}
		if (world == 1 && data.MapID != 1) || (world == 2 && data.MapID != 2) {
			mapWidget.Refresh()
			return
		}

		mx, my, ok := toImageMarker(data.Player.X, data.Player.Y)
		if !ok {
			mapWidget.Refresh()
			return
		}
		mapWidget.markerX = mx
		mapWidget.markerY = my
		mapWidget.showMarker = true
		mapWidget.Refresh()
	}

	teleportTo := func(name string, x, y float64) {
		world := getWorld()

		if world == 2 {
			data.MapID = 2
		} else {
			data.MapID = 1
		}
		data.Player.X = x
		data.Player.Y = y
		data.Player.Z = 0
		data.PlayerDirection = 0
		updateMarker()
		if name != "" {
			coordLabel.SetText(fmt.Sprintf("%s (%.1f, %.1f) on %s", name, x, y, worldSelect.Selected))
		} else {
			coordLabel.SetText(fmt.Sprintf("(%.1f, %.1f) on %s", x, y, worldSelect.Selected))
		}
	}

	locationsList := widget.NewList(
		func() int { return len(e.logic.getFilteredLocations()) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, o fyne.CanvasObject) {
			locs := e.logic.getFilteredLocations()
			if id < 0 || id >= len(locs) {
				return
			}
			o.(*widget.Label).SetText(locs[id].Name)
		},
	)

	locationsList.OnSelected = func(id widget.ListItemID) {
		locs := e.logic.getFilteredLocations()
		if id < 0 || id >= len(locs) {
			return
		}
		loc := locs[id]
		selectedLocationName = loc.Name


		world := getWorld()
		sm := settings.NewManager("")
		_ = sm.Load()
		x, y, ok := sm.GetMapPoint(world, loc.Name)
		if ok && mapWidget.imgWidth > 1 && mapWidget.imgHeight > 1 {
			calibLabel.SetText("Using saved map point (image-calified).")
			gameWidth, gameHeight := getGameSize()
			gx := float64(x) * float64(gameWidth-1) / float64(mapWidget.imgWidth-1)
			gy := float64(y) * float64(gameHeight-1) / float64(mapWidget.imgHeight-1)

			teleportTo(loc.Name, gx, gy)
			return
		}

		if loc.HasFallback {
			calibLabel.SetText("No point set for this location (using fallback).")

			teleportTo(loc.Name, loc.X, loc.Y)
			return
		}

		calibLabel.SetText("No point set for this location. Click 'Set Point', then click on the map.")
	}

	locationSearch.OnChanged = func(_ string) {
		e.logic.rebuildLocations(getWorld(), locationSearch.Text)
		locationsList.Refresh()
	}

	var updateMap func()
	currentMapPath := ""
	updateMap = func() {
		world := getWorld()
		nextMapPath := getMapImagePath(world)

		if currentMapPath != nextMapPath {
			mapWidget.SetImagePath(nextMapPath)
			mapWidget.ResetView()
			currentMapPath = nextMapPath
		}

		e.logic.rebuildLocations(world, locationSearch.Text)
		locationsList.Refresh()
		updateMarker()
		mapContainer.Refresh()
	}

	mapWidget.onClick = func(x, y int) {
		gameWidth, gameHeight := getGameSize()
		if mapWidget.imgWidth <= 0 || mapWidget.imgHeight <= 0 {
			return
		}

		world := getWorld()

		// If we're calibrating a named point, store it and teleport there
		if calibratingName != "" {
			name := calibratingName
			calibratingName = ""
			calibLabel.SetText("")
			sm := settings.NewManager("")
			_ = sm.Load()
			sm.SetMapPoint(world, name, float64(x), float64(y))
			if gameWidth > 1 && gameHeight > 1 && mapWidget.imgWidth > 1 && mapWidget.imgHeight > 1 {
				gx := float64(x) * float64(gameWidth-1) / float64(mapWidget.imgWidth-1)
				gy := float64(y) * float64(gameHeight-1) / float64(mapWidget.imgHeight-1)
				teleportTo(name, gx, gy)
			}
			return
		}

		gx := float64(x) * float64(gameWidth-1) / float64(mapWidget.imgWidth-1)
		gy := float64(y) * float64(gameHeight-1) / float64(mapWidget.imgHeight-1)

		if gx < 0 {
			gx = 0
		} else if gx > float64(gameWidth-1) {
			gx = float64(gameWidth - 1)
		}
		if gy < 0 {
			gy = 0
		} else if gy > float64(gameHeight-1) {
			gy = float64(gameHeight - 1)
		}

		targetX, targetY := gx, gy
		targetName := ""
		if snapToLandmark.Checked {
			// Prefer calibrated points (pixel-perfect on the map image)
			sm := settings.NewManager("")
			_ = sm.Load()
			points := sm.GetAllMapPoints(world)
			if len(points) > 0 {
				minDist := 1e18
				bestName := ""
				bestPx, bestPy := 0, 0
				for name, p := range points {
					dx := p["x"] - float64(x)
					dy := p["y"] - float64(y)
					dist := dx*dx + dy*dy
					if dist < minDist {
						minDist = dist
						bestName = name
						bestPx = int(p["x"])
						bestPy = int(p["y"])
					}
				}
				if bestName != "" && mapWidget.imgWidth > 1 && mapWidget.imgHeight > 1 && gameWidth > 1 && gameHeight > 1 {
					targetName = bestName
					targetX = float64(bestPx) * float64(gameWidth-1) / float64(mapWidget.imgWidth-1)
					targetY = float64(bestPy) * float64(gameHeight-1) / float64(mapWidget.imgHeight-1)
				}
			} else {
				// Fallback: built-in landmark coordinates
				name, nx, ny, found := findNearestLandmark(world, gx, gy)
				if found {
					targetX = nx
					targetY = ny
					targetName = name
				}
			}
		}

		teleportTo(targetName, targetX, targetY)
	}

	updateMap()
	worldSelect.OnChanged = func(_ string) {
		selectedLocationName = ""
		calibratingName = ""
		calibLabel.SetText("")
		locationsList.UnselectAll()
		updateMap()
	}

	resetViewBtn := widget.NewButton("Reset View", func() { mapWidget.ResetView() })
	setPointBtn := widget.NewButton("Set Point", func() {
		if selectedLocationName == "" {
			return
		}
		calibratingName = selectedLocationName
		calibLabel.SetText(fmt.Sprintf("Click on the map to set: %s", calibratingName))
	})
	clearPointBtn := widget.NewButton("Clear Point", func() {
		if selectedLocationName == "" {
			return
		}
		sm := settings.NewManager("")
		_ = sm.Load()
		sm.ClearMapPoint(getWorld(), selectedLocationName)
		calibLabel.SetText("Cleared saved map point (using fallback).")
	})
	addLocationBtn := widget.NewButton("Add", func() {
		name := strings.TrimSpace(locationSearch.Text)
		if name == "" {
			return
		}
		sm := settings.NewManager("")
		_ = sm.Load()
		sm.AddMapLocation(getWorld(), name)
		locationSearch.SetText("")
		e.logic.rebuildLocations(getWorld(), locationSearch.Text)
		locationsList.Refresh()
	})

	locationsScroll := container.NewVScroll(locationsList)
	locationsScroll.SetMinSize(fyne.NewSize(240, 400))
	locationsCard := widget.NewCard(
		"Locations",
		"",
		container.NewBorder(
			container.NewBorder(nil, nil, nil, addLocationBtn, locationSearch),
			container.NewHBox(setPointBtn, clearPointBtn, calibLabel),
			nil,
			nil,
			locationsScroll,
		),
	)

	mapSplit := container.NewHSplit(locationsCard, mapContainer)
	mapSplit.Offset = 0.28

	mapSection := container.NewVBox(
		widget.NewLabel("Click anywhere on the map to teleport (or pick a location on the left)."),
		container.NewHBox(worldSelect, resetViewBtn, snapToLandmark),
		mapSplit,
		coordLabel,
	)

	// --- Existing Data Cards Section ---
	cards := make([]fyne.CanvasObject, 0, 4)
	cards = append(cards, widget.NewCard("Player", "", container.NewVBox(
		inputs.NewLabeledIntEntryWithHint("Map ID", inputs.NewIntEntryWithData(&data.MapID), inputs.HintArgs{
			Align: inputs.NewAlign(fyne.TextAlignTrailing),
			Hints: &mapLookup,
		}),
		inputs.NewLabeledEntry("Position X", inputs.NewFloatEntryWithData(&data.Player.X)),
		inputs.NewLabeledEntry("Position Y", inputs.NewFloatEntryWithData(&data.Player.Y)),
		inputs.NewLabeledEntry("Position Z", inputs.NewFloatEntryWithData(&data.Player.Z)),
		inputs.NewLabeledEntry("Facing Direction", inputs.NewIntEntryWithData(&data.PlayerDirection)),
	)))
	cards = append(cards,
		widget.NewCard("GPS", "", container.NewVBox(
			inputs.NewLabeledEntry("World", e.newWorldSelect(&data.Gps.MapID)),
			inputs.NewLabeledEntry("Area ID", inputs.NewIntEntryWithData(&data.Gps.AreaID)),
			inputs.NewLabeledEntry("ID", inputs.NewIntEntryWithData(&data.Gps.GpsID)),
			inputs.NewLabeledEntry("Width", inputs.NewIntEntryWithData(&data.Gps.Width)),
			inputs.NewLabeledEntry("Height", inputs.NewIntEntryWithData(&data.Gps.Height)),
		)))
	if len(transport) >= 5 {
		bj := transport[3]
		cards = append(cards,
			widget.NewCard("Blackjack", "", container.NewVBox(
				inputs.NewLabeledEntry("Enabled", widget.NewCheckWithData("", binding.BindBool(&bj.Enabled))),
				inputs.NewLabeledEntry("World", e.newWorldSelectUnset(&bj.MapID)),
				inputs.NewLabeledEntry("Position X", inputs.NewFloatEntryWithData(&bj.Position.X)),
				inputs.NewLabeledEntry("Position Y", inputs.NewFloatEntryWithData(&bj.Position.Y)),
				inputs.NewLabeledEntry("Position Z", inputs.NewFloatEntryWithData(&bj.Position.Z)),
				inputs.NewLabeledEntry("Facing Direction", inputs.NewIntEntryWithData(&bj.Direction)),
			)))
		f := transport[4]
		cards = append(cards,
			widget.NewCard("Falcon", "", container.NewVBox(
				inputs.NewLabeledEntry("Enabled", widget.NewCheckWithData("", binding.BindBool(&f.Enabled))),
				inputs.NewLabeledEntry("World", e.newWorldSelectUnset(&f.MapID)),
				inputs.NewLabeledEntry("Position X", inputs.NewFloatEntryWithData(&f.Position.X)),
				inputs.NewLabeledEntry("Position Y", inputs.NewFloatEntryWithData(&f.Position.Y)),
				inputs.NewLabeledEntry("Position Z", inputs.NewFloatEntryWithData(&f.Position.Z)),
				inputs.NewLabeledEntry("Facing Direction", inputs.NewIntEntryWithData(&f.Direction)),
			)))
	}

	// Combine mapSection and cards into a single slice for NewVBox
	allVBoxItems := []fyne.CanvasObject{mapSection}
	allVBoxItems = append(allVBoxItems, cards...)

	leftCol := container.NewVScroll(container.NewVBox(allVBoxItems...))
	mapIDsCol := container.NewBorder(
		widget.NewLabel("Map IDs"), nil, nil, nil,
		container.NewVScroll(mapsTextBox),
	)
	findByNameCol := container.NewBorder(
		inputs.NewLabeledEntry("Find By Name:", e.search), nil, nil, nil,
		container.NewVScroll(e.results),
	)

	rightSplit := container.NewHSplit(mapIDsCol, findByNameCol)
	rightSplit.Offset = 0.40
	mainSplit := container.NewHSplit(leftCol, rightSplit)
	mainSplit.Offset = 0.35

	return widget.NewSimpleRenderer(mainSplit)
}

func (e *MapData) newWorldSelect(i *int) *widget.Select {
	s := widget.NewSelect([]string{"Balance", "Ruin"}, func(s string) {
		if s == "Balance" {
			*i = 1
		} else {
			*i = 2
		}
	})
	if *i == 1 {
		s.SetSelected("Balance")
	} else {
		s.SetSelected("Ruin")
	}
	return s
}

func (e *MapData) newWorldSelectUnset(i *int) *widget.Select {
	s := widget.NewSelect([]string{"-", "Balance", "Ruin"}, func(s string) {
		if s == "Balance" {
			*i = 1
		} else if s == "Ruin" {
			*i = 2
		} else {
			*i = -1
		}
	})
	if *i == 1 {
		s.SetSelected("Balance")
	} else if *i == 2 {
		s.SetSelected("Ruin")
	} else {
		s.SetSelected("-")
	}
	return s
}

// teleportDialog holds the state for the teleport dialog
type teleportDialog struct {
	data           *pr.MapData
	worldSelect    *widget.Select
	snapToLandmark *widget.Check
	coordLabel     *widget.Label
	mapWidget      *MapClickWidget
}

// getWorld returns the current world (1 or 2)
func (td *teleportDialog) getWorld() int {
	if td.worldSelect.Selected == "World of Ruin" {
		return 2
	}
	return 1
}

// getGameSize returns the game world dimensions
func (td *teleportDialog) getGameSize() (int, int) {
	gameWidth, gameHeight := td.data.Gps.Width, td.data.Gps.Height
	if gameWidth <= 0 {
		gameWidth = defaultWorldWidth
	}
	if gameHeight <= 0 {
		gameHeight = defaultWorldHeight
	}
	return gameWidth, gameHeight
}

// isValidMarker returns true if the marker should be shown
func (td *teleportDialog) isValidMarker(gameWidth, gameHeight int) bool {
	if td.mapWidget.imgWidth <= 0 || td.mapWidget.imgHeight <= 0 {
		return false
	}
	validPlayer := td.data.Player.X >= 0 && td.data.Player.Y >= 0 &&
		td.data.Player.X <= float64(gameWidth-1) && td.data.Player.Y <= float64(gameHeight-1)
	if !validPlayer {
		return false
	}
	world := td.getWorld()
	if (world == 1 && td.data.MapID != 1) || (world == 2 && td.data.MapID != 2) {
		return false
	}
	if gameWidth <= 1 || gameHeight <= 1 || td.mapWidget.imgWidth <= 1 || td.mapWidget.imgHeight <= 1 {
		return false
	}
	return true
}

// updateMarker updates the player marker position on the map
func (td *teleportDialog) updateMarker() {
	gameWidth, gameHeight := td.getGameSize()

	td.mapWidget.showMarker = false
	td.mapWidget.markerX, td.mapWidget.markerY = -1, -1

	if !td.isValidMarker(gameWidth, gameHeight) {
		td.mapWidget.Refresh()
		return
	}

	td.mapWidget.markerX = float64(td.mapWidget.imgWidth-1) * td.data.Player.X / float64(gameWidth-1)
	td.mapWidget.markerY = float64(td.mapWidget.imgHeight-1) * td.data.Player.Y / float64(gameHeight-1)
	td.mapWidget.showMarker = true
	td.mapWidget.Refresh()
}

// teleportTo sets the player position and updates the UI
func (td *teleportDialog) teleportTo(name string, x, y float64) {
	world := td.getWorld()
	if world == 2 {
		td.data.MapID = 2
	} else {
		td.data.MapID = 1
	}
	td.data.Player.X = x
	td.data.Player.Y = y
	td.data.Player.Z = 0
	td.data.PlayerDirection = 0
	td.updateMarker()

	if name != "" {
		td.coordLabel.SetText(fmt.Sprintf("%s (%.1f, %.1f) on %s", name, x, y, td.worldSelect.Selected))
	} else {
		td.coordLabel.SetText(fmt.Sprintf("(%.1f, %.1f) on %s", x, y, td.worldSelect.Selected))
	}
}

// handleMapClick processes a click on the map
func (td *teleportDialog) handleMapClick(x, y int) {
	gameWidth, gameHeight := td.getGameSize()
	if td.mapWidget.imgWidth <= 1 || td.mapWidget.imgHeight <= 1 || gameWidth <= 1 || gameHeight <= 1 {
		return
	}

	gx := float64(x) * float64(gameWidth-1) / float64(td.mapWidget.imgWidth-1)
	gy := float64(y) * float64(gameHeight-1) / float64(td.mapWidget.imgHeight-1)

	// Clamp coordinates
	gx = clampFloat(gx, 0, float64(gameWidth-1))
	gy = clampFloat(gy, 0, float64(gameHeight-1))

	targetX, targetY := gx, gy
	targetName := ""

	if td.snapToLandmark.Checked {
		world := td.getWorld()
		name, nx, ny, found := findNearestLandmark(world, gx, gy)
		if found {
			targetX = nx
			targetY = ny
			targetName = name
		}
	}

	td.teleportTo(targetName, targetX, targetY)
}

// clampFloat clamps a float64 value between min and max
func clampFloat(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (e *MapData) showTeleportDialog(data *pr.MapData) {
	td := &teleportDialog{
		data:           data,
		worldSelect:    widget.NewSelect([]string{"World of Balance", "World of Ruin"}, nil),
		snapToLandmark: widget.NewCheck("Snap to nearest landmark", nil),
		coordLabel:     widget.NewLabel(""),
		mapWidget:      NewMapClickWidget(getMapImagePath(1), nil, nil),
	}

	td.worldSelect.SetSelected("World of Balance")
	td.snapToLandmark.SetChecked(false)

	td.mapWidget.onClick = func(x, y int) { td.handleMapClick(x, y) }

	updateMap := func() {
		td.mapWidget.SetImagePath(getMapImagePath(td.getWorld()))
		td.updateMarker()
	}

	updateMap()
	td.worldSelect.OnChanged = func(_ string) { updateMap() }

	d := dialog.NewCustom("Click to Teleport", "Close",
		container.NewVBox(
			widget.NewLabel("Select World and click a location on the map to teleport."),
			td.worldSelect,
			td.snapToLandmark,
			td.mapWidget,
			td.coordLabel,
		), fyne.CurrentApp().Driver().AllWindows()[0])
	d.Show()
}
