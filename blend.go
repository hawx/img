package main

import (
	"github.com/nfnt/resize"
	"github.com/hawx/img/blend"
	"github.com/hawx/img/utils"
	"strings"
	"os"
	"fmt"
	"image"
	_ "image/png"
	_ "image/gif"
	_ "image/jpeg"
)

var cmdBlend = &Command{
	UsageLine: "blend <other> [options]",
	Short:     "blends two images together",
Long: `
  Blend takes an image file from STDIN (referred to as the 'base image') and
  another given as <other> (referred to as the 'blend image'), and blends them
  together using the method selected producing a result image which is printed
  to STDOUT.

    --modes          # List all available modes
    --opacity [n]    # Opacity of blend image layer (default: 1.0)
    --fit            # Fit the blend layer to the base layer, may result in loss of quality

    BASIC
    --normal         # Selects the blend image (default)
    --dissolve       # Randomly selects the blend image, depending on its opacity

    DARKEN
    --darken         # Selects the darkest value for each colour channel
    --multiply       # Multiples each colour channel
    --burn           # Darkens the base image to increase contrast
    --linear-burn    # Adds the blend colour to the base colour, then subtracts white
    --darker         # Selects the darkest colour by comparing the sum of channels

    LIGHTEN
    --lighten        # Selects the lightest value for each colour channel
    --screen         # Multiples the inverse of each colour channel
    --dodge          # Brightens the base image to decrease contrast
    --linear-dodge   # Adds the blend colour to the base colour
    --lighter        # Selects the lightest colour by comparing the sum of channels

    CONTRAST
    --overlay        # Multiplies or screens the colours, depending on the base colour
    --soft-light     # Darkens or lighten the colours, depending on the blend colour
    --hard-light     # Multiplies or screens the colours, depending on the blend colour
    --vivid-light    # Burns or dodges the colours, depending on the blend colour
    --linear-light   # Linear burns or dodges the colours, depending on the blend colour
    --pin-light      # Replaces the colours, depending on the blend colour
    --hard-mix       # Makes all pixels red, green, blue, white or black

    COMPARATIVE
    --difference     # Finds the absolute difference between the base and blend colour
    --exclusion      # Creates an effect similar to but lower in contrast than difference
    --subtraction    # Subtracts the blend colour from the base colour

    HSL
    --hue            # Uses just the hue of the blend colour
    --saturation     # Uses just the saturation of the blend colour
    --color          # Uses just the hue and saturation of the blend colour
    --luminosity     # Uses just the luminosity of the blend colour
`,
}

var blendModes, blendFit bool
var blendOpacity float64
var blendNormal, blendDissolve bool
var blendDarken, blendMultiply, blendBurn, blendLinearBurn, blendDarker bool
var blendLighten, blendScreen, blendDodge, blendLinearDodge, blendLighter bool
var blendOverlay, blendSoftLight, blendHardLight, blendVividLight bool
var blendLinearLight, blendPinLight, blendHardMix bool
var blendDifference, blendExclusion, blendAddition, blendSubtraction bool
var blendHue, blendSaturation, blendColor, blendLuminosity bool

func init() {
	cmdBlend.Run = runBlend

	cmdBlend.Flag.BoolVar(&blendModes, "modes", false, "")
  cmdBlend.Flag.Float64Var(&blendOpacity, "opacity", 1.0, "")
	cmdBlend.Flag.BoolVar(&blendFit, "fit", false, "")

	// BASIC
  cmdBlend.Flag.BoolVar(&blendNormal, "normal", false, "")
  cmdBlend.Flag.BoolVar(&blendDissolve, "dissolve", false, "")

	// DARKEN
  cmdBlend.Flag.BoolVar(&blendDarken, "darken", false, "")
  cmdBlend.Flag.BoolVar(&blendMultiply, "multiply", false, "")
  cmdBlend.Flag.BoolVar(&blendBurn, "burn", false, "")
	cmdBlend.Flag.BoolVar(&blendLinearBurn, "linear-burn", false, "")
  cmdBlend.Flag.BoolVar(&blendDarker, "darker", false, "")

	// LIGHTEN
  cmdBlend.Flag.BoolVar(&blendLighten, "lighten", false, "")
  cmdBlend.Flag.BoolVar(&blendScreen, "screen", false, "")
  cmdBlend.Flag.BoolVar(&blendDodge, "dodge", false, "")
	cmdBlend.Flag.BoolVar(&blendLinearDodge, "linear-dodge", false, "")
  cmdBlend.Flag.BoolVar(&blendLighter, "lighter", false, "")

	// CONTRAST
  cmdBlend.Flag.BoolVar(&blendOverlay, "overlay", false, "")
  cmdBlend.Flag.BoolVar(&blendSoftLight, "soft-light", false, "")
  cmdBlend.Flag.BoolVar(&blendHardLight, "hard-light", false, "")
	cmdBlend.Flag.BoolVar(&blendVividLight, "vivid-light", false, "")
	cmdBlend.Flag.BoolVar(&blendLinearLight, "linear-light", false, "")
	cmdBlend.Flag.BoolVar(&blendPinLight, "pin-light", false, "")
	cmdBlend.Flag.BoolVar(&blendHardMix, "hard-mix", false, "")

	// COMPARATIVE
  cmdBlend.Flag.BoolVar(&blendDifference, "difference", false, "")
  cmdBlend.Flag.BoolVar(&blendExclusion, "exclusion", false, "")
  cmdBlend.Flag.BoolVar(&blendAddition, "addition", false, "") // leave as alias
  cmdBlend.Flag.BoolVar(&blendSubtraction, "subtraction", false, "")

	// HSL
	cmdBlend.Flag.BoolVar(&blendHue, "hue", false, "")
  cmdBlend.Flag.BoolVar(&blendSaturation, "saturation", false, "")
  cmdBlend.Flag.BoolVar(&blendColor, "color", false, "")
  cmdBlend.Flag.BoolVar(&blendLuminosity, "luminosity", false, "")
}

