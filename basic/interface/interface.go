package main

import "fmt"

func main() {
	fmt.Println("hello")
	m := Mouse{
		name:  "哈哈哈",
		color: "红色",
	}
	d := Disk{
		name:  "呵呵呵",
		price: "23.56",
	}
	testInterface(&m)
	testInterface(&d)

	m.start()

	var u Usb
	u = &m
	u.stop()

	var u2 Usb
	u2 = &d
	u2.start()
}

type Usb interface {
	start()
	stop()
}
type Mouse struct {
	name  string
	color string
}
type Disk struct {
	name  string
	price string
}

func (m *Mouse) start() {
	fmt.Println(m.name, "start to work .....")
}

func (m Mouse) stop() {
	fmt.Println(m.name, "stop to work .....")
}

func (d *Disk) start() {
	fmt.Println(d.name, "start to work .....")
}

func (d *Disk) stop() {
	fmt.Println(d.name, "stop to work .....")
}
func testInterface(u Usb) {
	u.start()
	u.stop()
}
