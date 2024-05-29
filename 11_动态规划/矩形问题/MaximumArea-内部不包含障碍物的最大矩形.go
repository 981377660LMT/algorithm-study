// MaxAreaRectangleWithoutObstacle 最大矩形

package main

import (
	"fmt"
	"sort"
)

func main() {
	left, right, bottom, top := 0, 10, 0, 10
	xs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	ys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	bounding, area := MaximumArea(left, right, bottom, top, xs, ys)
	fmt.Println(bounding, area)
}

// !内部不包含障碍物的最大矩形
// 在一个 [l,r,b,t] 确定的大厅中，有 n 根柱子。要求找到大厅内部面积最大的一个矩形，内部不包含柱子。
// 返回矩形的(l, r, b, t)坐标和最大面积.
func MaximumArea(
	left, right, bottom, top int,
	xs, ys []int,
) (bounding [4]int, bestArea int) {
	xs, ys = append(xs[:0:0], xs...), append(ys[:0:0], ys...)
	n := int32(len(xs))
	for i := int32(0); i < n; i++ {
		xs[i] = max(xs[i], left)
		xs[i] = min(xs[i], right)
		ys[i] = max(ys[i], bottom)
		ys[i] = min(ys[i], top)
	}
	pts := make([][2]int, n+2)
	allX := make([]int, 0, n+2)
	allX = append(allX, left, right)
	allX = append(allX, xs...)
	unique(&allX)
	for i := int32(0); i < n; i++ {
		pts[i][0] = sort.SearchInts(allX, xs[i])
		pts[i][1] = ys[i]
	}
	pts[n][0], pts[n][1] = 0, bottom
	pts[n+1][0], pts[n+1][1] = 0, top
	n += 2
	sort.Slice(pts, func(i, j int) bool { return pts[i][1] < pts[j][1] })
	m := int32(len(allX))
	L, R, cnt := make([]int32, m), make([]int32, m), make([]int32, m)

	for i := int32(0); i < n; i++ {
		to := i
		for to+1 < n && pts[to+1][1] == pts[i][1] {
			to++
		}
		i = to
		for j := i + 1; j < n; j++ {
			cnt[pts[j][0]]++
		}
		ll, rr := int(left), int(left)
		{
			left := int32(0)
			for j := int32(1); j < m; j++ {
				if cnt[j] > 0 || j == m-1 {
					cand := allX[j] - allX[left]
					if rr-ll < cand {
						rr = allX[j]
						ll = allX[left]
					}
					L[j] = left
					R[left] = j
					left = j
				}
			}
		}
		for j := n - 1; j > i; j-- {
			x := pts[j][0]
			cnt[x]--
			if cnt[x] == 0 {
				L[R[x]] = L[x]
				R[L[x]] = R[x]
				cand := allX[R[x]] - allX[L[x]]
				if rr-ll < cand {
					ll = allX[L[x]]
					rr = allX[R[x]]
				}
			}
			cand := (rr - ll) * (pts[j][1] - pts[i][1])
			if cand > bestArea {
				bestArea = cand
				bounding[0], bounding[1], bounding[2], bounding[3] = ll, rr, pts[i][1], pts[j][1]
			}
		}
	}

	return
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

func unique(arr *[]int) {
	if len(*arr) <= 1 {
		return
	}
	sort.Ints(*arr)
	pos := 1
	for i := 1; i < len(*arr); i++ {
		if (*arr)[i] != (*arr)[i-1] {
			(*arr)[pos] = (*arr)[i]
			pos++
		}
	}
	*arr = (*arr)[:pos]
}
