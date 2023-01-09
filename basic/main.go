package main

import (
	"fmt"
	"strconv"
)

func main() {
	/*fmt.Println("hello 世界")
	var a int = 61
	var b int = 77
	fmt.Printf("a:%d,%b\n",a,a)
	fmt.Printf("b:%d,%b\n",b,b)
	fmt.Println("请输入一个字符串:")
	reader := bufio.NewReader(os.Stdin)
	sin,err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读到的字符串：%s\n",sin)
	}*/
	value := 0
	for i:=1;i<=100 ;i++  {
		value += i
	}
	fmt.Println(value)
	for i:=1;i<=10 ;i++  {
		for j:=1;j<=i ;j++  {
			fmt.Printf("%d * %d = %d\t",i,j,i*j)
			if j==i {
				fmt.Println()
			}
		}
	}
	s1 := []int{1,2,3}
	s2 := s1
	s2[1]=4
	fmt.Printf("%T\n",s1)
	fmt.Printf("%d,%d\n",cap(s1),len(s1))
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = append(s2, 5)
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Printf("%p,%p\n",s1,s2)
	fmt.Printf("%d,%d\n",cap(s1),len(s1))
	fmt.Printf("%d,%d\n",cap(s2),len(s2))
	attr := []int{1,2,3,4,5,6}
	sum(attr)
	fmt.Println(sum)
	var a int = 123
	fmt.Println(&a)
	fmt.Println(a)
	var p *int
	p = &a
	fmt.Println(p)
	fmt.Println(&p)
	fmt.Println(*p)

	b := [4]string{"a","b","c","d"}
	var c *[4]string
	c = &b
	fmt.Println(b,c)
	fmt.Println(*c)
	fmt.Println(b,c)
	fmt.Println(&b,&c)
	fmt.Println(*c)
	fmt.Printf("%p\n",c)
	(*c)[2]="asd";
	fmt.Println(b,*c)

	fmt.Println(strconv.Itoa(2)+"3")
}
func sum(attr []int) {
	fmt.Println(attr)
	attr[0]=100
	fmt.Println(attr)
	sum := 1
	for index,_ := range attr{
		sum += attr[index]
	}
	fmt.Printf("sum=%d\n",sum)
}
