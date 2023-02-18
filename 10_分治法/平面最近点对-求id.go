// https://www.luogu.com.cn/problem/P1429
// nlogn解法
// 2 ≤ n ≤ 100,000
// 平面最近点对

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const INF float64 = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	points := make([]Point2, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &points[i].x, &points[i].y)
	}

	calDist := func(p1, p2 Point2) float64 {
		return math.Sqrt(float64((p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y)))
	}
	minDist, _, _ := minDistPair(points, calDist)
	// Each value is a real number with at most 6 digits after the decimal point.
	fmt.Fprintf(out, "%.6f", minDist)
}

type I = float64

type Point2 struct {
	x, y I
}

type Point2WithID struct {
	Point2
	id int
}

func minDistPair(points []Point2, calDist func(p1, p2 Point2) float64) (minDist float64, pid1, pid2 int) {
	if len(points) == 2 {
		return calDist(points[0], points[1]), 0, 1
	}

	point2WithID := make([]Point2WithID, len(points))
	for i, p := range points {
		point2WithID[i] = Point2WithID{p, i}
	}

	sort.Slice(point2WithID, func(i, j int) bool {
		if point2WithID[i].x == point2WithID[j].x {
			return point2WithID[i].y < point2WithID[j].y
		}
		return point2WithID[i].x < point2WithID[j].x
	})

	var merge func(left, right int) (minDist float64, pid1, pid2 int)
	merge = func(left, right int) (minDist float64, pid1, pid2 int) {
		if left == right {
			return INF, point2WithID[left].id, point2WithID[left].id
		} else if left+1 == right {
			return calDist(point2WithID[left].Point2, point2WithID[right].Point2), point2WithID[left].id, point2WithID[right].id
		}

		mid := (left + right) >> 1
		leftMin, leftPid1, leftPid2 := merge(left, mid)
		rightMin, rightPid1, rightPid2 := merge(mid+1, right)

		// 获得两边区间最小的距离的点对距离后，用于求中间衔接长方形内是否有更小值
		minDist, pid1, pid2 = leftMin, leftPid1, leftPid2
		if rightMin < minDist {
			minDist, pid1, pid2 = rightMin, rightPid1, rightPid2
		}

		// !将有可能成为更小值的点加入,只要x轴小于等于当前求出的最小值，就有可能
		cands := make([]Point2WithID, 0)
		for i := left; i <= right; i++ {
			if math.Abs(float64(point2WithID[i].x-point2WithID[mid].x)) <= minDist {
				cands = append(cands, point2WithID[i])
			}
		}

		// 在中间点集里面两两计算两点距离
		for i := 0; i < len(cands); i++ {
			for j := i + 1; j < len(cands); j++ {
				dist := calDist(cands[i].Point2, cands[j].Point2)
				if dist < minDist {
					minDist, pid1, pid2 = dist, cands[i].id, cands[j].id
				}
			}
		}

		return
	}

	return merge(0, len(points)-1)
}
