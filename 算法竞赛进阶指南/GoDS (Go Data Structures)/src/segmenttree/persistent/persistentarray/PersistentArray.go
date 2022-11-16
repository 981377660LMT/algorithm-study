// !可持久化数组
// https://www.luogu.com.cn/problem/P3919
// 如题，你需要维护这样的一个长度为N的数组，支持如下几种操作
// 1.在某个历史版本上修改某一个位置上的值
// 2.访问某个历史版本上的某一位置的值
// 此外，每进行一次操作（对于操作2，即为生成一个完全一样的版本，不作任何改动)，
// 就会生成一个新的版本。版本编号即为当前操作的编号(从1开始编号，版本0表示初始状态数组)

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
)

// 存在海量小对象时禁用GC
func init() {
	debug.SetGCPercent(-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	nums := make([]E, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	root0 := Build(0, n-1, nums) // !初始版本

	versions := make([]*Node, 0, q+1)
	versions = append(versions, root0)
	for i := 0; i < q; i++ {
		var version, op, index, value int
		fmt.Fscan(in, &version, &op, &index)
		index--
		if op == 1 {
			fmt.Fscan(in, &value)
			newRoot := versions[version].Set(index, E(value))
			versions = append(versions, newRoot)
		} else {
			fmt.Fprintln(out, versions[version].Get(index))
			versions = append(versions, versions[version])
		}
	}
}

type E = int32

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
