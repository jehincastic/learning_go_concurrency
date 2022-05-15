package main

import "fmt"

// TODO: Build a Pipeline
// generator() -> square() -> print

// generator - convertes a list of integers to a channel
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range nums {
			out <- v
		}
		close(out)
	}()
	return out
}

// square - receive on inbound channel
// square the number
// output on outbound channel
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

func main() {
	// set up the pipeline

	// for n := range square(square(generator(1, 2, 3, 4, 5))) {
	// 	fmt.Println(n)
	// }

	dataInChannel := generator(1, 2, 3, 4, 5)
	dataOutChannel := square(dataInChannel)
	for data := range dataOutChannel {
		fmt.Println(data)
	}
	// run the last stage of pipeline
	// receive the values from square stage
	// print each one, until channel is closed.

}
