package main

import "fmt"

type ddd struct {
	a int
	v string
}

func main() {
	nums := make([]ddd, 0, 10)
	fmt.Println(nums[0])
}
