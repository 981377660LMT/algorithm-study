package functions

import (
	"fmt"

	"github.com/emirpasic/gods/utils"
)

func b() {
	strings := []interface{}{}
	strings = append(strings, "c")
	strings = append(strings, "a")
	strings = append(strings, "d")
	strings = append(strings, "b")
	utils.Sort(strings, utils.StringComparator)
	fmt.Println(strings)
}
