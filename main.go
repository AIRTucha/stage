package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

const executor = "/bin/sh"

func getArgs() (string, []string) {
	return os.Args[1], os.Args[2:]
}

func run(stringCmd string) {
	cmd := exec.Command(executor, "-c", stringCmd)

	stdout, _ := cmd.StdoutPipe()

	cmd.Start()

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		fmt.Println(scanner.Err())
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err.Error())
	}

	cmd.Wait()
}

func main() {
	rootPath := GetCurrentPath()
	arg, _ := getArgs()

	runAll := func() {
		t := ReadYaml(
			GetConfigPath(rootPath),
		)

		for _, cmd := range t[arg] {
			run(cmd)
		}
	}

	runAll()

	Watch(rootPath, runAll)
}
