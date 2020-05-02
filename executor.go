package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"
)

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
			fmt.Println("Called ", device)

			// golang stdlib /x/crypto/ssh doesnt currently fully support openssh config file.
			// instead we'll fork these out to host machine ssh
			var cmdOut bytes.Buffer
			cmd := exec.Command("ssh", device, "show version")
			cmd.Stdout = &cmdOut
			err := cmd.Run()

			fmt.Println(err)

			<-poolSemaphore
		}(device)
	}
	wg.Wait()
}
