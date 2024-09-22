package main

import (
	"fmt"
)

// TODO: golang记忆化sninppet
func main() {
	arr := []int{1, 2, 3, 4, 5}

	arr = Concat(arr, []int{6, 7, 8}, []int{9, 10, 11})
	fmt.Println(arr)

	type Duration int64

	var a, b Duration

	Min2(a, b)

}
