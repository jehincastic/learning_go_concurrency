package main

import (
	"fmt"
	"sync"
)

func main() {

	var balance int
	var wg sync.WaitGroup
	var mutex sync.Mutex

	deposit := func(amount int) {
		mutex.Lock()
		defer mutex.Unlock()
		balance += amount
	}

	withdrawal := func(amount int) {
		mutex.Lock()
		defer mutex.Unlock()
		balance -= amount
	}

	// make 100 deposits of $1
	// and 100 withdrawal of $1 concurrently.
	// run the program and check result.

	// TODO: fix the issue for consistent output.

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			deposit(1)
		}()
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			withdrawal(1)
		}()
	}

	wg.Wait()
	fmt.Println(balance)
}
