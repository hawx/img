package main

import (
	"github.com/hawx/img/brightness"
	"github.com/hawx/img/utils"
)

var cmdBrightness = &Command{
	UsageLine: "brightness [options]",
	Short:     "adjust image brightness",
Long: `
  Brightness takes an image from STDIN, adjusts the brightness and prints the
  result to STDOUT

    --by [n]        # Amount to adjust brightness by
    --ratio [n]     # Ratio to adjust brightness by (default: 1.2)
`,
}

var brightnessBy float64
var brightnessRatio float64

func init() {
	cmdBrightness.Run = runBrightness

	cmdBrightness.Flag.Float64Var(&brightnessBy, "by", 0.1, "")
	cmdBrightness.Flag.Float64Var(&brightnessRatio, "ratio", 1.2, "")
}

func runBrightness(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if utils.FlagVisited("by", cmd.Flag) {
		i = brightness.Adjust(i, utils.Adder(brightnessBy))
	} else {
		i = brightness.Adjust(i, utils.Multiplier(brightnessRatio))
	}

	utils.WriteStdout(i)
}
