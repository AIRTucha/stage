package main

import (
	"fmt"
	"path/filepath"
	"regexp"
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

				fmt.Println(event.Name)
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
			fmt.Printf("Can not watch %v due to %v", folderPath, err)
		}
	}
}

func filterSrc(paths []string, root string) []string {
	var files []string
	for _, file := range paths {
		isFine, err := filepath.Match(filepath.Join(root, "src/**/*"), file)
		if err != nil {
			fmt.Println(err)
		}
		if isFine {
			files = append(files, file)
		}
	}
	return files
}

func print(strs []string) {
	for _, str := range strs {
		fmt.Println((str))
	}
}

func globToRegexp(pattern string) string {
	slashExp, _ := regexp.Compile("\\/")
	doudbleStarExp, _ := regexp.Compile("\\*\\*\\/")
	starExp, _ := regexp.Compile("\\*")
	anyCharTemp, _ := regexp.Compile("\\!\\!\\!")

	return anyCharTemp.ReplaceAllLiteralString(
		starExp.ReplaceAllLiteralString(
			slashExp.ReplaceAllLiteralString(
				doudbleStarExp.ReplaceAllLiteralString(pattern, "!!!"),
				"\\/",
			),
			"[a-zA-Z0-9_.-]*",
		),
		".*",
	)
}

func match(pattern string, str string) bool {
	globExp, _ := regexp.Compile(globToRegexp(pattern))
	return globExp.MatchString(str)
}

func Watch(path string, stopPrev chan bool, onChange func(stopSignal chan bool)) {
	watcher := createWatcher()
	defer watcher.Close()

	go handleChanges(
		watcher,
		stopPrev,
		func(stopSignal chan bool) {
			onChange(stopSignal)
		},
	)

	addToWatcher(
		watcher,
		GetFoldersToWatch(
			path,
			func(str string) bool {
				isValid := match(filepath.Join(path, "src/**/*"), str)
				return isValid
			},
		),
	)

	waitClosing()
}
