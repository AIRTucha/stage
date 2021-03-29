package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
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

type RunningStatus int

const (
	success RunningStatus = 0
	stopped RunningStatus = 1
	failure RunningStatus = 2
)

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

func run(executor string, stringCmd string, stopSignal chan bool) RunningStatus {
	executionStatus := success

	cmd := exec.Command(executor, "-c", stringCmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	fmt.Println("\033[34m", fmt.Sprintf("> %s", stringCmd))
	cmd.Start()

	go func() {
		isStop := <-stopSignal
		if isStop {
			executionStatus = stopped

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
				syscall.Kill(-pgid, 15)
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
	}()

	go func() {
		scannerErr := bufio.NewScanner(stderr)
		for scannerErr.Scan() {
			fmt.Println("\033[31m", scannerErr.Text())
			executionStatus = failure
		}
	}()

	cmd.Wait()

	select {
	case stopSignal <- false:
	default:
	}
	return executionStatus
}

func main() {
	fmt.Println("\033[36m", "[stage] Starting version 0.0.4")

	rootPath := GetCurrentPath()
	arg := getArgs()
	isWatch := getIsWatch()
	stopSignal := make(chan bool)
	config, actions := ReadYaml(
		GetConfigPath(rootPath),
	)
	fmt.Println(config)
	runAll := func(stopSignal chan bool) {
		for i, cmd := range actions[arg] {
			switch executionStatus := run(config.engine, cmd, stopSignal); executionStatus {
			case failure:
				fmt.Println(
					"\033[33m",
					fmt.Sprintf("Command '%v' is interrupted due to error at step %v", arg, i+1),
				)
				return
			case stopped:
				return
			case success:
			}
		}
		if isWatch {
			fmt.Println("\033[36m", "[stage] Watching for new changes...")
		}
	}
	var waitGroup sync.WaitGroup

	externalStopSignal := make(chan bool)
	if isWatch {
		waitGroup.Add(1)
		fmt.Println("\033[36m", "[stage] Start watching for changes...")
		go Watch(rootPath, config.watch, stopSignal, externalStopSignal, runAll)
	}

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)
	go func() {
		<-interruptSignal
		select {
		case externalStopSignal <- true:
		default:
		}
		select {
		case stopSignal <- true:
		default:
		}
		if isWatch {
			waitGroup.Done()
		}
	}()

	runAll(stopSignal)

	waitGroup.Wait()
}
