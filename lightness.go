package main

import (
	"github.com/hawx/img/hsla"
	"github.com/hawx/img/utils"
)

var cmdLightness = &Command{
	UsageLine: "lightness [options]",
	Short:     "adjust image lightness",
Long: `
  Lightness takes an image from STDIN, adjusts the lightness and prints the
  result to STDOUT

    --by [n]         # Amount to adjust lightness by
    --ratio [n]      # Ratio to adjust lightness by (default: 1.2)
`,
}

var lightnessBy float64
var lightnessRatio float64

func init() {
	cmdLightness.Run = runLightness

	cmdLightness.Flag.Float64Var(&lightnessBy, "by", 0.1, "")
	cmdLightness.Flag.Float64Var(&lightnessRatio, "ratio", 1.2, "")
}

func runLightness(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if utils.FlagVisited("by", cmd.Flag) {
		i = hsla.Lightness(i, utils.Adder(lightnessBy))
	} else {
		i = hsla.Lightness(i, utils.Multiplier(lightnessRatio))
	}

	utils.WriteStdout(i)
}
