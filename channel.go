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
  blue or alpha) individually.

    --red             # Apply changes to red channel
    --green           # Apply changes to green channel
    --blue            # Apply changes to blue channel
    --alpha           # Apply changes to alpha channel

    --by [n]          # Amount to adjust value by
    --ratio [n]       # Ratio to adjust value by (default: 1.2)
`,
}

var channelRed, channelGreen, channelBlue, channelAlpha bool
var channelBy, channelRatio float64

func init() {
	cmdChannel.Run = runChannel

	cmdChannel.Flag.BoolVar(&channelRed, "red", false, "")
	cmdChannel.Flag.BoolVar(&channelGreen, "green", false, "")
	cmdChannel.Flag.BoolVar(&channelBlue, "blue", false, "")
	cmdChannel.Flag.BoolVar(&channelAlpha, "alpha", false, "")

	cmdChannel.Flag.Float64Var(&channelBy, "by", 0.1, "")
	cmdChannel.Flag.Float64Var(&channelRatio, "ratio", 1.2, "")
}

func runChannel(cmd *Command, args []string) {
	i := utils.ReadStdin()
	var adj utils.Adjuster

	if utils.FlagVisited("by", cmd.Flag) {
		adj = utils.Adder(channelBy)
	} else {
		adj = utils.Multiplier(channelRatio)
	}

	if channelRed   { i = channel.Red(i, adj) }
	if channelGreen { i = channel.Green(i, adj) }
	if channelBlue  { i = channel.Blue(i, adj) }
	if channelAlpha { i = channel.Alpha(i, adj) }

	utils.WriteStdout(i)
}
