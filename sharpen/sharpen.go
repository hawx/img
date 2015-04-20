// Package sharpen provides functions to sharpen an image.
package sharpen

import (
	"hawx.me/code/img/blur"
	"hawx.me/code/img/utils"

	"image"
	"image/color"
	"math"
)

// Sharpen takes an image and sharpens it by, essentially, unblurring it. It is
// currently extremely slow, so you are probably better off sticking to
// UnsharpMask.
func Sharpen(in image.Image, radius int, sigma float64) image.Image {
	// Copied from ImageMagick, obvs.
	//
	// Sharpens the image. Convolve the image with a Gaussian operator of the
	// given radius and standard deviation (sigma). For reasonable results radius
	// should be larger than sigma.
	//
	// Using a seperable kernel would be faster, but the negative weights cancel
	// out on the corners of the kernel producing often undesirable ringing in the
	// filtered result; this can be avoided by using a 2D gaussian shaped image
	// sharpening kernel instead.
	//
	// (-exp(-(u*u + v*v) / 2.0 * sigma*sigma)) / (2.0 * Pi * sigma*sigma)

	normalize := 0.0
	f := func(u, v int) float64 {
		usq := float64(u * u)
		vsq := float64(v * v)
		val := -math.Exp(-(usq+vsq)/(2.0*sigma*sigma)) / (2.0 * math.Pi * sigma * sigma)
		normalize += val
		return val
	}

	k := blur.NewKernel(radius*2+1, radius*2+1, f)
	k[radius+1][radius+1] = -2.0 * normalize

	return blur.Convolve(in, k, blur.CLAMP)
}

// UnsharpMask sharpens the given Image using the unsharp mask technique.
// Basically the image is blurred, then subtracted from the original for
// differences above the threshold value.
func UnsharpMask(in image.Image, radius int, sigma, amount, threshold float64) image.Image {
	blurred := blur.Gaussian(in, radius, sigma, blur.IGNORE)
	bounds := in.Bounds()
	out := image.NewRGBA(bounds)

	// Absolute difference between a and b, returns float64 between 0 and 1.
	diff := func(a, b float64) float64 {
		if a > b {
			return a - b
		}
		return b - a
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			ar, ag, ab, aa := utils.RatioRGBA(in.At(x, y))
			br, bg, bb, _ := utils.RatioRGBA(blurred.At(x, y))

			if diff(ar, br) >= threshold {
				ar = amount*(ar-br) + ar
			}

			if diff(ag, bg) >= threshold {
				ag = amount*(ag-bg) + ag
			}

			if diff(ab, bb) >= threshold {
				ab = amount*(ab-bb) + ab
			}

			out.Set(x, y, color.NRGBA{
				uint8(utils.Truncatef(ar * 255)),
				uint8(utils.Truncatef(ag * 255)),
				uint8(utils.Truncatef(ab * 255)),
				uint8(aa * 255),
			})
		}
	}

	return out
}
