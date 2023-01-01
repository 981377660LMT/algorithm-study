// manhattanmst
// F - Permutation Distance-曼哈顿最小生成树
// 曼哈顿距离最小生成树
// 给定 对1-n的排列perm n<=2e5
// !对每个i 求出 abs(perm[i]-perm[j])+abs(i-j) (i!=j) 的最小值
// 即对每个点求二维平面上的最近曼哈顿距离

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const INF int = 1e18

type Point struct{ x, y int }

func permutationDistance(perm []int) []int {
	n := len(perm)
	// !对每个点求二维平面上的最近曼哈顿距离
	points := make([]Point, n)
	for i := range perm {
		points[i] = Point{perm[i], i + 1}
	}

	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = INF
	}

	_, mstEdges := manhattanMST(points)
	for _, e := range mstEdges {
		v, w := e.v, e.w
		cost := abs(points[v].x-points[w].x) + abs(points[v].y-points[w].y)
		res[v] = min(res[v], cost)
		res[w] = min(res[w], cost)
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	perm := make([]int, n)
	for i := range perm {
		fmt.Fscan(in, &perm[i])
	}

	res := permutationDistance(perm)
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// https://github.dev/EndlessCheng/codeforces-go/blob/a0733fa7a046673ff42a058b0dca7852646fbf3b/copypasta/graph.go#L1835
func manhattanMST(points []Point) (mst int, mstEdges []struct{ v, w int }) {
	n := len(points)
	pointsCopy := make([]struct{ x, y, i int }, n)
	for i, p := range points {
		pointsCopy[i] = struct{ x, y, i int }{p.x, p.y, i}
	}

	type edge struct{ v, w, dis int }
	edges := []edge{}

	build := func() {
		sort.Slice(pointsCopy, func(i, j int) bool { a, b := pointsCopy[i], pointsCopy[j]; return a.x < b.x || a.x == b.x && a.y < b.y })

		// 离散化 y-x
		type pair struct{ v, i int }
		ps := make([]pair, n)
		for i, p := range pointsCopy {
			ps[i] = pair{p.y - p.x, i}
		}
		sort.Slice(ps, func(i, j int) bool { return ps[i].v < ps[j].v })
		kth := make([]int, n)
		k := 1
		kth[ps[0].i] = k
		for i := 1; i < n; i++ {
			if ps[i].v != ps[i-1].v {
				k++
			}
			kth[ps[i].i] = k
		}

		t := newFenwickTree(k + 1)
		for i := n - 1; i >= 0; i-- {
			p := pointsCopy[i]
			pos := kth[i]
			if j := t.query(pos); j != -1 {
				q := pointsCopy[j]
				dis := abs(p.x-q.x) + abs(p.y-q.y)
				edges = append(edges, edge{p.i, q.i, dis})
			}
			t.update(pos, p.x+p.y, i)
		}
	}
	build()
	for i := range pointsCopy {
		pointsCopy[i].x, pointsCopy[i].y = pointsCopy[i].y, pointsCopy[i].x
	}
	build()
	for i := range pointsCopy {
		pointsCopy[i].x = -pointsCopy[i].x
	}
	build()
	for i := range pointsCopy {
		pointsCopy[i].x, pointsCopy[i].y = pointsCopy[i].y, pointsCopy[i].x
	}
	build()

	sort.Slice(edges, func(i, j int) bool { return edges[i].dis < edges[j].dis })

	uf := newUnionFindArray(n)
	left := n - 1
	for _, e := range edges {
		if uf.union(e.v, e.w) {
			mst += e.dis
			mstEdges = append(mstEdges, struct{ v, w int }{e.v, e.w})
			left--
			if left == 0 {
				break
			}
		}
	}
	return
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func newUnionFindArray(n int) *unionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}
	return &unionFindArray{
		Rank:   rank,
		parent: parent,
	}
}

type unionFindArray struct {
	Rank   []int
	parent []int
}

func (ufa *unionFindArray) union(key1, key2 int) bool {
	root1, root2 := ufa.find(key1), ufa.find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.Rank[root1] > ufa.Rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.Rank[root2] += ufa.Rank[root1]
	return true
}

func (ufa *unionFindArray) find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

type fenwickTree struct {
	tree, idRec []int
}

func newFenwickTree(n int) *fenwickTree {
	tree := make([]int, n)
	idRec := make([]int, n)
	for i := range tree {
		tree[i], idRec[i] = math.MaxInt64, -1
	}
	return &fenwickTree{tree, idRec}
}

func (f *fenwickTree) update(pos, val, id int) {
	for ; pos > 0; pos &= pos - 1 {
		if val < f.tree[pos] {
			f.tree[pos], f.idRec[pos] = val, id
		}
	}
}

func (f *fenwickTree) query(pos int) int {
	minVal, minID := math.MaxInt64, -1
	for ; pos < len(f.tree); pos += pos & -pos {
		if f.tree[pos] < minVal {
			minVal, minID = f.tree[pos], f.idRec[pos]
		}
	}
	return minID
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
