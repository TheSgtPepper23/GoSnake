package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var YELLOW = color.RGBA{238, 175, 97, 255}

type UIElement struct {
	Font           *text.GoTextFaceSource
	AvailableWidth float64
}

// Draws a screen with text on it. It can have a title and a subtitle, the later being optional
func (u *UIElement) GameScreen(screen *ebiten.Image, hasSub bool, title, subtitle string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(u.AvailableWidth/2, 32)
	op.ColorScale.ScaleWithColor(YELLOW)
	op.LineSpacing = 32
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, title, &text.GoTextFace{
		Source: u.Font,
		Size:   32,
	}, op)

	if hasSub {
		op = &text.DrawOptions{}
		op.GeoM.Translate(u.AvailableWidth/2, 100)
		op.ColorScale.ScaleWithColor(YELLOW)
		op.LineSpacing = 10
		op.PrimaryAlign = text.AlignCenter
		text.Draw(screen, subtitle, &text.GoTextFace{
			Source: u.Font,
			Size:   10,
		}, op)
	}
}

func (u *UIElement) StatusBar(screen *ebiten.Image, score, speed string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(0, 0)
	op.ColorScale.ScaleWithColor(YELLOW)
	op.LineSpacing = 10
	op.PrimaryAlign = text.AlignStart
	text.Draw(screen, score, &text.GoTextFace{
		Source: u.Font,
		Size:   10,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(u.AvailableWidth, 0)
	op.ColorScale.ScaleWithColor(YELLOW)
	op.LineSpacing = 10
	op.PrimaryAlign = text.AlignEnd
	text.Draw(screen, speed, &text.GoTextFace{
		Source: u.Font,
		Size:   10,
	}, op)
}
