package main

import (
	"github.com/hawx/img/tint"
	"github.com/hawx/img/utils"
	"fmt"
	"image/color"
	"strconv"
	"regexp"
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

type localRGBA struct {
	R, G, B, A uint8
}

func (c *localRGBA) String() string {
	return fmt.Sprint(*c)
}

func (c *localRGBA) Set(value string) error {
	hex := regexp.MustCompile("#[0-9a-fA-F]{8}")
	// rgb := regexp.MustCompile("rgba\\((\\d{1,3},\\s*){3}\\d{1,3}\\)")

	if hex.MatchString(value) {
		// We have a hex number
		r, _ := strconv.ParseInt(value[1:3], 16, 16)
		g, _ := strconv.ParseInt(value[3:5], 16, 16)
		b, _ := strconv.ParseInt(value[5:7], 16, 16)
		a, _ := strconv.ParseInt(value[7:9], 16, 16)
		*c = localRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	// } else if rgb.MatchString(value) {
		// We have an rgba number
		//*c = localRGBA{255, 0, 0, 255}
	} else {
		return errors.New("unknown string format passed to with flag")
	}

	return nil
}

var tintWith localRGBA = localRGBA{255, 0, 0, 160}

func init() {
	cmdTint.Run = runTint

	cmdTint.Flag.Var(&tintWith, "with", "")
}

func runTint(cmd *Command, args []string) {
	i := utils.ReadStdin()

	tintColor := color.RGBA{
		uint8(tintWith.R),
		uint8(tintWith.G),
		uint8(tintWith.B),
		uint8(tintWith.A),
	}
	i = tint.Tint(i, tintColor)

	utils.WriteStdout(i)
}
