package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TheSgtPepper23/goSnake/files"
	"github.com/TheSgtPepper23/goSnake/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// every x frames the game will be updated
var arcadeFaceSource *text.GoTextFaceSource
var availableImages map[string]*ebiten.Image

func init() {
	//Load the font
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = s

	availableImages = make(map[string]*ebiten.Image)

	savedFiles, err := os.ReadDir("./assets")
	if err != nil {
		log.Fatal(err)
	}

	//Load the images
	for i := 0; i < len(savedFiles); i++ {
		if strings.HasSuffix(savedFiles[i].Name(), ".png") {
			tempImg, err := files.LoadImgFromFile(fmt.Sprintf("./assets/%s", savedFiles[i].Name()))
			if err != nil {
				log.Fatal(err)
			}
			availableImages[strings.Split(savedFiles[i].Name(), ".")[0]] = ebiten.NewImageFromImage(tempImg)
		}
	}
}

// Disminuir la velocidad máxima
func main() {
	ebiten.SetWindowSize(960, 780)
	ebiten.SetWindowTitle("SNAKE")

	game := &game.Game{}
	game.Initialize(arcadeFaceSource, &availableImages)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
