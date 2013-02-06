package main

import (
	"github.com/hawx/img/pixelate"
	"github.com/hawx/img/utils"
)

var cmdPixelate = &Command{
	UsageLine: "pixelate [options]",
	Short:     "pixelates image",
Long: `
  Pixelate takes an image from STDIN, pixelates it by averaging the colour in
  large rectangles, and prints the result to STDOUT

    --size [HxW]      # Size of pixel to pixelate with (default: 20x20)
`,
}

var pixelateSize utils.Dimension = utils.Dimension{20, 20}

func init() {
	cmdPixelate.Run = runPixelate

	cmdPixelate.Flag.Var(&pixelateSize, "size", "")
}

func runPixelate(cmd *Command, args []string) {
	i := utils.ReadStdin()
	i  = pixelate.Pixelate(i, pixelateSize.H, pixelateSize.W)
	utils.WriteStdout(i)
}
