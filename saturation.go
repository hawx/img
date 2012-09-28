package main

import (
	"github.com/hawx/img/hsla"
	"github.com/hawx/img/utils"
	"flag"
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

	cmd.Flag.Visit(func (f *flag.Flag) {
		if f == cmd.Flag.Lookup("by") {
			i = hsla.Saturation(i, func(i float64) float64 {
				return i + saturationBy
			})
		} else { // default is ratio
			i = hsla.Saturation(i, func(i float64) float64 {
				return i * saturationRatio
			})
		}
	})

	utils.WriteStdout(i)
}
