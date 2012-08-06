package utils

import (
	"strings"
	"strconv"
	"fmt"
)

type Pixel struct {
	H, W int
}

func (p *Pixel) String() string {
	return fmt.Sprint(*p)
}

func (p *Pixel) Set(value string) error {
	parts := strings.Split(value, "x")

	h, _ := strconv.Atoi(parts[0])
	w, _ := strconv.Atoi(parts[1])

	*p = Pixel{h, w}

	return nil
}
