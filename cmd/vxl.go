package cmd

import (
	"github.com/hawx/hadfield"
	"github.com/hawx/img/pixelate"
	"github.com/hawx/img/utils"
)

var (
	vxlHeight, vxlRows        int
	vxlFlip                   bool
	vxlTop, vxlLeft, vxlRight float64
)

func Vxl() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "vxl [options]",
		Short: "vxls image",
		Long: `
  Vxl takes an image from STDIN, pixelates it into isometric cubes (voxels-ish),
  and prints the result to STDOUT.

    --rows <num>        # Split into <num> rows
    --height <h>        # Height of cube to use (default: 20)

    --flip              # Flip orientation of cubes

    --top <ratio>       # Ratio to adjust lightness of top square
    --left <ratio>      # Ratio to adjust lightness of left part
    --right <ratio>     # Ratio to adjust lightness of right part
`,
	}

	cmd.Run = runVxl

	cmd.Flag.IntVar(&vxlRows, "rows", -1, "")
	cmd.Flag.IntVar(&vxlHeight, "height", 20, "")

	cmd.Flag.BoolVar(&vxlFlip, "flip", false, "")

	cmd.Flag.Float64Var(&vxlTop, "top", 1.0, "")
	cmd.Flag.Float64Var(&vxlLeft, "left", 2.0, "")
	cmd.Flag.Float64Var(&vxlRight, "right", 0.5, "")

	return cmd
}

func runVxl(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	if vxlRows > 0 {
		vxlHeight = utils.SizeForRows(i, vxlRows).H
	}

	i = pixelate.Vxl(i, vxlHeight, vxlFlip, vxlTop, vxlLeft, vxlRight)
	utils.WriteStdout(i, data)
}
