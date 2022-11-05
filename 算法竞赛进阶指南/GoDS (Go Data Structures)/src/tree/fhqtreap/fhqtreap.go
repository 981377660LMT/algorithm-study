// !163 普通平衡树 FHQ Treap https://www.bilibili.com/video/BV1kY4y1j7LC
// 164 文艺平衡树 FHQ Treap https://www.bilibili.com/video/BV1pd4y1D7Nu
// 169 可持久化平衡树 https://www.bilibili.com/video/BV1sB4y1L79D
package fhqtreap

import "time"

type node struct {
	left  *node
	right *node
	// 堆的随机权值
	priority int
	// 子树大小
	size int
	// 树结点的值
	value interface{}
}

func (cur *node) PushUp() {
	cur.size = cur.left.size + cur.right.size + 1
}

type FHQTreap struct {
	seed       uint
	root       *node
	comparator func(a, b interface{}) int
}

func NewFHQTreap(comparator func(a, b interface{}) int) *FHQTreap {
	return &FHQTreap{
		seed:       uint(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
	}
}

func (t *FHQTreap) fastRand() uint {
	t.seed ^= t.seed << 13
	t.seed ^= t.seed >> 17
	t.seed ^= t.seed << 5
	return t.seed
}
