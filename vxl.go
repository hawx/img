package main

import (
	"github.com/hawx/img/pixelate"
	"github.com/hawx/img/utils"
	"github.com/hawx/hadfield"
)

var cmdVxl = &hadfield.Command{
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

var vxlHeight, vxlRows int
var vxlFlip bool
var vxlTop, vxlLeft, vxlRight float64

func init() {
	cmdVxl.Run = runVxl

	cmdVxl.Flag.IntVar(&vxlRows,   "rows",   -1,    "")
	cmdVxl.Flag.IntVar(&vxlHeight, "height", 20,    "")

	cmdVxl.Flag.BoolVar(&vxlFlip,  "flip",   false, "")

	cmdVxl.Flag.Float64Var(&vxlTop, "top",   1.0,   "")
	cmdVxl.Flag.Float64Var(&vxlLeft, "left", 2.0,   "")
	cmdVxl.Flag.Float64Var(&vxlRight, "right", 0.5, "")
}

func runVxl(cmd *hadfield.Command, args []string) {
	i := utils.ReadStdin()

	if vxlRows > 0 {
		vxlHeight = utils.SizeForRows(i, vxlRows).H
	}

	i = pixelate.Vxl(i, vxlHeight, vxlFlip, vxlTop, vxlLeft, vxlRight)
	utils.WriteStdout(i)
}
