# img

A collection of image manipulation tools. Each tool takes an input file from
standard input, this needs to be in PNG, JPEG or GIF format. Every tool outputs
a PNG file, regardless of input type, to standard output.

This was tested and built with the latest version of Go. Install to
`$GOPATH/bin` with,

``` bash
$ go install github.com/hawx/img
```

You can then run `go help` and `go help [command]` for information.

To use as a Go library, run

``` bash
$ go get github.com/hawx/img
```

Then simply import the required libraries. You can read the documentation on
GoDoc: <http://godoc.org/github.com/hawx/img>.


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
  img = tint.Tint(img, color.NRGBA{131, 18, 19})

  output, _ := os.Create(os.Args[2])
  png.Encode(output, img)
}
```

This can then be run like `./example input.png output.png`.


## External Scripts

_img_ can be extended with external scripts. They must be named `img-something`,
and be somewhere on the `$PATH` to be found. They are then expected to respond
to the following flags:

- `--usage`: the usage string, eg. `"example [args]"`, the first word of
  which should be the command's name.

- `--short`: a short one-line description of the tool to show in `img
  help`.

- `--long`: a long, multiline description detailing arguments that the command
  can take also. Should follow format of other tools and have description
  indented two spaces, and arguments indented four spaces with a hash separating
  the description.

To show an example of an external script I have rewritten a lomo effect detailed
in [Lomography, UNIX Style][tao] in [shell][lomosh] and [Go][lomogo]


[lomosh]: https://gist.github.com/hawx/5047389
[lomogo]: https://gist.github.com/hawx/4566266
[tao]:    http://the.taoofmac.com/space/blog/2005/08/23/2359
