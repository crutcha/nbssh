package main

import (
	"bytes"
	//"encoding/json"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var terminalWidth uint16

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

// TODO: what if its windows?
func init() {
	winSize, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	terminalWidth = winSize.Col

	if err != nil {
		panic(err)
	}

	terminalWidth = winSize.Col
}

type Executor struct {
	command     string
	devices     []string
	concurrency int
}

type Result struct {
	CommandErr string `json:"commanderr"`
	Device     string `json:"device"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
}

func newExecutor(devices []string) *Executor {
	return &Executor{
		command:     *command,
		devices:     devices,
		concurrency: *concurrency,
	}
}

func (e *Executor) execute() {
	var wg sync.WaitGroup
	poolSemaphore := make(chan int, e.concurrency)
	defer close(poolSemaphore)

	for _, device := range e.devices {
		poolSemaphore <- 1

		wg.Add(1)
		go func(device string) {
			defer wg.Done()
			//fmt.Println("Called ", device)

			// golang stdlib /x/crypto/ssh doesnt currently fully support openssh config file.
			// instead we'll fork these out to host machine ssh
			ourCommand := "\"show version\""
			cmd := exec.Command("ssh", "192.168.10.1", ourCommand)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			cmdErr := cmd.Run()

			outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
			/*
				result := Result{
					CommandErr: fmt.Sprintf("%v", err),
					Device:     device,
					Stdout:     outStr,
					Stderr:     errStr,
				}
				resultOut, _ := json.Marshal(result)
				fmt.Println(string(resultOut))
			*/
			fmt.Println(strings.Repeat("#", int(terminalWidth)))
			fmt.Println(device)
			fmt.Println(strings.Repeat("#", int(terminalWidth)))
			fmt.Println(outStr)
			fmt.Printf(InfoColor, errStr)
			fmt.Println(cmdErr, "\n")

			<-poolSemaphore
		}(device)
	}
	wg.Wait()
}
