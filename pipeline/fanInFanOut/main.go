// generator() -> square() -> print

package main

import (
	"fmt"
	"sync"
)

func generator(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	// Implement fan-in
	// merge a list of channels to a single channel
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	for _, ch := range cs {
		go output(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	in := generator(2, 3)

	// TODO: fan out square stage to run two instances.
	sq1 := square(in)
	sq2 := square(in)

	// TODO: fan in the results of square stages.
	outCh := merge(sq1, sq2)
	for val := range outCh {
		fmt.Println(val)
	}
}
