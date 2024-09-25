package game

// Vec just has a position, this is the base for all the game elements
// dir indicates the direction, true is vertical and false is horizontal
type Vec struct {
	XPos float32
	YPos float32
	Dir  DIR
}

// easy way to compare two coords
func (c *Vec) Equals(c2 *Vec) bool {
	return c.XPos == c2.XPos && c.YPos == c2.YPos
}
