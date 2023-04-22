// https://kopricky.github.io/code/SegmentTrees/orthogonal_range_report.html

// 正交矩形报告(Orthogonal Range Report)
// 统计矩形内的点, 每次查询 O(logN^2 + 点数)

package main

import (
	"fmt"
	"sort"
	"time"
)

func main() {
	points := [][2]int{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8}, {9, 9}, {10, 10}}
	orr := NewOrthogonalRangeReport(points)
	fmt.Println(orr.Query(1, 3, 1, 3))

	n := 100000
	points = make([][2]int, n)
	for i := 0; i < n; i++ {
		points[i] = [2]int{i, i}
	}
	time1 := time.Now()
	orr = NewOrthogonalRangeReport(points)
	for i := 0; i < 100; i++ {
		orr.Query(0, n, 0, n)
	}
	fmt.Println(time.Since(time1))

}

// 正交范围报告(Orthogonal Range Report).
type OrthogonalRangeReport struct {
	n  int
	xs []int
	ys [][][2]int
}

func NewOrthogonalRangeReport(points [][2]int) *OrthogonalRangeReport {
	sz := len(points)
	n := 1
	for n < sz {
		n <<= 1
	}
	sorted := make([][3]int, sz)
	for i := 0; i < sz; i++ {
		sorted[i] = [3]int{points[i][0], points[i][1], i}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i][0] < sorted[j][0]
	})

	xs := make([]int, sz)
	ys := make([][][2]int, 2*n-1)
	for i := 0; i < sz; i++ {
		xs[i] = sorted[i][0]
		ys[i+n-1] = [][2]int{{sorted[i][1], sorted[i][2]}}
	}

	for i := n - 2; i >= 0; i-- {
		nums1 := ys[2*i+1]
		nums2 := ys[2*i+2]
		ys[i] = make([][2]int, len(nums1)+len(nums2))
		p := 0
		p1 := 0
		p2 := 0
		for p1 < len(nums1) && p2 < len(nums2) {
			if nums1[p1][0] < nums2[p2][0] {
				ys[i][p] = nums1[p1]
				p++
				p1++
			} else {
				ys[i][p] = nums2[p2]
				p++
				p2++
			}
		}
		for p1 < len(nums1) {
			ys[i][p] = nums1[p1]
			p++
			p1++
		}
		for p2 < len(nums2) {
			ys[i][p] = nums2[p2]
			p++
			p2++
		}
	}

	return &OrthogonalRangeReport{n: n, xs: xs, ys: ys}
}

func (orr *OrthogonalRangeReport) Query(x1, x2, y1, y2 int) []int {
	lxid := sort.SearchInts(orr.xs, x1)
	rxid := sort.SearchInts(orr.xs, x2)
	if lxid >= rxid {
		return []int{}
	}
	report := make([]int, 0)
	orr._query(lxid, rxid, y1, y2, &report, 0, 0, orr.n)
	return report
}

func (orr *OrthogonalRangeReport) _query(lxid, rxid, ly, ry int, report *[]int, k, l, r int) {
	if r <= lxid || rxid <= l {
		return
	}
	if lxid <= l && r <= rxid {
		ysk := orr.ys[k]
		st := sort.Search(len(ysk), func(i int) bool { return ysk[i][0] >= ly })
		ed := sort.Search(len(ysk), func(i int) bool { return ysk[i][0] > ry })
		for i := st; i < ed; i++ {
			*report = append(*report, ysk[i][1])
		}
	} else {
		orr._query(lxid, rxid, ly, ry, report, 2*k+1, l, (l+r)>>1)
		orr._query(lxid, rxid, ly, ry, report, 2*k+2, (l+r)>>1, r)
	}
}
