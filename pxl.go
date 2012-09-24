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

    --size [HxW]    # Size of pixel to pxl with (default: 20x20)
    --left          # Use only left triangles
    --right         # Use only right triangles
`,
}

var pxlSize utils.Pixel = utils.Pixel{20, 20}
var pxlLeft, pxlRight, pxlBoth bool

func init() {
	cmdPxl.Run = runPxl

	cmdPxl.Flag.Var(&pxlSize, "size", "")
	cmdPxl.Flag.BoolVar(&pxlLeft,  "left",  false, "")
	cmdPxl.Flag.BoolVar(&pxlRight, "right", false, "")
}

func runPxl(cmd *Command, args []string) {
	i := utils.ReadStdin()

	triangle := pixelate.BOTH
	if pxlLeft  { triangle = pixelate.LEFT }
	if pxlRight { triangle = pixelate.RIGHT }

	i  = pixelate.Pxl(i, triangle, pxlSize.H, pxlSize.W)
	utils.WriteStdout(i)
}
