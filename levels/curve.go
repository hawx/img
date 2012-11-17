package levels

import (
	"fmt"
	"strings"
	"strconv"
)

type Point struct {
	X, Y float64
}

func (p *Point) String() string {
	return fmt.Sprintf("%v:%v", p.X, p.Y)
}

func P(x, y float64) *Point {
	return &Point{x, y}
}

type Curve struct {
	Points []*Point
}

func (c *Curve) String() string {
	s := ""
	for i, p := range c.Points {
		if i != 0 { s += ", " }
		s += p.String()
	}

	return "Curve{" + s + "}"
}

// Value calculates the y-value for the given x-value by approximating a bezier
// curve through the points the Curve describes.
func (c *Curve) Value(x float64) float64 {
	// First use a linear connecting curve.
	// - Given a value x find the two Points either side.
	// - Now linear interpolate the value

	var left, right Point
	for i, p := range c.Points {
		if p.X > x {
			left = p
			right = c.Points[i + 1]
		}
	}
}

func C(ps [][]float64) *Curve {
	points := make([]*Point, len(ps))
	for i, p := range ps {
		points[i] = P(p[0], p[1])
	}

	return &Curve{points}
}


func ParseCurveString(s string) *Curve {
	pairs  := strings.Split(s, ":")
	points := make([]*Point, len(pairs))
	conv   := func(v string) (r float64) { r, _ = strconv.ParseFloat(v, 64); return }

	for i, p := range pairs {
		parts := strings.Split(p, ",")
		points[i] = P(conv(parts[0]), conv(parts[1]))
	}

	return &Curve{points}
}
