// https://www.luogu.com.cn/problem/solution/P3835
// TODO
// !FIXME

package persistentfhqtreap

import (
	"fmt"
	"time"
)

// https://www.cnblogs.com/dx123/p/16584604.html
func demo() {
	t := NewPersistentFHQTreap(func(a, b interface{}) int {
		return a.(int) - b.(int)
	}, 10)
	res := t.Add(&t.root, 1)
	fmt.Println(res)
	res = t.Add(&t.root, 1)
	fmt.Println(res)
}

type node struct {
	left, right int
	// 堆的随机权值
	priority uint
	// 子树大小
	size int
	// 树结点的值
	value interface{}
}

// 需要开 50倍长度的数组 (动态开点)
type PersistentFHQTreap struct {
	seed       uint
	nodeId     int // 从1开始
	root       int
	comparator func(a, b interface{}) int
	nodes      []node
}

func NewPersistentFHQTreap(comparator func(a, b interface{}) int, n int) *PersistentFHQTreap {
	return &PersistentFHQTreap{
		seed:       uint(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, (50*n + 10)),
	}
}

func (t *PersistentFHQTreap) Add(root *int, value interface{}) int {
	var x, y, z int
	t.split(*root, value, &x, &y, false)
	t.newNode(&z, value)
	t.root = t.merge(t.merge(x, z), y)
	return t.nodeId - 1
}

func (sl *PersistentFHQTreap) At(root int, index int) interface{} {
	n := sl.Len()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("%d index out of range: [%d,%d]", index, 0, n-1))
	}
	return sl.nodes[sl.kthNode(sl.root, index+1)].value
}

func (t *PersistentFHQTreap) Discard(root *int, value interface{}) int {
	var x, y, z int
	t.split(*root, value, &x, &z, false)
	t.split(x, value, &x, &y, true)
	y = t.merge(t.nodes[y].left, t.nodes[y].right)
	*root = t.merge(t.merge(x, y), z)
	return t.nodeId - 1
}

func (t *PersistentFHQTreap) BisectLeft(root *int, value interface{}) int {
	var x, y int
	t.split(*root, value, &x, &y, true)
	res := t.nodes[x].size
	*root = t.merge(x, y)
	return res
}

func (t *PersistentFHQTreap) BisectRight(root *int, value interface{}) int {
	var x, y int
	t.split(*root, value, &x, &y, false)
	res := t.nodes[x].size
	*root = t.merge(x, y)
	return res
}

func (t *PersistentFHQTreap) kthNode(root int, k int) int {
	if k <= t.nodes[t.nodes[root].left].size {
		return t.kthNode(t.nodes[root].left, k)
	}
	if k == t.nodes[t.nodes[root].left].size+1 {
		return root
	}
	return t.kthNode(t.nodes[root].right, k-t.nodes[t.nodes[root].left].size-1)
}

func (t *PersistentFHQTreap) split(root int, value interface{}, x, y *int, strictLess bool) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	if strictLess {
		if t.comparator(t.nodes[root].value, value) < 0 {
			// !动态开点
			t.nodeId++
			*x = t.nodeId
			t.nodes[*x] = t.nodes[root]
			t.split(t.nodes[root].right, value, &t.nodes[*x].right, y, strictLess)
			t.pushUp(*x)
		} else {
			t.nodeId++
			*y = t.nodeId
			t.nodes[*y] = t.nodes[root]
			t.split(t.nodes[root].left, value, x, &t.nodes[*y].left, strictLess)
			t.pushUp(*y)
		}
	} else {
		if t.comparator(t.nodes[root].value, value) <= 0 {
			t.nodeId++
			*x = t.nodeId
			t.nodes[*x] = t.nodes[root]
			t.split(t.nodes[root].right, value, &t.nodes[*x].right, y, strictLess)
			t.pushUp(*x)
		} else {
			t.nodeId++
			*y = t.nodeId
			t.nodes[*y] = t.nodes[root]
			t.split(t.nodes[root].left, value, x, &t.nodes[*y].left, strictLess)
			t.pushUp(*y)
		}
	}

}

func (sl *PersistentFHQTreap) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x + y
	}
	if sl.nodes[x].priority < sl.nodes[y].priority {
		sl.nodes[x].right = sl.merge(sl.nodes[x].right, y)
		sl.pushUp(x)
		return x
	}
	sl.nodes[y].left = sl.merge(x, sl.nodes[y].left)
	sl.pushUp(y)
	return y
}

func (t *PersistentFHQTreap) Len() int {
	return t.nodes[t.root].size
}

func (sl *PersistentFHQTreap) pushUp(root int) {
	sl.nodes[root].size = sl.nodes[sl.nodes[root].left].size + sl.nodes[sl.nodes[root].right].size + 1
}

func (t *PersistentFHQTreap) newNode(x *int, value interface{}) {
	t.nodeId++
	*x = t.nodeId
	t.nodes[*x].value = value
	t.nodes[*x].priority = t.fastRand()
	t.nodes[*x].size = 1
}

func (t *PersistentFHQTreap) fastRand() uint {
	t.seed ^= t.seed << 13
	t.seed ^= t.seed >> 17
	t.seed ^= t.seed << 5
	return t.seed
}
