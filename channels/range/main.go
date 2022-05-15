package main

import "fmt"

func main() {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < 6; i++ {
			// TODO: send iterator over channel
			ch <- i
		}
	}()

	// TODO: range over channel to recv values
	for val := range ch {
		fmt.Println(val)
	}
}
