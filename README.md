# img

A selection of image manipulation tools.

Requires Go 1.

## shuffle

Randomly shuffles pixels around the image. Use `-v` or `-h` to constrain it to
vertical or horizontal shuffling, respectively.

``` bash
$ go build shuffle.go
$ ./shuffle --help
$ ./shuffle -v < input.png > output.png
```

![Shuffle](http://github.com/hawx/img/raw/master/examples/shuffle.png)

## pixelate

Pixelates an image. Use `--size HxW` to set pixel size used.

``` bash
$ go build pixelate.go
$ ./pixelate --help
$ ./pixelate --size 10x50 < input.png > output.png
```

![Pixelate](http://github.com/hawx/img/raw/master/examples/pixelate.png)

## pxl

Implementation of the triangle filter from [pxl app][pxlapp], using the
algorithm loosely described by [revdancatt][rev].

``` bash
$ go build pxl.go
$ ./pxl --help
$ ./pxl 30 < input.png > output.png
```

![pxl](http://github.com/hawx/img/raw/master/examples/pxl.png)

## greyscale

Creates a greyscale version of an image.

``` bash
$ go build greyscale.go
$ ./greyscale --help
$ ./greyscale --average < input.png > output.png
```

![Greyscale](http://github.com/hawx/img/raw/master/examples/greyscale.png)

## colourpixels

Features a problematic pixelation algorithm where calculations overflow
producing incorrect (but generally pretty) results.

``` bash
$ go build colourpixels.go
$ ./colourpixels --help
$ ./colourpixels --size 20x30 < input.png > output.png
```

![colourpixels](http://github.com/hawx/img/raw/master/examples/colourpixels.png)

## contrast

Adjusts the contrast of the given image.

``` bash
$ go build contrast.go
$ ./contrast --help
$ ./contrast -25 < input.png > output.png
```

![contrast](http://github.com/hawx/img/raw/master/examples/contrast.png)


# Composition

These tools have been created to do one task each, and to use standard
input/output so that they can be easily composed. For example;

``` bash
$ < input.png ./colourpixels | ./pxl | ./greyscale --minimal > output.png
```

![Composed](http://github.com/hawx/img/raw/master/examples/composed.png)


[pxlapp]: http://kohlberger.net/apps/pxl
[rev]:    http://revdancatt.com/2012/03/31/the-pxl-effect-with-javascript-and-canvas-and-maths/
[cpng]:   https://github.com/wvanbergen/chunky_png
