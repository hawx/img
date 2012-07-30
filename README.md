# Glytch

A selection of image manipulation tools.


## shuffle

Randomly shuffles pixels around the image. Use `-v` or `-h` to constrain it to
vertical or horizontal shuffling, respectively.

Requires Go 1.

``` bash
$ go build shuffle.go
$ ./shuffle --help
$ ./shuffle -v < input.png > output.png
```

## pixelate

Pixelates an image. Use `--size HxW` to set pixel size used.

Requires Go 1.

``` bash
$ go build pixelate.go
$ ./pixelate --help
$ ./pixelate --size 10x50 < input.png > output.png
```

## pxl

Implementation of the triangle filter from [pxl app][pxlapp], using the
algorithm loosely described by [revdancatt][rev].

Requires Ruby and the [chunky_png][cpng] gem.

``` bash
$ ruby pxl.rb < input.png > output.png
```

[pxlapp]: http://kohlberger.net/apps/pxl
[rev]:    http://revdancatt.com/2012/03/31/the-pxl-effect-with-javascript-and-canvas-and-maths/
[cpng]:   https://github.com/wvanbergen/chunky_png
