package main

import "fmt"

func main() {
	type name struct {
		a int
	}

	res := make([]name, 2)
	res[0].a = 1

	foo := &res[0]
	foo.a = 2
	fmt.Println(foo, res)
}
