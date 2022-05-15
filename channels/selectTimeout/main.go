package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 1)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "one"
	}()

	// TODO: implement timeout for recv on channel ch
	select {
	case msg1 := <-ch:
		fmt.Println(msg1)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout")
	}
}
