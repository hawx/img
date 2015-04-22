package cmd

import (
	"hawx.me/code/hadfield"
	"hawx.me/code/img/pixelate"
	"hawx.me/code/img/utils"
)

var (
	pxlAlias, pxlCrop, pxlLeft, pxlRight, pxlBoth bool
	pxlSize                                       utils.Dimension = utils.Dimension{-1, -1}
	pxlRows, pxlCols                              int
)

func Pxl() *hadfield.Command {
	cmd := &hadfield.Command{
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

	cmd.Run = runPxl

	cmd.Flag.BoolVar(&pxlAlias, "alias", false, "")
	cmd.Flag.BoolVar(&pxlCrop, "crop", false, "")
	cmd.Flag.BoolVar(&pxlLeft, "left", false, "")
	cmd.Flag.BoolVar(&pxlRight, "right", false, "")

	cmd.Flag.Var(&pxlSize, "size", "")
	cmd.Flag.IntVar(&pxlRows, "rows", -1, "")
	cmd.Flag.IntVar(&pxlCols, "cols", -1, "")

	return cmd
}

func runPxl(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	triangle := pixelate.BOTH
	if pxlLeft {
		triangle = pixelate.LEFT
	}
	if pxlRight {
		triangle = pixelate.RIGHT
	}

	style := pixelate.FITTED
	if pxlCrop {
		style = pixelate.CROPPED
	}

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
