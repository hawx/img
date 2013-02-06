package main

import (
	"github.com/hawx/img/pixelate"
	"github.com/hawx/img/utils"
)

var cmdPxl = &Command{
	UsageLine: "pxl [options]",
	Short:     "pxls image",
Long: `
  Pxl takes an image from STDIN, pxls it by averaging the colour in large
  rectangles, and prints the result to STDOUT

    --alias         # Do not use antialiasing
    --left          # Use only left triangles
    --right         # Use only right triangles
    --size <HxW>    # Size of pixel to pxl with (default: 20x20)
`,
}

var pxlAlias bool
var pxlSize utils.Dimension = utils.Dimension{20, 20}
var pxlLeft, pxlRight, pxlBoth bool

func init() {
	cmdPxl.Run = runPxl

	cmdPxl.Flag.BoolVar(&pxlAlias, "alias", false, "")
	cmdPxl.Flag.BoolVar(&pxlLeft,  "left",  false, "")
	cmdPxl.Flag.BoolVar(&pxlRight, "right", false, "")
	cmdPxl.Flag.Var(&pxlSize, "size", "")
}

func runPxl(cmd *Command, args []string) {
	i := utils.ReadStdin()

	triangle := pixelate.BOTH
	if pxlLeft  { triangle = pixelate.LEFT }
	if pxlRight { triangle = pixelate.RIGHT }

	if pxlAlias {
		i = pixelate.AliasedPxl(i, pxlSize, triangle)
	} else {
		i = pixelate.Pxl(i, pxlSize, triangle)
	}

	utils.WriteStdout(i)
}
