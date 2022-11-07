// https://github.dev/EndlessCheng/codeforces-go/blob/6d127a66c2a11651e8d11783af687264780e82a8/copypasta/treap.go#L2
// https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e
// !https://www.luogu.com.cn/blog/specialflag/solution-p3369 (详细)
// !https://baobaobear.github.io/post/20191215-fhq-treap/ (支持区间加、区间删和区间求和 pushDown里定义)

// 163 普通平衡树 FHQ Treap https://www.bilibili.com/video/BV1kY4y1j7LC
// !164 文艺平衡树 FHQ Treap https://www.bilibili.com/video/BV1pd4y1D7Nu (区间翻转)
// 169 可持久化平衡树 https://www.bilibili.com/video/BV1sB4y1L79D

// FHQ Treap依靠分裂与合并两个操作来维护树的平衡,
// !这种操作方式支持维护序列、区间操作、可持久化等特性。

package fhqtreap

import (
	"fmt"
	"time"
)

func demo() {
	t := NewFHQTreap([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, func(a, b interface{}) int {
		return a.(int) - b.(int)
	})
	fmt.Println(t.InOrder())

	t.Reverse(0, 4)
	fmt.Println(t.InOrder())
}

type node struct {
	left, right int
	// 堆的随机权值
	priority uint
	// 子树大小
	size int
	// 树结点的值
	value interface{}

	// 有时候需要维护一些额外的信息(翻转区间等)
	lazy int
}

type FHQTreap struct {
	seed       uint
	nodeId     int // 从1开始
	root       int
	comparator func(a, b interface{}) int
	nodes      []node
}

func NewFHQTreap(nums []interface{}, comparator func(a, b interface{}) int) *FHQTreap {
	n := len(nums)
	res := &FHQTreap{
		seed:       uint(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, n+2),
		nodeId:     1,
	}

	// !build
	for i := range nums {
		res.root = res.merge(res.root, res.newNode(nums[i]))
	}
	return res
}

// !翻转闭区间 [left, right] 的值 (0-indexed)
//  0 <= left <= right < t.Len()
func (t *FHQTreap) Reverse(left, right int) {
	var x, y, z int
	t.split(t.root, right+1, &x, &z)
	t.split(x, left, &x, &y)
	t.nodes[y].lazy ^= 1
	t.root = t.merge(t.merge(x, y), z)
}

func (t *FHQTreap) InOrder() []interface{} {
	res := make([]interface{}, 0, t.Len())
	t.inOrder(t.root, &res)
	return res
}

func (t *FHQTreap) inOrder(root int, res *[]interface{}) {
	if root == 0 {
		return
	}
	t.pushDown(root) // !注意下传 lazy 标记
	t.inOrder(t.nodes[root].left, res)
	*res = append(*res, t.nodes[root].value)
	t.inOrder(t.nodes[root].right, res)
}

func (t *FHQTreap) Len() int {
	return t.nodes[t.root].size
}

// 按照k(排名)分裂
func (t *FHQTreap) split(root, k int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	t.pushDown(root)

	if leftSize := t.nodes[t.nodes[root].left].size; k > leftSize {
		k -= leftSize + 1
		*x = root
		t.split(t.nodes[root].right, k, &t.nodes[root].right, y)
	} else {
		*y = root
		t.split(t.nodes[root].left, k, x, &t.nodes[root].left)
	}

	t.pushUp(root)
}

// 合并左Treap的根结点x和右Treap的根结点y
func (t *FHQTreap) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x + y
	}

	// 随机值小的作为根结点
	if t.nodes[x].priority < t.nodes[y].priority {
		t.pushDown(x)
		t.nodes[x].right = t.merge(t.nodes[x].right, y)
		t.pushUp(x)
		return x
	}

	t.pushDown(y)
	t.nodes[y].left = t.merge(x, t.nodes[y].left)
	t.pushUp(y)
	return y
}

func (t *FHQTreap) pushDown(root int) {
	if t.nodes[root].lazy == 0 {
		return
	}

	// !交换左右子树
	t.nodes[root].left, t.nodes[root].right = t.nodes[root].right, t.nodes[root].left
	t.nodes[t.nodes[root].left].lazy ^= 1
	t.nodes[t.nodes[root].right].lazy ^= 1
	t.nodes[root].lazy = 0
}

func (t *FHQTreap) pushUp(root int) {
	t.nodes[root].size = t.nodes[t.nodes[root].left].size + t.nodes[t.nodes[root].right].size + 1
}

func (t *FHQTreap) newNode(value interface{}) int {
	t.nodeId++
	index := t.nodeId
	t.nodes[index].value = value
	t.nodes[index].priority = t.fastRand()
	t.nodes[index].size = 1
	return index
}

func (t *FHQTreap) fastRand() uint {
	t.seed ^= t.seed << 13
	t.seed ^= t.seed >> 17
	t.seed ^= t.seed << 5
	return t.seed
}
