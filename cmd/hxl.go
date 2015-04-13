package cmd

import (
	"github.com/hawx/hadfield"
	"github.com/hawx/img/pixelate"
	"github.com/hawx/img/utils"
)

var (
	hxlWidth, hxlCols int
)

func Hxl() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "hxl [options]",
		Short: "pixelates image into equilateral triangles",
		Long: `
  Hxl takes an image from STDIN, pixelates it into equilateral triangles
  forming hexagons, and prints the result to STDOUT

    --cols <num>     # Split into <num> columns
    --width <w>      # Width of the base of each triangle (default: 20)
`,
	}

	cmd.Run = runHxl

	cmd.Flag.IntVar(&hxlCols, "cols", -1, "")
	cmd.Flag.IntVar(&hxlWidth, "width", 20, "")

	return cmd
}

func runHxl(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	if hxlCols > 0 {
		hxlWidth = utils.SizeForCols(i, hxlCols).W
	}

	i = pixelate.Hxl(i, hxlWidth)
	utils.WriteStdout(i, data)
}
