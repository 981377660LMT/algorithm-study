// CF1620E Replace the Numbers-在线查询
// https://www.luogu.com.cn/problem/CF1620E
// 给出 q 个操作，操作分为两种：

// 1 x 在序列末尾插入数字 x。
// 2 x y 把序列中的所有 x 替换为 y。

// 求这个序列操作后的结果。
// 并查集，类似未来日记中的技巧

package main

import (
	"bufio"
	"fmt"
	"os"
)

type UfMap struct {
	parent map[int32]int32
}

func NewUfMap() *UfMap {
	return &UfMap{parent: make(map[int32]int32)}
}

func (uf *UfMap) Find(x int) int {
	x32 := int32(x)
	if _, ok := uf.parent[x32]; !ok {
		uf.parent[x32] = x32
		return x
	}
	for uf.parent[x32] != x32 {
		uf.parent[x32] = uf.parent[uf.parent[x32]]
		x32 = uf.parent[x32]
	}
	return int(x32)
}

func (uf *UfMap) UnionTo(child, parent int) bool {
	root1, root2 := int32(uf.Find(child)), int32(uf.Find(parent))
	if root1 == root2 {
		return false
	}
	uf.parent[root1] = root2
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	firstIndex := make(map[int]int) // 每个值出现的第一个位置
	values := make([]int, 0)        // 每个位置对应的值
	uf := NewUfMap()                // 维护哪些位置的值相同
	curIndex := 0
	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var x int
			fmt.Fscan(in, &x)

			values = append(values, x)
			if _, ok := firstIndex[x]; !ok {
				firstIndex[x] = curIndex
			} else {
				uf.UnionTo(curIndex, firstIndex[x])
			}
			curIndex++
		} else {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if x == y {
				continue
			}
			if _, ok := firstIndex[x]; !ok {
				continue
			}
			if _, ok := firstIndex[y]; !ok {
				firstIndex[y] = firstIndex[x]
				values[firstIndex[y]] = y
				delete(firstIndex, x)
				continue
			}
			uf.UnionTo(firstIndex[x], firstIndex[y])
			delete(firstIndex, x)
		}
	}

	res := make([]int, 0, curIndex)
	for i := 0; i < curIndex; i++ {
		res = append(res, values[uf.Find(i)])
	}
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}
