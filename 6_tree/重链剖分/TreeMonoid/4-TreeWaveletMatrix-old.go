package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// L番目の数字()
	AnUnusualKinginPakenKingdom()
}

func AnUnusualKinginPakenKingdom() {
	// https://atcoder.jp/contests/pakencamp-2022-day1/tasks/pakencamp_2022_day1_j
	// 路径上的中位数

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	tree := _NT(n)
	values := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		tree.AddEdge(a-1, b-1, 1)
		values[i] = c
	}
	tree.Build(0)
	twm := NewTreeWaveletMatrix(tree, values, false, -1, nil)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		x := twm.MedianPath(false, a, b, 0)
		y := twm.MedianPath(true, a, b, 0)
		fmt.Fprintln(out, (x+y)/2)
	}
}

func L番目の数字() {
	// https://atcoder.jp/contests/utpc2011/tasks/utpc2011_12
	// 路径上第k小的数

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	tree := _NT(n)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		tree.AddEdge(a-1, b-1, 1)
	}
	tree.Build(0)
	twm := NewTreeWaveletMatrix(tree, values, true, -1, nil)
	for i := 0; i < q; i++ {
		var a, b, k int
		fmt.Fscan(in, &a, &b, &k)
		a, b, k = a-1, b-1, k-1
		fmt.Fprintln(out, twm.KthPath(a, b, k, 0))
	}
}

const INF int = 1e18

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }
func inv(a E) E   { return -a }

type TreeWaveletMatrix struct {
	n        int
	tree     *_T
	wm       *WaveletMatrixSegments
	isVertex bool
}

// !树的路径查询, 维护的量需要满足群的性质，且必须要满足交换律.
//
//	data: 顶点的值, 或者边的值.(边的编号为两个定点中较深的那个点的编号).
//	isVertex: data是否为顶点的值以及查询的时候是否是顶点权值.
//	log: 如果要支持异或,则需要按照异或的值来决定值域.设为-1时表示不使用异或.
//	sumData: 如果要支持区间和,设为nil时表示不使用区间和.
func NewTreeWaveletMatrix(tree *_T, data []E, isVertex bool, log int, sumData []int) *TreeWaveletMatrix {
	res := &TreeWaveletMatrix{tree: tree, isVertex: isVertex, n: len(tree.Tree)}
	n := res.n
	A1 := make([]E, n)
	var S1 []E
	if len(sumData) > 0 {
		S1 = make([]E, n)
	}

	if isVertex {
		for v := 0; v < n; v++ {
			A1[tree.LID[v]] = data[v]
		}
		if len(sumData) == n {
			for v := 0; v < n; v++ {
				S1[tree.LID[v]] = E(sumData[v])
			}
		}
	} else {
		if len(sumData) != 0 {
			for e := 0; e < n-1; e++ {
				S1[tree.LID[tree.EidtoV(e)]] = sumData[e]
			}
		}
		for e := 0; e < n-1; e++ {
			A1[tree.LID[tree.EidtoV(e)]] = data[e]
		}
	}

	res.wm = NewWaveletMatrixSegments(A1, log, S1)
	return res
}

// s到t的路径上有多少个数在[a,b)之间.
func (tw *TreeWaveletMatrix) CountPath(s, t int, a, b, xorVal E) int {
	return tw.wm.CountRangeSegments(tw.getSegments(s, t), a, b, xorVal)
}

// 子树root中有多少个数在[a,b)之间.
func (tw *TreeWaveletMatrix) CountSubtree(root int, a, b, xorVal E) int {
	l, r := tw.tree.LID[root], tw.tree.RID[root]
	offset := 1
	if tw.isVertex {
		offset = 0
	}
	return tw.wm.CountRange(l+offset, r, a, b, xorVal)
}

// s到t的路径上第k小的数以及前k个数的和.
func (tw *TreeWaveletMatrix) KthValueAndSumPath(s, t, k int, xorVal E) (int, E) {
	return tw.wm.KthValueAndSumSegments(tw.getSegments(s, t), k, xorVal)
}

// 子树内第k小的数以及前k个数的和.
func (tw *TreeWaveletMatrix) KthValueAndSumSubtree(root, k int, xorVal E) (int, E) {
	l, r := tw.tree.LID[root], tw.tree.RID[root]
	offset := 1
	if tw.isVertex {
		offset = 0
	}
	return tw.wm.KthValueAndSum(l+offset, r, k, xorVal)
}

