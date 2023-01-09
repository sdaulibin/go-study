package main

import (
	"fmt"
	"sync"
)

var sum int

func main() {
	run()
}

func run() {
	var wg sync.WaitGroup
	wg.Add(110)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			add(i)
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			//计数器值减1
			defer wg.Done()
			go fmt.Println("和为:", readSum())
		}()
	}
	wg.Wait()
}

func add(i int) {
	sum += i
}

func readSum() int {
	b := sum
	return b
}
