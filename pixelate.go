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

    --cols <num>      # Split into <num> columns
    --rows <num>      # Split into <num> rows
    --size <HxW>      # Size of pixel to pixelate with (default: 20x20)
`,
}

var pixelateSize utils.Dimension = utils.Dimension{20, 20}
var pixelateRows, pixelateCols int

func init() {
	cmdPixelate.Run = runPixelate

	cmdPixelate.Flag.Var(&pixelateSize, "size", "")
	cmdPixelate.Flag.IntVar(&pixelateRows, "rows", -1, "")
	cmdPixelate.Flag.IntVar(&pixelateCols, "cols", -1, "")
}

func runPixelate(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if pixelateRows > 0 && pixelateCols > 0 {
		pixelateSize = utils.SizeForRowsAndCols(i, pixelateRows, pixelateCols)
	} else if pixelateRows > 0 {
		pixelateSize = utils.SizeForRows(i, pixelateRows)
	} else if pixelateCols > 0 {
		pixelateSize = utils.SizeForCols(i, pixelateCols)
	}

	i  = pixelate.Pixelate(i, pixelateSize)
	utils.WriteStdout(i)
}
