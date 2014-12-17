package main

import (
	"github.com/hawx/img/pixelate"
	"github.com/hawx/img/utils"
	"github.com/hawx/hadfield"
)

var cmdPxl = &hadfield.Command{
	Usage: "pxl [options]",
	Short: "pxls image",
Long: `
  Pxl takes an image from STDIN, pxls it by averaging the colour in large
  rectangles, and prints the result to STDOUT

    --alias         # Do not use antialiasing
    --crop          # Crop final image to exact triangles
    --left          # Use only left triangles
    --right         # Use only right triangles

    --cols <num>    # Split into <num> columns
    --rows <num>    # Split into <num> rows
    --size <HxW>    # Size of pixel to pxl with
`,
}

var pxlAlias, pxlCrop, pxlLeft, pxlRight, pxlBoth bool
var pxlSize utils.Dimension = utils.Dimension{-1, -1}
var pxlRows, pxlCols int

func init() {
	cmdPxl.Run = runPxl

	cmdPxl.Flag.BoolVar(&pxlAlias, "alias", false, "")
	cmdPxl.Flag.BoolVar(&pxlCrop,  "crop",  false, "")
	cmdPxl.Flag.BoolVar(&pxlLeft,  "left",  false, "")
	cmdPxl.Flag.BoolVar(&pxlRight, "right", false, "")

	cmdPxl.Flag.Var(&pxlSize,      "size",         "")
	cmdPxl.Flag.IntVar(&pxlRows,   "rows",  -1,    "")
	cmdPxl.Flag.IntVar(&pxlCols,   "cols",  -1,    "")
}

func runPxl(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	triangle := pixelate.BOTH
	if pxlLeft  { triangle = pixelate.LEFT }
	if pxlRight { triangle = pixelate.RIGHT }

	style := pixelate.FITTED
	if pxlCrop { style = pixelate.CROPPED }

	if pxlRows > 0 && pxlCols > 0 {
		pxlSize = utils.SizeForRowsAndCols(i, pxlRows, pxlCols)
	} else if pxlRows > 0 {
		pxlSize = utils.SizeForRows(i, pxlRows)
	} else if pxlCols > 0 {
		pxlSize = utils.SizeForCols(i, pxlCols)
	}

	// If no sizes given, guess
	if pxlSize.H == -1 && pxlSize.W == -1 {
		bounds := i.Bounds()

		if bounds.Dx() > bounds.Dy() {
			pxlSize = utils.SizeForCols(i, 20)
		} else {
			pxlSize = utils.SizeForRows(i, 20)
		}
	}

	if pxlAlias {
		i = pixelate.AliasedPxl(i, pxlSize, triangle, style)
	} else {
		i = pixelate.Pxl(i, pxlSize, triangle, style)
	}

	utils.WriteStdout(i, data)
}
