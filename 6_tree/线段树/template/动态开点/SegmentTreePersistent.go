// https://ei1333.github.io/library/structure/segment-tree/persistent-segment-tree.hpp
// 可持久化线段树.

// NewSegmentTreePersistent() *SegmentTreePersistent
// Build(leaves []E) *PNode
// Update(root *PNode, index int, value E) *PNode
// Query(root *PNode, start, end int) E

package main

import (
	"fmt"
	"time"
)

func main() {
	seg := NewSegmentTreePersistent()
	root1 := seg.Build(5, func(i int32) E { return i })
	fmt.Println(seg.Query(root1, 0, 5), root1)
	root2 := seg.Update(root1, 0, 10)
	fmt.Println(seg.Query(root2, 0, 5), root2, root1)

	time1 := time.Now()
	n := int32(2e5)
	big := seg.Build(n, func(i int32) E { return 0 })
	for i := int32(0); i < n; i++ {
		big = seg.Update(big, i, 1)
		seg.Query(big, 0, i)
	}
	time2 := time.Now()
	fmt.Println(time2.Sub(time1)) // 222.3792ms
}

type E = int32

func (*SegmentTreePersistent) e() E        { return 0 }
func (*SegmentTreePersistent) op(a, b E) E { return a + b }

type SegmentTreePersistent struct {
	size int32
}

func NewSegmentTreePersistent() *SegmentTreePersistent {
	return &SegmentTreePersistent{}
}

func (s *SegmentTreePersistent) Build(n int32, f func(i int32) E) *Node {
	s.size = n
	return s._build(0, s.size, f)
}

func (s *SegmentTreePersistent) Update(root *Node, index int32, value E) *Node {
	if index < 0 || index >= s.size {
		return root
	}
	return s._update(root, index, value, 0, s.size)
}

func (s *SegmentTreePersistent) Query(root *Node, start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > s.size {
		end = s.size
	}
	if start >= end {
		return s.e()
	}
	return s._query(root, start, end, 0, s.size)
}

func (s *SegmentTreePersistent) _build(l, r int32, f func(i int32) E) *Node {
	if l+1 >= r {
		return &Node{data: f(l)}
	}
	mid := (l + r) >> 1
	return s._merge(s._build(l, mid, f), s._build(mid, r, f))
}

func (s *SegmentTreePersistent) _merge(l, r *Node) *Node {
	return &Node{data: s.op(l.data, r.data), l: l, r: r}
}

func (s *SegmentTreePersistent) _update(root *Node, index int32, value E, l, r int32) *Node {
	if r <= index || index+1 <= l {
		return root
	}
	if index <= l && r <= index+1 {
		return &Node{data: value}
	}
	mid := (l + r) >> 1
	return s._merge(s._update(root.l, index, value, l, mid), s._update(root.r, index, value, mid, r))
}

func (s *SegmentTreePersistent) _query(root *Node, start, end int32, l, r int32) E {
	if r <= start || end <= l {
		return s.e()
	}
	if start <= l && r <= end {
		return root.data
	}
	mid := (l + r) >> 1
	return s.op(s._query(root.l, start, end, l, mid), s._query(root.r, start, end, mid, r))
}

type Node struct {
	data E
	l, r *Node
}

func (p *Node) String() string {
	leaves := []E{}
	var dfs func(*Node)
	dfs = func(root *Node) {
		if root == nil {
			return
		}
		if root.l == nil && root.r == nil {
			leaves = append(leaves, root.data)
			return
		}
		dfs(root.l)
		dfs(root.r)
	}
	dfs(p)
	return fmt.Sprintf("%v", leaves)
}
