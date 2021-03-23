package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func getArgs() (string, []string) {
	return os.Args[1], os.Args[2:]
}

func getPath() string {
	path, err := os.Getwd()
	if err != nil {
		panic("Can not get current path")
	}
	return path
}

type T = map[string]([]string)

func main() {

	data, err := ioutil.ReadFile(
		filepath.Join(getPath(), "stages.yaml"),
	)
	if err != nil {
		panic("stages.yaml file does not exist")
	}

	t := T{}

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		panic(err)
	}

	arg, _ := getArgs()

	for i, cmd := range t[arg] {
		stdout, err := exec.Command("/bin/sh", "-c", cmd).Output()
		if err != nil {
			panic(
				fmt.Sprintf(
					"Error in command '%v' at script number %v:\n%v",
					arg,
					i,
					err.Error(),
				),
			)
		}
		fmt.Printf(
			"%s",
			stdout,
		)

	}

}
