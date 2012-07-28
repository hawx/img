package main

import (
	"os"
	"math/rand"
	"image"
	"image/png"
)


func randBetween(a, b int) int {
	return rand.Intn(b - a)
}

func readImage(path string) image.Image {
	file, _ := os.Open(path)
	defer file.Close()
	img, _  := png.Decode(file)

	return img
}

func writeImage(img image.Image, path string) {
	file, _ := os.Create(path)
	defer file.Close()
	png.Encode(file, img)
}

func shuffle(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			xb := randBetween(b.Min.X, b.Max.X)
			yb := randBetween(b.Min.Y, b.Max.Y)

			a := img.At(x, y)
			b := img.At(xb, yb)

			o.Set(x, y, b)
			o.Set(xb, yb, a)
		}
	}

	return o
}

func verticalShuffle(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			yb := randBetween(b.Min.Y, b.Max.Y)

			a := img.At(x, y)
			b := img.At(x, yb)

			o.Set(x, y, b)
			o.Set(x, yb, a)
		}
	}

	return o
}

func horizontalShuffle(img image.Image) image.Image {
	b := img.Bounds()
	o := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			xb := randBetween(b.Min.X, b.Max.X)

			a := img.At(x, y)
			b := img.At(xb, y)

			o.Set(x, y, b)
			o.Set(xb, y, a)
		}
	}

	return o
}

func main() {
	i := readImage("image.png")
	i  = horizontalShuffle(i)
	writeImage(i, "output.png")
}
