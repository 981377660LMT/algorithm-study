// https://atcoder.jp/contests/abc287/tasks/abc287_g
// 有N种类型的卡，每种卡有10^100张。每种卡有分数和数量限制属性 (ai,bi)。
// 维护以下三种操作：
// 1 x y，将第x种卡的分数修改为y。
// 2 x y，将第x种卡的数量限制修改为y。
// 3 x，选x张卡，要求最大化分数和，且每种类型的卡的数量不得超过其数量限制。

// 解法:
// !动态开点权值线段树维护每种score的count+树上二分查询

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(-1)
}

// cards:每张卡片的(分数,数量限制)
func balanceUpdateQuery(cards [][2]int, queries [][]int) []int {
	res := []int{}
	tree := CreateSegmentTree(0, 1e9) // 0~1e9分数(权值线段树)
	for i := range cards {
		score, count := cards[i][0], cards[i][1]
		tree.Update(score, score, count)
	}

	for _, query := range queries {
		op := query[0]
		if op == 1 {
			index, newScore := query[1], query[2]
			preScore, preCount := cards[index][0], cards[index][1]
			tree.Update(preScore, preScore, -preCount)
			cards[index][0] = newScore
			tree.Update(newScore, newScore, preCount)
		} else if op == 2 {
			index, count := query[1], query[2]
			preScore, preCount := cards[index][0], cards[index][1]
			tree.Update(preScore, preScore, count-preCount)
			cards[index][1] = count
		} else {
			count := query[1]
			if tree.data.count < count {
				res = append(res, -1)
			} else {
				res = append(res, tree.Search(count))
			}
		}
	}

	return res
}

type E = struct{ left, sum, count int }
type Id = int // add count

func e(left, right int) E { return E{left: left} }
func id() Id              { return 0 }
func op(left, right E) E {
	return E{
		left: left.left,
		sum:  left.sum + right.sum, count: left.count + right.count,
	}
}
func mapping(parent Id, child E) E {
	return E{
		left: child.left,
		sum:  child.left*parent + child.sum, count: parent + child.count,
	}
}
func composition(parent, child Id) Id { return child } // 无需修改

// !树上二分查询最大分数和
func (o *Node) Search(remain int) int {
	if o.left == o.right {
		return o.left * remain
	}
	o.pushDown()
	if o.rightChild.data.count >= remain {
		return o.rightChild.Search(remain)
	}
	return o.leftChild.Search(remain-o.rightChild.data.count) + o.rightChild.data.sum
}

//
//
//
// 指定区间上下界建立线段树
func CreateSegmentTree(lower, upper int) *Node {
	root := newNode(lower, upper)
	return root
}

type Node struct {
	left, right           int
	leftChild, rightChild *Node
	data                  E
	lazy                  Id
}

func (o *Node) Update(left, right int, lazy Id) {
	if left <= o.left && o.right <= right {
		o.propagate(lazy)
		return
	}

	o.pushDown()
	mid := (o.left + o.right) >> 1
	if left <= mid {
		o.leftChild.Update(left, right, lazy)
	}
	if right > mid {
		o.rightChild.Update(left, right, lazy)
	}
	o.pushUp()
}

func newNode(left, right int) *Node {
	return &Node{left: left, right: right, lazy: id(), data: e(left, right)}
}

// op
func (o *Node) pushUp() {
	o.data = op(o.leftChild.data, o.rightChild.data)
}

func (o *Node) pushDown() {
	mid := (o.left + o.right) >> 1
	if o.leftChild == nil {
		o.leftChild = newNode(o.left, mid)
	}
	if o.rightChild == nil {
		o.rightChild = newNode(mid+1, o.right)
	}

	if o.lazy != id() {
		o.leftChild.propagate(o.lazy)
		o.rightChild.propagate(o.lazy)
		o.lazy = id()
	}
}

// mapping + composition
func (o *Node) propagate(lazy Id) {
	o.data = mapping(lazy, o.data)
	o.lazy = composition(lazy, o.lazy)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	cards := make([][2]int, n) // !卡片的(分数,拥有的枚数)
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		cards[i] = [2]int{a, b}
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([][]int, q) // !查询的(类型,参数1,参数2)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 3 {
			var count int
			fmt.Fscan(in, &count)
			queries[i] = []int{op, count}
		} else {
			var index, scoreOrCount int
			fmt.Fscan(in, &index, &scoreOrCount)
			index--
			queries[i] = []int{op, index, scoreOrCount}
		}
	}

	res := balanceUpdateQuery(cards, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}
