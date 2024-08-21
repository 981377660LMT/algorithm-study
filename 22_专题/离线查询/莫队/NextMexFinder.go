// NextMex(v)：返回容器中缺失的最小的v的倍数的数，并将这个数加入容器中。

package main

import "fmt"

func main() {
	M := NewNextMexFinder()
	fmt.Println(M.NextMex(4)) // 4
	fmt.Println(M.NextMex(2)) // 2
	fmt.Println(M.NextMex(2)) // 6
}

type NextMexFinder struct {
	next    map[int]int
	visited map[int]struct{}
}

func NewNextMexFinder() *NextMexFinder {
	return &NextMexFinder{next: map[int]int{}, visited: map[int]struct{}{0: {}}}
}

func (nmf *NextMexFinder) NextMex(v int) int {
	res := nmf.next[v]
	for {
		_, has := nmf.visited[res]
		if !has {
			break
		}
		res += v
	}
	nmf.visited[res] = struct{}{}
	nmf.next[v] = res
	return res
}
