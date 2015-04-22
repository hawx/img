package cmd

import (
	"hawx.me/code/hadfield"
	"hawx.me/code/img/crop"
	"hawx.me/code/img/utils"
)

var (
	cropSquare, cropCircle, cropTriangle bool
	cropSize                             int
	cropCentre, cropTop, cropTopRight, cropRight, cropBottomRight, cropBottom,
	cropBottomLeft, cropLeft, cropTopLeft bool
)

func Crop() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "crop [options]",
		Short: "crop an image",
		Long: `
  Crop takes an image fron STDIN, and returns a cropped version to STDOUT. By
  default it will use the largest size possible; that is, if the image is wider
  than tall it will use the height, and vica versa.

    --square               # Crop to a square (default)
    --circle               # Crop to a circle
    --triangle             # Crop to an equilateral triangle

    --size <pixels>        # Size to crop to (default: largest possible)

    --centre               # Centre the image
    --top                  # Centre the image to the top of the frame
    --top-right            # Centre the image to the top-right of the frame
    --right                # Centre the image to the right of the frame
    --bottom-right         # Centre the image to the bottom-right of the frame
    --bottom               # Centre the image to the bottom of the frame
    --bottom-left          # Centre the image to the bottom-left of the frame
    --left                 # Centre the image to the left of the frame
    --top-left             # Centre the image to the top-left of the frame
`,
	}

	cmd.Run = runCrop

	cmd.Flag.BoolVar(&cropSquare, "square", false, "")
	cmd.Flag.BoolVar(&cropCircle, "circle", false, "")
	cmd.Flag.BoolVar(&cropTriangle, "triangle", false, "")

	cmd.Flag.IntVar(&cropSize, "size", -1, "")

	cmd.Flag.BoolVar(&cropCentre, "centre", false, "")
	cmd.Flag.BoolVar(&cropTop, "top", false, "")
	cmd.Flag.BoolVar(&cropTopRight, "top-right", false, "")
	cmd.Flag.BoolVar(&cropRight, "right", false, "")
	cmd.Flag.BoolVar(&cropBottomRight, "bottom-right", false, "")
	cmd.Flag.BoolVar(&cropBottom, "bottom", false, "")
	cmd.Flag.BoolVar(&cropBottomLeft, "bottom-left", false, "")
	cmd.Flag.BoolVar(&cropLeft, "left", false, "")
	cmd.Flag.BoolVar(&cropTopLeft, "top-left", false, "")

	return cmd
}

func runCrop(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	direction := utils.Centre

	if cropTop {
		direction = utils.Top
	} else if cropTopRight {
		direction = utils.TopRight
	} else if cropRight {
		direction = utils.Right
	} else if cropBottomRight {
		direction = utils.BottomRight
	} else if cropBottom {
		direction = utils.Bottom
	} else if cropBottomLeft {
		direction = utils.BottomLeft
	} else if cropLeft {
		direction = utils.Left
	} else if cropTopLeft {
		direction = utils.TopLeft
	}

	if cropCircle {
		i = crop.Circle(i, cropSize, direction)
	} else if cropTriangle {
		i = crop.Triangle(i, cropSize, direction)
	} else {
		i = crop.Square(i, cropSize, direction)
	}

	utils.WriteStdout(i, data)
}
