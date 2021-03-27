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
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdinDone := make(chan bool, 1)
	stderrDone := make(chan bool, 1)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	fmt.Println("\033[34m", fmt.Sprintf("> %s", stringCmd))
	cmd.Start()

	go func() {
		isStop := <-stopSignal
		if isStop {
			isNotCrashed = false

			err := cmd.Process.Signal(os.Interrupt)
			if err != nil {
				Crash(err)
			}

			err = cmd.Process.Signal(syscall.SIGINT)
			if err != nil {
				Crash(err)
			}

			err = cmd.Process.Kill()
			if err != nil {
				Crash(err)
			}

			if pgid, err := syscall.Getpgid(cmd.Process.Pid); err == nil {
				syscall.Kill(-pgid, 15) // note the minus sign
			} else {
				Crash(err)
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println("\033[32m", scanner.Text())
		}
		stdinDone <- true
	}()

	go func() {
		scannerErr := bufio.NewScanner(stderr)
		for scannerErr.Scan() {
			fmt.Println("\033[31m", scannerErr.Text())
		}
		stderrDone <- true
	}()

	cmd.Wait()

	<-stdinDone
	<-stderrDone

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
