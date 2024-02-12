// https://maspypy.github.io/library/geo/count_points_in_triangles.hpp

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
)

// https://judge.yosupo.jp/problem/count_points_in_triangle
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	pointsA := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pointsA[i][0], &pointsA[i][1])
	}
	var m int
	fmt.Fscan(in, &m)
	pointsB := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &pointsB[i][0], &pointsB[i][1])
	}

	cp := NewCountPointsInTriangles(pointsA, pointsB)
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		fmt.Fprintln(out, cp.Query(a, b, c))
	}
}

const LIMIT int = 1e9 + 10 // LIMIT*LIMIT < 1<<63

// 三角形内的点计数.
// 给定点群A和点群B，每次询问给出三个点i,j,k，求三角形AiAjAk内部的B中点的个数。
// O(nmlogm)预处理,O(1)查询.
type CountPointsInTriangles struct {
	A, B      [][2]int
	newIndex  []int   // 极角排序
	samePoint []int   // A[i]与B中点重合的个数
	seg       [][]int // 线段A[i]A[j]上的B[k]个数
	tri       [][]int // OA[i]A[j] 内的B[k]个数
}

// 指定点集构建.
func NewCountPointsInTriangles(points1, points2 [][2]int) *CountPointsInTriangles {
	res := &CountPointsInTriangles{}
	points1 = append(points1[:0:0], points1...)
	points2 = append(points2[:0:0], points2...)
	res.A = points1
	res.B = points2
	res.build()
	return res
}

// 查询三角形 A[i]A[j]A[k] 内B中点的个数.
func (cp *CountPointsInTriangles) Query(i, j, k int) int {
	i, j, k = cp.newIndex[i], cp.newIndex[j], cp.newIndex[k]
	if i > j {
		i, j = j, i
	}
	if j > k {
		j, k = k, j
	}
	if i > j {
		i, j = j, i
	}
	pi, pj, pk := cp.A[i], cp.A[j], cp.A[k]
	d := det(pj[0]-pi[0], pj[1]-pi[1], pk[0]-pi[0], pk[1]-pi[1])
	if d == 0 {
		return 0
	}
	if d > 0 {
		return cp.tri[i][j] + cp.tri[j][k] - cp.tri[i][k] - cp.seg[i][k]
	}
	x := cp.tri[i][k] - cp.tri[i][j] - cp.tri[j][k]
	return x - cp.seg[i][j] - cp.seg[j][k] - cp.samePoint[j]
}

func (cp *CountPointsInTriangles) takeOrigin() [2]int {
	return [2]int{-LIMIT, rand.Intn(2*LIMIT) - LIMIT}
}

