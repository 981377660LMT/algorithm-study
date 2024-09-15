package main

import (
	"fmt"
)

func main() {
	nexts := []int32{1, 2, 3, 0, 5, 6, 4}
	start := int32(0)
	cycle := collectCycle(nexts, start)
	fmt.Println(cycle)
}

// 置换环找环.nexts数组中元素各不相同.
func collectCycle(nexts []int32, start int32) []int32 {
	cycle := []int32{}
	cur := start
	for {
		cycle = append(cycle, cur)
		cur = nexts[cur]
		if cur == start {
			break
		}
	}

	return cycle
}
