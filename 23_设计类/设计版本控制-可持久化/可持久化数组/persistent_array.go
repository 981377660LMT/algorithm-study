package main

import (
	"fmt"
	"math/bits"
)

func main() {
	arr0 := NewPersistentArray(10, func(i int32) int32 { return i })
	fmt.Println(arr0.GetAll())
	arr1 := arr0.Set(3, 100)
	fmt.Println(arr0.GetAll())
	fmt.Println(arr1.GetAll())
	fmt.Println(arr0.Get(3))
	fmt.Println(arr1.Get(3))
}

type V = int32

type pNode struct {
	key         V
	left, right *pNode
}

func (p *pNode) Copy() *pNode {
	return &pNode{key: p.key, left: p.left, right: p.right}
}

type PersistentArray struct {
	n    int32
	root *pNode
}

func NewPersistentArray(n int32, f func(int32) V) *PersistentArray {
	res := &PersistentArray{n: n}
	res.root = res._build(n, f)
	return res
}

func (p *PersistentArray) Set(index int32, value V) *PersistentArray {
	node := p.root
	if node == nil {
		return p._new(nil)
	}
	newNode := node.Copy()
	res := p._new(newNode)
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
	newNode.key = value
	return res
}

func (p *PersistentArray) Get(index int32) V {
	node := p.root
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

func (p *PersistentArray) Copy() *PersistentArray {
	if p.root == nil {
		return p._new(nil)
	} else {
		return p._new(p.root.Copy())
	}
}

func (p *PersistentArray) GetAll() []V {
	node := p.root
	if node == nil {
		return nil
	}
	res := make([]V, 0, p.n)
	q := []*pNode{node}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		res = append(res, cur.key)
		if cur.left != nil {
			q = append(q, cur.left)
		}
		if cur.right != nil {
			q = append(q, cur.right)
		}
	}
	return res
}

func (p *PersistentArray) Len() int32 { return p.n }

func (p *PersistentArray) _build(n int32, f func(int32) V) *pNode {
	if n == 0 {
		return nil
	}
	pool := make([]*pNode, n)
	for i := int32(0); i < n; i++ {
		pool[i] = &pNode{key: f(i)}
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

func (p *PersistentArray) _new(root *pNode) *PersistentArray {
	return &PersistentArray{n: p.n, root: root}
}
