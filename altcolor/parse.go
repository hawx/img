package altcolor

import (
	"image/color"
	"regexp"
	"strconv"
)

var (
	SHORT_HEX = regexp.MustCompile("#([A-Fa-f0-9])([A-Fa-f0-9])([A-Fa-f0-9])([A-Fa-f0-9])?")
	LONG_HEX  = regexp.MustCompile("#([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})?")

	RGB  = regexp.MustCompile("rgb\\(\\s*(\\d+)\\s*,\\s*(\\d+)\\s*,\\s*(\\d+)\\s*\\)")
	RGBA = regexp.MustCompile("rgba\\(\\s*(\\d+)\\s*,\\s*(\\d+)\\s*,\\s*(\\d+)\\s*,\\s*(\\d+)\\s*\\)")
)

func intToInt(s string) uint8 {
	r, _ := strconv.Atoi(s)
	return uint8(r)
}

func hexToInt(s string) uint8 {
	r, _ := strconv.ParseInt(s, 16, 16)
	return uint8(r)
}

func Parse(str string) color.Color {
	_parts := SHORT_HEX.FindAllStringSubmatch(str, 5)
	if _parts != nil {
		parts := _parts[0]
		if parts[4] == "" {
			parts[4] = "F"
		}
		return color.NRGBA{
			hexToInt(parts[1] + parts[1]),
			hexToInt(parts[2] + parts[2]),
			hexToInt(parts[3] + parts[3]),
			hexToInt(parts[4] + parts[4]),
		}
	}

	_parts = LONG_HEX.FindAllStringSubmatch(str, 5)
	if _parts != nil {
		parts := _parts[0]
		if parts[4] == "" {
			parts[4] = "FF"
		}
		return color.NRGBA{
			hexToInt(parts[1]),
			hexToInt(parts[2]),
			hexToInt(parts[3]),
			hexToInt(parts[4]),
		}
	}

	_parts = RGB.FindAllStringSubmatch(str, 4)
	if _parts != nil {
		parts := _parts[0]
		return color.NRGBA{
			intToInt(parts[1]),
			intToInt(parts[2]),
			intToInt(parts[3]),
			255,
		}
	}

	_parts = RGBA.FindAllStringSubmatch(str, 5)
	if _parts != nil {
		parts := _parts[0]
		return color.NRGBA{
			intToInt(parts[1]),
			intToInt(parts[2]),
			intToInt(parts[3]),
			intToInt(parts[4]),
		}
	}

	return nil
}
