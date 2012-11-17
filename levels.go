package main

import (
	"github.com/hawx/img/levels"
	"github.com/hawx/img/utils"
	"image"
)

var cmdLevels = &Command{
	UsageLine: "levels [options]",
	Short:     "adjust image levels (and stuff)",
	Long: `
  Levels adjusts the levels of an image. You can set the black/white point or
  give points in a curve to use. It can act on all colour channels (default
  behaviour), or just those you want.

    --red             # Only act on red channel
    --green           # Only act on green channel
    --blue            # Only act on blue channel

    --auto            # Auto adjust levels to fit
    --auto-black      # Auto adjust black point to fit
    --auto-white      # Auto adjust white point to fit

    --black [n]       # Set black point
    --white [n]       # Set white point

    --curve [c]       # Set curve. Argument is a list of 'point,value' pairs
                      # delimited by colons, eg. --curve 0,0:33,40:66,60:100,100
`,
}

var levelsRed, levelsGreen, levelsBlue bool
var levelsAuto, levelsAutoBlack, levelsAutoWhite bool
var levelsBlack, levelsWhite float64
var levelsCurve string

func init() {
	cmdLevels.Run = runLevels

	cmdLevels.Flag.BoolVar(&levelsRed, "red", false, "")
	cmdLevels.Flag.BoolVar(&levelsGreen, "green", false, "")
	cmdLevels.Flag.BoolVar(&levelsBlue, "blue", false, "")

	cmdLevels.Flag.BoolVar(&levelsAuto, "auto", false, "")
	cmdLevels.Flag.BoolVar(&levelsAutoBlack, "auto-black", false, "")
	cmdLevels.Flag.BoolVar(&levelsAutoWhite, "auto-white", false, "")

	cmdLevels.Flag.Float64Var(&levelsBlack, "black", 0, "")
	cmdLevels.Flag.Float64Var(&levelsWhite, "white", 100, "")

	cmdLevels.Flag.StringVar(&levelsCurve, "curve", "", "")
}

func runLevels(cmd *Command, args []string) {
	i := utils.ReadStdin()

	if !levelsRed && !levelsGreen && !levelsBlue {
		levelsRed = true; levelsGreen = true; levelsBlue = true
	}

	if levelsRed {
		i = runLevelsOnChannel(cmd, args, i, levels.RedChannel)
	}

	if levelsGreen {
		i = runLevelsOnChannel(cmd, args, i, levels.GreenChannel)
	}

	if levelsBlue {
		i = runLevelsOnChannel(cmd, args, i, levels.BlueChannel)
	}

	utils.WriteStdout(i)
}

func runLevelsOnChannel(cmd *Command, args []string, img image.Image,
	ch levels.Channel) image.Image {

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
