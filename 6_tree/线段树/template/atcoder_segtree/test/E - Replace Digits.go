/* eslint-disable no-console */

// https://atcoder.jp/contests/abl/tasks/abl_e

// 输入 n(≤2e5) 和 q(≤2e5)。
// 初始有一个长为 n 的字符串 s，
// !所有字符都是 1，s 的下标从 1 开始。
// 然后输入 q 个替换操作，每个操作输入 L,R (1≤L≤R≤n) 和 d (1≤d≤9)。
// !你需要把 s 的 [L,R] 内的所有字符替换为 d。
// !对每个操作，把替换后的 s 看成一个十进制数，输出这个数模 998244353 的结果。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

// https://atcoder.jp/contests/abl/tasks/abl_e

const N int = 2e5 + 5
const MOD int = 998244353

var pow10 [N]int
var pow10PreSum [N]int

type Operation struct {
	left, right, target int
}

func replaceDigits(n int, operations []Operation) []int {
	// !预处理pow10 和 pow10PreSum
	pow10[0] = 1
	pow10PreSum[0] = 1
	for i := 1; i <= n; i++ {
		pow10[i] = (pow10[i-1] * 10) % MOD
		pow10PreSum[i] = (pow10PreSum[i-1] + pow10[i]) % MOD
	}

	initNums := make([]E, n)
	for i := range initNums {
		initNums[i] = E{1, 1}
	}

	tree := NewLazySegTree(initNums)
	res := make([]int, len(operations))
	for i, op := range operations {
		tree.Update(op.left-1, op.right, Id(op.target))
		res[i] = tree.QueryAll().Sum
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	// !每个操作 1<=L<=R<=n 1<=d<=9
	operations := make([]Operation, q)
	for i := range operations {
		fmt.Fscan(in, &operations[i].left, &operations[i].right, &operations[i].target)
	}

	res := replaceDigits(n, operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// !线段树维护的值的类型
type E = struct {
	Sum, Length int
}

// !更新操作的值的类型/懒标记的值的类型
type Id = int

// !线段树维护的值的幺元.
//  alias: e
func (tree *LazySegTree) e() E { return E{} }

// !更新操作/懒标记的幺元
//  alias: id
func (tree *LazySegTree) id() Id { return -1 }

// !合并左右区间的值
//  alias: op
func (tree *LazySegTree) op(left, right E) E {
	return E{
		Sum:    (left.Sum*pow10[right.Length] + right.Sum) % MOD,
		Length: left.Length + right.Length,
	}
}

// !父结点的懒标记更新子结点的值
//  alias: mapping
func (tree *LazySegTree) mapping(lazy Id, data E) E {
	if lazy == -1 {
		return data
	}

	data.Sum = (int(lazy) * pow10PreSum[data.Length-1]) % MOD
	return data
}

// !合并父结点的懒标记和子结点的懒标记
//  alias: composition
func (tree *LazySegTree) composition(parentLazy, childLazy Id) Id {
	if parentLazy >= 0 {
		return parentLazy
	}
	return childLazy
}

//
//
//
//
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
