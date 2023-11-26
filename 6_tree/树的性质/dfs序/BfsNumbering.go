// https://maspypy.github.io/library/graph/ds/bfs_numbering.hpp
// 给一颗以1为根大小为n的有根树，每个点有点权，查询若干次，每次查询以x的子树，深度为y的这一层中，所有点权的和（节点x的深度为0）.
// !求每个root的子树中,绝对深度为dep的顶点的欧拉序/括号序的范围

// !ID[v]：每个顶点的欧拉序编号 (0-indexed)
// !FindRange(v, dep)：以v为顶点的子树中, `绝对深度`为dep的顶点的欧拉序的范围(左闭右开)

package main

import (
	"bufio"
	"fmt"
	"os"
)

//
//     4(0)
//     / \
//    /   \
//   3(2)  2(1)
//  /   \
// 0(3)  1(4)

func main() {
	// https://atcoder.jp/contests/abc202/tasks/abc202_e
	// !子树中特定深度的结点个数
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	parents := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &parents[i])
		parents[i]--
	}

	g := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		p := parents[i]
		g[p] = append(g[p], Edge{i + 1, 1})
	}
	B := NewBFSNumbering(g, 0)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var root, dep int
		fmt.Fscan(in, &root, &dep)
		root--
		left, right := B.FindRange(root, dep)
		fmt.Fprintln(out, right-left)
	}

}

type Edge struct{ to, weight int }
type BFSNumbering struct {
	Depth  []int // 每个点的绝对深度(0-based)
	Id     []int // 每个点的欧拉序起点编号(0-based)
	Parent []int // 不存在时为-1

	count           int
	root            int
	vs              []int
	lid, rid, depId []int
	lidSeq          []int
	graph           [][]Edge
}

func NewBFSNumbering(graph [][]Edge, root int) *BFSNumbering {
	res := &BFSNumbering{graph: graph, root: root}
	res.build()
	return res
}

// 查询root的子树中,`绝对深度`为dep的顶点的欧拉序/括号序的左闭右开区间[start, end).
//
//	0 <= start < end <= n.
func (b *BFSNumbering) FindRange(root, dep int) (start, end int) {
	if dep < b.Depth[root] || dep >= len(b.depId)-1 {
		return 0, 0
	}
	left1, right1 := b.lid[root], b.rid[root]
	left2, right2 := b.depId[dep], b.depId[dep+1]
	start = b.bs(left2-1, right2, left1)
	end = b.bs(left2-1, right2, right1)
	return
}

func (b *BFSNumbering) build() {
	n := len(b.graph)
	b.vs = make([]int, 0, n)
	b.Parent = make([]int, n)
	for i := range b.Parent {
		b.Parent[i] = -1
	}
	b.Id = make([]int, n)
	b.lid = make([]int, n)
	b.rid = make([]int, n)
	b.Depth = make([]int, n)
	b.bfs()
	b.dfs(b.root)
	d := maxs(b.Depth...)
	b.depId = make([]int, d+2)
	for i := 0; i < n; i++ {
		b.depId[b.Depth[i]+1]++
	}
	for i := 0; i < d+1; i++ {
		b.depId[i+1] += b.depId[i]
	}
	b.lidSeq = make([]int, 0, n)
	for i := 0; i < n; i++ {
		b.lidSeq = append(b.lidSeq, b.lid[b.vs[i]])
	}
}

func (b *BFSNumbering) bfs() {
	queue := []int{b.root}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		b.Id[v] = len(b.vs)
		b.vs = append(b.vs, v)
		for _, e := range b.graph[v] {
			if e.to == b.Parent[v] {
				continue
			}
			queue = append(queue, e.to)
			b.Parent[e.to] = v
			b.Depth[e.to] = b.Depth[v] + 1
		}
	}
}

func (b *BFSNumbering) dfs(v int) {
	b.lid[v] = b.count
	b.count++
	for _, e := range b.graph[v] {
		if e.to == b.Parent[v] {
			continue
		}
		b.dfs(e.to)
	}
	b.rid[v] = b.count
}

func (b *BFSNumbering) bs(left, right, x int) int {
	for left+1 < right {
		mid := (left + right) / 2
		if b.lidSeq[mid] >= x {
			right = mid
		} else {
			left = mid
		}
	}
	return right
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
