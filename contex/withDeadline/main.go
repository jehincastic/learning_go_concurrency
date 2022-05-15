package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

type data struct {
	result string
}

func main() {

	// TODO: set deadline for goroutine to return computational result.

	compute := func(ctx context.Context) <-chan data {
		ch := make(chan data)
		go func() {
			defer close(ch)
			deadline, ok := ctx.Deadline()
			if ok {
				if deadline.Sub(time.Now().Add(50*time.Millisecond)) < 0 {
					fmt.Println("time out")
					return
				}
			}
			// Simulate work.
			time.Sleep(50 * time.Millisecond)

			// Report result.
			select {
			case ch <- data{"123"}:
			case <-ctx.Done():
				return
			}
		}()
		return ch
	}

	// Wait for the work to finish. If it takes too long move on.
	deadline := time.Now().Add(1 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	ch := compute(ctx)
	d, ok := <-ch
	if ok {
		fmt.Printf("work complete: %s\n", d)
	} else {
		fmt.Println("work cancelled")
	}
}
