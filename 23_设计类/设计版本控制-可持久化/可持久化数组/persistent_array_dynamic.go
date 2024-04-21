// persistent_array_dynamic

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	arr0 := NewPersistentArray(10, func(i int32) int32 { return i })
	fmt.Println(arr0.ToList())
	arr1 := arr0.Set(3, 100)
	fmt.Println(arr0.ToList())
	fmt.Println(arr1.ToList())
	fmt.Println(arr0.Get(3))
	fmt.Println(arr1.Get(3))
}

type aNode[V any] struct {
	left, right *aNode[V]
	key         V
}

func newANode[V any](key V) *aNode[V] {
	return &aNode[V]{key: key}
}

func (node *aNode[V]) Copy() *aNode[V] {
	return &aNode[V]{key: node.key, left: node.left, right: node.right}
}

type PersistentArray[V any] struct {
	n    int32
	root *aNode[V]
}

func NewPersistentArray[V any](n int32, f func(int32) V) *PersistentArray[V] {
	res := &PersistentArray[V]{n: n}
	res.root = res.build(n, f)
	return res
}

func (pa *PersistentArray[V]) Set(index int32, v V) *PersistentArray[V] {
	node := pa.root
	if node == nil {
		return pa.newWithRoot(nil)
	}
	newNode := node.Copy()
	res := pa.newWithRoot(newNode)
	index++
	b := bits.Len32(uint32(index))
	for i := b - 2; i >= 0; i-- {
		if index>>i&1 == 1 {
			node = node.right
			newNode.right = node.Copy()
			newNode = newNode.right
		} else {
			node = node.left
			newNode.left = node.Copy()
			newNode = newNode.left
		}
	}
	newNode.key = v
	return res
}

func (pa *PersistentArray[V]) Get(index int32) V {
	node := pa.root
	index++
	b := bits.Len32(uint32(index))
	for i := b - 2; i >= 0; i-- {
		if index>>i&1 == 1 {
			node = node.right
		} else {
			node = node.left
		}
	}
	return node.key
}

func (pa *PersistentArray[V]) Copy() *PersistentArray[V] {
	if pa.root == nil {
		return &PersistentArray[V]{n: pa.n, root: nil}
	}
	return &PersistentArray[V]{n: pa.n, root: pa.root.Copy()}
}

func (pa *PersistentArray[V]) ToList() []V {
	node := pa.root
	if node == nil {
		return nil
	}
	res := make([]V, 0, pa.n)
	queue := []*aNode[V]{node}
	for len(queue) > 0 {
		node = queue[0]
		queue = queue[1:]
		res = append(res, node.key)
		if node.left != nil {
			queue = append(queue, node.left)
		}
		if node.right != nil {
			queue = append(queue, node.right)
		}
	}
	return res
}

func (pa *PersistentArray[V]) Len() int32 {
	return pa.n
}

func (pa *PersistentArray[V]) build(n int32, f func(int32) V) *aNode[V] {
	if n == 0 {
		return nil
	}
	pool := make([]*aNode[V], n)
	for i := int32(0); i < n; i++ {
		pool[i] = newANode(f(i))
	}
	for i := int32(1); i <= n; i++ {
		if 2*i-1 < n {
			pool[i-1].left = pool[2*i-1]
		}
		if 2*i < n {
			pool[i-1].right = pool[2*i]
		}
	}
	return pool[0]
}

func (pa *PersistentArray[V]) newWithRoot(root *aNode[V]) *PersistentArray[V] {
	res := &PersistentArray[V]{n: pa.n, root: root}
	return res
}
