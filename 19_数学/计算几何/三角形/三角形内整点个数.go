// https://maspypy.github.io/library/geo/count_points_in_triangles.hpp
// 查询三角形 A[i]A[j]A[k] 内B中点的个数.
// TODO: 有问题
package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func main() {
	ponts1 := [][2]int{{0, 13}, {13, 0}, {0, 0}}
	ponts2 := [][2]int{{1, 2}, {3, 4}, {3, 2}, {1, 1}}
	CP := NewCountPointsInTriangles(ponts1, ponts2)
	fmt.Println(CP.Query(0, 1, 2), CP._rk)
}

const INF int = 1e18

// 三角形内的点计数.
//  O(n^2*m) 预处理, O(1) 查询.
type CountPointsInTriangles struct {
	_A, _B    [][2]int
	_I, _rk   []int
	samePoint []int
	seg       [][]int // 线段A[i]A[j]上的B[k]个数
	tri       [][]int // OA[i]A[j] 内的B[k]个数
}

// 指定点集构建.
func NewCountPointsInTriangles(points1, points2 [][2]int) *CountPointsInTriangles {
	res := &CountPointsInTriangles{}
	points1 = append(points1[:0:0], points1...)
	points2 = append(points2[:0:0], points2...)
	res._A = points1
	res._B = points2
	res.build()
	return res
}

// 查询三角形 A[i]A[j]A[k] 内B中点的个数.
func (cp *CountPointsInTriangles) Query(i, j, k int) int {
	i, j, k = cp._rk[i], cp._rk[j], cp._rk[k]
	if i > j {
		i, j = j, i
	}
	if j > k {
		j, k = k, j
	}
	if i > j {
		i, j = j, i
	}
	pi, pj, pk := cp._A[i], cp._A[j], cp._A[k]
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
	n, m := len(cp._A), len(cp._B)
	for {
		O := [2]int{-INF, rand.Intn(2*INF) - INF}
		ok := true

		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if cp._A[i] == cp._A[j] {
					continue
				}
				pi, pj := cp._A[i], cp._A[j]
				if det(pi[0]-O[0], pi[1]-O[1], pj[0]-O[0], pj[1]-O[1]) == 0 {
					ok = false
				}
			}
		}

		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if cp._A[i] == cp._B[j] {
					continue
				}
				pi, pj := cp._A[i], cp._B[j]
				if det(pi[0]-O[0], pi[1]-O[1], pj[0]-O[0], pj[1]-O[1]) == 0 {
					ok = false
				}
			}
		}

		if ok {
			return O
		}
	}
}

func (cp *CountPointsInTriangles) build() {
	O := cp.takeOrigin() // 平移到原点
	for i := range cp._A {
		cp._A[i][0] -= O[0]
		cp._A[i][1] -= O[1]
	}
	for i := range cp._B {
		cp._B[i][0] -= O[0]
		cp._B[i][1] -= O[1]
	}
	n, m := len(cp._A), len(cp._B)
	cp._I = make([]int, n)
	cp._rk = make([]int, n)
	for i := range cp._I {
		cp._I[i] = i
	}
	sort.Slice(cp._I, func(i, j int) bool {
		return det(cp._A[cp._I[i]][0], cp._A[cp._I[i]][1], cp._A[cp._I[j]][0], cp._A[cp._I[j]][1]) > 0
	})
	for i := range cp._I {
		cp._rk[cp._I[i]] = i
	}
	cp._A = reArrage(cp._A, cp._I)
	cp.samePoint = make([]int, n)
	cp.seg = make([][]int, n)
	cp.tri = make([][]int, n)
	for i := range cp.seg {
		cp.seg[i] = make([]int, n)
		cp.tri[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if cp._A[i] == cp._B[j] {
				cp.samePoint[i]++
			}
		}
	}

	for i := 0; i < n; i++ {
		pi := cp._A[i]
		for j := i + 1; j < n; j++ {
			pj := cp._A[j]
			for k := 0; k < m; k++ {
				pk := cp._B[k]
				if det(pi[0], pi[1], pk[0], pk[1]) <= 0 {
					continue
				}
				if det(pj[0], pj[1], pk[0], pk[1]) >= 0 {
					continue
				}

				d := det(pk[0]-pi[0], pk[1]-pi[1], pj[0]-pi[0], pj[1]-pi[1])
				if d == 0 {
					cp.seg[i][j]++
				}
				if d < 0 {
					cp.tri[i][j]++
				}
			}
		}
	}
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
