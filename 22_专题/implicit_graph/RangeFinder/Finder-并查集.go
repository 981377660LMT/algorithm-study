// 寻找前驱后继/区间删除

package main

const INF int = 1e18

// 2612. 最少翻转操作数
// https://leetcode.cn/problems/minimum-reverse-operations/
func minReverseOperations(n int, p int, banned []int, k int) []int {
	finder := [2]*Finder{NewFinder(n), NewFinder(n)}

	for i := 0; i < n; i++ {
		finder[(i&1)^1].Erase(i, i+1)
	}
	for _, i := range banned {
		finder[i&1].Erase(i, i+1)
	}

	getRange := func(i int) (int, int) {
		return max(i-k+1, k-i-1), min(i+k-1, 2*n-k-i-1)
	}
	setUsed := func(u int) {
		finder[u&1].Erase(u, u+1)
	}

	findUnused := func(u int) int {
		left, right := getRange(u)
		pre, ok := finder[(u+k+1)&1].Prev(right)
		if ok && left <= pre && pre <= right {
			return pre
		}
		next, ok := finder[(u+k+1)&1].Next(left)
		if ok && left <= next && next <= right {
			return next
		}
		return -1
	}

	dist := OnlineBfs(n, p, setUsed, findUnused)
	res := make([]int, n)
	for i, d := range dist {
		if d == INF {
			res[i] = -1
		} else {
			res[i] = d
		}
	}
	return res
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

// 利用并查集寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
// 初始时,所有位置都未被访问过.
type Finder struct {
	n                int
	lParent, rParent []int
}

func NewFinder(n int) *Finder {
	lParent, rParent := make([]int, n+2), make([]int, n+2)
	for i := range lParent {
		lParent[i] = i
		rParent[i] = i
	}
	return &Finder{n: n, lParent: lParent, rParent: rParent}
}

// 找到x左侧第一个未被访问过的位置(包含x).
func (f *Finder) Prev(x int) (res int, ok bool) {
	res = f.lFind(x + 1)
	if res > 0 {
		res--
		ok = true
	}
	return
}

// x右侧第一个未被访问过的位置(包含x).
func (f *Finder) Next(x int) (res int, ok bool) {
	res = f.rFind(x + 1)
	if res < f.n+1 {
		res--
		ok = true
	}
	return
}

// 删除[left, right)区间内的元素.
//   0<=left<=right<=n.
func (f *Finder) Erase(left, right int) {
	if left >= right {
		return
	}
	left, right = left+1, right+1
	leftRoot, rightRoot := f.rFind(left), f.rFind(right)
	for rightRoot != leftRoot {
		f.rUnion(leftRoot, leftRoot+1)
		leftRoot = f.rFind(leftRoot + 1)
	}
	leftRoot, rightRoot = f.lFind(left-1), f.lFind(right-1)
	for rightRoot != leftRoot {
		f.lUnion(rightRoot, rightRoot-1)
		rightRoot = f.lFind(rightRoot - 1)
	}
}

func (f *Finder) lUnion(x, y int) {
	if x < y {
		x, y = y, x
	}
	rootX := f.lFind(x)
	rootY := f.lFind(y)
	if rootX == rootY {
		return
	}
	f.lParent[rootX] = rootY
}

func (f *Finder) rUnion(x, y int) {
	if x > y {
		x, y = y, x
	}
	rootX := f.rFind(x)
	rootY := f.rFind(y)
	if rootX == rootY {
		return
	}
	f.rParent[rootX] = rootY
}

func (f *Finder) lFind(x int) int {
	for x != f.lParent[x] {
		f.lParent[x] = f.lParent[f.lParent[x]]
		x = f.lParent[x]
	}
	return x
}

func (f *Finder) rFind(x int) int {
	for x != f.rParent[x] {
		f.rParent[x] = f.rParent[f.rParent[x]]
		x = f.rParent[x]
	}
	return x
}

// 在线bfs.
//   不预先给出图，而是通过两个函数 setUsed 和 findUnused 来在线寻找边.
//   setUsed(u)：将 u 标记为已访问。
//   findUnused(u)：找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `-1`。

func OnlineBfs(
	n int, start int,
	setUsed func(u int), findUnused func(cur int) (next int),
) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	queue := []int{start}
	setUsed(start)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for {
			next := findUnused(cur)
			if next == -1 {
				break
			}
			dist[next] = dist[cur] + 1 // weight
			queue = append(queue, next)
			setUsed(next)
		}
	}

	return
}
