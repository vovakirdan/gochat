package main

import "fmt"

type TextColor string

const (
	Reset TextColor = "\033[0m"
	Red TextColor = "\033[31m"
	Green TextColor = "\033[32m"
	Yellow TextColor = "\033[33m"
	Blue TextColor = "\033[34m"
	Magenta TextColor = "\033[35m"
	Cyan TextColor = "\033[36m"
	Gray TextColor = "\033[37m"
	White TextColor = "\033[97m"
)

func ColoredText(text string, color TextColor) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}