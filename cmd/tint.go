package cmd

import (
	"errors"
	"fmt"
	"image/color"

	"hawx.me/code/hadfield"
	"hawx.me/code/img/altcolor"
	"hawx.me/code/img/tint"
	"hawx.me/code/img/utils"
)

var (
	tintWith localNRGBA = localNRGBA{255, 0, 0, 160}
)

func Tint() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "tint [options]",
		Short: "tint an image with a colour",
		Long: `
  Tint takes an image and tints it using the specified colour, the result is
  printed to STDOUT

    --with [colour]       # Colour to tint with (default: #FF0000A0)
`,
	}

	cmd.Run = runTint

	cmd.Flag.Var(&tintWith, "with", "")

	return cmd
}

func runTint(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	tintColor := color.NRGBA{
		uint8(tintWith.R),
		uint8(tintWith.G),
		uint8(tintWith.B),
		uint8(tintWith.A),
	}
	i = tint.Tint(i, tintColor)

	utils.WriteStdout(i, data)
}

type localNRGBA struct {
	R, G, B, A uint8
}

func (c *localNRGBA) String() string {
	return fmt.Sprint(*c)
}

func (c *localNRGBA) Set(value string) error {
	col := altcolor.Parse(value)

	if col == nil {
		return errors.New(`unknown string format. Accepts:
  #FFF                  3-digit hexadecimal (#RGB)
  #FFFC                 4-digit hexadecimal (#RGBA)
  #F0CDBB               6-digit hexadecimal (#RRGGBB)
  #F0CDBBAC             8-digit hexadecimal (#RRGGBBAA)

  rgb(100,50,10)        Red, green, blue
  rgba(100,50,10,127)    Red, green, blue and alpha
`)
	}

	r, g, b, a := utils.NormalisedRGBA(col)
	*c = localNRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}

	return nil
}
