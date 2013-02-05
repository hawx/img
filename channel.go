package main

import (
	"github.com/hawx/img/channel"
	"github.com/hawx/img/utils"
)

var cmdChannel = &Command{
	UsageLine: "channel [options]",
	Short:     "adjust the value of each colour channel individually",
Long: `
  Channel allows you to adjust the value of each colour channel (red, green,
  blue or alpha) individually. Defaults to red, green and blue.

    --red             # Apply changes to red channel
    --green           # Apply changes to green channel
    --blue            # Apply changes to blue channel

    --hue             # Apply changes to hue channel
    --saturation      # Apply changes to saturation channel
    --lightness       # Apply changes to lightness channel
    --brightness      # Apply changes to brightness channel

    --alpha           # Apply changes to alpha channel

    --by [n]          # Amount to adjust value by
    --ratio [n]       # Ratio to adjust value by (default: 1.2)
`,
}

var channelRed, channelGreen, channelBlue bool
var channelHue, channelSaturation, channelLightness, channelBrightness bool
var channelAlpha bool
var channelBy, channelRatio float64

func init() {
	cmdChannel.Run = runChannel

	cmdChannel.Flag.BoolVar(&channelRed,        "red", false, "")
	cmdChannel.Flag.BoolVar(&channelGreen,      "green", false, "")
	cmdChannel.Flag.BoolVar(&channelBlue,       "blue", false, "")

	cmdChannel.Flag.BoolVar(&channelHue,        "hue", false, "")
	cmdChannel.Flag.BoolVar(&channelSaturation, "saturation", false, "")
	cmdChannel.Flag.BoolVar(&channelLightness,  "lightness", false, "")
	cmdChannel.Flag.BoolVar(&channelBrightness, "brightness", false, "")

	cmdChannel.Flag.BoolVar(&channelAlpha,      "alpha", false, "")

	cmdChannel.Flag.Float64Var(&channelBy,      "by", 0.1, "")
	cmdChannel.Flag.Float64Var(&channelRatio,   "ratio", 1.2, "")
}

func runChannel(cmd *Command, args []string) {
	i := utils.ReadStdin()
	var adj utils.Adjuster

	if utils.FlagVisited("by", cmd.Flag) {
		adj = utils.Adder(channelBy)
	} else {
		adj = utils.Multiplier(channelRatio)
	}

	if !(channelRed || channelGreen || channelBlue || channelHue ||
		channelSaturation || channelLightness || channelBrightness || channelAlpha) {
		channelRed   = true
		channelGreen = true
		channelBlue  = true
	}

	if channelRed        { i = channel.Red(i, adj) }
	if channelGreen      { i = channel.Green(i, adj) }
	if channelBlue       { i = channel.Blue(i, adj) }

	if channelHue        { i = channel.Hue(i, adj) }
	if channelSaturation { i = channel.Saturation(i, adj) }
	if channelLightness  { i = channel.Lightness(i, adj) }
	if channelBrightness { i = channel.Brightness(i, adj) }

	if channelAlpha      { i = channel.Alpha(i, adj) }

	utils.WriteStdout(i)
}
