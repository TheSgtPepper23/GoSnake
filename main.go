package main

import (
	"bytes"
	"log"

	"github.com/TheSgtPepper23/goSnake/files"
	"github.com/TheSgtPepper23/goSnake/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// every x frames the game will be updated
var arcadeFaceSource *text.GoTextFaceSource
var normalImage *ebiten.Image

func init() {
	//Load the font
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = s

	//Load the images
	normalImg, err := files.LoadImgFromFile("./assets/normal.png")
	if err != nil {
		log.Fatal("Error converting image")
	}
	normalImage = ebiten.NewImageFromImage(normalImg)
}

// Disminuir la velocidad máxima
func main() {
	ebiten.SetWindowSize(960, 780)
	ebiten.SetWindowTitle("SNAKE")

	game := &game.Game{}
	game.Initialize(arcadeFaceSource)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
