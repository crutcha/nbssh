package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
)

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
			err := cmd.Run()

			outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
			result := Result{
				CommandErr: fmt.Sprintf("%v", err),
				Device:     device,
				Stdout:     outStr,
				Stderr:     errStr,
			}
			resultOut, _ := json.Marshal(result)
			fmt.Println(string(resultOut))

			<-poolSemaphore
		}(device)
	}
	wg.Wait()
}
