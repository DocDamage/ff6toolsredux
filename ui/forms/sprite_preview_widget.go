package forms

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"ffvi_editor/io"
	"ffvi_editor/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SpritePreviewWidget displays a sprite with palette
type SpritePreviewWidget struct {
	widget.BaseWidget

	sprite         *models.FF6Sprite
	currentFrame   int
	scale          float32
	showGridToggle *widget.Check
	showGrid       bool

	// Rendered sprite
	spriteImage *canvas.Image
	container   *fyne.Container

	// Controls
	frameLabel  *widget.Label
	frameSlider *widget.Slider
	scaleLabel  *widget.Label
	scaleSlider *widget.Slider
	frameNext   *widget.Button
	framePrev   *widget.Button

	frameImages []image.Image
}

// NewSpritePreviewWidget creates a new sprite preview widget
func NewSpritePreviewWidget(sprite *models.FF6Sprite) *SpritePreviewWidget {
	p := &SpritePreviewWidget{
		sprite:       sprite,
		currentFrame: 0,
		scale:        2.0,
		showGrid:     false,
	}

	p.frameImages = decodeSpriteFrames(sprite)

	p.ExtendBaseWidget(p)
	return p
}

// CreateRenderer creates the renderer for the widget
func (p *SpritePreviewWidget) CreateRenderer() fyne.WidgetRenderer {
	p.showGridToggle = widget.NewCheck("Show Grid", func(checked bool) {
		p.showGrid = checked
		p.updatePreview()
	})

	p.frameLabel = widget.NewLabel("Frame 0/0")
	p.frameLabel.Alignment = fyne.TextAlignCenter

	frameCount := len(p.frameImages)
	if frameCount <= 0 {
		frameCount = 1
	}

	p.frameSlider = widget.NewSlider(0, float64(frameCount-1))
	p.frameSlider.OnChanged = func(f float64) {
		frames := len(p.frameImages)
		if frames <= 0 {
			frames = 1
		}

		p.currentFrame = int(f + 0.5)
		if p.currentFrame < 0 {
			p.currentFrame = 0
		}
		if p.currentFrame >= frames {
			p.currentFrame = frames - 1
		}
		p.frameLabel.SetText(fmt.Sprintf("Frame %d / %d", p.currentFrame+1, frames))
		p.updatePreview()
	}
	p.frameSlider.Step = 1

	p.framePrev = widget.NewButton("◀", func() {
		if p.currentFrame > 0 {
			p.currentFrame--
			p.frameSlider.SetValue(float64(p.currentFrame))
			p.updatePreview()
		}
	})

	p.frameNext = widget.NewButton("▶", func() {
		frames := len(p.frameImages)
		if frames <= 0 {
			frames = 1
		}
		if p.currentFrame < frames-1 {
			p.currentFrame++
			p.frameSlider.SetValue(float64(p.currentFrame))
			p.updatePreview()
		}
	})

	p.scaleLabel = widget.NewLabel("1x")
	p.scaleLabel.Alignment = fyne.TextAlignCenter

	p.scaleSlider = widget.NewSlider(1, 8)
	p.scaleSlider.OnChanged = func(f float64) {
		p.scale = float32(f)
		p.scaleLabel.SetText(fmt.Sprintf("%dx", int(f)))
		p.updatePreview()
	}
	p.scaleSlider.SetValue(2)
	p.scaleSlider.Step = 1

	// Build sprite image
	p.spriteImage = canvas.NewImageFromImage(p.renderSprite())
	p.frameLabel.SetText(fmt.Sprintf("Frame %d / %d", p.currentFrame+1, frameCount))

	// Controls container
	controlsBox := container.NewVBox(
		widget.NewLabel("Frame Controls"),
		p.frameLabel,
		container.NewHBox(
			p.framePrev,
			p.frameSlider,
			p.frameNext,
		),
		widget.NewSeparator(),
		widget.NewLabel("Display"),
		container.NewBorder(nil, nil, widget.NewLabel("Scale:"), nil, p.scaleSlider),
		p.scaleLabel,
		p.showGridToggle,
	)

	p.container = container.NewBorder(
		controlsBox,
		nil,
		nil,
		nil,
		p.spriteImage,
	)

	return widget.NewSimpleRenderer(p.container)
}

// updatePreview updates the preview display
func (p *SpritePreviewWidget) updatePreview() {
	if p.spriteImage == nil {
		return
	}
	p.spriteImage.Image = p.renderSprite()
	p.spriteImage.Refresh()
}

