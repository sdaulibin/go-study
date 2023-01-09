package main

import "fmt"

func main() {
	age := Age(25)
	// age.String()
	// (&age).Modify()
	// age.String()

	sm := Age.String
	sm(age)
}

type Age uint

func (age Age) String() {
	fmt.Println("the age is", age)
}

func (age *Age) Modify() {
	*age = Age(30)
}
