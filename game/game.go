package game

import (
	"fmt"
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	SNAKE_SIZE = 10
	WIDTH      = 320
	HEIGHT     = 260
	OFFSET     = 10
	MAXX       = (WIDTH / SNAKE_SIZE)
	MAXY       = (HEIGHT / SNAKE_SIZE)
)

var SPEEDS = [...]int{30, 20, 15, 12, 10, 6, 5, 4, 3, 2, 1}

type GameState int

const (
	NEW GameState = iota
	OVER
	GAME
	PAUSE
)

type Game struct {
	availableCells  []Vec
	uiDrawer        *UIElement
	availableImages map[string]*ebiten.Image
	snake           Snake
	food            Vec
	currentSpeed    int
	updateCounter   int
	state           GameState
	gameOver        bool
	movementBuff    bool
}

// TODO: remove this function
func (g Game) GetAvailbleCells() *[]Vec {
	return &g.availableCells
}

// Resets the game state
func (g *Game) Initialize(font *text.GoTextFaceSource, availableImages map[string]*ebiten.Image) {
	availableCells := make([]Vec, 0)
	for i := 0; i < MAXX; i++ {
		for j := (OFFSET / 10); j < MAXY; j++ {
			availableCells = append(availableCells, Vec{
				XPos: float32(i * SNAKE_SIZE),
				YPos: float32(j * SNAKE_SIZE),
			})
		}
	}

	g.availableCells = availableCells
	drawer := UIElement{
		Font:           font,
		AvailableWidth: WIDTH,
	}
	g.snake = NewSnake(3, WIDTH, HEIGHT, SNAKE_SIZE)
	g.currentSpeed = 0
	g.snake.SpeedX = 0
	g.snake.SpeedY = float32(SNAKE_SIZE)
	g.updateCounter = 0
	g.gameOver = false
	g.movementBuff = false
	g.placeFood()
	g.uiDrawer = &drawer
	g.availableImages = availableImages
}

func (g *Game) reset() {
	g.snake = NewSnake(3, WIDTH, HEIGHT, SNAKE_SIZE)
	g.currentSpeed = 0
	g.snake.SpeedX = 0
	g.snake.SpeedY = float32(SNAKE_SIZE)
	g.updateCounter = 0
	g.gameOver = false
	g.movementBuff = false
	g.placeFood()
}

func (g *Game) placeFood() {
	freeCells := make([]Vec, len(g.availableCells))
	copy(freeCells, g.availableCells)
	for _, node := range g.snake.Body {
		ind := DetermineIndex(node)
		freeCells[ind] = freeCells[len(freeCells)-1]
		freeCells = freeCells[:len(freeCells)-1]
	}
	g.food = freeCells[rand.IntN(len(freeCells))]
}

func (g *Game) Update() error {
	g.handleInput()
	switch g.state {
	case GAME:
		g.updateCounter++
		// updates every x frames to control snake speed
		if g.updateCounter == SPEEDS[g.currentSpeed] {
			g.updateCounter = 0
			for i := len(g.snake.Body) - 1; i > 0; i-- {
				g.snake.Body[i].XPos = g.snake.Body[i-1].XPos
				g.snake.Body[i].YPos = g.snake.Body[i-1].YPos
				g.snake.Body[i].Dir = g.snake.Body[i-1].Dir
			}
			g.snake.Body[0].YPos += g.snake.SpeedY
			g.snake.Body[0].XPos += g.snake.SpeedX

			// Bites its Body
			for j := 1; j < len(g.snake.Body); j++ {
				if g.snake.Body[0].Equals(&g.snake.Body[j]) {
					g.state = OVER
					break
				}
			}

			// hits a wall
			if g.snake.Body[0].XPos < 0 ||
				g.snake.Body[0].XPos > float32(WIDTH-SNAKE_SIZE) ||
				g.snake.Body[0].YPos < float32(OFFSET) ||
				g.snake.Body[0].YPos > float32(HEIGHT-SNAKE_SIZE) {
				g.state = OVER
			}

			// gets a food
			if g.snake.Body[0].Equals(&g.food) {
				g.snake.FeedSnake()
				if g.currentSpeed <= 10 && len(g.snake.Body)%10 == 0 {
					g.currentSpeed++
				}
				g.placeFood()
			}
			g.movementBuff = false
		}
	}
	return nil
}

func drawImage(target *ebiten.Image, toPrint *ebiten.Image, xPos, yPos, angle float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(op.GeoM.Apply(xPos, yPos))
	op.GeoM.Rotate(angle)
	target.DrawImage(toPrint, op)
}

