package main

import "os"

func GetIsWatch() bool {
	for _, arg := range os.Args {
		if arg == "-w" || arg == "--watch" {
			return true
		}
	}
	return false
}

func GetArgs() string {
	return os.Args[1]
}
