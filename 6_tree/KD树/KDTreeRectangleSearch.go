// https://beet-aizu.github.io/library/datastructure/kdtree.cpp

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
)

func init() {
	debug.SetGCPercent(-1)
}

// TLE
// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DSL_2_C
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	kd := NewKDTreeRectangleSearch()
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		kd.AddPoint(i, x, y)
	}
	kd.Build()

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		// sxi ≤ x ≤ txi and syi ≤ y ≤ tyi.
		var sx, tx, sy, ty int
		fmt.Fscan(in, &sx, &tx, &sy, &ty)
		res := kd.Find(sx, tx, sy, ty)
		sort.Ints(res)
		for _, v := range res {
			fmt.Fprintln(out, v)
		}
		fmt.Fprintln(out)
	}
}

type KDTreeRectangleSearch struct {
	points [][]int
	nodes  []Node
	nodeId int
}

type Node struct {
	pos         int
	left, right int
}

func NewKDTreeRectangleSearch() *KDTreeRectangleSearch {
	return &KDTreeRectangleSearch{}
}

func (k *KDTreeRectangleSearch) AddPoint(i, x, y int) {
	k.points = append(k.points, []int{i, x, y})
	k.nodes = append(k.nodes, Node{left: -1, right: -1, pos: -1})
}

func (k *KDTreeRectangleSearch) Build() int {
	return k.build(0, len(k.points), 0)
}

// 查找闭区间[x1, x2] * [y1, y2]内的点的索引
func (k *KDTreeRectangleSearch) Find(x1, x2, y1, y2 int) []int {
	res := []int{}
	k.find(0, x1, x2, y1, y2, 0, &res) // root
	return res
}

func (k *KDTreeRectangleSearch) find(root, x1, x2, y1, y2, depth int, res *[]int) {
	x := k.points[k.nodes[root].pos][1]
	y := k.points[k.nodes[root].pos][2]
	if x1 <= x && x <= x2 && y1 <= y && y <= y2 {
		*res = append(*res, k.points[k.nodes[root].pos][0])
	}

	if depth%2 == 0 {
		if k.nodes[root].left != -1 {
			if x1 <= x {
				k.find(k.nodes[root].left, x1, x2, y1, y2, depth+1, res)
			}
		}
		if k.nodes[root].right != -1 {
			if x <= x2 {
				k.find(k.nodes[root].right, x1, x2, y1, y2, depth+1, res)
			}
		}
	} else {
		if k.nodes[root].left != -1 {
			if y1 <= y {
				k.find(k.nodes[root].left, x1, x2, y1, y2, depth+1, res)
			}
		}
		if k.nodes[root].right != -1 {
			if y <= y2 {
				k.find(k.nodes[root].right, x1, x2, y1, y2, depth+1, res)
			}
		}
	}
}

func (k *KDTreeRectangleSearch) build(left, right, depth int) int {
	if left >= right {
		return -1
	}

	mid := (left + right) / 2
	t := k.nodeId
	k.nodeId++
	if depth%2 == 0 {
		sort.Slice(k.points[left:right], func(i, j int) bool {
			return k.points[left+i][1] < k.points[left+j][1]
		})
	} else {
		sort.Slice(k.points[left:right], func(i, j int) bool {
			return k.points[left+i][2] < k.points[left+j][2]
		})
	}

	k.nodes[t].pos = mid
	k.nodes[t].left = k.build(left, mid, depth+1)
	k.nodes[t].right = k.build(mid+1, right, depth+1)
	return t
}
