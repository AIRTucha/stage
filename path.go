package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		Crashf("Can not get current path due to %v", err.Error())
	}
	return path
}

func selectDirs(info []fs.FileInfo) []string {
	var files []string
	for _, file := range info {
		if file.IsDir() {
			files = append(files, file.Name())
		}
	}
	return files
}

func getNestedFolders(paths []string) []string {
	var nestedPath []string
	for _, path := range paths {
		nestedPath = append(nestedPath, GetAllSubFolders(path)...)
	}
	return nestedPath
}

func fullfilPaths(rootPath string, content []string) []string {
	var fulFilledContent []string
	for _, contentPath := range content {
		fulFilledContent = append(fulFilledContent, filepath.Join(rootPath, contentPath))
	}
	return fulFilledContent
}

func GetAllSubFolders(path string) []string {
	fileInfo, _ := ioutil.ReadDir(path)

	return append(
		getNestedFolders(
			fullfilPaths(
				path,
				selectDirs(fileInfo),
			),
		),
		path,
	)
}

func containFileToWatch(root string, paths []fs.FileInfo, matchPattern func(string) bool) bool {
	for _, path := range paths {
		if matchPattern(filepath.Join(root, path.Name())) {
			return true
		}
	}
	return false
}

func GetFoldersToWatch(path string, matchPattern func(string) bool) []string {
	allPaths := GetAllSubFolders(path)
	var selectedPaths []string
	for _, folderPath := range allPaths {
		dirContent, _ := ioutil.ReadDir(folderPath)
		if containFileToWatch(folderPath, dirContent, matchPattern) {
			selectedPaths = append(selectedPaths, folderPath)
		}
	}
	return selectedPaths
}
