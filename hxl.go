package main

import (
	"github.com/hawx/img/pixelate"
	"github.com/hawx/img/utils"
)

var cmdHxl = &Command{
	UsageLine: "hxl [options]",
	Short:     "pixelates image into equilateral triangles",
Long: `
  Hxl takes an image from STDIN, pixelates it into equilateral triangles
  forming hexagons, and prints the result to STDOUT

    --width [w]      # Width of the base of each triangle (default: 20)
`,
}

var hxlWidth int

func init() {
	cmdHxl.Run = runHxl

	cmdHxl.Flag.IntVar(&hxlWidth, "width", 20, "")
}

func runHxl(cmd *Command, args []string) {
	i := utils.ReadStdin()
	i  = pixelate.Hxl(i, hxlWidth)
	utils.WriteStdout(i)
}
