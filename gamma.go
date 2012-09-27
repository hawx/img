package main

import (
	"github.com/hawx/img/gamma"
	"github.com/hawx/img/utils"
)

var cmdGamma = &Command{
	UsageLine: "gamma [options]",
	Short:     "adjust image gamma",
Long: `
	Gamma takes an image from STDIN, adjusts the gamma and prints the result to
	STDOUT

	  --by [n]       # Amount to adjust gamma by (default: 1.8)
`,
}

var gammaBy float64

func init() {
	cmdGamma.Run = runGamma

	cmdGamma.Flag.Float64Var(&gammaBy, "by", 1.8, "")
}

func runGamma(cmd *Command, args []string) {
	i := utils.ReadStdin()
	i  = gamma.Adjust(i, gammaBy)
	utils.WriteStdout(i)
}
