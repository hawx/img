package main

import (
	"strings"
	"fmt"
	"os"
	"sync"
	"log"
	"text/template"
	"unicode/utf8"
	"unicode"
	"io"
	"flag"
)

// NOTE: this is mainly ripped from the go command see
// https://code.google.com/p/go/source/browse/src/cmd/go/main.go

// A Command is a single img command, like img hue or img pixelate.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'img help' output.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command wil do its own flag parsing.
	CustomFlags bool
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise it is a
// documentationpseudo-command such as ....
func (c *Command) Runnable() bool {
	return c.Run != nil
}

// Commands list the available commands and help topics.  The order here is the
// order in which they are printed by 'img help'.
var commands = []*Command{
	cmdBlend,
	cmdBrightness,
	cmdContrast,
	cmdGamma,
	cmdGreyscale,
	cmdHue,
	cmdHxl,
	cmdLightness,
	cmdPixelate,
	cmdPxl,
	cmdSaturation,
	cmdShuffle,
	cmdTint,
}

var exitStatus = 0
var exitMu sync.Mutex

func setExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			cmd.Run(cmd, args)
			exit()
			return
		}
	}

	fmt.Fprintf(os.Stderr, "img: unknown subcommand %q\nRun 'img help' for usage.\n", args[0])
	setExitStatus(2)
	exit()
}

var usageTemplate = `Usage: img [command] [arguments]

  Img is a set of image manipulation tools. They each take an image from STDIN
  and print the result to STDOUT (in some cases they may also require a second
  image, consult the help for the particular command).

  An example usage,

    $ img greyscale < input.png > output.png

  As standard input and output are used throughout, commands can be easily
  chained together using pipes (and parentheses for clarity),

    $ (img greyscale | img pxl | img contrast --by 0.05) < input.png > output.png

  Commands: {{range .}}{{if .Runnable}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}{{end}}

Use "img help [command]" for more information about a command.
`

var helpTemplate = `{{if .Runnable}}Usage: img {{.UsageLine}}
{{end}}{{.Long}}
`

	func tmpl(w io.Writer, text string, data interface{}) {
		t := template.New("top")
		t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
		template.Must(t.Parse(text))
		if err := t.Execute(w, data); err != nil {
			panic(err)
		}
	}

	func capitalize(s string) string {
		if s == "" {
			return s
		}
		r, n := utf8.DecodeRuneInString(s)
		return string(unicode.ToTitle(r)) + s[n:]
	}

	func printUsage(w io.Writer) {
		tmpl(w, usageTemplate, commands)
	}

	func usage() {
		printUsage(os.Stderr)
		os.Exit(2)
	}

	// help implements the 'help' command
	func help(args []string) {
		if len(args) == 0 {
			printUsage(os.Stdout)
			return
		}
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "usage: img help command\n\nToo many arguments given.\n")
			os.Exit(2)
		}

		arg := args[0]

		for _, cmd := range commands {
			if cmd.Name() == arg {
				tmpl(os.Stdout, helpTemplate, cmd)
				return
			}
		}

		fmt.Fprintf(os.Stderr, "Unknown help topic %#q. Run 'go help'.\n", arg)
		os.Exit(2)
	}


func exit() {
	os.Exit(exitStatus)
}
