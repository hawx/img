# img [![docs](http://godoc.org/github.com/hawx/img?status.png)](http://godoc.org/github.com/hawx/img)

A collection of image manipulation tools. Each tool takes an input file from
standard input, this needs to be in PNG, JPEG or GIF format. They output the
resulting image (by default in PNG format) to standard output.

To install run,

``` bash
$ go install github.com/hawx/img
```

You can then run `go help` and `go help [command]` for information.

> The aim of _img_ is not to be fast, if you want fast use
> [GraphicsMagick](http://www.graphicsmagick.org/), the aim is to be
> readable. Diving through endless files of C and C++ to find out how a certain
> effect is implemented is no fun, reading a single Go file is hopefully better.


## Example (Command Line)

Here is an example: First we convert the image to greyscale using the values
from the red colour channel, then boost the contrast slightly using a linear
function, and finally tint the image with a dark red.

``` bash
(img greyscale --red | \
  img contrast --linear --ratio 1.5 | \
  img tint --with '#83121344') < input.png > output.png
```

You can see here how easy it is to chain different tools together using pipes.


## Example (Go)

You can also use the _img_ libraries in Go code. We could rewrite the previous
example as,

``` go
// example.go
package main

import (
  "github.com/hawx/img/contrast"
  "github.com/hawx/img/greyscale"
  "github.com/hawx/img/tint"

  "image/png"
  "os"
)

func main() {
  input, _ := os.Open(os.Args[1])
  img,   _ := png.Decode(input)

  img = greyscale.Red(img)
  img = contrast.Linear(img, 1.5)
  img = tint.Tint(img, color.NRGBA{131, 18, 19, 255})

  output, _ := os.Create(os.Args[2])
  png.Encode(output, img)
}
```

This can then be compiled and run like `./example input.png output.png`.
