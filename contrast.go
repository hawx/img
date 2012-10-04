package main

import (
	"github.com/hawx/img/contrast"
	"github.com/hawx/img/utils"
)

var cmdContrast = &Command{
	UsageLine: "contrast [options]",
	Short:     "adjust image contrast",
Long: `
  Contrast takes an image from STDIN, adjusts the contrast and prints the result
  to STDOUT

    --ratio [n]         # Ratio to shift contrast by (default: 1.2)
`,
}

var contrastRatio float64

func init() {
	cmdContrast.Run = runContrast

	cmdContrast.Flag.Float64Var(&contrastRatio, "ratio", 1.2, "")
}

func runContrast(cmd *Command, args []string) {
	i := utils.ReadStdin()
	i  = contrast.Adjust(i, contrastRatio)
	utils.WriteStdout(i)
}
