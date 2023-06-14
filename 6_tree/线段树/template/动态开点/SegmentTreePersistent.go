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
	root1 := seg.Build([]E{1, 2, 3, 4, 5})
	fmt.Println(seg.Query(root1, 0, 5), root1)
	root2 := seg.Update(root1, 0, 10)
	fmt.Println(seg.Query(root2, 0, 5), root2, root1)

	time1 := time.Now()
	n := int(2e5)
	big := seg.Build(make([]E, n))
	for i := 0; i < n; i++ {
		big = seg.Update(big, i, 1)
		seg.Query(big, 0, i)
	}
	time2 := time.Now()
	fmt.Println(time2.Sub(time1)) // 222.3792ms
}

type E = int

func (*SegmentTreePersistent) e() E        { return 0 }
func (*SegmentTreePersistent) op(a, b E) E { return a + b }

//
//
//
type SegmentTreePersistent struct {
	size int
}

func NewSegmentTreePersistent() *SegmentTreePersistent {
	return &SegmentTreePersistent{}
}

func (s *SegmentTreePersistent) Build(leaves []E) *PNode {
	s.size = len(leaves)
	return s._build(0, s.size, leaves)
}

func (s *SegmentTreePersistent) Update(root *PNode, index int, value E) *PNode {
	if index < 0 || index >= s.size {
		return root
	}
	return s._update(root, index, value, 0, s.size)
}

func (s *SegmentTreePersistent) Query(root *PNode, start, end int) E {
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

func (s *SegmentTreePersistent) _build(l, r int, leaves []E) *PNode {
	if l+1 >= r {
		return &PNode{data: leaves[l]}
	}
	mid := (l + r) >> 1
	return s._merge(s._build(l, mid, leaves), s._build(mid, r, leaves))
}

func (s *SegmentTreePersistent) _merge(l, r *PNode) *PNode {
	return &PNode{data: s.op(l.data, r.data), l: l, r: r}
}

func (s *SegmentTreePersistent) _update(root *PNode, index int, value E, l, r int) *PNode {
	if r <= index || index+1 <= l {
		return root
	}
	if index <= l && r <= index+1 {
		return &PNode{data: value}
	}
	mid := (l + r) >> 1
	return s._merge(s._update(root.l, index, value, l, mid), s._update(root.r, index, value, mid, r))
}

func (s *SegmentTreePersistent) _query(root *PNode, start, end int, l, r int) E {
	if r <= start || end <= l {
		return s.e()
	}
	if start <= l && r <= end {
		return root.data
	}
	mid := (l + r) >> 1
	return s.op(s._query(root.l, start, end, l, mid), s._query(root.r, start, end, mid, r))
}

type PNode struct {
	data E
	l, r *PNode
}

func (p *PNode) String() string {
	leaves := []E{}
	var dfs func(*PNode)
	dfs = func(root *PNode) {
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
