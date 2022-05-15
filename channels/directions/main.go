package main

import "fmt"

// TODO: Implement relaying of message with Channel Direction

func genMsg(ch chan<- int) {
	defer close(ch)
	// send message on ch1
	for i := 0; i < 10; i++ {
		ch <- i + 1
	}
}

func relayMsg(ch1 <-chan int, ch2 chan<- int) {
	// recv message on ch1
	defer close(ch2)
	for v := range ch1 {
		ch2 <- v
	}
	// send it on ch2
}

func main() {
	// create ch1 and ch2
	ch1 := make(chan int)
	ch2 := make(chan int)
	// spine goroutine genMsg and relayMsg
	go genMsg(ch1)
	go relayMsg(ch1, ch2)
	// recv message on ch2
	for v := range ch2 {
		fmt.Printf("Received: %v\n", v)
	}
}
