package main

import (
	"github.com/hawx/img/hsla"
	"github.com/hawx/img/utils"
)

var cmdSaturation = &Command{
	UsageLine: "saturation [options]",
	Short:     "adjust image saturation",
Long: `
  Saturation takes an image from STDIN, adjusts the saturation and prints the
  result to STDOUT

    --by [n]       # Amount to adjust saturation by
    --ratio [n]    # Ratio to adjust saturation by (default: 1.2)
`,
}

var saturationBy float64
var saturationRatio float64

func init() {
	cmdSaturation.Run = runSaturation

	cmdSaturation.Flag.Float64Var(&saturationBy, "by", 0.1, "")
	cmdSaturation.Flag.Float64Var(&saturationRatio, "ratio", 1.2, "")
}

func runSaturation(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if utils.FlagVisited("by", cmd.Flag) {
		i = hsla.Saturation(i, utils.Adder(saturationBy))
	} else {
		i = hsla.Saturation(i, utils.Multiplier(saturationRatio))
	}

	utils.WriteStdout(i)
}
