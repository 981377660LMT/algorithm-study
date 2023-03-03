// 动态开点权值线段树
// Update
// Query/QueryAll
// Set/Get
// Merge
// Split

// https://www.luogu.com.cn/problem/P5494
// 给定一个多重集合,编号为1
// 支持以下操作：
// 0 p x y : 将可重集 p 中大于等于 x 且小于等于 y 的值移动到一个新的可重集中
//           新可重集编号为从 2 开始的正整数，是上一次产生的新可重集的编号+1
// 1 p t : 将可重集 t 中的所有元素移动到可重集 p 中,且清空可重集 t (数据保证在此后的操作中不会出现可重集 t)
// 2 p k q : 在 可重集 p 中加入k个q
// 3 p x y : 输出可重集 p 中大于等于 x 且小于等于 y 的值的个数
// 4 p k : 输出可重集 p 中第 k 小的值,不存在输出-1

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}
	gits := make([]*LazyNode, 1, q+1)
	gits = append(gits, CreateSegmentTree(0, 2e5))
	for num, count := range nums {
		num++
		gits[1].Update(num, num, count)
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var p, x, y int
			fmt.Fscan(in, &p, &x, &y)
			_, root2 := gits[p].Split(x, y)
			gits = append(gits, root2)
		} else if op == 1 {
			var p, t int
			fmt.Fscan(in, &p, &t)
			gits[p] = gits[p].Merge(gits[t])
		} else if op == 2 {
			var p, k, q int
			fmt.Fscan(in, &p, &k, &q)
			gits[p].Update(q, q, k)
		} else if op == 3 {
			var p, x, y int
			fmt.Fscan(in, &p, &x, &y)
			fmt.Fprintln(out, gits[p].Query(x, y))
		} else if op == 4 {
			var p, k int
			fmt.Fscan(in, &p, &k)
			if k > gits[p].QueryAll() {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, gits[p].kth(k))
			}
		}

	}
}

// 指定区间上下界建立权值线段树.
func CreateSegmentTree(lower, upper int) *LazyNode {
	root := &LazyNode{left: lower, right: upper}
	return root
}

type LazyNode struct {
	left, right           int
	sum                   int
	lazy                  int
	leftChild, rightChild *LazyNode
}

func (LazyNode) op(a, b int) int {
	return a + b
}

func (o *LazyNode) propagate(add int) {
	o.lazy += add                         // % mod
	o.sum += (o.right - o.left + 1) * add // % mod
}

func (o *LazyNode) pushDown() {
	m := (o.left + o.right) >> 1
	if o.leftChild == nil {
		o.leftChild = &LazyNode{left: o.left, right: m}
	}
	if o.rightChild == nil {
		o.rightChild = &LazyNode{left: m + 1, right: o.right}
	}
	if add := o.lazy; add != 0 {
		o.leftChild.propagate(add)
		o.rightChild.propagate(add)
		o.lazy = 0
	}
}

func (o *LazyNode) pushUp() {
	o.sum = o.op(o.leftChild.QueryAll(), o.rightChild.QueryAll())
}

// Build from array. [1,len(nums))]
func (o *LazyNode) Build(nums []int) {
	o.build(nums, 1, len(nums))
}

func (o *LazyNode) build(nums []int, left, right int) {
	o.left, o.right = left, right
	if left == right {
		o.sum = int(nums[left-1])
		return
	}
	m := (left + right) >> 1
	o.leftChild = &LazyNode{}
	o.leftChild.build(nums, left, m)
	o.rightChild = &LazyNode{}
	o.rightChild.build(nums, m+1, right)
	o.pushUp()
}

// [left, right]
func (o *LazyNode) Update(left, right int, add int) {
	if left <= o.left && o.right <= right {
		o.propagate(add)
		return
	}
	o.pushDown()
	m := (o.left + o.right) >> 1
	if left <= m {
		o.leftChild.Update(left, right, add)
	}
	if m < right {
		o.rightChild.Update(left, right, add)
	}
	o.pushUp()
}

// [left, right]
func (o *LazyNode) Query(left, right int) int {
	if o == nil || left > o.right || right < o.left {
		return 0
	}
	if left <= o.left && o.right <= right {
		return o.sum
	}
	o.pushDown()
	return o.op(o.leftChild.Query(left, right), o.rightChild.Query(left, right))
}

func (o *LazyNode) Set(pos int, val int) {
	if o.left == o.right {
		o.sum = val
		return
	}
	o.pushDown()
	m := (o.left + o.right) >> 1
	if pos <= m {
		o.leftChild.Set(pos, val)
	} else {
		o.rightChild.Set(pos, val)
	}
	o.pushUp()
}

func (o *LazyNode) Get(pos int) int {
	if o.left == o.right {
		return o.sum
	}
	o.pushDown()
	m := (o.left + o.right) >> 1
	if pos <= m {
		return o.leftChild.Get(pos)
	}
	return o.rightChild.Get(pos)
}

func (o *LazyNode) QueryAll() int {
	if o != nil {
		return o.sum
	}
	return 0
}

// 线段树合并
func (o *LazyNode) Merge(b *LazyNode) *LazyNode {
	if o == nil {
		return b
	}
	if b == nil {
		return o
	}
	if o.left == o.right {
		o.sum += b.sum
		return o
	}
	o.leftChild = o.leftChild.Merge(b.leftChild)
	o.rightChild = o.rightChild.Merge(b.rightChild)
	o.pushUp()
	return o
}

// 线段树分裂
//  将区间 [l,r] 从原树分离到 other 上, this 为原树的剩余部分
func (o *LazyNode) Split(left, right int) (this, other *LazyNode) {
	this, other = o.split(nil, left, right)
	return
}

func (o *LazyNode) split(b *LazyNode, l, r int) (*LazyNode, *LazyNode) {
	if o == nil || l > o.right || r < o.left {
		return o, nil
	}
	if l <= o.left && o.right <= r {
		return nil, o
	}
	if b == nil {
		b = &LazyNode{left: o.left, right: o.right}
	}
	o.leftChild, b.leftChild = o.leftChild.split(b.leftChild, l, r)
	o.rightChild, b.rightChild = o.rightChild.split(b.rightChild, l, r)
	o.pushUp()
	b.pushUp()
	return o, b
}

// 权值线段树求第 k 小
// 调用前需保证 1 <= k <= root.QueryAll()
func (o *LazyNode) kth(k int) int {
	if o.left == o.right {
		return o.left
	}
	if lc := o.leftChild.QueryAll(); k <= lc {
		return o.leftChild.kth(k)
	} else {
		return o.rightChild.kth(k - lc)
	}
}
