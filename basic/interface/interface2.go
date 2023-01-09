package main

import (
	"fmt"
	"log"
)

func main() {
	p := Person{
		name: "李彬",
		age:  36,
		address: Address{
			province: "山东",
			city:     "青岛",
		},
	}
	// fmt.Println(p)
	// fmt.Println(p.String())
	// printString(p)
	printString(p.address)

	printString(&p)
}

type Person struct {
	name    string
	age     int
	address Address
}

type Address struct {
	province string
	city     string
}

type Stringer interface {
	String() string
}

func (p *Person) String() string {
	return fmt.Sprintf("姓名 %s,年龄 %d", p.name, p.age)
}

func printString(s Stringer) {
	log.Println(s.String())
}

func (addr Address) String() string {
	return fmt.Sprintf("地址:省 %s,市 %s", addr.province, addr.city)
}
