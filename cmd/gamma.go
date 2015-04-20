package cmd

import (
	"github.com/hawx/hadfield"
	"hawx.me/code/img/gamma"
	"hawx.me/code/img/utils"
)

var (
	gammaAuto, gammaUndo bool
	gammaBy              float64
)

func Gamma() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "gamma [options]",
		Short: "adjust image gamma",
		Long: `
  Gamma takes an image from STDIN, adjusts the gamma and prints the result to
  STDOUT

    --auto         # Automatically alter gamma to "best" value (default)
    --by <n>       # Amount to adjust gamma by
    --undo         # Adjust by the reciprocal of the amount, instead
`,
	}

	cmd.Run = runGamma

	cmd.Flag.BoolVar(&gammaAuto, "auto", false, "")
	cmd.Flag.Float64Var(&gammaBy, "by", 1.8, "")
	cmd.Flag.BoolVar(&gammaUndo, "undo", false, "")

	return cmd
}

func runGamma(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	if gammaUndo {
		gammaBy = 1.0 / gammaBy
	}

	if utils.FlagVisited("by", cmd.Flag) {
		i = gamma.Adjust(i, gammaBy)
	} else {
		i = gamma.Auto(i)
	}

	utils.WriteStdout(i, data)
}
