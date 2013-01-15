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
	"os/exec"
	"bytes"
	"errors"
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
	tmpl(os.Stdout, helpTemplate, c)
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise it is a
// documentationpseudo-command such as ....
func (c *Command) Runnable() bool {
	return c.Run != nil
}


// An External is an external command, named "img-something".
type External struct {
	Name      string
	Path      string
	UsageLine string
	Short     string
	Long      string
}

func (e *External) Usage() {
	tmpl(os.Stdout, externalHelpTemplate, e)
	os.Exit(2)
}

func (e External) String() string {
	return "External{" + e.Name + "}"
}



// Commands list the available commands and help topics.  The order here is the
// order in which they are printed by 'img help'.
var commands = []*Command{
	cmdBlend,
	cmdBrightness,
	cmdChannel,
	cmdContrast,
	cmdGamma,
	cmdGreyscale,
	cmdHue,
	cmdHxl,
	cmdLevels,
	cmdLightness,
	cmdPixelate,
	cmdPxl,
	cmdSaturation,
	cmdShuffle,
	cmdTint,
}

// Will load these in main
var externals = []*External{}

var exitStatus = 0
var exitMu sync.Mutex
const PREFIX = "img-"

func setExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

func nameOfExternal(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func findExternalsIn(dir string) ([]string, error) {
	found := []string{}

	cmd := exec.Command("ls", dir)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return found, err
	}

	for _, possible := range strings.Split(out.String(), "\n") {
		if strings.HasPrefix(nameOfExternal(possible), "img-") {
			found = append(found, dir + "/" + possible)
		}
	}

	return found, nil
}

func runExternal(ext string, flag string) string {
	cmd := exec.Command(ext, flag)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// handle error
	}
	return out.String()
}

// Modified from Go's os/exec#LookPath
func lookupExternals() ([]*External, error) {
	found := []*External{}
	pathenv := os.Getenv("PATH")

	for _, dir := range strings.Split(pathenv, ":") {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}

		if exts, err := findExternalsIn(dir); err == nil {
			for _, ext := range exts {
				usage := runExternal(ext, "--usage")
				short := runExternal(ext, "--short")
				long  := runExternal(ext, "--long")
				name  := nameOfExternal(ext)[4:]

				found = append(found, &External{name, ext, usage, short, long})
			}
		}
	}

	return found, errors.New("executable file not found in $PATH")
}

func main() {
	externals, _ = lookupExternals()

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

	for _, ext := range externals {
		if ext.Name == args[0] {
			if args[1] == "-h" || args[1] == "--help" {
				ext.Usage()
			}

			// run the external command
			cmd := exec.Command(ext.Path, args[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				// handle error
			}
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

  Commands: {{range .Commands}}{{if .Runnable}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}{{end}}

  {{if .HasExternals}}External Commands: {{range .Externals}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}{{end}}
Use "img help [command]" for more information about a command.
`

var helpTemplate = `{{if .Runnable}}Usage: img {{.UsageLine}}
{{end}}{{.Long}}
`

var externalHelpTemplate = `Usage: img {{.UsageLine}}
{{.Long}}
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

type CommandsAndExternals struct {
	Commands  []*Command
	Externals []*External
}

func (c CommandsAndExternals) HasExternals() bool {
	return len(c.Externals) > 0
}


func printUsage(w io.Writer) {
	fmt.Println(externals)
	tmpl(w, usageTemplate, CommandsAndExternals{commands, externals})
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
		fmt.Fprintf(os.Stderr, "Usage: img help command\n\nToo many arguments given.\n")
		os.Exit(2)
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}

	for _, ext := range externals {
		if ext.Name == arg {
			tmpl(os.Stdout, externalHelpTemplate, ext)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q. Run 'go help'.\n", arg)
	os.Exit(2)
}


func exit() {
	os.Exit(exitStatus)
}
