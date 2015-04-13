package exif

import (
	"fmt"
	"github.com/hawx/img/exif"
)

func ExampleReading() {
	data := exif.Load("test.jpg")
	fmt.Println(exif.Get("UserComment"))
}

func ExampleModification() {
	data := exif.Load("test.jpg")
	exif.Set("UserComment", "Nice test photo")
	exif.Save()
}

func ExampleCopying() {
	exif.Load("test.jpg").Write("other.jpg")
}
