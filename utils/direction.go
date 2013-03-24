package utils

// Direction represents one of 9 directions. these being the 4 vertices, the 4
// sides, and the Centre of a square.
type Direction int

const (
	Centre Direction = iota
	Top
	TopRight
	Right
	BottomRight
	Bottom
	BottomLeft
	Left
	TopLeft
)
