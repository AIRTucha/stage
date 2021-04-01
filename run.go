package main

import (
	"bufio"
	"os"
	"os/exec"
	"syscall"
)

type RunningStatus int

const (
	success RunningStatus = 0
	stopped RunningStatus = 1
	failure RunningStatus = 2
)

func Run(executor string, stringCmd string, stopSignal chan bool) RunningStatus {
	executionStatus := success

	cmd := exec.Command(executor, "-c", stringCmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	PrintAction(stringCmd)

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
			PrintStdIn(scanner.Text())
		}
	}()

	go func() {
		scannerErr := bufio.NewScanner(stderr)
		for scannerErr.Scan() {
			PrintStrErr(scannerErr.Text())
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

func RunAll(
	arg string,
	isWatch bool,
	config Config,
	actions map[string][]string,
) func(stopSignal chan bool) {
	return func(stopSignal chan bool) {
		for i, cmd := range actions[arg] {
			switch executionStatus := Run(config.engine, cmd, stopSignal); executionStatus {
			case failure:
				PrintError(
					"Command '%v' is interrupted due to error at step %v",
					arg,
					i+1,
				)
				return
			case stopped:
				return
			case success:
			}
		}
		PrintInfo("Action '%v' is finished.", arg)
		if isWatch {
			PrintInfo("Watching for new changes...")
		}
	}
}
