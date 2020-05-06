package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
)

func main() {
	fmt.Println(kingpin.Parse())
	matchingDevices := queryDevices()

	fmt.Println("Executing against: ", matchingDevices)

	//executor := newExecutor(matchingDevices)
	//executor.execute()
}
