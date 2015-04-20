package cmd

import (
	"github.com/hawx/hadfield"
	"hawx.me/code/img/pixelate"
	"hawx.me/code/img/utils"
)

var (
	pixelateCrop               bool
	pixelateSize               utils.Dimension = utils.Dimension{20, 20}
	pixelateRows, pixelateCols int
)

func Pixelate() *hadfield.Command {
	cmd := &hadfield.Command{
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

	cmd.Run = runPixelate

	cmd.Flag.BoolVar(&pixelateCrop, "crop", false, "")
	cmd.Flag.Var(&pixelateSize, "size", "")
	cmd.Flag.IntVar(&pixelateRows, "rows", -1, "")
	cmd.Flag.IntVar(&pixelateCols, "cols", -1, "")

	return cmd
}

func runPixelate(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

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

	i = pixelate.Pixelate(i, pixelateSize, style)
	utils.WriteStdout(i, data)
}