// renderSprite renders the current frame with the palette
func (p *SpritePreviewWidget) renderSprite() image.Image {
	if len(p.frameImages) == 0 {
		// Return empty image
		return image.NewRGBA(image.Rect(0, 0, 128, 128))
	}

	if p.currentFrame < 0 || p.currentFrame >= len(p.frameImages) {
		p.currentFrame = 0
	}

	frame := p.frameImages[p.currentFrame]
	srcBounds := frame.Bounds()

	// Calculate scaled dimensions
	scaledW := srcBounds.Dx() * int(p.scale)
	scaledH := srcBounds.Dy() * int(p.scale)

	// Create output image
	dst := image.NewRGBA(image.Rect(0, 0, scaledW, scaledH))

	// Fill with background (magenta for transparency)
	bgColor := color.RGBA{R: 255, G: 0, B: 255, A: 255}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	// Draw palette colors
	// Scale and draw sprite
	for y := 0; y < srcBounds.Dy(); y++ {
		for x := 0; x < srcBounds.Dx(); x++ {
			c := frame.At(srcBounds.Min.X+x, srcBounds.Min.Y+y)
			r, g, b, a := c.RGBA()

			if a > 0 { // Only draw non-transparent pixels
				// Fill scaled pixel
				for sy := 0; sy < int(p.scale); sy++ {
					for sx := 0; sx < int(p.scale); sx++ {
						dx := x*int(p.scale) + sx
						dy := y*int(p.scale) + sy
						if dx < scaledW && dy < scaledH {
							dst.SetRGBA(dx, dy, color.RGBA{
								R: uint8(r >> 8),
								G: uint8(g >> 8),
								B: uint8(b >> 8),
								A: uint8(a >> 8),
							})
						}
					}
				}
			}
		}
	}

	// Draw grid if enabled
	if p.showGrid {
		p.drawGrid(dst, int(p.scale))
	}

	return dst
}

// drawGrid draws a tile grid on the sprite
func (p *SpritePreviewWidget) drawGrid(img *image.RGBA, scale int) {
	gridColor := color.RGBA{R: 0, G: 255, B: 0, A: 128}
	tileSize := 8 * scale

	bounds := img.Bounds()

	// Vertical lines
	for x := 0; x < bounds.Dx(); x += tileSize {
		for y := 0; y < bounds.Dy(); y++ {
			img.SetRGBA(x, y, gridColor)
		}
	}

	// Horizontal lines
	for y := 0; y < bounds.Dy(); y += tileSize {
		for x := 0; x < bounds.Dx(); x++ {
			img.SetRGBA(x, y, gridColor)
		}
	}
}

// SetSprite updates the sprite to display
func (p *SpritePreviewWidget) SetSprite(sprite *models.FF6Sprite) {
	p.sprite = sprite
	p.frameImages = decodeSpriteFrames(sprite)

	p.currentFrame = 0
	if p.frameSlider != nil {
		frames := len(p.frameImages)
		if frames <= 0 {
			frames = 1
		}
		p.frameSlider.Max = float64(frames - 1)
		p.frameSlider.SetValue(0)
		p.frameLabel.SetText(fmt.Sprintf("Frame %d / %d", p.currentFrame+1, frames))
	}
	p.updatePreview()
}

// SetFrame sets the current frame
func (p *SpritePreviewWidget) SetFrame(frameIdx int) {
	if frameIdx >= 0 && frameIdx < len(p.frameImages) {
		p.currentFrame = frameIdx
		if p.frameSlider != nil {
			p.frameSlider.SetValue(float64(frameIdx))
		}
		p.updatePreview()
	}
}

// GetCurrentFrame returns the current frame index
func (p *SpritePreviewWidget) GetCurrentFrame() int {
	return p.currentFrame
}

// decodeSpriteFrames splits sprite data into per-frame images using the converter.
func decodeSpriteFrames(sprite *models.FF6Sprite) []image.Image {
	if sprite == nil || sprite.Width == 0 || sprite.Height == 0 || sprite.Frames <= 0 {
		return nil
	}

	bytesPerFrame := (sprite.Width * sprite.Height) / 2
	if bytesPerFrame <= 0 || len(sprite.Data) < bytesPerFrame {
		return nil
	}

	frameCount := sprite.Frames
	maxFrames := len(sprite.Data) / bytesPerFrame
	if maxFrames > 0 && frameCount > maxFrames {
		frameCount = maxFrames
	}
	if frameCount <= 0 {
		return nil
	}

	conv := io.NewFF6SpriteConverter()
	frames := make([]image.Image, 0, frameCount)

	for i := 0; i < frameCount; i++ {
		start := i * bytesPerFrame
		end := start + bytesPerFrame
		if end > len(sprite.Data) {
			break
		}

		frameSprite := &models.FF6Sprite{
			ID:        sprite.ID,
			Name:      sprite.Name,
			Type:      sprite.Type,
			Data:      append([]byte(nil), sprite.Data[start:end]...),
			Palette:   sprite.Palette,
			Width:     sprite.Width,
			Height:    sprite.Height,
			Frames:    1,
			FrameRate: sprite.FrameRate,
		}

		img, err := conv.DecodeFF6Sprite(frameSprite)
		if err != nil {
			continue
		}
		frames = append(frames, img)
	}

	return frames
}