// s到t的路径上第k小的数.
func (tw *TreeWaveletMatrix) KthPath(s, t, k int, xorVal E) E {
	return tw.wm.KthSegments(tw.getSegments(s, t), k, xorVal)
}

// 子树内第k小的数.
func (tw *TreeWaveletMatrix) KthSubtree(root, k int, xorVal E) E {
	l, r := tw.tree.LID[root], tw.tree.RID[root]
	offset := 1
	if tw.isVertex {
		offset = 0
	}
	return tw.wm.Kth(l+offset, r, k, xorVal)
}

// s到t的路径上的中位数.
func (tw *TreeWaveletMatrix) MedianPath(upper bool, s, t int, xorVal E) E {
	return tw.wm.MedianSegments(upper, tw.getSegments(s, t), xorVal)
}

// 子树内的中位数.
func (tw *TreeWaveletMatrix) MedianSubtree(upper bool, root int, xorVal E) E {
	l, r := tw.tree.LID[root], tw.tree.RID[root]
	offset := 1
	if tw.isVertex {
		offset = 0
	}
	return tw.wm.Median(upper, l+offset, r, xorVal)
}

// s到t的路径上第k1小到第k2小的数的和.
func (tw *TreeWaveletMatrix) SumPath(s, t, k1, k2 int, xorVal E) E {
	return tw.wm.SumSegments(tw.getSegments(s, t), k1, k2, xorVal)
}

// 子树内第k1小到第k2小的数的和.
func (tw *TreeWaveletMatrix) SumSubtree(root, k1, k2 int, xorVal E) E {
	l, r := tw.tree.LID[root], tw.tree.RID[root]
	offset := 1
	if tw.isVertex {
		offset = 0
	}
	return tw.wm.Sum(l+offset, r, k1, k2, xorVal)
}

// s到t的路径上所有数的和.
func (tw *TreeWaveletMatrix) SumAllPath(s, t int) E {
	return tw.wm.SumAllSegments(tw.getSegments(s, t))
}

// 子树内所有数的和.
func (tw *TreeWaveletMatrix) SumAllSubtree(root int) E {
	l, r := tw.tree.LID[root], tw.tree.RID[root]
	offset := 1
	if tw.isVertex {
		offset = 0
	}
	return tw.wm.SumAll(l+offset, r)
}

func (tw *TreeWaveletMatrix) getSegments(s, t int) [][2]int {
	segments := tw.tree.GetPathDecomposition(s, t, tw.isVertex)
	for i := range segments {
		seg := &segments[i]
		if seg[0] > seg[1] {
			seg[0], seg[1] = seg[1], seg[0]
		}
		seg[1]++
	}
	return segments
}

type WaveletMatrixSegments struct {
	n, log int
	mid    []int
	bv     []*BitVector
	preSum [][]int
	unit   E
}

// log:如果要支持异或,则需要按照异或的值来决定值域
//
//	设为-1时表示不使用异或
//
// sumData:如果要支持区间和,则需要传入前缀和数组
//
//	设为nil时表示不使用区间和
func NewWaveletMatrixSegments(nums []E, log int, sumData []E) *WaveletMatrixSegments {
	res := &WaveletMatrixSegments{}
	res.build(nums, log, sumData)
	return res
}

// 返回区间 [left, right) 中 范围在 [a, b) 中的 元素的个数.
func (wm *WaveletMatrixSegments) CountRange(left, right, a, b, xor int) int {
	return wm.prefixCount(left, right, b, xor) - wm.prefixCount(left, right, a, xor)
}

func (wm *WaveletMatrixSegments) CountRangeSegments(segments [][2]int, a, b, xor int) int {
	res := 0
	for _, seg := range segments {
		res += wm.CountRange(seg[0], seg[1], a, b, xor)
	}
	return res
}

// 返回区间 [left, right) 中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果)
//
//	如果k < 0, 返回 (-1, 0); 如果k >= right-left, 返回 (-1, 区间 op 的结果)
func (wm *WaveletMatrixSegments) KthValueAndSum(left, right, k, xor int) (int, E) {
	if k < 0 {
		return -1, 0
	}
	if right-left <= k {
		return -1, wm.get(wm.log, left, right)
	}
	res, sum := 0, wm.unit
	count := 0
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		kf := f*(right-left-r0+l0) + (f^1)*(r0-l0)
		if count+kf > k {
			if f == 0 {
				left, right = l0, r0
			} else {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			}
		} else {
			var s E
			if f == 0 {
				s = wm.get(d, l0, r0)
			} else {
				s = wm.get(d, wm.mid[d]-l0+left, wm.mid[d]-r0+right)
			}
			count += kf
			res |= 1 << d
			sum = op(sum, s)
			if f == 0 {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			} else {
				left, right = l0, r0
			}
		}
	}
	sum = op(sum, wm.get(0, left, left+k-count))
	return res, sum
}