func runBlend(cmd *Command, args []string) {
	if blendModes { printModes() }

	a := utils.ReadStdin()

	path := args[0]
	file, _ := os.Open(path)
	defer file.Close()
	b, _, _ := image.Decode(file)
	var f (func(a, b image.Image) image.Image)

	b = blend.Fade(b, blendOpacity)

	if blendFit {
		ab := a.Bounds()
		bb := b.Bounds()

		// Need to work this out better, see
		// http://www.codinghorror.com/blog/2007/07/better-image-resizing.html
		if bb.Dx() < ab.Dx() || bb.Dy() < ab.Dy() {
			// b is going to get BIGGER
			b = resize.Resize(uint(ab.Dx()), uint(ab.Dy()), b, resize.Bilinear)
		} else {
			// b is going to get SMALLER
			b = resize.Resize(uint(ab.Dx()), uint(ab.Dy()), b, resize.Bicubic)
		}
	}

	if blendNormal {
		f = blend.Normal
	} else if blendDissolve {
		f = blend.Dissolve

	} else if blendDarken {
		f = blend.Darken
	} else if blendMultiply {
		f = blend.Multiply
	} else if blendBurn {
		f = blend.Burn
	} else if blendLinearBurn {
		f = blend.LinearBurn
	} else if blendDarker {
		f = blend.Darker

	} else if blendLighten {
		f = blend.Lighten
	} else if blendScreen {
		f = blend.Screen
	} else if blendDodge {
		f = blend.Dodge
	} else if blendLinearDodge {
		f = blend.LinearDodge
	} else if blendLighter {
		f = blend.Lighter

	} else if blendOverlay {
		f = blend.Overlay
	} else if blendSoftLight {
		f = blend.SoftLight
	} else if blendHardLight {
		f = blend.HardLight
	} else if blendVividLight {
		f = blend.VividLight
	} else if blendLinearLight {
		f = blend.LinearLight
	} else if blendPinLight {
		f = blend.PinLight
	} else if blendHardMix {
		f = blend.HardMix

	} else if blendDifference {
		f = blend.Difference
	} else if blendExclusion {
		f = blend.Exclusion
	} else if blendAddition {
		f = blend.Addition
	} else if blendSubtraction {
		f = blend.Subtraction

	} else if blendHue {
		f = blend.Hue
	} else if blendSaturation {
		f = blend.Saturation
	} else if blendColor {
		f = blend.Color
	} else if blendLuminosity {
		f = blend.Luminosity

	} else {
		f = blend.Normal
	}

	utils.WriteStdout(f(a, b))
}

func printModes() {
	modes := []string{
		"normal", "dissolve",
		"darken", "multiply", "burn", "linear-burn", "darker",
		"lighten", "screen", "dodge", "linear-dodge", "lighter",
		"overlay", "soft-light", "hard-light", "vivid-light", "linear-light",
		"pin-light", "hard-mix",
		"difference", "exclusion", "addition", "subtraction",
		"hue", "saturation", "color", "luminosity",
	}

	msg := strings.Join(modes, "\n")

	fmt.Fprintf(os.Stdout, msg + "\n")
	os.Exit(0)
}
