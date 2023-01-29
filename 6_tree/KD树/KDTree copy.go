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
	"sort"
)

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
		minDist := kdtree.FindNearest(points[i], float64(2*n))
		fmt.Fprint(out, minDist, " ")
	}
}

type Point []int

type KDTree struct {
	dim     int
	calDist func(p1, p2 Point) float64
	root    *KDTreeNode
}

type KDTreeNode struct {
	point Point
	left  *KDTreeNode
	right *KDTreeNode
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

	pointsCopy := make([]Point, len(points))
	copy(pointsCopy, points)
	res.root = res.build(pointsCopy, 0)
	return res
}

// 查找距离point最近的点(不包含与point重合的点), 返回距离
//  upperBound: 从上面传过来的已知最近点，或者看做裁剪范围
func (kdtree *KDTree) FindNearest(point Point, upperBound float64) float64 {
	return kdtree.findNearest(kdtree.root, point, 0, upperBound)
}

func (kdtree *KDTree) build(points []Point, depth int) *KDTreeNode {
	if len(points) == 0 {
		return nil
	}

	axis := depth % kdtree.dim
	sort.Slice(points, func(i, j int) bool {
		return points[i][axis] < points[j][axis]
	})
	mid := len(points) / 2 // !中位数,可以用nth_element优化到O(nlogn)建树

	res := &KDTreeNode{point: points[mid]}
	leftPoints := points[:mid]
	rightPoints := points[mid+1:]
	res.left = kdtree.build(leftPoints, depth+1)
	res.right = kdtree.build(rightPoints, depth+1)
	return res
}

func (kdtree *KDTree) findNearest(node *KDTreeNode, target Point, depth int, upperBound float64) float64 {
	if node == nil {
		return upperBound
	}

	dist := kdtree.calDist(node.point, target)
	if dist == 0 { // !移除自己(重合时)
		dist = upperBound
	}

	if dist < upperBound {
		upperBound = dist
	}

	axis := depth % kdtree.dim
	near, far := node.left, node.right
	if target[axis] > node.point[axis] {
		near, far = far, near
	}

	upperBound = kdtree.findNearest(near, target, depth+1, upperBound)
	if upperBound > math.Abs(float64(node.point[axis]-target[axis])) {
		// TODO
		// 1.这里是不是应该这样写?
		// 2.怎么支持返回pid?
		cand := kdtree.findNearest(far, target, depth+1, upperBound)
		if cand < upperBound {
			upperBound = cand
		}
	}

	return upperBound
}
