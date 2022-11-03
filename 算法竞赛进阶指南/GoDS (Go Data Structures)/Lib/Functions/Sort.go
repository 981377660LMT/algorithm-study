package main

import (
	"fmt"

	"github.com/emirpasic/gods/utils"
)

func main() {
	strings := []interface{}{}
	strings = append(strings, "c")
	strings = append(strings, "a")
	strings = append(strings, "d")
	strings = append(strings, "b")
	utils.Sort(strings, utils.StringComparator)
	fmt.Println(strings)
}
