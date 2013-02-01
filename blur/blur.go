package blur

import (
	"github.com/hawx/img/utils"
	"image"
	"image/color"
	"math"
	"os"
)

type Style float64

const (
	// Ignore edges, may leave them semi-transparent
	IGNORE Style = iota
	// Clamp edges, may leave them looking unblurred
	CLAMP
	// Wrap edges, may change colour of edges
	WRAP
)


func abs(num int) int {
	if num < 0 { return -num }
	return num
}

func correct(num int) bool {
	return (num > 0) && (num % 2 != 0)
}

// A Kernel is a 2-dimensional array of ratios. A 1-dimensional, horizontal or
// vertical, kernel can easily be defined by a 2-dimensional array in the
// obvious manner. The weights are taken as row by column, so kernel[0]
// references the first row and kernel[i][0] (for all i) is the first column.
type Kernel [][]float64

func NewHorizontalKernel(width int, f func(x int) float64) Kernel {
	if !correct(width) {
		utils.Warn("Error: kernel size wrong!")
		os.Exit(2)
	}

	mx := (width - 1) / 2
	k  := [][]float64{ make([]float64, width) }

	for x := 0; x < width; x++ {
		k[0][x] = f(mx - x)
	}

	return k
}

func NewVerticalKernel(height int, f func(y int) float64) Kernel {
	if !correct(height) {
		utils.Warn("Error: kernel size wrong!")
		os.Exit(2)
	}

	my := (height - 1) / 2
	k  := make([][]float64, height)

	for y := 0; y < height; y++ {
		k[y] = []float64{ f(my - y) }
	}

	return k
}

// NewKernel creates a new Kernel of the dimensions given, it is populated by
// the given function which itself is passed the signed x and y offsets from the
// mid point.
func NewKernel(height, width int, f func(x, y int) float64) Kernel {
	if !correct(width) || !correct(height) {
		utils.Warn("Error: kernel size wrong!")
		os.Exit(2)
		// should return error really!
	}

	mx := (width - 1) / 2
	my := (height - 1) / 2
	k := make([][]float64, height)

	for y := 0; y < height; y++ {
		k[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			k[y][x] = f(mx - x, my - y)
		}
	}

	return k
}

// Normalised returns a copy of the Kernel where the sum of all entries is 1.
func (k Kernel) Normalised() Kernel {
	total := 0.0
	for y := 0; y < k.Height(); y++ {
		for x := 0; x < k.Width(); x++ {
			total += k[y][x]
		}
	}

	nk := make([][]float64, k.Height())

	for y := 0; y < k.Height(); y++ {
		nk[y] = make([]float64, k.Width())
		for x := 0; x < k.Width(); x++ {
			nk[y][x] = k[y][x] / total
		}
	}

	return nk
}

// Height returns the height of the Kernel.
func (k Kernel) Height() int {
	return len(k)
}

// Width returns the width of the Kernel.
func (k Kernel) Width() int {
	if k.Height() > 0 { return len(k[0]) }
	return 0
}

// Mid returns the centre Point of the Kernel.
func (k Kernel) Mid() image.Point {
	return image.Pt((k.Width() - 1) / 2, (k.Height() - 1) / 2)
}

func Convolve(in image.Image, weights Kernel, style Style) image.Image {
	bnds := in.Bounds()
	mid := weights.Mid()
	o := image.NewRGBA(bnds)

	for y := bnds.Min.Y; y < bnds.Max.Y; y++ {
		for x := bnds.Min.X; x < bnds.Max.X; x++ {
			var r, g, b, a, offset float64

			for oy := 0; oy < weights.Height(); oy++ {
				for ox := 0; ox < weights.Width(); ox++ {
					factor := weights[oy][ox]
					pt := image.Pt(x + ox - mid.X, y + oy - mid.Y)

					if pt == weights.Mid() {
						// Ignore!

					} else if pt.In(bnds) {
						or,og,ob,oa := in.At(pt.X, pt.Y).RGBA()

						r += float64(or) * factor
						g += float64(og) * factor
						b += float64(ob) * factor
						a += float64(oa) * factor

					} else {
						switch style {
						case CLAMP:
							offset += factor

						case WRAP:
							if pt.X >= bnds.Max.X {
								pt.X = pt.X - bnds.Max.X
							} else if pt.X < bnds.Min.X {
								pt.X = bnds.Dx() + pt.X
							}

							if pt.Y >= bnds.Max.Y {
								pt.Y = pt.Y - bnds.Max.Y
							} else if pt.Y < bnds.Min.Y {
								pt.Y = bnds.Dy() + pt.Y
							}

							or,og,ob,oa := in.At(pt.X, pt.Y).RGBA()

							r += float64(or) * factor
							g += float64(og) * factor
							b += float64(ob) * factor
							a += float64(oa) * factor
						}
					}
				}
			}

			if offset != 0 && style == CLAMP {
				or,og,ob,oa := in.At(x,y).RGBA()
				r += float64(or) * offset
				g += float64(og) * offset
				b += float64(ob) * offset
				a += float64(oa) * offset
			}

			o.Set(x, y, color.RGBA{
				uint8(utils.Truncatef(r / 255)),
				uint8(utils.Truncatef(g / 255)),
				uint8(utils.Truncatef(b / 255)),
				uint8(utils.Truncatef(a / 255)),
			})
		}
	}

	return o
}

// Perform a convolution with two Kernels in succession.
func Convolve2(in image.Image, a, b Kernel, style Style) image.Image {
	return Convolve(Convolve(in, a, style), b, style)
}


func Box(in image.Image, size utils.Pixel, style Style) image.Image {
	f := func(n int) float64 { return 1.0 }

	tall := NewVerticalKernel(size.H, f).Normalised()
	wide := NewHorizontalKernel(size.W, f).Normalised()

	return Convolve2(in, tall, wide, style)
}

func Gaussian(in image.Image, size utils.Pixel, sigma float64, style Style) image.Image {
	f := func(n int) float64 {
		return math.Exp( -float64(n*n) / (2 * sigma*sigma) )
	}

	tall := NewVerticalKernel(size.H, f).Normalised()
	wide := NewHorizontalKernel(size.W, f).Normalised()

	return Convolve2(in, tall, wide, style)
}
