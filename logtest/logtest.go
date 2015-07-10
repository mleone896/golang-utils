package main

import (
	"log"
	"os"
	"time"
)

var progName = os.Args[0]

func loopAndPrint(count int) {

	for {

		log.Printf("%s Im just sitting here looping at %d just to confirm fluentd picks up any container\n", progName, count)
		count++
		time.Sleep(2 * time.Second)
		loopAndPrint(count)
	}
}

func run() {
	loopAndPrint(0)

}

func main() {

	run()

}
