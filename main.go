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

func main() {
	availableImages := make(map[string]*ebiten.Image)

	defFont, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}

	savedFiles, err := os.ReadDir("./assets")
	if err != nil {
		log.Fatal(err)
	}

	// Load the images
	for i := 0; i < len(savedFiles); i++ {
		if strings.HasSuffix(savedFiles[i].Name(), ".png") {
			tempImg, err := files.LoadImgFromFile(fmt.Sprintf("./assets/%s", savedFiles[i].Name()))
			if err != nil {
				log.Fatal(err)
			}
			availableImages[strings.Split(savedFiles[i].Name(), ".")[0]] = ebiten.NewImageFromImage(tempImg)
		}
	}

	ebiten.SetWindowSize(960, 780)
	ebiten.SetWindowTitle("SNAKE")

	game := &game.Game{}
	game.Initialize(defFont, availableImages)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
