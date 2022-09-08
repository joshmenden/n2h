package log

import "fmt"

var (
	// https://golangbyexample.com/print-output-text-color-console/
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorWhite = "\033[37m"
	colorBlue  = "\033[34m"
	colorsMap  = map[string]string{
		"green": colorGreen,
		"blue":  colorBlue,
	}
)

func Status(message, emoji string) {
	fmt.Printf("%s%s %s%s\n", colorGreen, emoji, message, colorReset)
}

func Error(err error) {
	fmt.Printf("%s🚨 Error: %v%s\n", colorRed, err, colorReset)
	panic(err)
}

func Linkify(text, link string) string {
	return fmt.Sprintf("\033]8;;%s\a%s\033]8;;\a", link, text)
}

func Colorify(message, color string) string {
	return fmt.Sprintf("%s%s%s", colorsMap[color], message, colorReset)
}

func SubStatus(message string) {
	fmt.Printf("\t- %s%s%s\n", colorWhite, message, colorReset)
}
