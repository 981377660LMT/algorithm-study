// No.899-妖精集合
// https://yukicoder.me/problems/no/899/editorial
// 给定一棵树,每个点有若干个妖精
// 把食物放在某个点v,距离这个点<=2的所有妖精都会集合到这个点
// !给定q个查询,问把食物放在qi时,有多少个妖精会集合到这个点

// 1. 父亲和父亲的父亲的妖精都会集合到这个点
// 2. 子树中特定深度的结点个数(dep=curDep,curDep+1,curDep+2)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	adjList := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adjList[u] = append(adjList[u], Edge{v, 1})
		adjList[v] = append(adjList[v], Edge{u, 1})
	}
	B := NewBFSNumbering(adjList, 0)

	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	leaves := make([]E, n)
	for i := 0; i < n; i++ {
		id := B.Id[i]
		leaves[id] = E{size: 1, sum: nums[i]}
	}
	seg := NewLazySegTree(leaves)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var root int
		fmt.Fscan(in, &root)

		p := B.Parent[root]
		pp := -1
		if p != -1 {
			pp = B.Parent[p]
		}

		res := 0
		// 移动非子树部分
		if pp >= 0 {
			res += seg.Get(B.Id[pp]).sum
			seg.Set(B.Id[pp], E{size: 1, sum: 0})
		}
		if p >= 0 {
			res += seg.Get(B.Id[p]).sum
			seg.Set(B.Id[p], E{size: 1, sum: 0})
			left, right := B.FindRange(p, B.Depth[p]+1) // !兄弟结点
			res += seg.Query(left, right).sum
			seg.Update(left, right, Id{mul: 0, add: 0})
		}

		// 移动子树部分
		for d := 0; d < 3; d++ {
			left, right := B.FindRange(root, B.Depth[root]+d)
			res += seg.Query(left, right).sum
			seg.Update(left, right, Id{mul: 0, add: 0})
		}

		fmt.Fprintln(out, res)
		seg.Set(B.Id[root], E{size: 1, sum: res})
	}
}

const INF = 1e18

// RangeAffineRangeSum (这里用来将区间全部置为0)
type E = struct{ size, sum int }
type Id = struct{ mul, add int }

func (*LazySegTree) e() E   { return E{size: 1} }
func (*LazySegTree) id() Id { return Id{mul: 1} }
func (*LazySegTree) op(left, right E) E {
	return E{
		size: left.size + right.size,
		sum:  (left.sum + right.sum),
	}
}

func (*LazySegTree) mapping(lazy Id, data E) E {
	return E{
		size: data.size,
		sum:  (data.sum*lazy.mul + data.size*lazy.add),
	}
}

func (*LazySegTree) composition(parentLazy, childLazy Id) Id {
	return Id{
		mul: (parentLazy.mul * childLazy.mul),
		add: (parentLazy.mul*childLazy.add + parentLazy.add),
	}
}

// !template
type LazySegTree struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id
}

func NewLazySegTree(leaves []E) *LazySegTree {
	tree := &LazySegTree{}
	n := len(leaves)
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = leaves[i]
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//   0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Query(left, right int) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}
func (tree *LazySegTree) QueryAll() E {
	return tree.data[1]
}

// 更新切片[left:right]的值
//   0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Update(left, right int, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := 1; i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree) MinLeft(right int, predicate func(data E) bool) int {
	if right == 0 {
		return 0
	}
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown((right - 1) >> i)
	}
	res := tree.e()
	for {
		right--
		for right > 1 && right&1 != 0 {
			right >>= 1
		}
		if !predicate(tree.op(tree.data[right], res)) {
			for right < tree.size {
				tree.pushDown(right)
				right = right<<1 | 1
				if predicate(tree.op(tree.data[right], res)) {
					res = tree.op(tree.data[right], res)
					right--
				}
			}
			return right + 1 - tree.size
		}
		res = tree.op(tree.data[right], res)
		if (right & -right) == right {
			break
		}
	}
	return 0
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree) MaxRight(left int, predicate func(data E) bool) int {
	if left == tree.n {
		return tree.n
	}
	left += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}
	res := tree.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data[left])) {
			for left < tree.size {
				tree.pushDown(left)
				left <<= 1
				if predicate(tree.op(res, tree.data[left])) {
					res = tree.op(res, tree.data[left])
					left++
				}
			}
			return left - tree.size
		}
		res = tree.op(res, tree.data[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return tree.n
}

// 单点查询(不需要 pushUp/op 操作时使用)
func (tree *LazySegTree) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *LazySegTree) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTree) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *LazySegTree) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *LazySegTree) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
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

// 查询root的子树中,深度为dep的顶点的欧拉序/括号序的左闭右开区间[left, right).
//  0 <= left < right <= n.
//  dep 是绝对深度.
func (b *BFSNumbering) FindRange(root, dep int) (left, right int) {
	if dep < b.Depth[root] || dep >= len(b.depId)-1 {
		return 0, 0
	}
	left1, right1 := b.lid[root], b.rid[root]
	left2, right2 := b.depId[dep], b.depId[dep+1]
	left = b.bs(left2-1, right2, left1)
	right = b.bs(left2-1, right2, right1)
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
