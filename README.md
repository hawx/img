# img

A selection of image manipulation tools.

Requires Go 1. Install to `$GOPATH/bin` with,

``` bash
$ go install github.com/hawx/img
```

Use `go help` and `go help [command]` for information.

## shuffle

Randomly shuffles pixels around the image. Use `-v` or `-h` to constrain it to
vertical or horizontal shuffling, respectively.

``` bash
$ img shuffle --vertical < input.png > output.png
```

![Shuffle](http://github.com/hawx/img/raw/master/examples/shuffle.jpg)

## pixelate

Pixelates an image. Use `--size HxW` to set pixel size used.

``` bash
$ img pixelate --size 10x50 < input.png > output.png
```

![Pixelate](http://github.com/hawx/img/raw/master/examples/pixelate.jpg)

## pxl

Implementation of the triangle filter from [pxl app][pxlapp], using the
algorithm loosely described by [revdancatt][rev].

``` bash
$ img pxl --size 30x30 < input.png > output.png
```

![pxl](http://github.com/hawx/img/raw/master/examples/pxl.jpg)

## hxl

An (almost; that is I'm not sure this is exactly the same) implementation of the
equilateral triangle filter from [pxl app][pxlapp].

``` bash
$ img hxl --width 50 < input.png > output.png
```

![hxl](http://github.com/hawx/img/raw/master/examples/hxl.jpg)

## greyscale

Creates a greyscale version of an image.

``` bash
$ img greyscale --average < input.png > output.png
```

![Greyscale](http://github.com/hawx/img/raw/master/examples/greyscale.jpg)

## contrast

Adjusts the contrast of the given image.

``` bash
$ img contrast --by -25 < input.png > output.png
```

![contrast](http://github.com/hawx/img/raw/master/examples/contrast.jpg)

## brightness

Adjusts the brightness of the given image.

``` bash
$ img brightness --by -25 < input.png > output.png
```

![brightness](http://github.com/hawx/img/raw/master/examples/brightness.jpg)

## hue, saturation and lightness

Adjust the hue, saturation and lightness of the an image.

``` bash
$ img hue --by -30 < input.png > output.png
$ img saturation --by 0.3 < input.png > output.png
$ img lightness --by -0.07 < input.png > output.png
```

![hsl](http://github.com/hawx/img/raw/master/examples/hsl.jpg)

## blend

Allows you to blend two images together using a blend mode. Takes one image from
STDIN (the base image, imagine the bottom layer in photoshop) and one image as
an argument (the blend image, the layer above).

``` bash
$ img blend --screen blend.png --opacity 0.3 < input.png > output.png
```

![blend](http://github.com/hawx/img/raw/master/examples/blend.jpg)

# Composition

These tools have been created to do one task each, and to use standard
input/output so that they can be easily composed. For example;

``` bash
$ (./shuffle --horizontal | ./hxl | ./hue --by -20) < input.png > output.png
```

![Composed](http://github.com/hawx/img/raw/master/examples/composed.jpg)


# Notes on using the img package in go

You can easily use img in programmatically as well,

``` bash
$ go get github.com/hawx/img
$ cat > greyscaler.go
package main

import (
  "github.com/hawx/img/greyscale"
  "os"
  "image/png"
)

func main() {
  file, _ := os.Open(os.Args[1])
  img,  _ := png.Decode(file)

  img = greyscale.Average(img)

  out, _ := os.Create(os.Args[2])
  png.Encode(out, img)
}

$ go build greyscaler.go
$ ./greyscaler input.png output.png
```

To view documentation run `godoc -http=:8080` then navigate to
<http://localhost:8080/pkg/github.com/hawx/img/>, or see it on [GoPkgDoc][docs].


[pxlapp]: http://kohlberger.net/apps/pxl
[rev]:    http://revdancatt.com/2012/03/31/the-pxl-effect-with-javascript-and-canvas-and-maths/
[docs]:   http://go.pkgdoc.org/github.com/hawx/img
