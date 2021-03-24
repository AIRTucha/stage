package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

func waitClosing() {
	done := make(chan bool, 1)

	fmt.Println("Watch for changes...")
	<-done
	fmt.Println("Stop wattching.")
}

func handleChanges(watcher *fsnotify.Watcher, onChange func()) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				onChange()
			}
		case err, ok := <-watcher.Errors:
			if ok {
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

func Watch(path string, onChange func()) {
	watcher := createWatcher()
	defer watcher.Close()

	go handleChanges(watcher, onChange)

	addToWatcher(
		watcher,
		append(
			GetAllSubFolders(path),
			path,
		),
	)

	waitClosing()
}
