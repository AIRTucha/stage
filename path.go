package main

import (
	"io/fs"
	"io/ioutil"
	"os"
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

func GetAllSubFolders(path string) []string {
	fileInfo, err := ioutil.ReadDir(path)

	if err != nil {
		Crash("Can not get all subforlder due to ", err)
	}

	return selectDirs(fileInfo)
}
