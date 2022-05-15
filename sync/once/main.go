package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	load := func() {
		fmt.Println("Run only once initialization")
	}

	// flag := false
	// var mu sync.Mutex
	var once sync.Once

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(j int) {
			defer wg.Done()

			//TODO: modify so that load function gets called only once.
			// mu.Lock()
			// if !flag {
			//	load(j + 1)
			// 	flag = true
			// }
			// mu.Unlock()
			once.Do(load)
		}(i)
	}
	wg.Wait()
}
