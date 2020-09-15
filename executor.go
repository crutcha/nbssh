package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
)

var terminalWidth int

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
	width, _, err := terminal.GetSize(int(os.Stdout.Fd()))

	// github actions doesnt like us doing stdin/stdout stuff...
	// maybe we should just set a sane default here
	if err != nil {
		fmt.Println("Unable to determine terminal size!")
	}

	terminalWidth = width
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
			cmd := exec.Command("ssh", "-oStrictHostKeyChecking=accept-new", deviceString, *command)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			cmdErr := cmd.Run()

			var displayedOutput bytes.Buffer
			outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
			displayedOutput.WriteString(fmt.Sprintf(BannerColor, strings.Repeat("#", int(terminalWidth))))
			displayedOutput.WriteString(fmt.Sprintf(BannerColor, device))
			displayedOutput.WriteString(fmt.Sprintf(BannerColor, strings.Repeat("#", int(terminalWidth))))
			displayedOutput.WriteString(fmt.Sprintln(outStr))

			if errStr != "" {
				displayedOutput.WriteString(fmt.Sprintf(WarningColor, errStr))
			}
			if cmdErr != nil {
				displayedOutput.WriteString(fmt.Sprintf(WarningColor, cmdErr))
			}

			fmt.Println(displayedOutput.String())
			<-poolSemaphore
		}(device)
	}
	wg.Wait()
}
