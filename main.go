package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	SNAKE_SIZE  = 10
	WIDTH       = 320
	HEIGHT      = 260
	GAME_OFFSET = 10
	MAXX        = (WIDTH / SNAKE_SIZE)
	MAXY        = (HEIGHT / SNAKE_SIZE)
)

// every x frames the game will be updated
var SPEEDS = [...]int{20, 15, 12, 10, 6, 5, 4, 3, 2, 1}
var YELLOW = color.RGBA{238, 175, 97, 255}
var arcadeFaceSource *text.GoTextFaceSource

var availableCells = make([]Vec, 0)

type GameState int

const (
	NEW GameState = iota
	OVER
	GAME
	PAUSE
)

type DIR int

const (
	UP DIR = iota
	DOWN
	LEFT
	RIGHT
)

// Vec just has a position, this is the base for all the game elements
// dir indicates the direction, true is vertical and false is horizontal
type Vec struct {
	xPos float32
	yPos float32
	dir  DIR
}

// easy way to compare two coords
func (c *Vec) Equals(c2 *Vec) bool {
	return c.xPos == c2.xPos && c.yPos == c2.yPos
}

// the snake has an array of coords as a body and a vertical and horizontal speed, those are in the snake to be able to access them inside feedSnake
type Snake struct {
	body   []Vec
	speedY float32
	speedX float32
}

// Creates a new snake of the indicated size. The snake its created vertically with the head in the middle and the rest
// of the body following upwards
func newSnake(size int) Snake {
	tempBody := []Vec{
		{
			xPos: (WIDTH / 2) - SNAKE_SIZE,
			yPos: (HEIGHT / 2) - SNAKE_SIZE,
			dir:  DOWN,
		},
	}

	for i := 0; i < size-1; i++ {
		tempBody = append(tempBody, Vec{
			xPos: (WIDTH / 2) - SNAKE_SIZE,
			yPos: float32((HEIGHT / 2) - (SNAKE_SIZE * (i + 2))),
			dir:  DOWN,
		})
	}
	return Snake{
		body: tempBody,
	}
}

// increases the size of the snake adding a new coord at the end.
// change it so when the new coord is added it doesnt apear sideways
func (s *Snake) feedSnake() {
	var xPos, yPos float32
	if s.body[len(s.body)-1].dir == DOWN || s.body[len(s.body)-1].dir == UP {
		xPos = s.body[len(s.body)-1].xPos
		yPos = s.body[len(s.body)-1].yPos + SNAKE_SIZE
	} else {
		yPos = s.body[len(s.body)-1].yPos
		xPos = s.body[len(s.body)-1].xPos + SNAKE_SIZE
	}
	s.body = append(s.body, Vec{
		xPos: xPos,
		yPos: yPos,
		dir:  s.body[len(s.body)-1].dir,
	})
}

type Game struct {
	snake          Snake
	currentSpeed   int
	food           Vec
	updateCounter  int
	gameOver       bool
	state          GameState
	availableCells *[]Vec
	movementBuff   bool
}

// Draws a screen with text on it. It can have a title and a subtitle, the later being optional
func gameScreen(screen *ebiten.Image, hasSub bool, title, subtitle string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(WIDTH/2, 32)
	op.ColorScale.ScaleWithColor(YELLOW)
	op.LineSpacing = 32
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, title, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   32,
	}, op)

	if hasSub {
		op = &text.DrawOptions{}
		op.GeoM.Translate(WIDTH/2, 100)
		op.ColorScale.ScaleWithColor(YELLOW)
		op.LineSpacing = 10
		op.PrimaryAlign = text.AlignCenter
		text.Draw(screen, subtitle, &text.GoTextFace{
			Source: arcadeFaceSource,
			Size:   10,
		}, op)
	}
}

func statusBar(screen *ebiten.Image, score, speed string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(0, 0)
	op.ColorScale.ScaleWithColor(YELLOW)
	op.LineSpacing = 10
	op.PrimaryAlign = text.AlignStart
	text.Draw(screen, score, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   10,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(WIDTH, 0)
	op.ColorScale.ScaleWithColor(YELLOW)
	op.LineSpacing = 10
	op.PrimaryAlign = text.AlignEnd
	text.Draw(screen, speed, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   10,
	}, op)
}

// places a new food in any of the tiles where the snake its NOT
func (g *Game) placeFood() {
	freeCells := make([]*Vec, 0)
	remaining := len(g.snake.body)
	found := false
	for i := 0; i < len(*g.availableCells); i++ {
		found = false
		//Can finish early if all the positions of the snake are found. Not all the tiles will be available, but its faster
		if remaining == 0 {
			break
		}
		for j := 0; j < len(g.snake.body); j++ {
			//crazy pointers
			if g.snake.body[j].Equals(&(*g.availableCells)[i]) {
				remaining--
				found = true
				break
			}
		}
		if !found {
			//crazy pointers
			freeCells = append(freeCells, &(*g.availableCells)[i])
		}
	}

	g.food = *freeCells[rand.IntN(len(freeCells))]
}

