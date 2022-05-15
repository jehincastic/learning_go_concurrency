package main

import (
	"fmt"
	"sync"
)

var sharedRsc = make(map[string]interface{})

func main() {
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	cv := sync.NewCond(mu)

	wg.Add(1)
	go func() {
		defer wg.Done()

		//TODO: suspend goroutine until sharedRsc is populated.
		cv.L.Lock()
		for len(sharedRsc) == 0 {
			cv.Wait()
		}
		cv.L.Unlock()
		fmt.Println(sharedRsc["rsc1"])
	}()

	// writes changes to sharedRsc
	sharedRsc["rsc1"] = "foo"
	cv.L.Lock()
	fmt.Println("Updated Map")
	cv.Signal()
	cv.L.Unlock()
	wg.Wait()
}
