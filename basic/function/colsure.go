package main

import "fmt"

func main() {
	fmt.Println("hello")
	c := colsure()
	fmt.Println(c())
	fmt.Println(c())
	fmt.Println(c())
}

func colsure() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}
