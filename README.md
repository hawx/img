# img

A selection of image manipulation tools. Each tool requires an input file from
stdin, either `.png`, `.jpeg` or `.gif`. They always output a `.png` file,
regardless of input type, using stdout.

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

![Shuffle](http://hawx.github.com/img/examples/shuffle.jpg)

## pixelate

Pixelates an image. Use `--size HxW` to set pixel size used.

``` bash
$ img pixelate --size 10x50 < input.png > output.png
```

![Pixelate](http://hawx.github.com/img/examples/pixelate.jpg)

## pxl

Implementation of the triangle filter from [pxl app][pxlapp], using the
algorithm loosely described by [revdancatt][rev].

``` bash
$ img pxl --size 30x30 < input.png > output.png
```

![pxl](http://hawx.github.com/img/examples/pxl.jpg)

## hxl

An (almost; that is I'm not sure this is exactly the same) implementation of the
equilateral triangle filter from [pxl app][pxlapp].

``` bash
$ img hxl --width 50 < input.png > output.png
```

![hxl](http://hawx.github.com/img/examples/hxl.jpg)

## greyscale

Creates a greyscale version of an image.

``` bash
$ img greyscale --average < input.png > output.png
```

![Greyscale](http://hawx.github.com/img/examples/greyscale.jpg)

## contrast

Adjusts the contrast of the given image.

``` bash
$ img contrast --by -25 < input.png > output.png
```

![contrast](http://hawx.github.com/img/examples/contrast.jpg)

## brightness

Adjusts the brightness of the given image.

``` bash
$ img brightness --by -25 < input.png > output.png
```

![brightness](http://hawx.github.com/img/examples/brightness.jpg)

## hue, saturation and lightness

Adjust the hue, saturation and lightness of the an image.

``` bash
$ img hue --by -30 < input.png > output.png
$ img saturation --by 0.3 < input.png > output.png
$ img lightness --by -0.07 < input.png > output.png
```

![hsl](http://hawx.github.com/img/examples/hsl.jpg)

## blend

Allows you to blend two images together using a blend mode. Takes one image from
STDIN (the base image, imagine the bottom layer in photoshop) and one image as
an argument (the blend image, the layer above).

``` bash
$ img blend --screen blend.png --opacity 0.3 < input.png > output.png
```

![blend](http://hawx.github.com/img/examples/blend-modes.jpg)

## levels

Allows you to alter the levels of an image. You can set (or auto set) white and
black points, or pass a curve to use, along with the channels to act on.

``` bash
$ img levels --red --green --curve "0,20 50,40 100,100" < input.png > output.png
```

![levels](http://hawx.github.com/img/examples/levels.jpg)

# Composition

These tools have been created to do one task each, and to use standard
input/output so that they can be easily composed. For example;

``` bash
$ (img shuffle --horizontal | img hxl | img hue --by -20) < input.png > output.png
```

![Composed](http://hawx.github.com/img/examples/composed.jpg)



# External scripts

It is possible to extend `img` with external scripts. They must be named
`img-something`, and be somewhere on your `$PATH`. They must also respond to the
flags `--usage`, `--short` and `--long` so they can be integrated into `img
help`.

As an example of an external script, here is `img-lomo`, that applies a
lomography effect detailed in [Lomography, UNIX Style][tao].  It makes use of
[imagemagick][im] to generate a mask, then blends the images with img.

``` bash
#!/usr/bin/env sh

function usage {
  echo "lomo [options]"
}

function short {
  echo "applies a simple lomo effect to the image."
}

function long {
  echo "  Applies a simple lomo effect to the image, boosting its saturation and
  composing with a black edged mask."
}

case "$1" in
  --usage ) usage
    exit
    ;;
  --short ) short
    exit
    ;;
  --long ) long
    exit
    ;;
esac

# Using the method described on:
# http://the.taoofmac.com/space/blog/2005/08/23/2359

# generate a mask using imagemagick
convert -size 80x60 xc:black \
        -fill white -draw 'rectangle 1,1 78,58' \
        -gaussian 7x15 +matte lomo_mask.png
mogrify -resize 800x600 -gaussian 0x5 lomo_mask.png

# then put it together with img
(
  img contrast --ratio 1.2 |
  img saturation --ratio 1.2 |
  img blend --multiply --fit lomo_mask.png
)

# delete the mask
rm lomo_mask.png
```

See [this gist](https://gist.github.com/4566266) for a pure Go rewrite of this.

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
[tao]:    http://the.taoofmac.com/space/blog/2005/08/23/2359
[im]:     http://www.imagemagick.org
