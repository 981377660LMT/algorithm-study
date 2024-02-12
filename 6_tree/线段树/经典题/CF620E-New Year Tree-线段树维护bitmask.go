package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

// !CF620E-New Year Tree-线段树维护bitmask
// https://www.luogu.com.cn/problem/CF620E
// 给定一棵n个节点的树,每个节点有一种颜色.
// 有q次操作：
// 1 u x: 把u为根的子树中的所有节点的颜色都变成x
// 2 u :查询u为根的子树中有多少种颜色
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &colors[i])
	}

	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	lid, rid := make([]int, n), make([]int, n)
	dfn := 0
	var dfs func(cur, pre int)
	dfs = func(cur, pre int) {
		lid[cur] = dfn
		dfn++
		for _, next := range tree[cur] {
			if next != pre {
				dfs(next, cur)
			}
		}
		rid[cur] = dfn
	}
	dfs(0, -1)

	leaves := make([]int, n)
	for i := 0; i < n; i++ {
		leaves[lid[i]] = 1 << colors[i]
	}
	seg := NewLazySegTreeFrom(leaves)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var u, x int
			fmt.Fscan(in, &u, &x)
			u--
			seg.Update(lid[u], rid[u], x)
		} else {
			var u int
			fmt.Fscan(in, &u)
			u--
			res := seg.Query(lid[u], rid[u])
			fmt.Fprintln(out, bits.OnesCount(uint(res)))
		}
	}
}

const INF = 1e18

// RangeAssignRangeCountColor

type E = int
type Id = int

func (*LazySegTree) e() E   { return 0 }
func (*LazySegTree) id() Id { return -1 }
func (*LazySegTree) op(left, right E) E {
	return left | right
}
func (*LazySegTree) mapping(f Id, g E) E {
	if f == -1 {
		return g
	}
	return 1 << f
}
func (*LazySegTree) composition(f, g Id) Id {
	if f == -1 {
		return g
	}
	return f
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// !template
type LazySegTree struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id
}

func NewLazySegTree(n int, f func(int) E) *LazySegTree {
	tree := &LazySegTree{}
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
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

func NewLazySegTreeFrom(leaves []E) *LazySegTree {
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
//
//	0<=left<=right<=len(tree.data)
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
func (tree *LazySegTree) GetAll() []E {
	res := make([]E, tree.n)
	copy(res, tree.data[tree.size:tree.size+tree.n])
	return res
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
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
				if tmp := tree.op(tree.data[right], res); predicate(tmp) {
					res = tmp
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
				if tmp := tree.op(res, tree.data[left]); predicate(tmp) {
					res = tmp
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

func (tree *LazySegTree) String() string {
	var sb []string
	sb = append(sb, "[")
	for i := 0; i < tree.n; i++ {
		if i != 0 {
			sb = append(sb, ", ")
		}
		sb = append(sb, fmt.Sprintf("%v", tree.Get(i)))
	}
	sb = append(sb, "]")
	return strings.Join(sb, "")
}
