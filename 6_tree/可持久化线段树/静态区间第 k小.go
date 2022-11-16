// https://www.luogu.com.cn/problem/P3834

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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	findKth := useKthTree(nums)
	for i := 0; i < q; i++ {
		var left, right, k int
		fmt.Fscan(in, &left, &right, &k)
		fmt.Fprintln(out, findKth(left, right, k))
	}
}

// !给定数组 查询闭区间内第k小的数
//  1 <= left <= right <= n
//  k >= 1
func useKthTree(nums []int) func(left, right, k int) int {
	// !1.离散化
	// sorted(set(nums))
	set := make(map[int]struct{}, len(nums))
	for _, num := range nums {
		set[num] = struct{}{}
	}
	allNums := make([]int, 0, len(set))
	for num := range set {
		allNums = append(allNums, num)
	}
	sort.Ints(allNums)

	mp := make(map[int]int, len(allNums))
	for i, num := range allNums {
		mp[num] = i + 1
	}

	n := len(nums)
	roots := make([]*Node, n+1)
	roots[0] = Build(1, n)
	for i, num := range nums {
		roots[i+1] = roots[i].Update(mp[num], 1)
	}

	return func(left, right, k int) int {
		return allNums[roots[right].kth(roots[left-1], k)-1]
	}
}

type Data = int

func e() Data                          { return 0 }
func op(leftData, rightData Data) Data { return leftData + rightData }

type Node struct {
	left, right           int
	leftChild, rightChild *Node
	count                 Data
}

func newNode(left, right int) *Node { return &Node{left: left, right: right} }
func (o *Node) pushUp()             { o.count = op(o.leftChild.count, o.rightChild.count) } // op
func (o *Node) propagate(lazy int)  { o.count += lazy }                                     // mapping

// usage:
//  roots := make([]*Node, 1, maxVersion+1)
//  roots[0] = Build(1, len(nums), nums)
func Build(left, right int) *Node {
	if left == right {
		return newNode(left, right)
	}

	node := newNode(left, right)
	mid := (left + right) >> 1
	node.leftChild = Build(left, mid)
	node.rightChild = Build(mid+1, right)
	// node.pushUp()
	return node
}

// usage:
//  newRoot := roots[version].Update(index, value)
//  roots = append(roots, newRoot)
//  1 <= index <= len(nums)
// !注意：为了拷贝一份 Node，这里的接收器不是指针
func (o Node) Update(index int, count int) *Node {
	if o.left == o.right {
		o.propagate(count)
		return &o
	}

	mid := (o.left + o.right) >> 1
	if index <= mid {
		o.leftChild = o.leftChild.Update(index, count)
	} else {
		o.rightChild = o.rightChild.Update(index, count)
	}
	o.pushUp()
	return &o
}

// 查询区间 [left,right] 中第 k 小在整个数组上的名次（从 1 开始）.
//  1 <= left <= right <= len(nums)
// !注意返回的是（排序去重后的数组的）`下标`，不是元素值.
//
// usage:
//     roots := make([]*Node, 1, maxVersion+1)
//     roots[0] = Build(1, len(A))
//     !roots[i+1] = roots[i].update(kth[i], 1)  // kth[i] 为 A[i] 离散化后的值（从 1 开始） , 1 表示增加1个计数
//     roots[right].kth(roots[left-1], k)  // 类似前缀和 [left,right] 1<=left<=right<=n
func (o *Node) kth(oldNode *Node, k int) int {
	if o.left == o.right {
		return o.left
	}

	leftCount := o.leftChild.count - oldNode.leftChild.count // 值域左半边的数字个数
	if k <= leftCount {
		return o.leftChild.kth(oldNode.leftChild, k)
	}
	return o.rightChild.kth(oldNode.rightChild, k-leftCount)
}