// 如果k < 0, 返回 (-1, 0); 如果k >= segments总长, 返回 (-1, 区间 op 的结果)
func (wm *WaveletMatrixSegments) KthValueAndSumSegments(segments [][2]int, k, xor int) (int, E) {
	if k < 0 {
		return -1, 0
	}
	totalLen := 0
	for _, seg := range segments {
		totalLen += seg[1] - seg[0]
	}
	if k >= totalLen {
		return -1, wm.SumAllSegments(segments)
	}
	count := 0
	sum := wm.unit
	res := 0
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		c := 0
		for _, seg := range segments {
			L, R := seg[0], seg[1]
			l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
			c += f*(R-L-r0+l0) + (f^1)*(r0-l0)
		}
		if count+c > k {
			for i := range segments {
				seg := &segments[i]
				L, R := seg[0], seg[1]
				l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
				if f == 0 {
					seg[0], seg[1] = l0, r0
				} else {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				}
			}
		} else {
			count += c
			res |= 1 << d
			for i := range segments {
				seg := &segments[i]
				L, R := seg[0], seg[1]
				l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
				var s E
				if f == 0 {
					s = wm.get(d, l0, r0)
				} else {
					s = wm.get(d, wm.mid[d]-l0+L, wm.mid[d]-r0+R)
				}
				sum = op(sum, s)
				if f == 0 {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				} else {
					seg[0], seg[1] = l0, r0
				}
			}
		}
	}

	for _, seg := range segments {
		L, R := seg[0], seg[1]
		t := min(R-L, k-count)
		sum = op(sum, wm.get(0, L, L+t))
		count += t
	}
	return res, sum
}

// 如果不存在,返回-1.
func (wm *WaveletMatrixSegments) Kth(left, right, k, xor int) E {
	if k < 0 || k >= right-left {
		return -1
	}
	count := 0
	res := 0
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		c := f*(right-left-r0+l0) + (f^1)*(r0-l0)
		if count+c > k {
			if f == 0 {
				left, right = l0, r0
			} else {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			}
		} else {
			count += c
			res |= 1 << d
			if f == 0 {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			} else {
				left, right = l0, r0
			}
		}
	}
	return res
}

// 如果不存在,返回-1.
func (wm *WaveletMatrixSegments) KthSegments(segments [][2]int, k, xor int) E {
	if k < 0 {
		return -1
	}
	totalLen := 0
	for _, seg := range segments {
		totalLen += seg[1] - seg[0]
	}
	if k >= totalLen {
		return -1
	}
	count := 0
	res := 0
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		c := 0
		for i := range segments {
			seg := &segments[i]
			L, R := seg[0], seg[1]
			l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
			c += f*(R-L-r0+l0) + (f^1)*(r0-l0)
		}
		if count+c > k {
			for i := range segments {
				seg := &segments[i]
				L, R := seg[0], seg[1]
				l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
				if f == 0 {
					seg[0], seg[1] = l0, r0
				} else {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				}
			}
		} else {
			count += c
			res |= 1 << d
			for i := range segments {
				seg := &segments[i]
				L, R := seg[0], seg[1]
				l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
				if f == 0 {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				} else {
					seg[0], seg[1] = l0, r0
				}
			}
		}
	}

	return res
}

// 区间中位数.
//
//	upper: true表示上中位数, false表示下中位数.
func (wm *WaveletMatrixSegments) Median(upper bool, left, right, xor int) E {
	n := right - left
	var k int
	if upper {
		k = n / 2
	} else {
		k = (n - 1) / 2
	}
	return wm.Kth(left, right, k, xor)
}

func (wm *WaveletMatrixSegments) MedianSegments(upper bool, segments [][2]int, xor int) E {
	n := 0
	for _, seg := range segments {
		n += seg[1] - seg[0]
	}
	var k int
	if upper {
		k = n / 2
	} else {
		k = (n - 1) / 2
	}
	return wm.KthSegments(segments, k, xor)
}

func (wm *WaveletMatrixSegments) Sum(left, right, k1, k2, xor int) E {
	return wm.prefixSum(left, right, k2, xor) - wm.prefixSum(left, right, k1, xor)
}

