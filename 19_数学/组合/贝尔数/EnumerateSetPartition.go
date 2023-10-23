package main

import "fmt"

func main() {
	setPartition(3)
}

func setPartition(n int) {
	groups := [][]int{} // 或者用一个 roots 数组表示集合的根节点（代表元）
	var f func(int)
	f = func(p int) {
		if p == n {
			// do groups ...
			fmt.Println(groups)
			return
		}
		groups = append(groups, []int{p})
		f(p + 1)
		groups = groups[:len(groups)-1]
		for i := range groups {
			groups[i] = append(groups[i], p)
			f(p + 1)
			groups[i] = groups[i][:len(groups[i])-1]
		}
	}

	f(0)
}
