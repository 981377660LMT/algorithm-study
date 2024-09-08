// yuki1054 Union add query-带权并查集
// https://yukicoder.me/problems/no/1054
// 1 a b: 合并a和b所在的集合.
// 2 a b: 将a所在的集合的所有元素的值加上b.
// 3 a: 输出点a的值.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	groupAdd := make([]int, n)
	values := make([]int, n)
	groups := make([][]int32, n)
	for i := int32(0); i < n; i++ {
		groups[i] = []int32{i}
	}
	uf := NewUnionFindArraySimple32(n)

	union := func(a, b int32) {
		uf.Union(a, b, func(big, small int32) {
			for _, v := range groups[small] {
				groups[big] = append(groups[big], v)
				values[v] += groupAdd[small] - groupAdd[big] // inv
			}
			groups[small] = nil
		})
	}

	addGroup := func(a int32, v int) {
		root := uf.Find(a)
		groupAdd[root] += v // op
	}

	get := func(a int32) int {
		root := uf.Find(a)
		return values[a] + groupAdd[root] // op
	}

	for i := int32(0); i < q; i++ {
		var op, a, b int32
		fmt.Fscan(in, &op, &a, &b)
		if op == 1 {
			a, b = a-1, b-1
			union(a, b)
		} else if op == 2 {
			a -= 1
			addGroup(a, int(b))
		} else {
			a -= 1
			fmt.Fprintln(out, get(a))
		}
	}
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32, beforeMerge func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	if beforeMerge != nil {
		beforeMerge(root1, root2)
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	root := key
	for u.data[root] >= 0 {
		root = u.data[root]
	}
	for key != root {
		key, u.data[key] = u.data[key], root
	}
	return root
}

func (u *UnionFindArraySimple32) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}
