package main

import (
	"fmt"
)

func main() {
	nexts := []int32{1, 2, 3, 0, 5, 6, 4}
	start := int32(0)
	cycle := collectCycle(func(i int32) int32 { return nexts[i] }, start)
	fmt.Println(cycle)
}

// 置换环找环.
// mapFn: 映射函数, 例如 nexts[i] = i+1, 0->1, 1->2, 2->3, 3->0.
// start: 起始点, 例如 0.
func collectCycle(mapFn func(int32) int32, start int32) []int32 {
	cycle := []int32{}
	cur := start
	for {
		cycle = append(cycle, cur)
		cur = mapFn(cur)
		if cur == start {
			break
		}
	}

	return cycle
}
