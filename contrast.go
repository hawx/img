package main

import (
	"github.com/hawx/img/contrast"
	"github.com/hawx/img/utils"
)

var cmdContrast = &Command{
	UsageLine: "contrast [options]",
	Short:     "adjust image contrast",
Long: `
  Contrast takes a png file from STDIN, adjusts the contrast and prints the result
  to STDOUT

    --by [n]         # Amount to shift contrast by (default: 15.0)
`,
}

var contrastBy float64

func init() {
	cmdContrast.Run = runContrast

	cmdContrast.Flag.Float64Var(&contrastBy, "by", 15.0, "")
}

func runContrast(cmd *Command, args []string) {
	i := utils.ReadStdin()
	i  = contrast.Adjust(i, contrastBy)
	utils.WriteStdout(i)
}
