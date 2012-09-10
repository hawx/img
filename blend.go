package main

import (
	"./blend"
	"./utils"
	"strings"
	"flag"
	"os"
	"fmt"
	"image"
	"image/png"
)

func subhead(s string) string {
	return "\033[90m" + s + "\033[0m"
}

func printHelp() {
	msg := "Usage: compose <other> [opts]\n" +
		"\n" +
		"  Takes a png file from STDIN and blends it with the png file at <other>\n" +
		"  using the method chosen. The result is printed to STDOUT.\n" +
		"\n" +
		"  Note: 'blend colour' refers to the colour taken from <other>, while 'base \n" +
		"  color' refers to the colour taken from STDIN.\n" +
		"\n" +
		"  --opacity      # Opacity of blended image layer (default: 0.5)\n" +
		"\n" +
		subhead("  BASIC\n") +
		"  --normal       # Paints pixels using <other> (default)\n" +
		"  --dissolve     # Paints pixels from <other> randomly, depending on opacity\n" +
		"\n" +
		subhead("  DARKEN\n") +
		"  --darken       # Selects the darkest value for each colour channel\n" +
		"  --multiply     # Multiplies each colour channel\n" +
		"  --burn         # Darkens the base colour to increase contrast\n" +
		"  --darker       # Selects the darkest colour by comparing the sum of channels\n" +
		"\n" +
		subhead("  LIGHTEN\n") +
		"  --lighten      # Selects the lightest value for each colour channel\n" +
		"  --screen       # Multiples the inverse of each colour channel\n" +
		"  --dodge        # Brightens the base colour to decrease contrast\n" +
		"  --lighter      # Selects the lightest colour by comparing the sum of channels\n" +
		"\n" +
		subhead("  CONTRAST\n") +
		"  --overlay      # Multiplies or screens the colours, depending on the base colour\n" +
		"  --soft-light   # Darkens or lightens the colours, depending on the blend colour\n" +
		"  --hard-light   # Multiplies or screens the colours, depending on the blend colour\n" +
		"\n" +
		subhead("  COMPARATIVE\n") +
		"  --difference   # Finds the absolute difference between the base and blend colour\n" +
		"  --exclusion    # Creates an effect similar to, but lower in contrast, than difference\n" +
		"  --addition     # Adds the blend colour to the base colour\n" +
		"  --subtraction  # Subtracts the blend colour from the base colour\n" +
		"\n" +
		subhead("  HSL\n") +
		"  --hue          # Uses just the hue of the blend colour\n" +
		"  --saturation   # Uses just the saturation of the blend colour\n" +
		"  --color        # Uses just the hue and saturation of the blend colour\n" +
		"  --luminosity   # Uses just the luminosity of the blend colour\n" +
		"\n" +
		"  --modes        # Lists all available modes\n" +
		"  --help         # Display this help message\n" +
		"\n"

	fmt.Fprintf(os.Stderr, msg)
	os.Exit(0)
}

func printModes() {
	modes := []string{
		"normal", "dissolve",
		"darken", "multiply", "burn", "darker",
		"lighten", "screen", "dodge", "lighter",
		"overlay", "soft-light", "hard-light",
		"difference", "exclusion", "addition", "subtraction",
		"hue", "saturation", "color", "luminosity",
	}

	msg := strings.Join(modes, "\n")

	fmt.Fprintf(os.Stdout, msg + "\n")
	os.Exit(0)
}


var help         = flag.Bool("help", false, "")
var modes        = flag.Bool("modes", false, "")
var opacity      = flag.Float64("opacity", 1.0, "")

// BASIC
var normalM      = flag.Bool("normal", false, "")
var dissolveM    = flag.Bool("dissolve", false, "")

// DARKEN
var	darkenM      = flag.Bool("darken", false, "")
var	multiplyM    = flag.Bool("multiply", false, "")
var	burnM        = flag.Bool("burn", false, "")
var darkerM      = flag.Bool("darker", false, "")

// LIGHTEN
var	lightenM     = flag.Bool("lighten", false, "")
var	screenM      = flag.Bool("screen", false, "")
var	dodgeM       = flag.Bool("dodge", false, "")
var lighterM     = flag.Bool("lighter", false, "")

// CONTRAST
var	overlayM     = flag.Bool("overlay", false, "")
var	softLightM   = flag.Bool("soft-light", false, "")
var	hardLightM   = flag.Bool("hard-light", false, "")

// COMPARATIVE
var	differenceM  = flag.Bool("difference", false, "")
var exclusionM   = flag.Bool("exclusion", false, "")
var	additionM    = flag.Bool("addition", false, "")
var	subtractionM = flag.Bool("subtraction", false, "")

// HSL
var	hueM         = flag.Bool("hue", false, "")
var	saturationM  = flag.Bool("saturation", false, "")
var	colorM       = flag.Bool("color", false, "")
var	luminosityM  = flag.Bool("luminosity", false, "")



func main() {
	flag.Parse()
	if *help { printHelp() }
	if *modes { printModes() }

	a := utils.ReadStdin()

	path := flag.Args()[0]
	file, _ := os.Open(path)
	b, _ := png.Decode(file)
	var f (func(a, b image.Image) image.Image)

	b = blend.Fade(b, *opacity)

	if *normalM {
		f = blend.Normal
	} else if *dissolveM {
		f = blend.Dissolve

	} else if *darkenM {
		f = blend.Darken
	} else if *multiplyM {
		f = blend.Multiply
	} else if *burnM {
		f = blend.Burn
	} else if *darkerM {
		f = blend.Darker

	} else if *lightenM {
		f = blend.Lighten
	} else if *screenM {
		f = blend.Screen
	} else if *dodgeM {
		f = blend.Dodge
	} else if *lighterM {
		f = blend.Lighter

	} else if *overlayM {
		f = blend.Overlay
	} else if *softLightM {
		f = blend.SoftLight
	} else if *hardLightM {
		f = blend.HardLight

	} else if *differenceM {
		f = blend.Difference
	} else if *exclusionM {
		f = blend.Exclusion
	} else if *additionM {
		f = blend.Addition
	} else if *subtractionM {
		f = blend.Subtraction

	} else if *hueM {
		f = blend.Hue
	} else if *saturationM {
		f = blend.Saturation
	} else if *colorM {
		f = blend.Color
	} else if *luminosityM {
		f = blend.Luminosity

	} else {
		f = blend.Normal
	}

	utils.WriteStdout(f(a, b))
}
