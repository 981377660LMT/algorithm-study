// カッコ列をグラフにする。各頂点の範囲を表す配列 LR も作る。
// 全体を表す根ノードも作って、N+1頂点。

// (())() → 得到(n//2)+1 个结点,其中0号结点是虚拟根结点
// graph: [[1 3] [2] [] []] (有向邻接表)
// leftRight: [[0 6] [0 3] [1 2] [4 5]] (每个顶点的括号范围,左闭右开)
//
//           0 [0,6)
//          / \
//   [0,4) 1   3 [4,6)
//        /
// [1,3) 2

// 有效的括号序列形成的树(括号树)，结合LCA

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// fmt.Println(BracketTree("(())()"))
	// https://yukicoder.me/problems/no/1778
	// 给定一个有效的括号序列,每次可以删除一段匹配的括号
	// 给定q个查询,每个查询形如 [start1,start2]
	// 求包含这两段括号的区间中最靠内部的一段区间

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	var s string
	fmt.Fscan(in, &s)

	tree, leftRight := BracketTree(s)

	idToRoot := make([]int, n) // 每个起点/终点位置对应的树节点 [0, n)
	for i := 1; i < len(leftRight); i++ {
		left, right := leftRight[i][0], leftRight[i][1]
		idToRoot[left] = i
		idToRoot[right-1] = i
	}

	lca := NewLCA((n / 2) + 1)
	for i := 0; i < len(tree); i++ {
		for _, j := range tree[i] {
			lca.AddEdge(i, j)
		}
	}
	lca.Build(0)

	for i := 0; i < q; i++ {
		var start, end int // 1<=start<end<=n
		fmt.Fscan(in, &start, &end)
		start--
		end--
		root1, root2 := idToRoot[start], idToRoot[end]
		lca_ := lca.LCA(root1, root2)
		if lca_ == 0 { // 不被包含
			fmt.Fprintln(out, -1)
		} else {
			left, right := leftRight[lca_][0], leftRight[lca_][1]
			fmt.Fprintln(out, left+1, right)
		}
	}
}

// BracketTree :有效括号序列转换成树.
//  (())() → 得到 `(n/2)+1` 个结点,其中0号结点是虚拟根结点.
//  tree: [[1 3] [2] [] []] (有向邻接表)
//  leftRight: [[0 6) [0 4) [1 3) [4 6)] (每个顶点的括号范围,左闭右开)
//
//           0 [0,6)
//          / \
//   [0,4) 1   3 [4,6)
//        /
// [1,3) 2
func BracketTree(s string) (tree [][]int, leftRight [][2]int) {
	n := len(s) / 2
	tree = make([][]int, n+1)
	leftRight = make([][2]int, n+1)
	now, nxt := 0, 1
	leftRight[0] = [2]int{0, len(s)}
	par := make([]int, n+1)
	for i := range par {
		par[i] = -1
	}

	for i := range s {
		if s[i] == '(' {
			tree[now] = append(tree[now], nxt)
			par[nxt] = now
			leftRight[nxt][0] = i
			now = nxt
			nxt++
		}
		if s[i] == ')' {
			leftRight[now][1] = i + 1
			now = par[now]
		}
	}

	return
}

type LCAHLD struct {
	Depth, Parent      []int
	Tree               [][]int
	dfn, top, heavySon []int
	dfnId              int
}

func NewLCA(n int) *LCAHLD {
	tree := make([][]int, n)
	dfn := make([]int, n)      // vertex => dfn
	top := make([]int, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n)    // 深度
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	return &LCAHLD{
		Tree:     tree,
		dfn:      dfn,
		top:      top,
		Depth:    depth,
		Parent:   parent,
		heavySon: heavySon,
	}
}

// 添加无向边 u-v.
func (hld *LCAHLD) AddEdge(u, v int) {
	hld.Tree[u] = append(hld.Tree[u], v)
	hld.Tree[v] = append(hld.Tree[v], u)
}

func (hld *LCAHLD) Build(root int) {
	hld.build(root, -1, 0)
	hld.markTop(root, root)
}

func (hld *LCAHLD) LCA(u, v int) int {
	for {
		if hld.dfn[u] > hld.dfn[v] {
			u, v = v, u
		}
		if hld.top[u] == hld.top[v] {
			return u
		}
		v = hld.Parent[hld.top[v]]
	}
}

func (hld *LCAHLD) Dist(u, v int) int {
	return hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)]
}

func (hld *LCAHLD) build(cur, pre, dep int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, next := range hld.Tree[cur] {
		if next != pre {
			nextSize := hld.build(next, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	hld.Depth[cur] = dep
	hld.heavySon[cur] = heavySon
	hld.Parent[cur] = pre
	return subSize
}

func (hld *LCAHLD) markTop(cur, top int) {
	hld.top[cur] = top
	hld.dfn[cur] = hld.dfnId
	hld.dfnId++
	if hld.heavySon[cur] != -1 {
		hld.markTop(hld.heavySon[cur], top)
		for _, next := range hld.Tree[cur] {
			if next != hld.heavySon[cur] && next != hld.Parent[cur] {
				hld.markTop(next, next)
			}
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
