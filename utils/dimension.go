package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Dimension represents a rectangle with Height and Width.
type Dimension struct {
	H, W int
}

func (d *Dimension) String() string {
	return fmt.Sprint("%vx%v", *d)
}

func (d *Dimension) Set(value string) error {
	parts := strings.Split(value, "x")

	h,e := strconv.Atoi(parts[0])

	if e != nil {
		return errors.New("error parsing height, expect HxW where H and W are integers")
	}

	w,f := strconv.Atoi(parts[1])

	if f != nil {
		return errors.New("error parsing width, expect HxW where H and W are integers")
	}

	*d = Dimension{h, w}

	return nil
}
