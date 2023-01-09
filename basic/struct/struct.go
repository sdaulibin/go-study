package main

import "fmt"

func main() {
	fmt.Println("hello")
	var s  person
	fmt.Println(s)
	s.name="李彬"
	s.age=35
	s.sex="男"
	s.address="青岛市"
	book := book{
		bookName:  "哈哈哈",
		bookPrice: "34.5",
	}
	s.book = &book
	fmt.Println(*s.book)

	fmt.Println(s)
	var sp *person = &s
	fmt.Println(sp,&sp,*sp,s)
	fmt.Printf("%p\n",&s)
	(*sp).age=34
	fmt.Println(sp,&sp,*sp,s)
	fmt.Printf("%p\n",&s)

	sp2 := new(person)
	sp2.age=44
	fmt.Println(sp2,*sp2,&sp2)

	s.book.bookName="呵呵呵"
	fmt.Println(*sp.book)
	fmt.Println(book)
	fmt.Println(s.book)

	s.work()
}

type person struct {
	name string
	age int
	sex string
	address string
	book *book
}

type book struct {
	bookName string
	bookPrice string
}

func (p person) work()  {
	fmt.Println(p.name,"is working")
}
