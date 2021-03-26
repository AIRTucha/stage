package main

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
)

func waitClosing() {
	done := make(chan bool, 1)
	<-done
}

func handleChanges(
	watcher *fsnotify.Watcher,
	initialStopSignal chan bool,
	onChange func(stopSignal chan bool),
) {
	stopSignal := initialStopSignal
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("Detected change")
				select {
				case stopSignal <- true:
				default:
				}

				go func() {
					time.AfterFunc(time.Second/4, func() {

						stopSignal = make(chan bool)
						onChange(stopSignal)
					})

				}()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			} else {
				Crash(err)
			}
		}
	}
}

func createWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Crash("Can not watch for files due to ", err)
	}

	return watcher
}

func addToWatcher(watcher *fsnotify.Watcher, paths []string) {
	for _, folderPath := range paths {
		if err := watcher.Add(folderPath); err != nil {
			Crashf("Can not watch %v due to %v", folderPath, err)
		}
	}
}

func Watch(path string, stopPrev chan bool, onChange func(stopSignal chan bool)) {
	watcher := createWatcher()
	defer watcher.Close()

	go handleChanges(
		watcher,
		stopPrev,
		func(stopSignal chan bool) {
			// func() { stopSignal <- true }()
			onChange(stopSignal)
		},
	)

	addToWatcher(
		watcher,
		append(
			GetAllSubFolders(path),
			path,
		),
	)

	waitClosing()
}
