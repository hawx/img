package cmd

import (
	"github.com/hawx/hadfield"
	"github.com/hawx/img/channel"
	"github.com/hawx/img/utils"
)

var (
	channelRed, channelGreen, channelBlue                              bool
	channelHue, channelSaturation, channelLightness, channelBrightness bool
	channelAlpha                                                       bool
	channelBy, channelRatio                                            float64
)

func Channel() *hadfield.Command {
	cmd := &hadfield.Command{
		Usage: "channel [options]",
		Short: "adjust the value of each colour channel individually",
		Long: `
  Channel allows you to adjust the value of each colour channel (red, green,
  blue or alpha) individually. All values are scaled to the range [0,1].
  Defaults to acting on the red, green and blue channels.

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

	cmd.Run = runChannel

	cmd.Flag.BoolVar(&channelRed, "red", false, "")
	cmd.Flag.BoolVar(&channelGreen, "green", false, "")
	cmd.Flag.BoolVar(&channelBlue, "blue", false, "")

	cmd.Flag.BoolVar(&channelHue, "hue", false, "")
	cmd.Flag.BoolVar(&channelSaturation, "saturation", false, "")
	cmd.Flag.BoolVar(&channelLightness, "lightness", false, "")
	cmd.Flag.BoolVar(&channelBrightness, "brightness", false, "")

	cmd.Flag.BoolVar(&channelAlpha, "alpha", false, "")

	cmd.Flag.Float64Var(&channelBy, "by", 0.1, "")
	cmd.Flag.Float64Var(&channelRatio, "ratio", 1.2, "")

	return cmd
}

func runChannel(cmd *hadfield.Command, args []string) {
	i, data := utils.ReadStdin()
	var adj utils.Adjuster

	if utils.FlagVisited("by", cmd.Flag) {
		adj = utils.Adder(channelBy)
	} else {
		adj = utils.Multiplier(channelRatio)
	}

	if !(channelRed || channelGreen || channelBlue || channelHue ||
		channelSaturation || channelLightness || channelBrightness || channelAlpha) {
		channelRed = true
		channelGreen = true
		channelBlue = true
	}

	if channelRed {
		i = channel.Adjust(i, adj, channel.Red)
	}
	if channelGreen {
		i = channel.Adjust(i, adj, channel.Green)
	}
	if channelBlue {
		i = channel.Adjust(i, adj, channel.Blue)
	}

	if channelHue {
		i = channel.Adjust(i, adj, channel.Hue)
	}
	if channelSaturation {
		i = channel.Adjust(i, adj, channel.Saturation)
	}
	if channelLightness {
		i = channel.Adjust(i, adj, channel.Lightness)
	}
	if channelBrightness {
		i = channel.Adjust(i, adj, channel.Brightness)
	}

	if channelAlpha {
		i = channel.Adjust(i, adj, channel.Alpha)
	}

	utils.WriteStdout(i, data)
}
