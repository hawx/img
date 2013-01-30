package main

import (
	"github.com/hawx/img/blur"
	"github.com/hawx/img/utils"
)

var cmdBlur = &Command{
	UsageLine: "blur [options]",
	Short:     "blur an image",
Long: `
  Blur takes an image from STDIN, and prints a blurred version to STDOUT.

    --radius <r>                 # Set radius of blur (default: 2.0)
    --size <height>x<width>      # Set size of blur

    --box                        # Perform box blur
    --gaussian <sigma>           # Perform gaussian blur (default: 5.0)
`,
}

var blurRadius int
var blurSize utils.Pixel

var blurBox bool
var blurGaussian float64

func init() {
	cmdBlur.Run = runBlur

	cmdBlur.Flag.IntVar(&blurRadius, "radius", 2.0, "")
	cmdBlur.Flag.Var(&blurSize, "size", "")

	cmdBlur.Flag.BoolVar(&blurBox, "box", false, "")
	cmdBlur.Flag.Float64Var(&blurGaussian, "gaussian", 5.0, "")
}

func runBlur(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if !utils.FlagVisited("size", cmd.Flag) {
		diameter := blurRadius * 2 + 1
		blurSize = utils.Pixel{diameter, diameter}
	}

	if blurBox {
		i = blur.Box(i, blurSize)
	} else {
		i = blur.Gaussian(i, blurSize, blurGaussian)
	}

	utils.WriteStdout(i)
}
