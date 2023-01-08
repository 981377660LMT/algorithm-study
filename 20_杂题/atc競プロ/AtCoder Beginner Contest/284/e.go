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

const TARGET int = 1e6

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	adjList := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
	}

	type Item struct {
		cur     int
		visited *Node
	}

	res := 1
	queue := make([]Item, 0, n)
	visited := Build(0, n-1, make([]int, n))
	visited = visited.Set(0, 1)
	queue = append(queue, Item{0, visited})

	for len(queue) != 0 {
		item := queue[0]
		queue = queue[1:]
		cur := item.cur
		visited := item.visited
		for _, next := range adjList[cur] {
			if visited.Get(next) == 1 {
				continue
			}
			res++
			if res == TARGET {
				fmt.Fprintln(out, TARGET)
				return
			}
			queue = append(queue, Item{next, visited.Set(next, 1)})
		}
	}

	fmt.Fprintln(out, res)
}

type E = int

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

func (o *Node) Get(index int) int {
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
