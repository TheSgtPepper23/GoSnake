package game

// the snake has an array of coords as a body and a vertical and horizontal speed, those are in the snake to be able to access them inside feedSnake
type Snake struct {
	Body     []Vec
	SpeedY   float32
	SpeedX   float32
	NodeSize int
}

// Creates a new snake of the indicated size. The snake its created vertically with the head in the middle and the rest
// of the Body following upwards
func NewSnake(size, width, height, snakeSize int) Snake {
	tempBody := []Vec{
		{
			XPos: float32((width / 2) - snakeSize),
			YPos: float32((height / 2) - snakeSize),
			Dir:  DOWN,
		},
	}

	for i := 0; i < size-1; i++ {
		tempBody = append(tempBody, Vec{
			XPos: float32((width / 2) - snakeSize),
			YPos: float32((height / 2) - (snakeSize * (i + 2))),
			Dir:  DOWN,
		})
	}
	return Snake{
		Body:     tempBody,
		NodeSize: snakeSize,
	}
}

// increases the size of the snake adding a new coord at the end.
// change it so when the new coord is added it doesnt apear sideways
func (s *Snake) FeedSnake() {
	var xPos, yPos float32
	if s.Body[len(s.Body)-1].Dir == DOWN || s.Body[len(s.Body)-1].Dir == UP {
		xPos = s.Body[len(s.Body)-1].XPos
		yPos = s.Body[len(s.Body)-1].YPos + float32(s.NodeSize)
	} else {
		yPos = s.Body[len(s.Body)-1].YPos
		xPos = s.Body[len(s.Body)-1].XPos + float32(s.NodeSize)
	}
	s.Body = append(s.Body, Vec{
		XPos: xPos,
		YPos: yPos,
		Dir:  s.Body[len(s.Body)-1].Dir,
	})
}
