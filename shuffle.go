package main

import (
	"github.com/hawx/img/shuffle"
	"github.com/hawx/img/utils"
)

var cmdShuffle = &Command{
	UsageLine: "shuffle [options]",
	Short:     "shuffles pixels of the image",
Long: `
  Shuffle takes an image, shuffles the pixels of the image, then prints the
  result to STDOUT

    --horizontal     # Use horizontal shuffling only
    --vertical       # Use vertical shuffling only
`,
}

var shuffleVertical, shuffleHorizontal bool

func init() {
	cmdShuffle.Run = runShuffle

	cmdShuffle.Flag.BoolVar(&shuffleVertical, "vertical", false, "")
	cmdShuffle.Flag.BoolVar(&shuffleHorizontal, "horizontal", false, "")
}

func runShuffle(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if (shuffleVertical && !shuffleHorizontal) {
		i = shuffle.Vertically(i)
	} else if (shuffleHorizontal && !shuffleVertical) {
		i = shuffle.Horizontally(i)
	} else {
		i = shuffle.Shuffle(i)
	}

	utils.WriteStdout(i)
}
