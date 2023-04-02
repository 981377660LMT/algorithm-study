// 覆盖区间的方案数/区间覆盖的方案数
// 给定[0,X]内的n个`闭区间`
// 选择其中若干个,使得这些区间能够覆盖[0,X]内的所有点
// 求方案数
// n<=1e5

// !按终点排序,dp[i][mex]表示前i个区间覆盖,最左端没有被覆盖的点为mex的方案数
// 选区间的问题都可以用数据结构优化
//  把区间按照右端点排序 dp[i][mex] 然后用区间更新

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func main() {

	bf := func(intervals [][2]int, x int) int {
		n := len(intervals)
		res := 0
		for s := 1; s < 1<<n; s++ {
			cur := make(map[int]struct{})
			for i := 0; i < n; i++ {
				if s>>i&1 == 1 {
					for v := intervals[i][0]; v <= intervals[i][1]; v++ {
						cur[v] = struct{}{}
					}
				}
			}
			flag := true
			for v := 0; v <= x; v++ {
				if _, ok := cur[v]; !ok {
					flag = false
					break
				}
			}
			if flag {
				res++
			}
		}
		return res
	}
	for i := 0; i < 10; i++ {
		n := rand.Intn(12) + 1
		x := rand.Intn(1000)
		intervals := make([][2]int, n)
		for i := range intervals {
			s := rand.Intn(1000)
			e := rand.Intn(1000)
			if s > e {
				s, e = e, s
			}
			intervals[i] = [2]int{s, e}
		}
		if solve(intervals, x) != bf(intervals, x) {
			fmt.Println(intervals, x)
			fmt.Println(solve(intervals, x), bf(intervals, x))
			break
		}
	}
	fmt.Println("pass")

	for i := 0; i < 10; i++ {
		n := rand.Intn(1e5) + 1
		x := rand.Intn(1e5)
		intervals := make([][2]int, n)
		for i := range intervals {
			s := rand.Intn(1e5)
			e := rand.Intn(1e5)
			if s > e {
				s, e = e, s
			}
			intervals[i] = [2]int{s, e}
		}
		time1 := time.Now()
		solve(intervals, x)
		fmt.Println(time.Since(time1))
	}
}

const MOD int = 1e9 + 7

// 按照区间右端点排序,dp[i][mex]表示前i个区间,未被覆盖的最左端点为mex的方案数.
func solve(intervals [][2]int, x int) int {
	sort.Slice(intervals, func(i, j int) bool { return intervals[i][1] < intervals[j][1] })
	maxEnd := max(x, intervals[len(intervals)-1][1])
	leaves := make([]E, maxEnd+2) // mex: 0-maxEnd+1
	for i := range leaves {
		leaves[i] = E{size: 1}
	}
	dp := NewLazySegTree(leaves)
	dp.Set(0, E{sum: 1, size: 1})
	for _, interval := range intervals {
		s, e := interval[0], interval[1]
		dp.Update(0, s, Id{mul: 2, add: 0}) // mex为之前的值,乘以2
		sum_ := dp.Query(s, e+2).sum
		dp.Update(e+1, e+2, Id{mul: 1, add: sum_}) // mex 为 e+1
	}
	return dp.Query(x+1, maxEnd+2).sum
}

// RangeAffineRangeSum (这里用来做区间乘法)
type E = struct{ size, sum int }
type Id = struct{ mul, add int }

func (*LazySegTree) e() E   { return E{size: 1} }
func (*LazySegTree) id() Id { return Id{mul: 1} }
func (*LazySegTree) op(left, right E) E {
	return E{
		size: left.size + right.size,
		sum:  (left.sum + right.sum) % MOD,
	}
}

func (*LazySegTree) mapping(lazy Id, data E) E {
	return E{
		size: data.size,
		sum:  (data.sum*lazy.mul + data.size*lazy.add) % MOD,
	}
}

func (*LazySegTree) composition(parentLazy, childLazy Id) Id {
	return Id{
		mul: (parentLazy.mul * childLazy.mul) % MOD,
		add: (parentLazy.mul*childLazy.add + parentLazy.add) % MOD,
	}
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
