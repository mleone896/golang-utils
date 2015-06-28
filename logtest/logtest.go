package main

import (
	"log"
	"time"
)

func loopAndPrint(count int) {

	for {

		log.Printf("Im just sitting here looping at %d\n", count)
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
