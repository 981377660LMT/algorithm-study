package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	Retain(&arr, func(index int) bool { return arr[index]%2 == 0 })
	fmt.Println(arr)
}

func Retain[T any](arr *[]T, f func(index int) bool) {
	ptr := 0
	n := len(*arr)
	for i := 0; i < n; i++ {
		if f(i) {
			(*arr)[ptr] = (*arr)[i]
			ptr++
		}
	}
	*arr = (*arr)[:ptr:ptr]
}
