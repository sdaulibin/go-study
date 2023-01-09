package main

import (
	"fmt"
	"sync"
	"time"
)

var m sync.Mutex

var set = make(map[int]bool,0)

func print(num int)  {
	m.Lock()
	if _,ok := set[num];ok {
		fmt.Println(num)
	}
	set[num] = true
	m.Unlock()
}

func main()  {
	for i := 0; i<10; i++  {
		go print(100)
	}
	time.Sleep(time.Second)
}