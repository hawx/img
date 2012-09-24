package main

import (
	"github.com/hawx/img/hsla"
	"github.com/hawx/img/utils"
)

var cmdHue = &Command{
	UsageLine: "hue [options]",
	Short:     "adjust image hue",
Long: `
  Hue takes an image from STDIN, adjusts the hue and prints the result to
  STDOUT.

    --by [n]         # Amount to adjust hue by (default: 60.0)
`,
}

var hueBy float64

func init() {
	cmdHue.Run = runHue

	cmdHue.Flag.Float64Var(&hueBy, "by", 60.0, "")
}

func runHue(cmd *Command, args []string) {
	i := utils.ReadStdin()
	i  = hsla.Hue(i, hueBy)
	utils.WriteStdout(i)
}
