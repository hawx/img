package main

import (
	"github.com/hawx/img/pixelate"
	"github.com/hawx/img/utils"
)

var cmdVxl = &Command{
	UsageLine: "vxl [options]",
	Short:     "vxls image",
Long: `
  Vxl takes an image from STDIN, pixelates it into isometric cubes (voxels-ish),
  and prints the result to STDOUT.

    --rows <num>        # Split into <num> rows
    --height <h>        # Height of cube to use (default: 20)
`,
}

var vxlHeight, vxlRows int

func init() {
	cmdVxl.Run = runVxl

	cmdVxl.Flag.IntVar(&vxlRows, "rows", -1, "")
	cmdVxl.Flag.IntVar(&vxlHeight, "height", 20, "")
}

func runVxl(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if vxlRows > 0 {
		vxlHeight = utils.SizeForRows(i, vxlRows).H
	}

	i = pixelate.Vxl(i, vxlHeight)
	utils.WriteStdout(i)
}