func (wm *WaveletMatrixSegments) SumSegments(segments [][2]int, k1, k2, xor int) E {
	return wm.prefixSumSegments(segments, k2, xor) - wm.prefixSumSegments(segments, k1, xor)
}

func (wm *WaveletMatrixSegments) SumAll(left, right int) E {
	return wm.get(wm.log, left, right)
}

func (wm *WaveletMatrixSegments) SumAllSegments(segments [][2]int) E {
	res := wm.unit
	for _, seg := range segments {
		res = op(res, wm.get(wm.log, seg[0], seg[1]))
	}
	return res
}

// 返回使得 check(count,prefixSum) 为 true 的最大 (count, prefixSum) 对.
//
//	!(即区间内小于 val 的数的个数count和 和 prefixSum 满足 check函数, 找到这样的最大的 (count, prefixSum).
//	eg: val = 5 => 即区间内值域在 [0,5) 中的数的和满足 check 函数.
func (wm *WaveletMatrixSegments) MaxRight(left, right, xor int, check func(count int, preSum E) bool) (int, E) {
	if tmp := wm.get(wm.log, left, right); check(right-left, tmp) {
		return right - left, tmp
	}
	count := 0
	res := wm.unit
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		c := f*(right-left-r0+l0) + (f^1)*(r0-l0)
		var s E
		if f == 0 {
			s = wm.get(d, l0, r0)
		} else {
			s = wm.get(d, left+wm.mid[d]-l0, right+wm.mid[d]-r0)
		}
		if tmp := op(res, s); check(count+c, tmp) {
			count += c
			res = tmp
			if f == 0 {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			} else {
				left, right = l0, r0
			}
		} else {
			if f == 0 {
				left, right = l0, r0
			} else {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			}
		}
	}
	k := wm.binarySearch(func(k int) bool {
		return check(count+k, op(res, wm.get(0, left, left+k)))
	}, 0, right-left)
	count += k
	res = op(res, wm.get(0, left, left+k))
	return count, res
}

func (w *WaveletMatrixSegments) build(nums []E, log int, sumData []E) {
	numsCopy := make([]E, len(nums))
	max_ := 1
	for i, v := range nums {
		numsCopy[i] = v
		if v > max_ {
			max_ = v
		}
	}
	if log == -1 {
		log = bits.Len(uint(max_))
	}

	makeSum := sumData != nil
	sumData = append(sumData[:0:0], sumData...)

	w.unit = e()
	n := len(numsCopy)
	mid := make([]int, log)
	bv := make([]*BitVector, log)
	for i := 0; i < log; i++ {
		bv[i] = NewBitVector(n)
	}

	var preSum [][]E
	if makeSum {
		preSum = make([][]E, log+1)
		for i := range preSum {
			preSum[i] = make([]E, n+1)
			for j := range preSum[i] {
				preSum[i][j] = w.unit
			}
		}
	}

	a0, a1 := make([]E, n), make([]E, n)
	s0, s1 := make([]E, n), make([]E, n)
	for d := log - 1; d >= -1; d-- {
		p0, p1 := 0, 0
		if makeSum {
			for i := 0; i < n; i++ {
				preSum[d+1][i+1] = op(preSum[d+1][i], sumData[i])
			}
		}
		if d == -1 {
			break
		}
		for i := 0; i < n; i++ {
			f := (numsCopy[i] >> d) & 1
			if f == 0 {
				if makeSum {
					s0[p0] = sumData[i]
				}
				a0[p0] = numsCopy[i]
				p0++
			} else {
				if makeSum {
					s1[p1] = sumData[i]
				}
				bv[d].Set(i)
				a1[p1] = numsCopy[i]
				p1++
			}
		}
		mid[d] = p0
		bv[d].Build()
		numsCopy, a0 = a0, numsCopy
		sumData, s0 = s0, sumData
		for i := 0; i < p1; i++ {
			numsCopy[p0+i] = a1[i]
			if makeSum {
				sumData[p0+i] = s1[i]
			}
		}
	}

	w.n, w.log = n, log
	w.mid, w.bv, w.preSum = mid, bv, preSum
}

func (wm *WaveletMatrixSegments) binarySearch(f func(E) bool, ok, ng int) int {
	for abs(ok-ng) > 1 {
		x := (ok + ng) >> 1
		if f(x) {
			ok = x
		} else {
			ng = x
		}
	}
	return ok
}

