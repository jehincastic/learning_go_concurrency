package main

import (
	"fmt"
	"time"
)

func fun(s string) {
	for i := 0; i < 3; i++ {
		fmt.Println(s)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Direct call
	fun("direct call")

	// TODO: write goroutine with different variants for function call.

	// goroutine function call
	go fun("goroutine call")

	// goroutine with anonymous function
	go func(s string) {
		fun(s)
	}("goroutine anonymous call")

	// goroutine with function value call
	fv := fun
	go fv("goroutine function value call")

	// wait for goroutines to end
	// go routines does not run after main function ends
	time.Sleep(3 * time.Second)

	fmt.Println("done..")
}
