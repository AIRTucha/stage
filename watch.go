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
	debounce int,
	externalStopSignal chan bool,
	validatePath func(string) bool,
	onChange func(stopSignal chan bool),
) {
	var timer *time.Timer
	stopSignal := initialStopSignal
	for {
		select {
		case stop := <-externalStopSignal:
			if stop {
				fmt.Println("\033[36m", "\n [stage] Stop process...")
				stopSignal <- true
			}
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			fileName := event.Name
			if event.Op&fsnotify.Write == fsnotify.Write && validatePath(fileName) {
				select {
				case stopSignal <- true:
				default:
				}
				fmt.Println("\033[36m", fmt.Sprintf("File is %v canged.", GetRelativeToRoot(fileName)))
				go func() {
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(time.Millisecond*time.Duration(debounce), func() {
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
func removeWatcher(watcher *fsnotify.Watcher, paths []string) {
	for _, folderPath := range paths {
		if err := watcher.Remove(folderPath); err != nil {
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

func findMatchAny(patterns []string, rootPath string, path string) bool {
	for _, pattern := range patterns {
		if match(filepath.Join(rootPath, pattern), path) {
			return true
		}
	}
	return false
}

func Watch(path string, patterns []string, debounce int, stopPrev chan bool, externalStopSignal chan bool, onChange func(chan bool)) {
	watcher := createWatcher()
	defer watcher.Close()

	matchPatterns := func(str string) bool {
		return findMatchAny(patterns, path, str)
	}

	foldersToWatch := GetFoldersToWatch(
		path,
		func(str string) bool {
			isValid := matchPatterns(str)
			return isValid
		},
	)
	go handleChanges(
		watcher,
		stopPrev,
		debounce,
		externalStopSignal,
		matchPatterns,
		func(stopSignal chan bool) {
			removeWatcher(watcher, foldersToWatch)
			foldersToWatch = GetFoldersToWatch(
				path,
				func(str string) bool {
					isValid := matchPatterns(str)
					return isValid
				},
			)
			addToWatcher(watcher, foldersToWatch)
			onChange(stopSignal)
		},
	)

	addToWatcher(
		watcher,
		foldersToWatch,
	)

	waitClosing()
}
