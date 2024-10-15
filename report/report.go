package report

import "fmt"

var Reset = "\033[0m"
var Red = "\033[31;1m"
var Green = "\033[32;1m"
var Yellow = "\033[33;1m"
var Blue = "\033[34;1m"
var Magenta = "\033[35;1m"
var Cyan = "\033[36;1m"
var Gray = "\033[37;1m"
var White = "\033[97;1m"

func Errorf(format string, rest ...any) {
	fmt.Printf("%sError:%s ", Red, Reset)

	fmt.Printf(format, rest...)
}

func Warnf(format string, rest ...any) {
	fmt.Printf("%sWarning:%s ", Yellow, Reset)

	fmt.Printf(format, rest...)
}

func Notef(format string, rest ...any) {
	fmt.Printf("%sNote:%s ", Blue, Reset)

	fmt.Printf(format, rest...)
}
