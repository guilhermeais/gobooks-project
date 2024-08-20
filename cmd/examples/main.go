package main

import (
	"fmt"
	"time"
)

func worker(workerId int, data chan int) {
	for data := range data {
		fmt.Printf("Worker %d got %d\n", workerId, data)
		time.Sleep((time.Second))
	}
}

func main() {
	dataChannel := make(chan int)

	const WORKER_COUNT = 3
	for workerId := 0; workerId < WORKER_COUNT; workerId++ {
		go worker(workerId, dataChannel)
	}

	const DATA_COUNT = 10
	for i := 0; i < DATA_COUNT; i++ {
		dataChannel <- i
	}
}