// TODO: arreglar el margen de los strokes, se debe de considerar el ancho
// de la línea.
func drawSnakeNode(node, prev, next Vec, screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen, node.XPos,
		node.YPos, float32(SNAKE_SIZE),
		float32(SNAKE_SIZE),
		color.White,
		false,
	)

	if node.YPos-SNAKE_SIZE != next.YPos && node.YPos-SNAKE_SIZE != prev.YPos {
		// Draw top line
		vector.StrokeLine(screen,
			node.XPos,
			node.YPos,
			node.XPos+SNAKE_SIZE,
			node.YPos,
			1,
			color.RGBA{0, 0, 0, 255},
			false,
		)
	}

	if node.YPos+SNAKE_SIZE != next.YPos && node.YPos+SNAKE_SIZE != prev.YPos {
		// Draw bottom line
		vector.StrokeLine(screen,
			node.XPos,
			node.YPos+SNAKE_SIZE,
			node.XPos+SNAKE_SIZE,
			node.YPos+SNAKE_SIZE,
			1,
			color.RGBA{0, 0, 0, 255},
			false,
		)
	}

	if node.XPos-SNAKE_SIZE != next.XPos && node.XPos-SNAKE_SIZE != prev.XPos {
		// Draw left line
		vector.StrokeLine(screen,
			node.XPos,
			node.YPos,
			node.XPos,
			node.YPos+SNAKE_SIZE,
			1,
			color.RGBA{255, 0, 0, 255},
			false,
		)
	}

	if node.XPos+SNAKE_SIZE != next.XPos && node.XPos+SNAKE_SIZE != prev.XPos {
		// Draw right line
		vector.StrokeLine(screen,
			node.XPos+SNAKE_SIZE,
			node.YPos,
			node.XPos+SNAKE_SIZE,
			node.YPos+SNAKE_SIZE,
			1,
			color.RGBA{0, 0, 255, 255},
			false,
		)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{106, 13, 131, 255})
	// displays the snake, food and grid. Also in pause so the player can look at the game state under the pause letters
	if g.state == GAME || g.state == PAUSE {
		// Draws the food
		vector.DrawFilledRect(
			screen,
			g.food.XPos,
			g.food.YPos,
			float32(SNAKE_SIZE),
			float32(SNAKE_SIZE),
			color.RGBA{251, 144, 98, 255},
			false,
		)
		// Draws the snake body
		for i := 0; i < len(g.snake.Body); i++ {
			var prev Vec
			var next Vec
			if i != 0 {
				prev = g.snake.Body[i-1]
			}
			if i < len(g.snake.Body)-1 {
				next = g.snake.Body[i+1]
			}
			fmt.Println(prev)
			drawSnakeNode(g.snake.Body[i], prev, next, screen)
		}

		g.uiDrawer.StatusBar(screen, fmt.Sprintf("SCORE: %d", len(g.snake.Body)), fmt.Sprintf("SPEED: %d", g.currentSpeed+1))
	}

	// they are down here so theyre drawn over the game
	if g.state == PAUSE {
		g.uiDrawer.GameScreen(screen, false, "PAUSE", "")
	}

	if g.state == NEW {
		g.uiDrawer.GameScreen(screen, true, "GO SNAKE", "PRESS SPACE TO START")
	}

	if g.state == OVER {
		g.uiDrawer.GameScreen(screen, true, "GAME OVER", fmt.Sprintf("SCORE: %d", len(g.snake.Body)))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func (g *Game) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) && g.snake.SpeedX != 0 && !g.movementBuff {
		g.snake.SpeedX = 0
		g.snake.SpeedY = -float32(SNAKE_SIZE)
		g.movementBuff = true
		g.snake.Body[0].Dir = UP
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) && g.snake.SpeedX != 0 && !g.movementBuff {
		g.snake.SpeedX = 0
		g.snake.SpeedY = float32(SNAKE_SIZE)
		g.movementBuff = true
		g.snake.Body[0].Dir = DOWN
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) && g.snake.SpeedY != 0 && !g.movementBuff {
		g.snake.SpeedX = -float32(SNAKE_SIZE)
		g.snake.SpeedY = 0
		g.movementBuff = true
		g.snake.Body[0].Dir = LEFT
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) && g.snake.SpeedY != 0 && !g.movementBuff {
		g.snake.SpeedX = float32(SNAKE_SIZE)
		g.snake.SpeedY = 0
		g.movementBuff = true
		g.snake.Body[0].Dir = RIGHT
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch g.state {
		case GAME:
			g.state = PAUSE
		case OVER:
			g.reset()
			fallthrough
		case NEW:
			fallthrough
		case PAUSE:
			g.state = GAME
		}
	}
}

func DetermineIndex(target Vec) int {
	target.YPos -= OFFSET
	target.YPos /= SNAKE_SIZE
	target.XPos /= SNAKE_SIZE
	return (int(target.XPos) * (MAXY - 1)) + int(target.YPos)
}
