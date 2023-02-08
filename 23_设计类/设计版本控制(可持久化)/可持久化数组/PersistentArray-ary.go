// !可持久化数组
// Persistent Array
// Reference: <https://qiita.com/hotman78/items/9c643feae1de087e6fc5>
//            <https://ei1333.github.io/luzhiled/snippets/structure/persistent-array.html>
// - (2^LOG)-ary tree-based
// - Fully persistent
// - `get(root, k)`:  get k-th element
// - `set(root, k, data)`: make new array whose k-th element is updated to data

package main

import (
	"fmt"
	"runtime/debug"
	"time"
)

// 单组测试时禁用gc
func init() {
	debug.SetGCPercent(-1)
}

func main() {
	time1 := time.Now()
	arr := NewPersistentArray(3, -1)
	root := arr.Build(make([]int, 1e5))
	for i := 0; i < 1e5; i++ {
		arr.Set(root, i, i)
		arr.Get(root, i)
	}
	fmt.Println(time.Since(time1))
}

type E = int

type AryNode struct {
	data     E
	children []*AryNode
}

type PersistentArray struct {
	Root *AryNode // 初始版本
	null E        // 当index越界/不存在时返回的值
	log  int
}

// aryLog :一般取3或4
// null :当index越界/不存在时返回的值
func NewPersistentArray(aryLog int, null E) *PersistentArray {
	res := &PersistentArray{log: aryLog, null: null}
	return res
}

func (p *PersistentArray) Build(nums []E) *AryNode {
	for i := 0; i < len(nums); i++ {
		p.Root = p.Set(p.Root, i, nums[i])
	}
	return p.Root
}

// 0<=index<len(nums) 如果index越界，返回p.null
func (p *PersistentArray) Get(version *AryNode, index int) E {
	if version == nil {
		return p.null
	}
	if index == 0 {
		return version.data
	}
	return p.Get(version.children[index&((1<<p.log)-1)], index>>p.log)
}

func (p *PersistentArray) Set(version *AryNode, index int, data E) *AryNode {
	newNode := AryNode{data: p.null, children: make([]*AryNode, 1<<p.log)}
	if version != nil {
		copy(newNode.children, version.children)
		newNode.data = version.data
	}

	version = &newNode
	if index == 0 {
		version.data = data
	} else {
		ptr := p.Set(version.children[index&((1<<p.log)-1)], index>>p.log, data)
		version.children[index&((1<<p.log)-1)] = ptr
	}
	return version
}
