package main

import (
	"os"
	"os/signal"
	"sync"
	"time"
)

func waitInterupt() {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)
	<-interruptSignal
}

func stopWaiting(
	isWatch bool,
	waitGroup *sync.WaitGroup,
) {
	if isWatch {
		time.AfterFunc(time.Second*2, func() {
			waitGroup.Done()
		})
	}
}

func HandleInterupt(
	isWatch bool,
	stopSignal chan bool,
	externalStopSignal chan bool,
	waitGroup *sync.WaitGroup,
) {
	waitInterupt()
	Notify(externalStopSignal, true)
	Notify(stopSignal, true)
	stopWaiting(isWatch, waitGroup)
}
