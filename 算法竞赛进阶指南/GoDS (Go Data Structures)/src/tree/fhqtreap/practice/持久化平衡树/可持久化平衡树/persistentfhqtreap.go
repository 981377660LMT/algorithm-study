// https://www.luogu.com.cn/problem/solution/P3835
// 可持久化平衡树
// TODO
// !FIXME

package main

import (
	"fmt"
	"time"
)

// 您需要写一种数据结构（可参考题目标题），来维护一个可重整数集合，
// 其中需要提供以下操作（ 对于各个以往的历史版本 ）：
// 1、 插入 x
// 2、 删除 x（若有多个相同的数，应只删除一个，如果没有请忽略该操作）
// 3、 查询 x 的排名（排名定义为比当前数小的数的个数 +1）
// 4、查询排名为 x 的数
// 5、 求 x 的前驱（前驱定义为小于 x，且最大的数，如不存在输出 −2^31+1 ）
// 6、求 x 的后继（后继定义为大于 x，且最小的数，如不存在输出 2^31 −1 ）
// 和原本平衡树不同的一点是，每一次的任何操作都是基于某一个历史版本，
// 同时生成一个新的版本。（操作3, 4, 5, 6即保持原版本无变化）

// 每个版本的编号即为操作的序号（版本0即为初始状态，空树）
// https://www.cnblogs.com/dx123/p/16584604.html

// !哪里有问题
func main() {
	// in := bufio.NewReader(os.Stdin)
	// out := bufio.NewWriter(os.Stdout)
	// defer out.Flush()

	// var n int
	// fmt.Fscan(in, &n)
	// t := NewPersistentFHQTreap(func(a, b interface{}) int {
	// 	return a.(int) - b.(int)
	// }, n)

	// roots := make([]int, n+5) // 保存每个版本的根节点
	// for i := 0; i < n; i++ {
	// 	var version, op, x int
	// 	fmt.Fscan(in, &version, &op, &x)
	// 	roots[i] = roots[version]
	// 	switch op {
	// 	case 1:
	// 		newVersion := t.Add(&roots[i], x)
	// 		fmt.Println(newVersion, op, "asa")
	// 		roots[i] = newVersion
	// 	case 2:
	// 		newVersion := t.Discard(&roots[i], x)
	// 		roots[i] = newVersion
	// 	case 3:
	// 		fmt.Fprintln(out, t.At(roots[i], x))
	// 	case 4:
	// 		fmt.Fprintln(out, t.BisectLeft(&roots[i], x))
	// 	case 5:
	// 		fmt.Fprintln(out, t.lower(&roots[i], x))
	// 	case 6:
	// 		fmt.Fprintln(out, t.upper(&roots[i], x))
	// 	}
	// }
	sl := NewPersistentFHQTreap(func(a, b interface{}) int {
		return a.(int) - b.(int)
	}, 1000)
	roots := make([]int, 1000) // 保存每个版本的根节点
	roots[0] = sl.Add(&roots[0], 1)
	roots[0] = sl.Add(&roots[0], 2)
	roots[0] = sl.Add(&roots[0], 3)
	roots[1] = sl.Discard(&roots[0], 2)
	roots[1] = sl.Add(&roots[1], 4)
	fmt.Println(sl.At(roots[1], 1))
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
	z = t.newNode(value)
	t.root = t.merge(t.merge(x, z), y)
	return t.nodeId
}

func (sl *PersistentFHQTreap) At(root int, index int) interface{} {
	n := sl.Len()
	if index < 0 {
		index += n
	}
	// if index < 0 || index >= n {
	// 	panic(fmt.Sprintf("%d index out of range: [%d,%d]", index, 0, n-1))
	// }
	return sl.nodes[sl.kthNode(sl.root, index+1)].value
}

func (t *PersistentFHQTreap) Discard(root *int, value interface{}) int {
	var x, y, z int
	t.split(*root, value, &x, &z, false)
	t.split(x, value, &x, &y, true)
	y = t.merge(t.nodes[y].left, t.nodes[y].right)
	*root = t.merge(t.merge(x, y), z)
	return t.nodeId
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
	if k == t.nodes[t.nodes[root].left].size+1 {
		return root
	}
	if k <= t.nodes[t.nodes[root].left].size {
		return t.kthNode(t.nodes[root].left, k)
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

func (t *PersistentFHQTreap) lower(root *int, value int) int {
	var x, y int
	t.split(*root, value, &x, &y, true)
	if x == 0 {
		return -2147483647
	}
	s := t.nodes[x].size
	res := t.nodes[t.kthNode(x, s)].value.(int)
	*root = t.merge(x, y)
	return res
}

func (t *PersistentFHQTreap) upper(root *int, value int) int {
	var x, y int
	t.split(*root, value, &x, &y, false)
	if y == 0 {
		return 2147483647
	}
	res := t.nodes[t.kthNode(y, 1)].value.(int)
	*root = t.merge(x, y)
	return res
}

func (t *PersistentFHQTreap) Len() int {
	return t.nodes[t.root].size
}

func (sl *PersistentFHQTreap) pushUp(root int) {
	sl.nodes[root].size = sl.nodes[sl.nodes[root].left].size + sl.nodes[sl.nodes[root].right].size + 1
}

func (t *PersistentFHQTreap) newNode(value interface{}) int {
	t.nodeId++
	id := t.nodeId
	t.nodes[id].value = value
	t.nodes[id].priority = t.fastRand()
	t.nodes[id].size = 1
	return id
}

func (t *PersistentFHQTreap) fastRand() uint {
	t.seed ^= t.seed << 13
	t.seed ^= t.seed >> 17
	t.seed ^= t.seed << 5
	return t.seed
}
