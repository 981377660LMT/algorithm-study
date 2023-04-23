//
// Vantage Point Tree (vp tree)
//
// Description
//   Vantage point tree is a metric tree.
//   Each tree node has a point, radius, and two childs.
//   The points of left descendants are contained in the ball B(p,r)
//   and the points of right descendants are exluded from the ball.
//
//   We can find k-nearest neighbors of a given point p efficiently
//   by pruning search.
//
//   The data structure is independently proposed by J. Uhlmann and
//   P. N. Yianilos.
//
// Complexity:
//   !Construction: O(n log n).
//   !Search: O(log n)
//
//   In my implementation, its construction is few times slower than kd tree
//   and its search is bit faster than kd tree.
//
// References
//   J. Uhlmann (1991):
//   Satisfying General Proximity/Similarity Queries with Metric Trees.
//   Information Processing Letters, vol. 40, no. 4, pp. 175--179.
//
//   Peter N. Yianilos (1993):
//   Data structures and algorithms for nearest neighbor search in general metric spaces.
//   in Proceedings of the 4th Annual ACM-SIAM Symposium on Discrete algorithms,
//   Society for Industrial and Applied Mathematics Philadelphia, PA, USA. pp. 311--321.

// 近邻搜索之制高点树（VP-Tree）
// https://www.cxyzjd.com/article/y459541195/102846739
// 以图搜图

// TODO 有问题 不要使用

package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

func main() {
	tree := NewVantagePointTree([][2]V{{3, 4}, {1, 2}, {7, 8}, {5, 6}, {9, 10}}, func(p1, p2 [2]V) V {
		return (p1[0]-p2[0])*(p1[0]-p2[0]) + (p1[1]-p2[1])*(p1[1]-p2[1])
	})
	res := tree.KNN([2]V{1, 1}, 5)
	fmt.Println(res)
}

type V = int
type VantagePointTree struct {
	root     *VNode
	aux      []pair
	calDist2 func(p1, p2 [2]V) V // 计算两点距离的平方
	maxPq    *Heap
}

func NewVantagePointTree(points [][2]V, calDist2 func(p1, p2 [2]V) V) *VantagePointTree {
	res := &VantagePointTree{calDist2: calDist2}
	aux := make([]pair, len(points))
	for i := 0; i < len(points); i++ {
		aux[i] = pair{second: points[i]}
	}
	res.aux = aux
	res.maxPq = NewHeap(func(a, b H) bool {
		if a.first == b.first {
			if a.second.point[0] == b.second.point[0] {
				return a.second.point[1] > b.second.point[1]
			}
			return a.second.point[0] > b.second.point[0]
		}
		return a.first > b.first
	}, nil)
	root := res._build(0, len(points))
	res.root = root
	return res
}

// 搜索k个最近邻点.
func (vp *VantagePointTree) KNN(point [2]V, k int) [][2]V {
	vp._knn(vp.root, point, k)
	var res [][2]V
	for vp.maxPq.Len() > 0 {
		res = append(res, vp.maxPq.Pop().second.point)
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

func (vp *VantagePointTree) _build(l, r int) *VNode {
	if l == r {
		return nil
	}
	rand := rand.Intn(r-l) + l
	vp.aux[l], vp.aux[rand] = vp.aux[rand], vp.aux[l]
	p := vp.aux[l].second
	l++
	if l == r {
		return &VNode{point: p}
	}
	for i := l; i < r; i++ {
		vp.aux[i].first = vp.calDist2(p, vp.aux[i].second)
	}
	m := (l + r) >> 1
	sort.Slice(vp.aux, func(i, j int) bool {
		if vp.aux[i].first == vp.aux[j].first {
			if vp.aux[i].second[0] == vp.aux[j].second[0] {
				return vp.aux[i].second[1] < vp.aux[j].second[1]
			}
			return vp.aux[i].second[0] < vp.aux[j].second[0]
		}
		return vp.aux[i].first < vp.aux[j].first
	})
	return &VNode{point: p, th: V(math.Sqrt(float64(vp.aux[m].first))), l: vp._build(l, m), r: vp._build(m, r)}
}

func (vp *VantagePointTree) _knn(t *VNode, p [2]V, k int) {
	if t == nil {
		return
	}
	d := abs(p[0]-t.point[0]) + abs(p[1]-t.point[1])
	if vp.maxPq.Len() < k {
		vp.maxPq.Push(H{first: d, second: t})
	} else if vp.maxPq.data[0].first > d {
		vp.maxPq.Pop()
		vp.maxPq.Push(H{first: d, second: t})
	}
	if t.l == nil && t.r == nil {
		return
	}
	if d < t.th {
		vp._knn(t.l, p, k)
		if t.th-d <= vp.maxPq.data[0].first {
			vp._knn(t.r, p, k)
		}
	} else {
		vp._knn(t.r, p, k)
		if d-t.th <= vp.maxPq.data[0].first {
			vp._knn(t.l, p, k)
		}
	}
}

func abs(a V) V {
	if a < 0 {
		return -a
	}
	return a
}

type VNode struct {
	point [2]V
	th    V
	l, r  *VNode
}

type pair struct {
	first  V
	second [2]V
}

type H = struct {
	first  V
	second *VNode
}

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}

		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
