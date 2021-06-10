package main

import (
	"sync"
)

func main() {
	PrintInfo("Starting version 0.0.9")

	var waitGroup sync.WaitGroup

	rootPath := GetCurrentPath()
	arg := GetArgs()
	isWatch := GetIsWatch()
	stopSignal := make(chan bool)
	config, actions := ReadYaml(
		GetConfigPath(rootPath),
	)
	runAll := RunAll(
		arg,
		isWatch,
		config,
		actions,
	)
	externalStopSignal := make(chan bool)

	go Watch(
		isWatch,
		&waitGroup,
		config,
		rootPath,
		stopSignal,
		externalStopSignal,
		runAll,
	)
	go HandleInterupt(
		isWatch,
		stopSignal,
		externalStopSignal,
		&waitGroup,
	)

	runAll(stopSignal)

	waitGroup.Wait()

	ResetInputStyle()
}
