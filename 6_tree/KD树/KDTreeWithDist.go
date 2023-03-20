// 静态KD树查询每个点的最近点(不包含自己)
// !注意查找最近点kdtree复杂度最坏会退化到O(n)
// KDTree查找最近点的原理，就是在搜索过程中先近后远，
// 然后搜索较远分支时，用已经搜索到的最近距离直接成片的剪枝
// 从上面传过来的已知最近点，或者看做裁剪范围
// https://baike.baidu.com/item/%E9%82%BB%E8%BF%91%E7%AE%97%E6%B3%95/1151153?fromtitle=knn&fromid=3479559&fr=aladdin

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime/debug"
	"sort"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		points[i] = Point{x, i}
	}

	kdtree := NewKDTree(points, func(p1, p2 Point) float64 {
		return math.Abs(float64(p1[0]-p2[0])) + math.Abs(float64(p1[1]-p2[1]))
	})

	for i := 0; i < n; i++ {
		minDist, _ := kdtree.FindNearest(points[i], float64(2*n), false)
		fmt.Fprint(out, minDist, " ")
	}
}

type Point []int

type PointWithID struct {
	Point
	id int
}

// 指定点集与距离计算函数，构造KDTree
func NewKDTree(points []Point, calDist func(p1, p2 Point) float64) *KDTree {
	if len(points) == 0 {
		return nil
	}

	res := &KDTree{
		dim:     len(points[0]),
		calDist: calDist,
	}

	pointsWithID := make([]PointWithID, len(points))
	for i, point := range points {
		pointsWithID[i] = PointWithID{point, i}
	}
	res.root = res.build(pointsWithID, 0)
	return res
}

// !查找距离point最近的点, 返回距离和id
//  upperBound: 从上面传过来的已知最近点，或者看做裁剪范围
//    如果不存在距离小于upperBound的点，返回upperBound和-1
//  allowOverlap: 是否统计距离为0的点(重合)
//    如果allowOverlap为false，则不会统计距离为0的点
func (kdtree *KDTree) FindNearest(point Point, upperBound float64, allowOverlap bool) (float64, int) {
	return kdtree.findNearest(kdtree.root, point, 0, upperBound, allowOverlap)
}

func (kdtree *KDTree) build(pointsWithID []PointWithID, depth int) *KDTreeNode {
	if len(pointsWithID) == 0 {
		return nil
	}

	axis := depth % kdtree.dim
	sort.Slice(pointsWithID, func(i, j int) bool {
		return pointsWithID[i].Point[axis] < pointsWithID[j].Point[axis]
	})
	mid := len(pointsWithID) / 2 // !中位数,可以用nth_element优化到O(nlogn)建树

	res := &KDTreeNode{pointWithId: pointsWithID[mid]}
	leftPoints := pointsWithID[:mid]
	rightPoints := pointsWithID[mid+1:]
	res.left = kdtree.build(leftPoints, depth+1)
	res.right = kdtree.build(rightPoints, depth+1)
	return res
}

func (kdtree *KDTree) findNearest(node *KDTreeNode, target Point, depth int, upperBound float64, allowOverlap bool) (float64, int) {
	dist := kdtree.calDist(node.pointWithId.Point, target)

	if !allowOverlap && dist == 0 { // !移除自己(重合时)
		dist = upperBound
	}

	resId := -1
	if dist < upperBound {
		upperBound = dist
		resId = node.pointWithId.id
	}

	axis := depth % kdtree.dim
	near, far := node.left, node.right
	if target[axis] > node.pointWithId.Point[axis] {
		near, far = far, near
	}

	if near != nil {
		distCand1, idCand1 := kdtree.findNearest(near, target, depth+1, upperBound, allowOverlap)
		if distCand1 < upperBound {
			upperBound = distCand1
			resId = idCand1
		}
	}

	// !复杂度靠这个剪枝
	if far != nil && upperBound > math.Abs(float64(node.pointWithId.Point[axis]-target[axis])) {
		distCand2, idCand2 := kdtree.findNearest(far, target, depth+1, upperBound, allowOverlap)
		if distCand2 < upperBound {
			upperBound = distCand2
			resId = idCand2
		}
	}

	return upperBound, resId
}

type KDTree struct {
	dim     int
	calDist func(p1, p2 Point) float64
	root    *KDTreeNode
}

type KDTreeNode struct {
	pointWithId PointWithID
	left        *KDTreeNode
	right       *KDTreeNode
}
