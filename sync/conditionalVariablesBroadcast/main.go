package main

import (
	"fmt"
	"sync"
	"time"
)

var sharedRsc = make(map[string]interface{})

func main() {
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	cv := sync.NewCond(&mu)

	wg.Add(1)
	go func() {
		defer wg.Done()

		//TODO: suspend goroutine until sharedRsc is populated.
		cv.L.Lock()
		for len(sharedRsc) < 1 {
			cv.Wait()
		}
		cv.L.Unlock()
		fmt.Println("From 1", sharedRsc["rsc1"])
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		//TODO: suspend goroutine until sharedRsc is populated.
		cv.L.Lock()
		for len(sharedRsc) < 2 {
			cv.Wait()
		}
		cv.L.Unlock()
		fmt.Println("From 2", sharedRsc["rsc2"])
	}()

	// writes changes to sharedRsc
	time.Sleep(time.Second)
	cv.L.Lock()
	sharedRsc["rsc1"] = "foo"
	sharedRsc["rsc2"] = "bar"
	cv.Signal()
	cv.Signal()
	cv.L.Unlock()
	fmt.Println("Unlocked")
	wg.Wait()
}
