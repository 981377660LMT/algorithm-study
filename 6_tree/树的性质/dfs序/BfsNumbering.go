// bfs序编号.
// https://maspypy.github.io/library/graph/ds/bfs_numbering.hpp
// !求每个root的子树中,绝对深度为dep的顶点的bfs序的范围.

// !ID[v]：每个顶点的bfs序 (0-indexed)
// !GetRange(v, dep)：以v为顶点的子树中, `绝对深度`为depth的顶点的bfs序的范围(左闭右开)

package main

import (
	"bufio"
	"fmt"
	"os"
)

// 顶点(bfs序)
//
//        0(0)
//        / \
//       /   \
//      /     \
//     1(1)    2(2)
//    / |  \
//   /  |   \
//  3(3) 4(4) 6(5)
//  |
//  5(6)

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

	g := make([][][2]int, n)
	for i := 0; i < n-1; i++ {
		p := parents[i]
		g[p] = append(g[p], [2]int{i + 1, 1})
	}
	B := NewBFSNumbering(g, 0)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var root, dep int
		fmt.Fscan(in, &root, &dep)
		root--
		left, right := B.GetRange(root, dep)
		fmt.Fprintln(out, right-left)
	}

}

// BFS序编号.
type BFSNumbering struct {
	Id       []int32 // 每个点的bfs序编号(0-based)
	Depth    []int32 // 每个点的绝对深度(0-based)
	Parent   []int32 // 不存在时为-1
	Lid, Rid []int32 // 每个点欧拉序的左右区间

	graph       [][][2]int
	root        int32
	dfn         int32
	bfsOrder    []int32 // 按照bfs序遍历的顶点
	depthPreSum []int32
	lidSeq      []int32
}

func NewBFSNumbering(graph [][][2]int, root int) *BFSNumbering {
	res := &BFSNumbering{graph: graph, root: int32(root)}
	res.build()
	return res
}

// 查询root的子树中,`绝对深度`为depth的顶点的欧拉序/括号序的左闭右开区间[start, end).
//
//	0 <= start < end <= n.
func (b *BFSNumbering) GetRange(root, depth int) (start, end int) {
	if depth < int(b.Depth[root]) || depth >= len(b.depthPreSum)-1 {
		return 0, 0
	}
	left1, right1 := b.Lid[root], b.Rid[root]
	left2, right2 := b.depthPreSum[depth], b.depthPreSum[depth+1]
	start = int(b.bs(left2-1, right2, left1))
	end = int(b.bs(left2-1, right2, right1))
	return
}

func (b *BFSNumbering) build() {
	n := len(b.graph)
	b.Id = make([]int32, n)
	b.Depth = make([]int32, n)
	b.Lid = make([]int32, n)
	b.Rid = make([]int32, n)
	b.Parent = make([]int32, n)
	for i := range b.Parent {
		b.Parent[i] = -1
	}
	b.bfsOrder = make([]int32, 0, n)
	b.bfs()
	b.dfs(b.root)
	d := int32(-1)
	for _, v := range b.Depth {
		if v > d {
			d = v
		}
	}
	b.depthPreSum = make([]int32, d+2)
	for i := 0; i < n; i++ {
		b.depthPreSum[b.Depth[i]+1]++
	}
	for i := int32(0); i < d+1; i++ {
		b.depthPreSum[i+1] += b.depthPreSum[i]
	}
	b.lidSeq = make([]int32, 0, n)
	for i := 0; i < n; i++ {
		b.lidSeq = append(b.lidSeq, b.Lid[b.bfsOrder[i]])
	}
}

func (b *BFSNumbering) bfs() {
	queue := make([]int32, len(b.graph))
	head, tail := 0, 0
	queue[tail] = b.root
	tail++
	for head < tail {
		v := queue[head]
		head++
		b.Id[v] = int32(len(b.bfsOrder))
		b.bfsOrder = append(b.bfsOrder, v)
		for _, e := range b.graph[v] {
			next := int32(e[0])
			if next == b.Parent[v] {
				continue
			}
			queue[tail] = next
			tail++
			b.Parent[next] = v
			b.Depth[next] = b.Depth[v] + 1
		}
	}
}

func (b *BFSNumbering) dfs(v int32) {
	b.Lid[v] = b.dfn
	b.dfn++
	for _, e := range b.graph[v] {
		next := int32(e[0])
		if next != b.Parent[v] {
			b.dfs(next)
		}
	}
	b.Rid[v] = b.dfn
}

func (b *BFSNumbering) bs(left, right, x int32) int32 {
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
