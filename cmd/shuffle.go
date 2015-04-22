package cmd

import (
	"hawx.me/code/hadfield"
	"hawx.me/code/img/shuffle"
	"hawx.me/code/img/utils"
)

var (
	shuffleVertical, shuffleHorizontal bool
)

func Shuffle() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "shuffle [options]",
		Short: "shuffles pixels of the image",
		Long: `
  Shuffle takes an image, shuffles the pixels of the image, then prints the
  result to STDOUT

    --horizontal     # Use horizontal shuffling only
    --vertical       # Use vertical shuffling only
`,
	}

	cmd.Run = runShuffle

	cmd.Flag.BoolVar(&shuffleVertical, "vertical", false, "")
	cmd.Flag.BoolVar(&shuffleHorizontal, "horizontal", false, "")

	return cmd
}

func runShuffle(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	if shuffleVertical && !shuffleHorizontal {
		i = shuffle.Vertically(i)
	} else if shuffleHorizontal && !shuffleVertical {
		i = shuffle.Horizontally(i)
	} else {
		i = shuffle.Shuffle(i)
	}

	utils.WriteStdout(i, data)
}
