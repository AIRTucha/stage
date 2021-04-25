package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type RunningStatus int

const (
	success RunningStatus = 0
	stopped RunningStatus = 1
	failure RunningStatus = 2
)

func printFromReader(
	reader *io.ReadCloser,
	print func(string),
	status *RunningStatus,
) {
	scannerErr := bufio.NewScanner(*reader)
	for scannerErr.Scan() {
		print(scannerErr.Text())
		if status != nil {
			*status = failure
		}
	}
}

func interruptProcess(cmd *exec.Cmd) {
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

func listenForStop(
	cmd *exec.Cmd,
	stopSignal chan bool,
	status *RunningStatus,
) {
	isStop := <-stopSignal
	if isStop {
		*status = stopped
		interruptProcess(cmd)
	}
}

func Run(executor string, stringCmd string, stopSignal chan bool) RunningStatus {
	executionStatus := success

	cmd := exec.Command(executor, "-c", stringCmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Start()

	go listenForStop(cmd, stopSignal, &executionStatus)

	go printFromReader(&stdout, PrintStdIn, nil)
	go printFromReader(&stderr, PrintStrErr, &executionStatus)

	cmd.Wait()

	Notify(stopSignal, false)

	return executionStatus
}

func RunAll(
	arg string,
	isWatch bool,
	config Config,
	actions map[string][]string,
) func(stopSignal chan bool) {
	return func(stopSignal chan bool) {
		for _, cmdStr := range actions[arg] {
			PrintAction(cmdStr)
		}

		cmd := strings.Join(actions[arg], "; ")
		switch executionStatus := Run(config.engine, cmd, stopSignal); executionStatus {
		case failure:
			PrintError(
				"Command '%v' is interrupted due to error",
				arg,
			)
			return
		case stopped:
			return
		case success:
		}
		PrintInfo("Action '%v' is finished.", arg)
		if isWatch {
			PrintInfo("Watching for new changes...")
		}
	}
}
