package main

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
)

var terminalWidth uint16

const (
	InfoColor    = "\033[1;34m%s\033[0m\n"
	NoticeColor  = "\033[1;36m%s\033[0m\n"
	WarningColor = "\033[1;33m%s\033[0m\n"
	ErrorColor   = "\033[1;31m%s\033[0m\n"
	DebugColor   = "\033[0;36m%s\033[0m\n"
	BannerColor  = "\033[0;32m%s\033[0m\n"
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

			// Gather user data. No need to inquire about SSH keys, if no password is passed in, SSH keys
			// will automatically be checked.
			var passString string
			if *username == "" {
				currentUser, _ := user.Current()
				username = &currentUser.Username
			}
			if *password == "" {
				passString = ""
			} else {
				passString = fmt.Sprintf(":%s", *password)
			}

			// golang stdlib /x/crypto/ssh doesnt currently fully support openssh config file.
			// instead we'll fork these out to host machine ssh
			deviceString := fmt.Sprintf("%s%s@%s", *username, passString, device)
			cmd := exec.Command("ssh", deviceString, *command)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			cmdErr := cmd.Run()

			outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
			fmt.Printf(BannerColor, strings.Repeat("#", int(terminalWidth)))
			fmt.Printf(BannerColor, device)
			fmt.Printf(BannerColor, strings.Repeat("#", int(terminalWidth)))
			fmt.Println(outStr)

			if errStr != "" {
				fmt.Printf(WarningColor, errStr)
			}
			if cmdErr != nil {
				fmt.Printf(WarningColor, cmdErr, "\n")
			}

			<-poolSemaphore
		}(device)
	}
	wg.Wait()
}
