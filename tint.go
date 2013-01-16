package main

import (
	"github.com/hawx/img/tint"
	"github.com/hawx/img/utils"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"errors"
)

var cmdTint = &Command{
	UsageLine: "tint [options]",
	Short:     "tint an image with a colour",
Long: `
  Tint takes an image and tints it using the specified colour, the result is
  printed to STDOUT

    --with [colour]       # Colour to tint with (default: #FF0000A0)
`,
}

type localNRGBA struct {
	R, G, B, A uint8
}

func (c *localNRGBA) String() string {
	return fmt.Sprint(*c)
}

func (c *localNRGBA) Set(value string) error {
	if strings.HasPrefix(value, "#") {
		parseHex := func(s string) uint8 {
			r, _ := strconv.ParseInt(s, 16, 16)
			return uint8(r)
		}

		if len(value) == 4 {
			r := parseHex(value[1:2] + value[1:2])
			g := parseHex(value[2:3] + value[2:3])
			b := parseHex(value[3:4] + value[3:4])
			a := uint8(255)
			*c = localNRGBA{r, g, b, a}

		} else if len(value) == 5 {
			r := parseHex(value[1:2] + value[1:2])
			g := parseHex(value[2:3] + value[2:3])
			b := parseHex(value[3:4] + value[3:4])
			a := parseHex(value[4:5] + value[4:5])
			*c = localNRGBA{r, g, b, a}

		} else if len(value) == 7 {
			r := parseHex(value[1:3])
			g := parseHex(value[3:5])
			b := parseHex(value[5:7])
			a := uint8(255)
			*c = localNRGBA{r, g, b, a}

		} else if len(value) == 9 {
			r := parseHex(value[1:3])
			g := parseHex(value[3:5])
			b := parseHex(value[5:7])
			a := parseHex(value[7:9])
			*c = localNRGBA{r, g, b, a}

		} else {
			return errors.New(`unknown hexadecimal format. Accepts:
  #FFF                  3-digit hexadecimal (#RGB)
  #FFFC                 4-digit hexadecimal (#RGBA)
  #F0CDBB               6-digit hexadecimal (#RRGGBB)
  #F0CDBBAC             8-digit hexadecimal (#RRGGBBAA)
`)
		}

	} else if strings.HasPrefix(value, "rgb(") {
		parts := strings.Split(value[4:len(value)-1], ",")
		r, _ := strconv.Atoi(parts[0])
		g, _ := strconv.Atoi(parts[1])
		b, _ := strconv.Atoi(parts[2])
		*c = localNRGBA{uint8(r), uint8(g), uint8(b), 255}

	} else if strings.HasPrefix(value, "rgba(") {
		parts := strings.Split(value[5:len(value)-1], ",")
		r, _ := strconv.Atoi(parts[0])
		g, _ := strconv.Atoi(parts[1])
		b, _ := strconv.Atoi(parts[2])
		a, _ := strconv.ParseFloat(parts[3], 32)
		*c = localNRGBA{uint8(r), uint8(g), uint8(b), uint8(a * 255)}

	} else {
		return errors.New(`unknown string format. Accepts:
  #FFF                  3-digit hexadecimal (#RGB)
  #FFFC                 4-digit hexadecimal (#RGBA)
  #F0CDBB               6-digit hexadecimal (#RRGGBB)
  #F0CDBBAC             8-digit hexadecimal (#RRGGBBAA)

  rgb(100,50,10)        Red, green, blue (must not contain spaces)
  rgba(100,50,10,.5)    Red, green, blue and alpha (must not contain spaces)
`)
	}

	return nil
}

var tintWith localNRGBA = localNRGBA{255, 0, 0, 160}

func init() {
	cmdTint.Run = runTint

	cmdTint.Flag.Var(&tintWith, "with", "")
}

func runTint(cmd *Command, args []string) {
	i := utils.ReadStdin()

	tintColor := color.NRGBA{
		uint8(tintWith.R),
		uint8(tintWith.G),
		uint8(tintWith.B),
		uint8(tintWith.A),
	}
	i = tint.Tint(i, tintColor)

	utils.WriteStdout(i)
}
