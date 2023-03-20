// 将全部环存在KD-Tree上，然后遍历玩具，找到距离当前玩具最近的点(环)，判断一下

package main

import (
	"math"
	"runtime/debug"
	"sort"
)

func init() {
	debug.SetGCPercent(-1)
}

func circleGame(toys [][]int, circles [][]int, r int) int {
	points := make([]Point, len(circles))
	for i, circle := range circles {
		points[i] = Point{circle[0], circle[1]}
	}

	kdTree := NewKDTree(points, func(p1, p2 Point) float64 {
		dx, dy := float64(p1[0]-p2[0]), float64(p1[1]-p2[1])
		return math.Sqrt(dx*dx + dy*dy)
	})

	res := 0
	for _, toy := range toys {
		minDist, _ := kdTree.FindNearest(Point{toy[0], toy[1]}, float64(r), true)
		if minDist+float64(toy[2]) <= float64(r) {
			res++
		}
	}

	return res
}

type Point []int

type PointWithID struct {
	Point
	id int
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