// 返回区间 [left, right) 中 范围在 [0, x) 中的 元素的个数.
func (wm *WaveletMatrixSegments) prefixCount(left, right, x, xor int) int {
	if x == 0 {
		return 0
	}
	if x >= 1<<wm.log {
		return right - left
	}
	count := 0
	for d := wm.log - 1; d >= 0; d-- {
		add := (x >> d) & 1
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		kf := f*(right-left-r0+l0) + (f^1)*(r0-l0)
		if add == 1 {
			count += kf
			if f == 1 {
				left, right = l0, r0
			} else {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			}
		} else if f == 0 {
			left, right = l0, r0
		} else {
			left += wm.mid[d] - l0
			right += wm.mid[d] - r0
		}
	}
	return count
}

func (wm *WaveletMatrixSegments) prefixSum(left, right, k, xor int) E {
	_, res := wm.KthValueAndSum(left, right, k, xor)
	return res
}

func (wm *WaveletMatrixSegments) prefixSumSegments(segments [][2]int, k, xor int) E {
	_, res := wm.KthValueAndSumSegments(segments, k, xor)
	return res
}

func (wm *WaveletMatrixSegments) get(d, l, r int) E {
	return op(inv(wm.preSum[d][l]), wm.preSum[d][r])
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

type BitVector struct {
	data [][2]int
}

func NewBitVector(n int) *BitVector {
	return &BitVector{data: make([][2]int, (n+63)>>5)}
}

func (bv *BitVector) Set(i int) {
	bv.data[i>>5][0] |= 1 << (i & 31)
}

func (bv *BitVector) Build() {
	for i := 0; i < len(bv.data)-1; i++ {
		bv.data[i+1][1] = bv.data[i][1] + bits.OnesCount(uint(bv.data[i][0]))
	}
}

func (bv *BitVector) Rank(k int, f int) int {
	a, b := bv.data[k>>5][0], bv.data[k>>5][1]
	ret := b + bits.OnesCount(uint(a&((1<<(k&31))-1)))
	if f == 1 {
		return ret
	}
	return k - ret
}

type _T struct {
	Tree                 [][][2]int // (next, weight)
	Edges                [][3]int   // (u, v, w)
	Depth, DepthWeighted []int
	Parent               []int
	LID, RID             []int // 欧拉序[in,out)
	IdToNode             []int
	top, heavySon        []int
	timer                int
}

func _NT(n int) *_T {
	tree := make([][][2]int, n)
	lid := make([]int, n)
	rid := make([]int, n)
	IdToNode := make([]int, n)
	top := make([]int, n)   // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n) // 深度
	depthWeighted := make([]int, n)
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	edges := make([][3]int, 0, n-1)
	for i := range parent {
		parent[i] = -1
	}

	return &_T{
		Tree:          tree,
		Depth:         depth,
		DepthWeighted: depthWeighted,
		Parent:        parent,
		LID:           lid,
		RID:           rid,
		IdToNode:      IdToNode,
		top:           top,
		heavySon:      heavySon,
		Edges:         edges,
	}
}

// 添加无向边 u-v, 边权为w.
func (tree *_T) AddEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.Tree[v] = append(tree.Tree[v], [2]int{u, w})
	tree.Edges = append(tree.Edges, [3]int{u, v, w})
}

// 添加有向边 u->v, 边权为w.
func (tree *_T) AddDirectedEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.Edges = append(tree.Edges, [3]int{u, v, w})
}

// root:0-based
//
//	当root设为-1时，会从0开始遍历未访问过的连通分量
func (tree *_T) Build(root int) {
	if root != -1 {
		tree.build(root, -1, 0, 0)
		tree.markTop(root, root)
	} else {
		for i := 0; i < len(tree.Tree); i++ {
			if tree.Parent[i] == -1 {
				tree.build(i, -1, 0, 0)
				tree.markTop(i, i)
			}
		}
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *_T) Id(root int) (int, int) {
	return tree.LID[root], tree.RID[root]
}

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *_T) Eid(u, v int) int {
	if tree.LID[u] > tree.LID[v] {
		return tree.LID[u]
	}
	return tree.LID[v]
}

// 较深的那个点作为边的编号.
func (tree *_T) EidtoV(eid int) int {
	e := tree.Edges[eid]
	u, v := e[0], e[1]
	if tree.Parent[u] == v {
		return u
	}
	return v
}

func (tree *_T) LCA(u, v int) int {
	for {
		if tree.LID[u] > tree.LID[v] {
			u, v = v, u
		}
		if tree.top[u] == tree.top[v] {
			return u
		}
		v = tree.Parent[tree.top[v]]
	}
}