func (g *Game) Update() error {
	g.handleInput()
	switch g.state {
	case GAME:
		g.updateCounter++

		//updates every x frames to control snake speed
		if g.updateCounter == SPEEDS[g.currentSpeed] {
			g.updateCounter = 0
			for i := len(g.snake.body) - 1; i > 0; i-- {
				g.snake.body[i].xPos = g.snake.body[i-1].xPos
				g.snake.body[i].yPos = g.snake.body[i-1].yPos
				g.snake.body[i].dir = g.snake.body[i-1].dir
			}
			g.snake.body[0].yPos += g.snake.speedY
			g.snake.body[0].xPos += g.snake.speedX
			if g.snake.speedX == -10 {
				g.snake.body[0].dir = LEFT
			} else if g.snake.speedX == 10 {
				g.snake.body[0].dir = RIGHT
			} else if g.snake.speedY == -10 {
				g.snake.body[0].dir = UP
			} else {
				g.snake.body[0].dir = DOWN
			}

			//Bites its body
			for j := 1; j < len(g.snake.body); j++ {
				if g.snake.body[0].Equals(&g.snake.body[j]) {
					g.state = OVER
					break
				}
			}

			//hits a wall
			if g.snake.body[0].xPos < 0 || g.snake.body[0].xPos > WIDTH-SNAKE_SIZE || g.snake.body[0].yPos < GAME_OFFSET || g.snake.body[0].yPos > HEIGHT-SNAKE_SIZE {
				g.state = OVER
			}

			//gets a food
			if g.snake.body[0].Equals(&g.food) {
				g.snake.feedSnake()
				if g.currentSpeed <= 10 && len(g.snake.body)%10 == 0 {
					g.currentSpeed++
				}
				g.placeFood()
			}
			g.movementBuff = false
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{106, 13, 131, 255})
	//displays the snake, food and grid. Also in pause so the player can look at the game state under the pause letters
	if g.state == GAME || g.state == PAUSE {
		//Draws the food
		vector.DrawFilledRect(screen, g.food.xPos, g.food.yPos, SNAKE_SIZE, SNAKE_SIZE, color.RGBA{251, 144, 98, 255}, false)
		//Draws the snake body
		var dirColor color.RGBA
		for i := 0; i < len(g.snake.body); i++ {
			switch g.snake.body[i].dir {
			case UP:
				dirColor = color.RGBA{255, 0, 0, 255}
			case DOWN:
				dirColor = color.RGBA{0, 255, 0, 255}
			case LEFT:
				dirColor = color.RGBA{0, 0, 255, 255}
			case RIGHT:
				dirColor = color.RGBA{255, 255, 255, 255}
			}

			vector.DrawFilledRect(screen, g.snake.body[i].xPos, g.snake.body[i].yPos, SNAKE_SIZE, SNAKE_SIZE, dirColor, true)
		}

		//draws the grid over the other two elements
		for j := 0; j < len(*g.availableCells); j++ {
			vector.StrokeRect(screen, (*g.availableCells)[j].xPos, (*g.availableCells)[j].yPos, SNAKE_SIZE, SNAKE_SIZE, 1, color.RGBA{0, 102, 204, 255}, false)
		}

		statusBar(screen, fmt.Sprintf("SCORE: %d", len(g.snake.body)), fmt.Sprintf("SPEED: %d", g.currentSpeed+1))
	}

	//they are down here so theyre drawn over the game
	if g.state == PAUSE {
		gameScreen(screen, false, "PAUSE", "")
	}

	if g.state == NEW {
		gameScreen(screen, true, "GO SNAKE", "PRESS SPACE TO START")
	}

	if g.state == OVER {
		gameScreen(screen, true, "GAME OVER", fmt.Sprintf("SCORE: %d", len(g.snake.body)))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

// Resets the game state
func (g *Game) initialize() {
	g.snake = newSnake(3)
	g.currentSpeed = 0
	g.snake.speedX = 0
	g.snake.speedY = SNAKE_SIZE
	g.updateCounter = 0
	g.gameOver = false
	g.availableCells = &availableCells
	g.movementBuff = false
	g.placeFood()
}

func (g *Game) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) && g.snake.speedX != 0 && !g.movementBuff {
		g.snake.speedX = 0
		g.snake.speedY = -SNAKE_SIZE
		g.movementBuff = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) && g.snake.speedX != 0 && !g.movementBuff {
		g.snake.speedX = 0
		g.snake.speedY = SNAKE_SIZE
		g.movementBuff = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) && g.snake.speedY != 0 && !g.movementBuff {
		g.snake.speedX = -SNAKE_SIZE
		g.snake.speedY = 0
		g.movementBuff = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) && g.snake.speedY != 0 && !g.movementBuff {
		g.snake.speedX = SNAKE_SIZE
		g.snake.speedY = 0
		g.movementBuff = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch g.state {
		case GAME:
			g.state = PAUSE
		case OVER:
			g.initialize()
			fallthrough
		case NEW:
			fallthrough
		case PAUSE:
			g.state = GAME
		}
	}
}

// Disminuir la velocidad máxima
func main() {
	ebiten.SetWindowSize(960, 780)
	ebiten.SetWindowTitle("SNAKE")
	//Populate the board
	for i := 0; i < MAXX; i++ {
		for j := (GAME_OFFSET / 10); j < MAXY; j++ {
			availableCells = append(availableCells, Vec{xPos: float32(i * SNAKE_SIZE), yPos: float32(j * SNAKE_SIZE)})
		}
	}
	//Load the font
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = s
	game := &Game{}
	game.initialize()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
