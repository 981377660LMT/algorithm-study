// !可持久化队列
// https://judge.yosupo.jp/problem/persistent_queue
// 初始的队列版本为-1
// 你需要维护这样的一个队列，支持如下几种操作
// 0 versioni x => 在版本 versioni 上入队x(append) 版本+1
// 1 versioni => 在版本 versioni 上出队(popleft) 版本+1 输出出队的元素
// q<=5e5 -1<=versioni<i

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	// queue := NewPersistentQueue(50)
	// queue.Append(0, 1)
	// queue.Append(1, 2)
	// fmt.Println(queue)
	// queue.PopLeft(2)
	// fmt.Println(queue)
	// queue.PopLeft(3)
	// fmt.Println(queue)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	queue := NewPersistentQueue(q)
	for i := 0; i < q; i++ {
		var op, version, x int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &version, &x)
			version++
			queue.Append(version, x)
		} else {
			fmt.Fscan(in, &version)
			version++
			_, x := queue.PopLeft(version)
			fmt.Fprintln(out, x)
		}
	}
}

type E = int

type PersistentQueue struct {
	MaxVersion int      // !初始版本号为0
	roots      []*Node  // !记录每个版本的数组
	bounds     [][2]int // !记录每个版本的数组的左右边界(左闭右开)
}

func NewPersistentQueue(maxQuery int) *PersistentQueue {
	res := &PersistentQueue{
		roots:      make([]*Node, 0, maxQuery+1),
		bounds:     make([][2]int, 0, maxQuery+1),
		MaxVersion: -1,
	}

	initNums := make([]E, maxQuery)
	root0 := Build(0, len(initNums)-1, initNums)
	res.addVersion(root0, 0, 0)
	return res
}

// version>=0
func (q *PersistentQueue) Append(version int, value E) (newVersion int) {
	root, left, right := q.roots[version], q.bounds[version][0], q.bounds[version][1]
	newRoot := root.Set(right, value)
	newVersion = q.addVersion(newRoot, left, right+1)
	return
}

// version>=0
func (q *PersistentQueue) PopLeft(version int) (newVersion int, popped E) {
	if q.Len(version) == 0 {
		panic(fmt.Sprintf("queue is empty in version %d", version))
	}
	root, left, right := q.roots[version], q.bounds[version][0], q.bounds[version][1]
	popped = root.Get(left)
	newRoot := root.Set(left, 0)
	newVersion = q.addVersion(newRoot, left+1, right)
	return
}

// version>=0
func (q *PersistentQueue) Get(version, index int) E {
	root, left, right := q.roots[version], q.bounds[version][0], q.bounds[version][1]
	len := right - left
	if index < 0 {
		index += len
	}
	if index < 0 || index >= len {
		panic(fmt.Sprintf("index %d out of range in version %d", index, version))
	}
	return root.Get(left + index)
}

// version>=0
func (q *PersistentQueue) Len(version int) int {
	left, right := q.bounds[version][0], q.bounds[version][1]
	return right - left
}

// print queue from version 0 to maxVersion
func (q *PersistentQueue) String() string {
	sb := strings.Builder{} // PersistentQueue
	sb.WriteString("PersistentQueue:\n")
	for i := 0; i <= q.MaxVersion; i++ {
		sb.WriteString(fmt.Sprintf("version %d: ", i))
		sb.WriteString("[")
		elements := make([]string, 0, q.bounds[i][1]-q.bounds[i][0])
		for j := q.bounds[i][0]; j < q.bounds[i][1]; j++ {
			elements = append(elements, fmt.Sprintf("%v", q.roots[i].Get(j)))
		}
		sb.WriteString(strings.Join(elements, ","))
		sb.WriteString("]")
		sb.WriteString("\n")
	}
	return sb.String()
}

func (q *PersistentQueue) addVersion(root *Node, left, right int) int {
	q.roots = append(q.roots, root)
	q.bounds = append(q.bounds, [2]int{left, right})
	q.MaxVersion++
	return q.MaxVersion
}

type Node struct {
	left, right           int
	size                  int
	value                 E
	leftChild, rightChild *Node
}

func Build(left, right int, nums []E) *Node {
	node := &Node{left: left, right: right}
	if left == right {
		node.value = nums[left]
		node.size = 1
		return node
	}
	mid := (left + right) >> 1
	node.leftChild = Build(left, mid, nums)
	node.rightChild = Build(mid+1, right, nums)
	node.pushUp()
	return node
}

func (o *Node) Get(index int) E {
	if o.left == o.right {
		return o.value
	}
	mid := (o.left + o.right) >> 1
	if index <= mid {
		return o.leftChild.Get(index)
	}
	return o.rightChild.Get(index)
}

func (o Node) Set(index int, value E) *Node {
	// !修改时拷贝一个新节点(这里用值作为接收者，已经隐式拷贝了一份root结点o)
	if o.left == o.right {
		o.value = value
		return &o
	}

	mid := (o.left + o.right) >> 1
	if index <= mid {
		o.leftChild = o.leftChild.Set(index, value)
	} else {
		o.rightChild = o.rightChild.Set(index, value)
	}

	return &o
}

func (o *Node) pushUp() {
	o.size = o.leftChild.size + o.rightChild.size
}

func (o *Node) String() string {
	res := o.dfs()
	return fmt.Sprintf("PersistentArray %v", res)
}

func (o *Node) dfs() []E {
	res := make([]E, 0, o.size)
	if o.left == o.right {
		res = append(res, o.value)
		return res
	}
	res = append(res, o.leftChild.dfs()...)
	res = append(res, o.rightChild.dfs()...)
	return res
}
