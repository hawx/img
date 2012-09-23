package main

import (
	"github.com/hawx/img/brightness"
	"github.com/hawx/img/utils"
)

var cmdBrightness = &Command{
	UsageLine: "brightness [options]",
	Short:     "adjust image brightness",
Long: `
  Brightness takes a png file from STDIN, adjusts the brightness and prints the
  result to STDOUT

    --by [n]        # Amount to adjust brightness by (default: 15.0)
`,
}

var brightnessBy float64

func init() {
	cmdBrightness.Run = runBrightness

	cmdBrightness.Flag.Float64Var(&brightnessBy, "by", 20.0, "")
}

func runBrightness(cmd *Command, args []string) {
	i := utils.ReadStdin()
	i  = brightness.Adjust(i, brightnessBy)
	utils.WriteStdout(i)
}
