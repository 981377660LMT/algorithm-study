// F - Confluence
// 学生上学合流

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// https://atcoder.jp/contests/abc183/tasks/abc183_f
	// 给定n个元素和q个查询
	// 1 a b 将a和b合并
	// 2 x y 查询时已经和x合流的学生中,y班的学生有多少人

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	class := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &class[i])
		class[i]--
	}

	ufa := NewUnionFindArray(n)
	mps := make([]map[int]int, n)
	for i := 0; i < n; i++ {
		mps[i] = make(map[int]int)
		mps[i][class[i]] = 1
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var a, b int
			fmt.Fscan(in, &a, &b)
			a, b = a-1, b-1
			if !ufa.IsConnected(a, b) {
				ufa.UnionWithCallback(a, b, func(big, small int) {
					for k, v := range mps[small] {
						mps[big][k] += v
					}
				})
			}
		} else {
			var x, y int
			fmt.Fscan(in, &x, &y)
			x, y = x-1, y-1
			leader := ufa.Find(x)
			mp := mps[leader]
			fmt.Fprintln(out, mp[y])
		}
	}
}

func NewUnionFindArray(n int) *UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	cb(root2, root1)
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
