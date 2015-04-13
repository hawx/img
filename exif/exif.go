// Package exif provides a go interface for exiftool.
package exif

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const EXIFTOOL = "exiftool"

// Set of all writable tags
var writable = map[string]bool{}

// Set of tags exiftool reports are writable, but aren't/shouldn't be!
var stupidUnwritableWritableTags = map[string]bool{
	"PreviewImageLength": true,
	"PreviewImageStart":  true,
	"ThumbnailLength":    true,
	"ThumbnailOffset":    true,
	"Directory":          true,
	"FileName":           true,
}

// This needs working out! Though this seems to be working for now...
func shellEscape(s string) string {
	// return strings.Replace(s, " ", "\\ ", -1)
	return s
}

type Exif struct {
	path    string
	data    map[string]string
	changed map[string]bool // a poor man's set
}

// Exists returns true if "exiftool" is installed anywhere on the user's $PATH.
func Exists() bool {
	_, err := exec.LookPath(EXIFTOOL)
	return err == nil
}

// New creates an empty Exif object. To initialize with the exif data of an
// existing image use Load.
func New() *Exif {
	return &Exif{
		path:    "",
		data:    map[string]string{},
		changed: map[string]bool{},
	}
}

func getWritableTags() {
	if len(writable) == 0 {
		// Get the writable tags then
		out, _ := exec.Command("exiftool", "-listw").Output()
		for _, tag := range strings.Split(string(out), " ") {
			key := strings.TrimSpace(tag)
			if _, ok := stupidUnwritableWritableTags[key]; !ok {
				writable[key] = true
			}
		}
	}
}

// Load creates a new Exif object, populated with the exif data of the file at
// the path given. It will silently fail if any errors are encountered.
func Load(path string) *Exif {
	if !Exists() {
		return New()
	}

	out, err := exec.Command("exiftool", "-s", path).Output()
	if err != nil {
		return New()
	}

	getWritableTags()

	lines := strings.Split(string(out), "\n")
	exif := New()
	exif.path = path

	for _, line := range lines[:len(lines)-1] {
		parts := strings.SplitN(line, ":", 2)
		key := strings.Replace(parts[0], " ", "", -1)
		if _, ok := writable[key]; ok {
			val := strings.TrimSpace(parts[1])
			if !strings.Contains(val, "use -b option to extract") {
				exif.data[key] = val
			}
		}
	}

	return exif
}

// Decode can be used for unnamed files, for example STDIN. It loads the exif
// data from the file given.
func Decode(r io.Reader) *Exif {
	if !Exists() {
		return New()
	}

	// Create a temporary file for exiftool to read
	tmp, _ := ioutil.TempFile("", "img-exif-")
	path := tmp.Name()
	// Make sure the file is deleted after
	defer func() { os.Remove(path) }()
	// Copy the file to it
	io.Copy(tmp, r)
	// Load normally
	return Load(path)
}

// Get returns the value for the key given. The key should be in CamelCase form.
func (e Exif) Get(key string) string {
	val, ok := e.data[key]
	if !ok {
		return ""
	}
	return val
}

// Set changes the value of the key. The key should be in CamelCase form.
func (e Exif) Set(key, val string) {
	if _, ok := writable[key]; ok {
		e.changed[key] = true
		e.data[key] = val
	}
}

// Keys returns a list of all exif data keys present.
func (e Exif) Keys() []string {
	keys := make([]string, len(e.data))
	i := 0
	for k, _ := range e.data {
		keys[i] = k
		i++
	}
	return keys
}

func (e Exif) String() string {
	if len(e.data) == 0 {
		return "(Exif:empty)"
	}
	s := ""
	for k, v := range e.data {
		s += k + " = " + v + "\n"
	}
	return s[:len(s)-1]
}

// Save will write the changed (and only the changed) exif data to the path that
// was initially given.
func (e Exif) Save() error {
	if e.path == "" {
		return errors.New("exif.Exif does not have path to save to")
	}
	args := make([]string, len(e.changed)+2)
	args[0] = "-overwrite_original"

	i := 1
	for k, _ := range e.changed {
		args[i] = "-" + k + "=" + shellEscape(e.Get(k))
		i++
	}

	args[i] = e.path
	return exec.Command("exiftool", args...).Run()
}

// Write will write to the path given all of the exif data (changed or unchanged).
func (e Exif) Write(path string) error {
	if !Exists() {
		return errors.New("exif package requires 'exiftool' to load/write exif data")
	}

	args := make([]string, len(e.data)+3)
	args[0] = "-overwrite_original"

	i := 1
	for k, v := range e.data {
		args[i] = "-" + k + "=" + shellEscape(v)
		i++
	}

	args[i] = path
	return exec.Command("exiftool", args...).Run()
}
