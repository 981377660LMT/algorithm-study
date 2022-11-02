package main

import "fmt"

func main() {
	slice1 := make([]struct{}, 10)
	fmt.Println(slice1)
	slice1 = append(slice1, struct{}{})
	fmt.Println(slice1)
}
