package levels

import (
//	"github.com/hawx/img/utils"
	"fmt"
	"strings"
	"strconv"
)

const (
	VALUE_DELIM = ","
	POINT_DELIM = " "
)


type Point struct {
	X, Y float64
}

func (p *Point) String() string {
	return fmt.Sprintf("%v" + VALUE_DELIM + "%v", p.X, p.Y)
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
		if i != 0 { s += POINT_DELIM }
		s += p.String()
	}

	return "Curve{" + s + "}"
}

// Value calculates the y-value for the given x-value by approximating a bezier
// curve through the points the Curve describes.
//
// BUG(r): Only uses linear interpolation.
func (c *Curve) Value(x float64) float64 {
	// find the points either side of x
	var left, right *Point

	for i := 0; i < len(c.Points) - 1; i++ {
		if c.Points[i].X == x * 100.0 {
			return c.Points[i].X / 100.0

		} else if c.Points[i+1].X == x * 100.0 {
			return c.Points[i+1].X / 100.0

		} else if c.Points[i].X > x {
			left = c.Points[i]
			right = c.Points[i+1]
		}
	}

	// find the gradient
	m := (left.Y - right.Y) / (left.X - right.X)

	// and finally, find the value
	return (left.Y / 100.0) + (m * (x - (left.X / 100.0)))
}

func C(ps [][]float64) *Curve {
	points := make([]*Point, len(ps))
	for i, p := range ps {
		points[i] = P(p[0], p[1])
	}

	return &Curve{points}
}


func ParseCurveString(s string) *Curve {
	pairs  := strings.Split(s, POINT_DELIM)
	points := make([]*Point, len(pairs))
	conv   := func(v string) (r float64) { r, _ = strconv.ParseFloat(v, 64); return }

	for i, p := range pairs {
		parts := strings.Split(p, VALUE_DELIM)
		points[i] = P(conv(parts[0]), conv(parts[1]))
	}

	return &Curve{points}
}
