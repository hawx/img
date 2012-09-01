# img

A selection of image manipulation tools.

Requires Go 1. Build all tools into `built/` by running,

``` bash
$ ./build
```

All tools respond to the `--help` flag, use it to get information on options
that are supported.

## Conversion with 'from' and 'to'

`from` takes an input image (jpeg, gif or png) and outputs a `.png` file.

``` bash
$ ./from image.jpg > image.png
```

`to` takes a png image and outputs some other file (jpeg or png).

``` bash
$ image.png < ./to image.jpg
```

These allow you to use other filetypes as input, and get different output:

``` bash
$ ./from input.jpg | ./greyscale | ./pxl | ./to output.jpg
```

## shuffle

Randomly shuffles pixels around the image. Use `-v` or `-h` to constrain it to
vertical or horizontal shuffling, respectively.

``` bash
$ ./shuffle --vertical < input.png > output.png
```

![Shuffle](http://github.com/hawx/img/raw/master/examples/shuffle.png)

## pixelate

Pixelates an image. Use `--size HxW` to set pixel size used.

``` bash
$ ./pixelate --size 10x50 < input.png > output.png
```

![Pixelate](http://github.com/hawx/img/raw/master/examples/pixelate.png)

## pxl

Implementation of the triangle filter from [pxl app][pxlapp], using the
algorithm loosely described by [revdancatt][rev].

``` bash
$ ./pxl --size 30x30 < input.png > output.png
```

![pxl](http://github.com/hawx/img/raw/master/examples/pxl.png)

## hxl

An (almost; that is I'm not sure this is exactly the same) implementation of the
equilateral triangle filter from [pxl app][pxlapp].

``` bash
$ ./hxl --width 50 < input.png > output.png
```

![hxl](http://github.com/hawx/img/raw/master/examples/hxl.png)

## greyscale

Creates a greyscale version of an image.

``` bash
$ ./greyscale --average < input.png > output.png
```

![Greyscale](http://github.com/hawx/img/raw/master/examples/greyscale.png)

## colourpixels

Features a problematic pixelation algorithm where calculations overflow
producing incorrect (but generally pretty) results.

``` bash
$ ./colourpixels --size 20x30 < input.png > output.png
```

![colourpixels](http://github.com/hawx/img/raw/master/examples/colourpixels.png)

## contrast

Adjusts the contrast of the given image.

``` bash
$ ./contrast --by -25 < input.png > output.png
```

![contrast](http://github.com/hawx/img/raw/master/examples/contrast.png)

## brightness

Adjusts the brightness of the given image.

``` bash
$ ./brightness --by -25 < input.png > output.png
```

![brightness](http://github.com/hawx/img/raw/master/examples/brightness.png)

## hue

Adjusts the hue of the given image.

``` bash
$ ./hue --by -1.34 < input.png > output.png
```

![hue](http://github.com/hawx/img/raw/master/examples/hue.png)

## saturation

Adjusts the saturation of the given image.

``` bash
$ ./saturation --by 0.2 < input.png > output.png
```

![saturation](http://github.com/hawx/img/raw/master/examples/saturation.png)

## blend

Allows you to blend two images together using a blend mode. Takes one image from
STDIN (the base image, imagine the bottom layer in photoshop) and one image as
an argument (the blend image, the layer above).

``` bash
$ < input.png ./blend --screen blend.png --opacity 0.3 > output.png
```

![blend](http://github.com/hawx/img/raw/master/examples/blend.png)

# Composition

These tools have been created to do one task each, and to use standard
input/output so that they can be easily composed. For example;

``` bash
$ < input.png ./colourpixels | ./pxl | ./greyscale --minimal > output.png
```

![Composed](http://github.com/hawx/img/raw/master/examples/composed.png)


[pxlapp]: http://kohlberger.net/apps/pxl
[rev]:    http://revdancatt.com/2012/03/31/the-pxl-effect-with-javascript-and-canvas-and-maths/
