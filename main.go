package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// colorReset := "\033[0m"

// colorRed := "\033[31m"
// colorGreen := "\033[32m"
// colorYellow := "\033[33m"
// colorBlue := "\033[34m"
// colorPurple := "\033[35m"
// colorCyan := "\033[36m"
// colorWhite := "\033[37m"

const executor = "/bin/sh"

func getIsWatch() bool {
	for _, arg := range os.Args {
		if arg == "-w" || arg == "--watch" {
			return true
		}
	}
	return false
}

func getArgs() string {
	return os.Args[1]
}

func run(stringCmd string, stopSignal chan bool) bool {
	isNotCrashed := true

	cmd := exec.Command(executor, "-c", stringCmd)

	stdinDone := make(chan bool, 1)
	stderrDone := make(chan bool, 1)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	fmt.Println("\033[34m", fmt.Sprintf("> %s", stringCmd))
	cmd.Start()
	// "\033[35m", "Start", cmd.Process.Pid)

	go func() {
		// "\033[35m", "Wait for kill", cmd.Process.Pid)
		isStop := <-stopSignal
		if isStop {
			// exec.CommandContext(
			isNotCrashed = false

			err := cmd.Process.Signal(syscall.SIGINT)
			fmt.Println("\033[35m", "Int error: ", err)
			err = cmd.Process.Signal(os.Interrupt)
			fmt.Println("\033[35m", "OS Int error: ", err)
			err = cmd.Process.Signal(syscall.SIGTERM)
			fmt.Println("\033[35m", "Term error: ", err)
			cmd.Process.Signal(os.Kill)
			fmt.Println("Exit status", cmd.ProcessState.Exited())
			fmt.Println("\033[35m", "Kill error: ", err)
			err = cmd.Process.Kill()

		}
		// "\033[35m", "Stop wait for kill", cmd.Process.Pid)
	}()

	go func() {
		// "\033[35m", "Start scan", cmd.Process.Pid)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println("\033[32m", scanner.Text())
		}
		stdinDone <- true
		// "\033[35m", "Scan finished", cmd.Process.Pid)
	}()

	go func() {
		// "\033[35m", "Start scan err", cmd.Process.Pid)
		scannerErr := bufio.NewScanner(stderr)
		for scannerErr.Scan() {
			fmt.Println("\033[31m", scannerErr.Text())
		}
		stderrDone <- true
		// "\033[35m", "Scan finished err", cmd.Process.Pid)
	}()

	// "\033[35m", "Before wait", cmd.Process.Pid)
	cmd.Wait()

	// "\033[35m", "After wait", cmd.Process.Pid)

	<-stdinDone
	<-stderrDone

	// "\033[35m", "After std chs", cmd.Process.Pid)
	select {
	case stopSignal <- false:
	default:
	}
	return isNotCrashed
}

func main() {
	fmt.Println("\033[36m", "[stage] Starting version 0.0.2")

	rootPath := GetCurrentPath()
	arg := getArgs()
	isWatch := getIsWatch()
	stopSignal := make(chan bool)

	runAll := func(stopSignal chan bool) {
		t := ReadYaml(
			GetConfigPath(rootPath),
		)
		for _, cmd := range t[arg] {
			if !run(cmd, stopSignal) {
				// fmt.Println(
				// 	"\033[33m",
				// 	fmt.Sprintf("Command '%v' is interrupted due to error at step %v", arg, i+1),
				// )
				return
			}
		}
		if isWatch {
			fmt.Println("\033[36m", "[stage] Watching for new changes...")
		}
	}

	if isWatch {
		fmt.Println("\033[36m", "[stage] Start watching for changes...")
		go Watch(rootPath, stopSignal, runAll)
	}

	runAll(stopSignal)
	done := make(chan bool, 1)
	<-done
}
