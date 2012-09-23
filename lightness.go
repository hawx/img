package main

import (
	"github.com/hawx/img/hsla"
	"github.com/hawx/img/utils"
)

var cmdLightness = &Command{
	UsageLine: "lightness [options]",
	Short:     "adjust image lightness",
Long: `
  Lightness takes a png file from STDIN, adjusts the lightness and prints the
  result to STDOUT

    --by [n]         # Amount to adjust lightness by (default: 0.1)
`,
}

var lightnessBy float64

func init() {
	cmdLightness.Run = runLightness

	cmdLightness.Flag.Float64Var(&lightnessBy, "by", 0.1, "")
}

func runLightness(cmd *Command, args []string) {
	i := utils.ReadStdin()
	i  = hsla.Lightness(i, lightnessBy)
	utils.WriteStdout(i)
}
