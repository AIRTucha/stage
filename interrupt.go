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

func notify(
	boolChan chan bool,
	boolVal bool,
) {
	select {
	case boolChan <- boolVal:
	default:
	}
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
	notify(externalStopSignal, true)
	notify(stopSignal, true)
	stopWaiting(isWatch, waitGroup)
}
