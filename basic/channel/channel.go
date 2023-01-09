package main

import "fmt"

func main() {
	ch := make(chan string)
	go func() {
		fmt.Println("我是李彬")
		ch <- "哈哈"
	}()

	fmt.Println("我是main")
	v := <-ch
	fmt.Println("接收到的chan中的值为：", v)
}
