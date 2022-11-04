package functions

import "fmt"

func a() {
	slice1 := []int{1, 2, 3}
	slice2 := slice1
	slice2[0] = 4
	fmt.Println(slice1) // [4 2 3]
}
