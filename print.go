package main

import (
	"fmt"
)

func PrintAction(stringCmd string) {
	fmt.Println("\033[34m", fmt.Sprintf("> %s", stringCmd))
}

func PrintStdIn(str string) {
	fmt.Println("\033[32m", str)
}

func PrintStrErr(errStr string) {
	fmt.Println("\033[31m", errStr)
}

func PrintInfo(format string, a ...interface{}) {
	fmt.Println(
		"\033[36m",
		fmt.Sprintf(
			"[stage] %v",
			fmt.Sprintf(format, a...),
		),
	)
}

func PrintError(format string, a ...interface{}) {
	fmt.Println(
		"\033[33m",
		fmt.Sprintf(format, a...),
	)
}

func ResetInputStyle() {
	fmt.Println("\033[0m")
}
