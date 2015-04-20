package cmd

import (
	"os"

	"github.com/hawx/hadfield"
	"hawx.me/code/img/blur"
	"hawx.me/code/img/utils"
)

var (
	blurRadius   int
	blurStyle    string
	blurBox      bool
	blurGaussian float64
)

func Blur() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "blur [options]",
		Short: "blur an image",
		Long: `
  Blur takes an image from STDIN, and prints a blurred version to STDOUT.

    --radius <r>             # Set radius of blur (default: 2.0)
    --style <option>         # Either clamp, ignore or wrap (default: ignore)

    --box                    # Perform box blur
    --gaussian <sigma>       # Perform gaussian blur (default: 5.0)
`,
	}

	cmd.Run = runBlur

	cmd.Flag.IntVar(&blurRadius, "radius", 2.0, "")
	cmd.Flag.StringVar(&blurStyle, "style", "ignore", "")

	cmd.Flag.BoolVar(&blurBox, "box", false, "")
	cmd.Flag.Float64Var(&blurGaussian, "gaussian", 5.0, "")

	return cmd
}

func runBlur(cmd *hadfield.Command, args []string) {
	var style blur.Style

	switch blurStyle {
	case "clamp":
		style = blur.CLAMP
	case "ignore":
		style = blur.IGNORE
	case "wrap":
		style = blur.WRAP
	default:
		utils.Warn("--style must be one of 'clamp', 'ignore' or 'wrap'")
		os.Exit(2)
	}

	i, data := utils.ReadStdin()

	if blurBox {
		i = blur.Box(i, blurRadius, style)
	} else {
		i = blur.Gaussian(i, blurRadius, blurGaussian, style)
	}

	utils.WriteStdout(i, data)
}
