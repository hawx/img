package main

import (
	"github.com/hawx/img/pixelate"
	"github.com/hawx/img/utils"
	"github.com/hawx/hadfield"
)

var cmdPixelate = &hadfield.Command{
	Usage: "pixelate [options]",
	Short: "pixelates image",
Long: `
  Pixelate takes an image from STDIN, pixelates it by averaging the colour in
  large rectangles, and prints the result to STDOUT

    --crop            # Crop final image to edges of pixelation
    --cols <num>      # Split into <num> columns
    --rows <num>      # Split into <num> rows
    --size <HxW>      # Size of pixel to pixelate with (default: 20x20)
`,
}

var pixelateCrop bool
var pixelateSize utils.Dimension = utils.Dimension{20, 20}
var pixelateRows, pixelateCols int

func init() {
	cmdPixelate.Run = runPixelate

	cmdPixelate.Flag.BoolVar(&pixelateCrop, "crop", false, "")
	cmdPixelate.Flag.Var(&pixelateSize,     "size", "")
	cmdPixelate.Flag.IntVar(&pixelateRows,  "rows", -1, "")
	cmdPixelate.Flag.IntVar(&pixelateCols,  "cols", -1, "")
}

func runPixelate(cmd *hadfield.Command, args []string) {
	i := utils.ReadStdin()

	// Default
	style := pixelate.FITTED

	if pixelateCrop {
		style = pixelate.CROPPED
	}

	if pixelateRows > 0 && pixelateCols > 0 {
		pixelateSize = utils.SizeForRowsAndCols(i, pixelateRows, pixelateCols)
	} else if pixelateRows > 0 {
		pixelateSize = utils.SizeForRows(i, pixelateRows)
	} else if pixelateCols > 0 {
		pixelateSize = utils.SizeForCols(i, pixelateCols)
	}

	i  = pixelate.Pixelate(i, pixelateSize, style)
	utils.WriteStdout(i)
}
