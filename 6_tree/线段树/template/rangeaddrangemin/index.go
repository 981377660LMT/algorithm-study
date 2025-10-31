// RangeAddRangeMin/RangeMinRangeAdd
// 更快的 RangeAddRangeMin 线段树(比一般的懒更新线段树快 40% 左右)

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	test()
}

const INF int = 1e18

type summax = struct {
	sum, max int
}

type summin = struct {
	sum, min int
}

type SegmentTreeRangeAddRangeMin struct {
	n    int
	lazy int
	seg  *SegmentTreeGeneric[summin]
}

func NewSegmentTreeRangeAddRangeMin(n int, f func(int) int) *SegmentTreeRangeAddRangeMin {
	res := &SegmentTreeRangeAddRangeMin{}
	res.build(n, f)
	return res
}

// [l, r)
func (st *SegmentTreeRangeAddRangeMin) Query(l, r int) int {
	if l < 0 {
		l = 0
	}
	if r > st.n {
		r = st.n
	}
	if l >= r {
		return INF
	}
	res := st.seg.Query(l, r).min
	for l += st.seg.size; l > 0; l >>= 1 {
		if l&1 == 1 {
			l--
			res += st.seg.seg[l].sum
		}
	}
	return res + st.lazy
}

// [l, r)
func (st *SegmentTreeRangeAddRangeMin) UpdateRange(l, r int, x int) {
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
func (st *SegmentTreeRangeAddRangeMin) UpdatePrefix(r int, x int) {
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
func (st *SegmentTreeRangeAddRangeMin) UpdateSuffix(l int, x int) {
	if l < 0 {
		l = 0
	}
	if l >= st.n {
		return
	}
	t := st.seg.Get(l).sum + x
	st.seg.Set(l, summin{t, t})
}

func (st *SegmentTreeRangeAddRangeMin) UpdateAll(x int) {
	st.lazy += x
}

func (st *SegmentTreeRangeAddRangeMin) Set(i int, x int) {
	if i < 0 || i >= st.n {
		return
	}
	cur := st.Query(i, i+1)
	st.UpdateRange(i, i+1, x-cur)
}

func (st *SegmentTreeRangeAddRangeMin) Update(i int, x int) {
	if i < 0 || i >= st.n {
		return
	}
	cur := st.Query(i, i+1)
	if cur > x {
		st.UpdateRange(i, i+1, x-cur)
	}
}
func (st *SegmentTreeRangeAddRangeMin) build(n int, f func(int) int) {
	st.lazy = 0
	st.n = n
	pre := 0
	st.seg = NewSegmentTreeGeneric(
		n,
		func(i int) summin {
			t := f(i) - pre
			pre += t
			return summin{t, t}
		},
		func() summin { return summin{min: 2 * INF} },
		func(a, b summin) summin {
			a.min = min(a.min, a.sum+b.min)
			a.sum += b.sum
			return a
		},
	)
}

// SegmentTreeGeneric
type SegmentTreeGeneric[E any] struct {
	n, size int
	seg     []E
	e       func() E
	op      func(a, b E) E
}

func NewSegmentTreeGeneric[E any](
	n int, f func(int) E,
	e func() E, op func(a, b E) E,
) *SegmentTreeGeneric[E] {
	res := &SegmentTreeGeneric[E]{e: e, op: op}
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
func (st *SegmentTreeGeneric[E]) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeGeneric[E]) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeGeneric[E]) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTreeGeneric[E]) Query(start, end int) E {
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
func (st *SegmentTreeGeneric[E]) QueryAll() E { return st.seg[1] }
func (st *SegmentTreeGeneric[E]) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

func test() {
	rand.Seed(time.Now().UnixNano())

	for testCase := 0; testCase < 10000; testCase++ {
		n := rand.Intn(100) + 1 // 数组长度1~100
		seg := NewSegmentTreeRangeAddRangeMin(n, func(i int) int { return 0 })
		bf := NewBruteForce(n)

		ops := rand.Intn(200) + 1 // 操作次数1~20次
		for op := 0; op < ops; op++ {
			action := rand.Intn(4)
			if action == 0 {
				// 区间加操作
				l := rand.Intn(n)
				r := l + rand.Intn(n-l+1)
				x := rand.Intn(20) - 10 // -10到9之间的随机数
				seg.UpdateRange(l, r, x)
				bf.UpdateRange(l, r, x)
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
				// set
				i := rand.Intn(n)
				x := rand.Intn(20) - 10 // -10到9之间的
				seg.Set(i, x)
				bf.Set(i, x)
			} else if action == 3 {
				// update
				i := rand.Intn(n)
				x := rand.Intn(20) - 10 // -10到9之间的随机数
				seg.Update(i, x)
				bf.Update(i, x)
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

func (bf *BruteForce) UpdateRange(l, r, x int) {
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
		return INF
	}

	minVal := INF
	currentSum := 0
	for i := 0; i < bf.n; i++ {
		currentSum += bf.diff[i]
		if i >= l && i < r {
			if currentSum < minVal {
				minVal = currentSum
			}
		}
	}
	return minVal
}

func (bf *BruteForce) Set(i, x int) {
	if i < 0 || i >= bf.n {
		return
	}
	current := bf.Query(i, i+1)
	diff := x - current
	if diff != 0 {
		bf.UpdateRange(i, i+1, diff)
	}
}

func (bf *BruteForce) Update(i int, x int) {
	if i < 0 || i >= bf.n {
		return
	}
	current := bf.Query(i, i+1)
	if current > x {
		bf.UpdateRange(i, i+1, x-current)
	}
}
