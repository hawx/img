package cmd

import (
	"image"

	"github.com/hawx/hadfield"
	"github.com/hawx/img/channel"
	"github.com/hawx/img/levels"
	"github.com/hawx/img/utils"
)

var (
	levelsRed, levelsGreen, levelsBlue           bool
	levelsAuto, levelsAutoBlack, levelsAutoWhite bool
	levelsBlack, levelsWhite                     float64
	levelsCurve                                  string
)

func Levels() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "levels [options]",
		Short: "adjust image levels",
		Long: `
  Levels adjusts the levels of an image. You can set the black/white point or
  give points in a curve to use. It can act on all colour channels (default
  behaviour), or just those you want.

    --red             # Act on red channel
    --green           # Act on green channel
    --blue            # Act on blue channel

    --auto            # Auto adjust levels to fit
    --auto-black      # Auto adjust black point to fit
    --auto-white      # Auto adjust white point to fit

    --black [n]       # Set black point
    --white [n]       # Set white point

    --curve [c]       # Set curve. Argument is a list of 'point,value' pairs
                      # delimited by spaces, eg. --curve "0,0 33,40 66,60 100,100"
`,
	}

	cmd.Run = runLevels

	cmd.Flag.BoolVar(&levelsRed, "red", false, "")
	cmd.Flag.BoolVar(&levelsGreen, "green", false, "")
	cmd.Flag.BoolVar(&levelsBlue, "blue", false, "")

	cmd.Flag.BoolVar(&levelsAuto, "auto", false, "")
	cmd.Flag.BoolVar(&levelsAutoBlack, "auto-black", false, "")
	cmd.Flag.BoolVar(&levelsAutoWhite, "auto-white", false, "")

	cmd.Flag.Float64Var(&levelsBlack, "black", 0, "")
	cmd.Flag.Float64Var(&levelsWhite, "white", 100, "")

	cmd.Flag.StringVar(&levelsCurve, "curve", "", "")

	return cmd
}

func runLevels(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()

	if !levelsRed && !levelsGreen && !levelsBlue {
		levelsRed = true
		levelsGreen = true
		levelsBlue = true
	}

	if levelsRed {
		i = runLevelsOnChannel(cmd, args, i, channel.Red)
	}

	if levelsGreen {
		i = runLevelsOnChannel(cmd, args, i, channel.Green)
	}

	if levelsBlue {
		i = runLevelsOnChannel(cmd, args, i, channel.Blue)
	}

	utils.WriteStdout(i, data)
}

func runLevelsOnChannel(cmd *hadfield.Command, args []string, img image.Image,
	ch channel.Channel) image.Image {

	if levelsAuto {
		img = levels.Auto(img, ch)

	} else if levelsAutoBlack {
		img = levels.AutoBlack(img, ch)

	} else if levelsAutoWhite {
		img = levels.AutoWhite(img, ch)

	} else if utils.FlagVisited("black", cmd.Flag) {
		img = levels.SetBlack(img, ch, levelsBlack)

	} else if utils.FlagVisited("white", cmd.Flag) {
		img = levels.SetWhite(img, ch, levelsWhite)

	} else if utils.FlagVisited("curve", cmd.Flag) {
		img = levels.SetCurve(img, ch, levels.ParseCurveString(levelsCurve))
	}

	return img
}
