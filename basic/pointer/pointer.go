package main

import "fmt"

func main() {
	a := 10
	b := &a
	fmt.Printf("%v %p\n", a, &a)
	fmt.Printf("%v %p\n", b, b)
	c := *b
	fmt.Printf("type of b: %T\n", b)
	fmt.Printf("type of c: %T\n", c)
	fmt.Printf("value of c: %v\n", c)

	x := 10
	modify1(x)
	fmt.Printf("the value of x : %v\n", x)
	modify2(&x)
	fmt.Printf("the value of x : %v\n", x)
}
func modify1(x int) {
	x = 100
}

func modify2(x *int) {
	*x = 100
}
