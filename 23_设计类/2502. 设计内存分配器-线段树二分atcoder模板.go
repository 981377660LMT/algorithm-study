package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	p2894()
}

// P2894 [USACO08FEB] Hotel G
// https://www.luogu.com.cn/problem/P2894
// 第一行输入 n,m，n 代表有 n 个房间 (1≤n≤50,000)，编号为 1∼n，开始都为空房，m 表示以下有 m 行操作 (1≤m<50,000)，以下每行先输入一个数 i ，表示一种操作：
// 若 i 为 1，表示查询房间，再输入一个数 x，表示在 1,2,...,n 房间中找到长度为 x 的连续空房，输出连续 x 个房间中左端的房间号，尽量让这个房间号最小，若找不到长度为 x 的连续空房，输出 0。若找得到，在这 x 个空房间中住上人。
// 若 i 为 2，表示退房，再输入两个数 x,y 代表房间号 x∼x+y−1 退房，即让房间为空。
func p2894() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)

	seg := NewLazySegTree32(n, func(i int32) E { return E{pre0: 1, suf0: 1, max0: 1, len: 1} })
	for i := int32(0); i < m; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var x int32
			fmt.Fscan(in, &x)
			end := seg.MaxRight(0, func(data E) bool { return data.max0 < x })
			if end == n {
				fmt.Fprintln(out, 0)
			} else {
				start := end + 1 - x
				fmt.Fprintln(out, start+1)
				seg.Update(start, start+x, 1)
			}
		} else {
			var x, y int32
			fmt.Fscan(in, &x, &y)
			seg.Update(x-1, x-1+y, 0)
		}
	}
}

type E = struct {
	pre0, suf0, max0 int32
	len              int32
}

type Id = int8

func (*LazySegTree32) e() E   { return E{} }
func (*LazySegTree32) id() Id { return -1 }
func (*LazySegTree32) op(left, right E) E {
	res := E{pre0: left.pre0, suf0: right.suf0, len: left.len + right.len}
	if left.pre0 == left.len {
		res.pre0 += right.pre0
	}
	if right.suf0 == right.len {
		res.suf0 += left.suf0
	}
	res.max0 = max32(max32(left.max0, right.max0), left.suf0+right.pre0)
	return res
}
func (*LazySegTree32) mapping(f Id, g E) E {
	if f == 0 {
		l := g.len
		return E{pre0: l, suf0: l, max0: l, len: l}
	}
	if f == 1 {
		return E{len: g.len}
	}
	return g
}

func (*LazySegTree32) composition(f, g Id) Id {
	return f
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
func max32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}

// !template
type LazySegTree32 struct {
	n    int32
	size int32
	log  int32
	data []E
	lazy []Id
}

func NewLazySegTree32(n int32, f func(int32) E) *LazySegTree32 {
	tree := &LazySegTree32{}
	tree.n = n
	tree.log = int32(bits.Len32(uint32(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := int32(0); i < n; i++ {
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

func NewLazySegTree32From(leaves []E) *LazySegTree32 {
	tree := &LazySegTree32{}
	n := int32(len(leaves))
	tree.n = n
	tree.log = int32(bits.Len32(uint32(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := int32(0); i < n; i++ {
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
func (tree *LazySegTree32) Query(left, right int32) E {
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
func (tree *LazySegTree32) QueryAll() E {
	return tree.data[1]
}
func (tree *LazySegTree32) GetAll() []E {
	for i := int32(1); i < tree.size; i++ {
		tree.pushDown(i)
	}
	res := make([]E, tree.n)
	copy(res, tree.data[tree.size:tree.size+tree.n])
	return res
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree32) Update(left, right int32, f Id) {
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
	for i := int32(1); i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree32) MinLeft(right int32, predicate func(data E) bool) int32 {
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
func (tree *LazySegTree32) MaxRight(left int32, predicate func(data E) bool) int32 {
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
func (tree *LazySegTree32) Get(index int32) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *LazySegTree32) Set(index int32, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := int32(1); i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTree32) pushUp(root int32) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}

func (tree *LazySegTree32) pushDown(root int32) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}

func (tree *LazySegTree32) propagate(root int32, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func (tree *LazySegTree32) String() string {
	var sb []string
	sb = append(sb, "[")
	for i := int32(0); i < tree.n; i++ {
		if i != 0 {
			sb = append(sb, ", ")
		}
		sb = append(sb, fmt.Sprintf("%v", tree.Get(i)))
	}
	sb = append(sb, "]")
	return strings.Join(sb, "")
}
