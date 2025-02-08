// RangeAddRangeMax/RangeMaxRangeAdd
// 更快的 RangeAddRangeMax 线段树(比一般的懒更新线段树快 40% 左右)

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	test()
}

type SegmentTreeRangeAddRangeMax struct {
	n    int
	lazy int
	seg  *segmentTree
}

func NewSegmentTreeRangeAddRangeMax(n int, f func(int) int) *SegmentTreeRangeAddRangeMax {
	res := &SegmentTreeRangeAddRangeMax{}
	res.build(n, f)
	return res
}

// [l, r)
func (st *SegmentTreeRangeAddRangeMax) Query(l, r int) int {
	if l < 0 {
		l = 0
	}
	if r > st.n {
		r = st.n
	}
	if l >= r {
		return -INF
	}
	res := st.seg.Query(l, r).max
	for l += st.seg.size; l > 0; l >>= 1 {
		if l&1 == 1 {
			l--
			res += st.seg.seg[l].sum
		}
	}
	return res + st.lazy
}

// [l, r)
func (st *SegmentTreeRangeAddRangeMax) Update(l, r int, x int) {
	if l < 0 {
		l = 0
	}
	if r > st.n {
		r = st.n
	}
	if l >= r {
		return
	}
	st.UpdateSuffix(l, x)
	st.UpdateSuffix(r, -x)
}

// [0, r)
func (st *SegmentTreeRangeAddRangeMax) UpdatePrefix(r int, x int) {
	if r > st.n {
		r = st.n
	}
	if r <= 0 {
		return
	}
	st.lazy += x
	st.UpdateSuffix(r, -x)
}

// [l, n)
func (st *SegmentTreeRangeAddRangeMax) UpdateSuffix(l int, x int) {
	if l < 0 {
		l = 0
	}
	if l >= st.n {
		return
	}
	t := st.seg.Get(l).sum + x
	st.seg.Set(l, E{t, t})
}

func (st *SegmentTreeRangeAddRangeMax) UpdateAll(x int) {
	st.lazy += x
}

func (st *SegmentTreeRangeAddRangeMax) build(n int, f func(int) int) {
	st.lazy = 0
	st.n = n
	pre := 0
	st.seg = newSegmentTree(n, func(i int) E {
		t := f(i) - pre
		pre += t
		return E{t, t}
	})
}

const INF int = 1e18

type E = struct{ sum, max int }

func (*segmentTree) e() E { return E{max: -INF} }
func (*segmentTree) op(a, b E) E {
	a.max = max(a.max, a.sum+b.max)
	a.sum += b.sum
	return a
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

type segmentTree struct {
	n, size int
	seg     []E
}

func newSegmentTree(n int, f func(int) E) *segmentTree {
	res := &segmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *segmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *segmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *segmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}

func test() {
	rand.Seed(time.Now().UnixNano())

	for testCase := 0; testCase < 10000; testCase++ {
		n := rand.Intn(100) + 1 // 数组长度1~100
		seg := NewSegmentTreeRangeAddRangeMax(n, func(i int) int { return 0 })
		bf := NewBruteForce(n)

		ops := rand.Intn(20) + 1 // 操作次数1~20次
		for op := 0; op < ops; op++ {
			action := rand.Intn(5)
			if action == 0 {
				// 区间加操作
				l := rand.Intn(n)
				r := l + rand.Intn(n-l+1)
				x := rand.Intn(20) - 10 // -10到9之间的随机数
				seg.Update(l, r, x)
				bf.Update(l, r, x)
			} else if action == 1 {
				// 区间查询操作
				l := rand.Intn(n)
				r := l + rand.Intn(n-l+1)
				segRes := seg.Query(l, r)
				bfRes := bf.Query(l, r)
				if segRes != bfRes {
					fmt.Printf("测试用例失败！用例编号：%d\n", testCase)
					fmt.Printf("数组长度：%d\n", n)
					fmt.Printf("查询区间：[%d, %d)\n", l, r)
					fmt.Printf("线段树结果：%d，暴力结果：%d\n", segRes, bfRes)
					panic("结果不一致")
				}
			} else if action == 2 {
				// 更新前缀
				r := rand.Intn(n)
				x := rand.Intn(20) - 10 // -10到9之间的随机数
				seg.UpdatePrefix(r, x)
				bf.Update(0, r, x)
			} else if action == 3 {
				// 更新后缀
				l := rand.Intn(n)
				x := rand.Intn(20) - 10 // -10到9之间的随机数
				seg.UpdateSuffix(l, x)
				bf.Update(l, n, x)
			} else if action == 4 {
				// 更新全部
				x := rand.Intn(20) - 10 // -10到9之间的随机数
				seg.UpdateAll(x)
				bf.Update(0, n, x)
			}
		}
	}
	fmt.Println("所有测试用例通过！")
}

// 暴力方法实现
type BruteForce struct {
	diff []int // 差分数组，diff[i] 表示对i位置的差分
	n    int   // 数组长度
}

func NewBruteForce(n int) *BruteForce {
	return &BruteForce{
		diff: make([]int, n+1), // 使用n+1以便处理右边界
		n:    n,
	}
}

func (bf *BruteForce) Update(l, r, x int) {
	if l < 0 {
		l = 0
	}
	if r > bf.n {
		r = bf.n
	}
	if l >= r {
		return
	}
	bf.diff[l] += x
	bf.diff[r] -= x
}

func (bf *BruteForce) Query(l, r int) int {
	if l < 0 {
		l = 0
	}
	if r > bf.n {
		r = bf.n
	}
	if l >= r {
		return -INF
	}

	maxVal := -INF
	currentSum := 0
	for i := 0; i < bf.n; i++ {
		currentSum += bf.diff[i]
		if i >= l && i < r {
			if currentSum > maxVal {
				maxVal = currentSum
			}
		}
	}
	return maxVal
}