func (cp *CountPointsInTriangles) build() {
	O := cp.takeOrigin() // 平移到原点
	for i := range cp.A {
		cp.A[i][0] -= O[0]
		cp.A[i][1] -= O[1]
	}
	for i := range cp.B {
		cp.B[i][0] -= O[0]
		cp.B[i][1] -= O[1]
	}
	N, M := len(cp.A), len(cp.B)
	order := angleArgSort(cp.A)
	cp.A = reArrage(cp.A, order)
	cp.newIndex = make([]int, N)
	for i := 0; i < N; i++ {
		cp.newIndex[order[i]] = i
	}

	order = angleArgSort(cp.B)
	cp.B = reArrage(cp.B, order)
	cp.samePoint = make([]int, N)
	cp.seg = make([][]int, N)
	cp.tri = make([][]int, N)
	for i := range cp.seg {
		cp.seg[i] = make([]int, N)
		cp.tri[i] = make([]int, N)
	}

	// point
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			if cp.A[i] == cp.B[j] {
				cp.samePoint[i]++
			}
		}
	}

	m := 0
	for j := 0; j < N; j++ {
		// OA[i]A[j], B[k]
		for m < M && det(cp.A[j][0], cp.A[j][1], cp.B[m][0], cp.B[m][1]) < 0 {
			m++
		}

		C := make([][2]int, m)
		for k := 0; k < m; k++ {
			C[k] = [2]int{cp.B[k][0] - cp.A[j][0], cp.B[k][1] - cp.A[j][1]}
		}
		order := make([]int, m)
		for i := 0; i < m; i++ {
			order[i] = i
		}
		sort.Slice(order, func(i, j int) bool {
			return det(C[order[i]][0], C[order[i]][1], C[order[j]][0], C[order[j]][1]) > 0
		})
		C = reArrage(C, order)
		rank := make([]int, m)
		for k := 0; k < m; k++ {
			rank[order[k]] = k
		}
		bit := NewBitArray(m)

		k := m
		for i := j - 1; i >= 0; i-- {
			for k > 0 && det(cp.A[i][0], cp.A[i][1], cp.B[k-1][0], cp.B[k-1][1]) > 0 {
				k--
				bit.Add(rank[k], 1)
			}
			p := [2]int{cp.A[i][0] - cp.A[j][0], cp.A[i][1] - cp.A[j][1]}
			lb := binarySearch(func(n int) bool {
				if n == 0 {
					return true
				}
				return det(C[n-1][0], C[n-1][1], p[0], p[1]) > 0
			}, 0, m+1)
			ub := binarySearch(func(n int) bool {
				if n == 0 {
					return true
				}
				return det(C[n-1][0], C[n-1][1], p[0], p[1]) >= 0
			}, 0, m+1)
			cp.seg[i][j] += bit.QueryRange(lb, ub)
			cp.tri[i][j] += bit.QueryPrefix(lb)
		}
	}

}

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n     int
	total int
	data  []int
}

func NewBitArray(n int) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func NewBitArrayFrom(n int, f func(i int) int) *BITArray {
	total := 0
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &BITArray{n: n, total: total, data: data}
}

func (b *BITArray) Add(index int, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BITArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *BITArray) QueryAll() int {
	return b.total
}

func (b *BITArray) MaxRight(check func(index, preSum int) bool) int {
	i := 0
	s := 0
	k := 1
	for 2*k <= b.n {
		k *= 2
	}
	for k > 0 {
		if i+k-1 < b.n {
			t := s + b.data[i+k-1]
			if check(i+k, t) {
				i += k
				s = t
			}
		}
		k >>= 1
	}
	return i
}

// 0/1 树状数组查找第 k(0-based) 个1的位置.
// UpperBound.
func (b *BITArray) Kth(k int) int {
	return b.MaxRight(func(index, preSum int) bool { return preSum <= k })
}

func (b *BITArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}

type V = int

// 极角排序，返回值为点的下标
func angleArgSort(points [][2]V) []int {
	lower, origin, upper := []int{}, []int{}, []int{}
	O := [2]V{0, 0}
	for i, p := range points {
		if p == O {
			origin = append(origin, i)
		} else if p[1] < 0 || (p[1] == 0 && p[0] > 0) {
			lower = append(lower, i)
		} else {
			upper = append(upper, i)
		}
	}

	sort.Slice(lower, func(i, j int) bool {
		oi, oj := lower[i], lower[j]
		pi, pj := points[oi], points[oj]
		return pi[0]*pj[1]-pi[1]*pj[0] > 0
	})
	sort.Slice(upper, func(i, j int) bool {
		oi, oj := upper[i], upper[j]
		pi, pj := points[oi], points[oj]
		return pi[0]*pj[1]-pi[1]*pj[0] > 0
	})

	res := lower
	res = append(res, origin...)
	res = append(res, upper...)
	return res
}

func det(a1, a2, b1, b2 int) int {
	return a1*b2 - a2*b1
}

func reArrage(nums [][2]int, order []int) [][2]int {
	res := make([][2]int, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}

func binarySearch(check func(int) bool, ok, ng int) int {
	for abs(ok-ng) > 1 {
		x := (ng + ok) / 2
		if check(x) {
			ok = x
		} else {
			ng = x
		}
	}
	return ok
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
