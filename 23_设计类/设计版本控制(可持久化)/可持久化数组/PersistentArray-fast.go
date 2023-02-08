// !可持久化数组
// https://www.luogu.com.cn/problem/P3919
// 如题，你需要维护这样的一个长度为N的数组，支持如下几种操作
// 1.在某个历史版本上修改某一个位置上的值
// 2.访问某个历史版本上的某一位置的值
// 此外，每进行一次操作（对于操作2，即为生成一个完全一样的版本，不作任何改动)，
// 就会生成一个新的版本。版本编号即为当前操作的编号(从1开始编号，版本0表示初始状态数组)

package main

import (
	"fmt"
	"runtime/debug"
)

// 单组测试时禁用gc
func init() {
	debug.SetGCPercent(-1)
}

func main() {
	nums := NewPersistentArray([]E{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	v1 := nums.Set(0, 0, 100) // 初始版本号为0
	v2 := nums.Set(v1, 1, 200)
	nums.Set(v2, 2, 300)
	fmt.Println(nums.Get(0, 0), nums.Get(v1, 1), nums.Get(v2, 2))
}

type E = int

type PersistentArray struct {
	CurVersion int // 当前版本号从0开始
	versions   []*Node
}

func NewPersistentArray(nums []E) *PersistentArray {
	root := Build(0, len(nums)-1, nums)
	return &PersistentArray{versions: []*Node{root}}
}

// 将版本version上的index位置的值修改为value, 返回新版本号
//  0 <= version <= curVersion
//  0 <= index < len(nums)
func (o *PersistentArray) Set(version, index int, value E) int {
	preRoot := o.versions[version]
	newRoot := preRoot.Set(index, value)
	o.versions = append(o.versions, newRoot)
	o.CurVersion++
	return o.CurVersion
}

// 返回版本version上的index位置的值
//  0 <= version <= curVersion
//  0 <= index < len(nums)
func (o *PersistentArray) Get(version, index int) E {
	root := o.versions[version]
	return root.Get(index)
}

type Node struct {
	left, right           int
	value                 E
	leftChild, rightChild *Node
}

func Build(left, right int, nums []E) *Node {
	node := &Node{left: left, right: right}
	if left == right {
		node.value = nums[left]
		return node
	}
	mid := (left + right) >> 1
	node.leftChild = Build(left, mid, nums)
	node.rightChild = Build(mid+1, right, nums)
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

func (o *Node) String() string {
	res := o.dfs()
	return fmt.Sprintf("PersistentArray %v", res)
}

func (o *Node) dfs() []E {
	res := []E{}
	if o.left == o.right {
		res = append(res, o.value)
		return res
	}
	res = append(res, o.leftChild.dfs()...)
	res = append(res, o.rightChild.dfs()...)
	return res
}
