package ui

import "image/color"

type themeMode string

const (
	backgroundMode themeMode = "background"
	foregroundMode themeMode = "theme"
	comboMode      themeMode = "combo"
)

type theme struct {
	name       string
	mode       themeMode
	background color.Color
	foreground color.Color
}

var defaultTheme = theme{
	name:       "Default",
	mode:       comboMode,
	background: hexColor("#142019"),
	foreground: hexColor("#eef5ef"),
}

var themes = []theme{
	{name: "Red Background", mode: backgroundMode, background: hexColor("#3b1115"), foreground: hexColor("#f8e9e9")},
	{name: "Red Theme", mode: foregroundMode, background: hexColor("#101713"), foreground: hexColor("#f06a6a")},
	{name: "Red Combo", mode: comboMode, background: hexColor("#3b1115"), foreground: hexColor("#ffb0ad")},
	{name: "Orange Background", mode: backgroundMode, background: hexColor("#3c2110"), foreground: hexColor("#fff0e4")},
	{name: "Orange Theme", mode: foregroundMode, background: hexColor("#101713"), foreground: hexColor("#f28b50")},
	{name: "Orange Combo", mode: comboMode, background: hexColor("#3c2110"), foreground: hexColor("#ffc29b")},
	{name: "Yellow Background", mode: backgroundMode, background: hexColor("#3a3210"), foreground: hexColor("#fff9de")},
	{name: "Yellow Theme", mode: foregroundMode, background: hexColor("#101713"), foreground: hexColor("#ead04f")},
	{name: "Yellow Combo", mode: comboMode, background: hexColor("#3a3210"), foreground: hexColor("#fff0a1")},
	{name: "Green Background", mode: backgroundMode, background: hexColor("#102d1c"), foreground: hexColor("#e8f7eb")},
	{name: "Green Theme", mode: foregroundMode, background: hexColor("#101713"), foreground: hexColor("#62c985")},
	{name: "Green Combo", mode: comboMode, background: hexColor("#102d1c"), foreground: hexColor("#a8edbb")},
	{name: "Blue Background", mode: backgroundMode, background: hexColor("#10263d"), foreground: hexColor("#e8f2ff")},
	{name: "Blue Theme", mode: foregroundMode, background: hexColor("#101713"), foreground: hexColor("#62a9ef")},
	{name: "Blue Combo", mode: comboMode, background: hexColor("#10263d"), foreground: hexColor("#a9d2ff")},
	{name: "Purple Background", mode: backgroundMode, background: hexColor("#29183e"), foreground: hexColor("#f5eaff")},
	{name: "Purple Theme", mode: foregroundMode, background: hexColor("#101713"), foreground: hexColor("#be83f1")},
	{name: "Purple Combo", mode: comboMode, background: hexColor("#29183e"), foreground: hexColor("#ddb5ff")},
}

func hexColor(value string) color.Color {
	red := parseHex(value[1:3])
	green := parseHex(value[3:5])
	blue := parseHex(value[5:7])
	return color.RGBA{R: red, G: green, B: blue, A: 0xff}
}

func parseHex(value string) uint8 {
	var result uint8
	for _, digit := range value {
		result *= 16
		switch {
		case digit >= '0' && digit <= '9':
			result += uint8(digit - '0')
		case digit >= 'a' && digit <= 'f':
			result += uint8(digit-'a') + 10
		}
	}
	return result
}
