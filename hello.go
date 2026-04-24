package main

import (
	"fmt"
	"sync"
)

// RunGoroutines starts two goroutines that print messages to the console.
func RunGoroutines() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("Hello from goroutine 1")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Hello from goroutine 2")
	}()

	wg.Wait()
}
