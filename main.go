package main

import (
	"github.com/hawx/img/utils"
	"github.com/hawx/hadfield"
	"bytes"
	"path/filepath"
	"os/exec"
	"os"
	"strings"
	"flag"
)

type External struct {
	Path, Usage, Short, Long  string
}

func (e External) String() string {
	return "External{" + e.Path + "}"
}

func (e *External) Name() string {
	return filepath.Base(e.Path)[4:]
}

func (e *External) Data() interface{} {
	return map[string]interface{}{
		"Callable": e.Callable(),
		"Category": e.Category(),
		"Usage":    e.Usage,
		"Short":    e.Short,
		"Long":     e.Long,
		"Name":     e.Name(),
	}
}

func (e *External) Category() string {
	return "External"
}

func (e *External) Callable() bool {
	return true
}

func (e *External) Call(cmd hadfield.Interface, templates hadfield.Templates, args []string) {
	// args[0] is set to the executable's name, so we can safely replace it with
	// the output type. This is always going to be something, so needs to be
	// checked for, removed, and respected!
	args[0] = string(utils.Output)

	ex := exec.Command(e.Path, args...)
	ex.Stdin  = os.Stdin
	ex.Stdout = os.Stdout
	ex.Stderr = os.Stderr
	err := ex.Run()
	if err != nil {
		os.Exit(2)
	}
	return
}

func findExternalsIn(dir string) ([]string, error) {
	found := []string{}

	dirs, _ := filepath.Glob(dir + "/" + "*")
	for _, possible := range dirs {
		if strings.HasPrefix(filepath.Base(possible), "img-") {
			found = append(found, possible)
		}
	}

	return found, nil
}

func runExternal(ext string, flags... string) string {
	cmd := exec.Command(ext, flags...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// handle
	}
	return out.String()
}

func lookupExternals() hadfield.Commands {
	found := hadfield.Commands{}
	pathenv := os.Getenv("PATH")
	output := string(utils.Output)

	for _, dir := range strings.Split(pathenv, ":") {
		if dir == "" {
			dir = "."
		}

		if exts, err := findExternalsIn(dir); err == nil {
			for _, ext := range exts {
				usage := runExternal(ext, output, "--usage")
				short := runExternal(ext, output, "--short")
				long  := runExternal(ext, output, "--long")

				found = append(found, &External{ext, usage, short, long})
			}
		}
	}

	return found
}

// Commands list the available commands and help topics. The order here is the
// order in which they are printed by 'img help'.
var commands = hadfield.Commands{
	cmdBlend,
	cmdBlur,
	cmdChannel,
	cmdContrast,
	cmdCrop,
	cmdGamma,
	cmdGreyscale,
	cmdHxl,
	cmdLevels,
	cmdPixelate,
	cmdPxl,
	cmdSharpen,
	cmdShuffle,
	cmdTint,
	cmdVxl,
}

var templates = hadfield.Templates{
Usage: `Usage: img [command] [arguments]

  Img is a set of image manipulation tools. They each take an image from STDIN
  and print the result to STDOUT (in some cases they may also require a second
  image, consult the help for the particular command).

  An example usage,

    $ img greyscale < input.png > output.png

  As standard input and output are used throughout, commands can be easily
  chained together using pipes (and parentheses for clarity),

    $ (img greyscale | img pxl | img contrast --by 0.05) < input.png > output.png

  Commands: {{range .}}{{if .Callable}}{{if category . "Command"}}
    {{.Name | printf "%-15s"}} # {{.Short | trim}}{{end}}{{end}}{{end}}

  External Commands: {{range .}}{{if .Callable}}{{if category . "External"}}
    {{.Name | printf "%-15s"}} # {{.Short | trim}}{{end}}{{end}}{{end}}

Use "img help [command]" for more information about a command.
`,
Help: `{{if .Callable}}Usage: img {{.Usage}}
{{end}}{{.Long}}
`,
}

var builtIn = []string{
	"blend", "blur", "channel", "contrast", "crop", "gamma", "greyscale", "hxl",
	"levels", "pixelate", "pxl", "sharpen", "shuffle", "tint", "vxl",
}

func isRunningBuiltin(args []string) bool {
	for _, v := range builtIn {
		if args[0] == v { return true }
	}

	return false
}

func main() {
	var jpeg, png, tiff bool
	flag.BoolVar(&jpeg, "jpg",  false, "")
	flag.BoolVar(&jpeg, "jpeg", false, "")
	flag.BoolVar(&png,  "png",  false, "")
	flag.BoolVar(&tiff, "tiff", false, "")
	flag.BoolVar(&tiff, "tif",  false, "")

	flag.Parse()
	if jpeg { utils.Output = utils.JPEG }
	if png  { utils.Output = utils.PNG }
	if tiff { utils.Output = utils.TIFF }

	if !isRunningBuiltin(flag.Args()) {
		externals := lookupExternals()
		commands = append(commands, externals...)
	}

	hadfield.Run(commands, templates)
}