func (tree *_T) RootedLCA(u, v int, root int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
}

func (tree *_T) Dist(u, v int, weighted bool) int {
	if weighted {
		return tree.DepthWeighted[u] + tree.DepthWeighted[v] - 2*tree.DepthWeighted[tree.LCA(u, v)]
	}
	return tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)]
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
func (tree *_T) KthAncestor(root, k int) int {
	if k > tree.Depth[root] {
		return -1
	}
	for {
		u := tree.top[root]
		if tree.LID[root]-k >= tree.LID[u] {
			return tree.IdToNode[tree.LID[root]-k]
		}
		k -= tree.LID[root] - tree.LID[u] + 1
		root = tree.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (tree *_T) Jump(from, to, step int) int {
	if step == 1 {
		if from == to {
			return -1
		}
		if tree.IsInSubtree(to, from) {
			return tree.KthAncestor(to, tree.Depth[to]-tree.Depth[from]-1)
		}
		return tree.Parent[from]
	}
	c := tree.LCA(from, to)
	dac := tree.Depth[from] - tree.Depth[c]
	dbc := tree.Depth[to] - tree.Depth[c]
	if step > dac+dbc {
		return -1
	}
	if step <= dac {
		return tree.KthAncestor(from, step)
	}
	return tree.KthAncestor(to, dac+dbc-step)
}

func (tree *_T) CollectChild(root int) []int {
	res := []int{}
	for _, e := range tree.Tree[root] {
		next := e[0]
		if next != tree.Parent[root] {
			res = append(res, next)
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *_T) GetPathDecomposition(u, v int, vertex bool) [][2]int {
	up, down := [][2]int{}, [][2]int{}
	for {
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			down = append(down, [2]int{tree.LID[tree.top[v]], tree.LID[v]})
			v = tree.Parent[tree.top[v]]
		} else {
			up = append(up, [2]int{tree.LID[u], tree.LID[tree.top[u]]})
			u = tree.Parent[tree.top[u]]
		}
	}
	edgeInt := 1
	if vertex {
		edgeInt = 0
	}
	if tree.LID[u] < tree.LID[v] {
		down = append(down, [2]int{tree.LID[u] + edgeInt, tree.LID[v]})
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		up = append(up, [2]int{tree.LID[u], tree.LID[v] + edgeInt})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

func (tree *_T) GetPath(u, v int) []int {
	res := []int{}
	composition := tree.GetPathDecomposition(u, v, true)
	for _, e := range composition {
		a, b := e[0], e[1]
		if a <= b {
			for i := a; i <= b; i++ {
				res = append(res, tree.IdToNode[i])
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, tree.IdToNode[i])
			}
		}
	}
	return res
}

// 以root为根时,结点v的子树大小.
func (tree *_T) SubtreeSize(v, root int) int {
	if root == -1 {
		return tree.RID[v] - tree.LID[v]
	}
	if v == root {
		return len(tree.Tree)
	}
	x := tree.Jump(v, root, 1)
	if tree.IsInSubtree(v, x) {
		return tree.RID[v] - tree.LID[v]
	}
	return len(tree.Tree) - tree.RID[x] + tree.LID[x]
}

// child 是否在 root 的子树中 (child和root不能相等)
func (tree *_T) IsInSubtree(child, root int) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

func (tree *_T) ELID(u int) int {
	return 2*tree.LID[u] - tree.Depth[u]
}

func (tree *_T) ERID(u int) int {
	return 2*tree.RID[u] - tree.Depth[u] - 1
}

func (tree *_T) build(cur, pre, dep, dist int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, e := range tree.Tree[cur] {
		next, weight := e[0], e[1]
		if next != pre {
			nextSize := tree.build(next, cur, dep+1, dist+weight)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	tree.Depth[cur] = dep
	tree.DepthWeighted[cur] = dist
	tree.heavySon[cur] = heavySon
	tree.Parent[cur] = pre
	return subSize
}

func (tree *_T) markTop(cur, top int) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.IdToNode[tree.timer] = cur
	tree.timer++
	if tree.heavySon[cur] != -1 {
		tree.markTop(tree.heavySon[cur], top)
		for _, e := range tree.Tree[cur] {
			next := e[0]
			if next != tree.heavySon[cur] && next != tree.Parent[cur] {
				tree.markTop(next, next)
			}
		}
	}
	tree.RID[cur] = tree.timer
}
