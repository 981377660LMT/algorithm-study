package skiplist

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type msNode[T any] struct {
	key        T
	nodeWeight int
	next       []msNodeNext[T]
}

type msNodeNext[T any] struct {
	pointer  *msNode[T]
	distance int
}

type SkipListMultiset[T any] struct {
	head   *msNode[T]
	comp   func(a, b T) bool
	height int
	size   int
}

func NewSkipListMultiset[T any](comp func(a, b T) bool) *SkipListMultiset[T] {
	s := &SkipListMultiset[T]{
		comp:   comp,
		height: 1,
	}
	s.head = &msNode[T]{}
	s.head.next = []msNodeNext[T]{{}}
	return s
}

func testJump() bool {
	return rand.Intn(64) < 30
}

func (s *SkipListMultiset[T]) insertNode(cur *msNode[T], level int, res **msNode[T], key T) int {
	distance := 0
	for {
		nxt := cur.next[level].pointer
		dis := cur.next[level].distance
		if nxt == nil || s.comp(key, nxt.key) {
			break
		} else if s.comp(nxt.key, key) {
			distance += dis + nxt.nodeWeight
			cur = nxt
		} else {
			nxt.nodeWeight++
			*res = nil
			return distance + dis
		}
	}

	if level > 0 {
		below := s.insertNode(cur, level-1, res, key)
		if *res == nil || below < 0 {
			return -1
		}
		cur.next[level].distance++
		if !testJump() {
			return -1
		}
		oldPointer := cur.next[level].pointer
		oldDist := cur.next[level].distance
		*res = cloneLevel(*res, oldPointer, oldDist-below-1)
		cur.next[level] = msNodeNext[T]{
			pointer:  *res,
			distance: below,
		}
		return distance + below
	}

	oldNext := cur.next[level].pointer
	newNode := &msNode[T]{
		key:        key,
		nodeWeight: 1,
	}
	newNode.next = []msNodeNext[T]{{pointer: oldNext, distance: 0}}
	*res = newNode
	cur.next[level].pointer = newNode
	return distance
}

func cloneLevel[T any](n *msNode[T], oldPointer *msNode[T], newDistance int) *msNode[T] {
	n.next = append(n.next, msNodeNext[T]{pointer: oldPointer, distance: newDistance})
	return n
}

func (s *SkipListMultiset[T]) Insert(key T) {
	var res *msNode[T]
	dist := s.insertNode(s.head, s.height-1, &res, key)
	if res == nil {
		s.size++
		return
	}
	s.size++
	if dist < 0 || !testJump() {
		return
	}
	s.height++
	s.head.next = append(s.head.next, msNodeNext[T]{pointer: res, distance: dist})
	res.next = append(res.next, msNodeNext[T]{pointer: nil, distance: s.size - dist - res.nodeWeight})
	return
}

func (s *SkipListMultiset[T]) Erase(key T) bool {
	removed := false
	for level := s.height - 1; level >= 0; level-- {
		cur := s.head
		for {
			nxt := cur.next[level].pointer
			if nxt == nil || s.comp(key, nxt.key) {
				break
			}
			if s.comp(nxt.key, key) {
				cur = nxt
			} else {
				cur.next[level].distance--
				nxt.nodeWeight--
				if nxt.nodeWeight == 0 {
					cur.next[level].pointer = nxt.next[level].pointer
				}
				removed = true
				break
			}
		}
	}
	if removed {
		s.size--
	}
	return removed
}

func (s *SkipListMultiset[T]) EraseCount(key T, count int) {
	for i := 0; i < count; i++ {
		if !s.Erase(key) {
			break
		}
	}
}

func (s *SkipListMultiset[T]) Find(key T) *msNode[T] {
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil {
			if s.comp(cur.next[level].pointer.key, key) {
				cur = cur.next[level].pointer
			} else if !s.comp(key, cur.next[level].pointer.key) {
				return cur.next[level].pointer
			} else {
				break
			}
		}
	}
	return nil
}

func (s *SkipListMultiset[T]) Count(key T) int {
	if node := s.Find(key); node != nil {
		return node.nodeWeight
	}
	return 0
}

func (s *SkipListMultiset[T]) Size() int {
	return s.size
}

func (s *SkipListMultiset[T]) Empty() bool {
	return s.size == 0
}

func (s *SkipListMultiset[T]) Rank(key T) int {
	cur := s.head
	order := 0
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil && s.comp(cur.next[level].pointer.key, key) {
			order += cur.next[level].distance + cur.next[level].pointer.nodeWeight
			cur = cur.next[level].pointer
		}
	}
	return order
}

func (s *SkipListMultiset[T]) Kth(k int) *msNode[T] {
	if k < 0 || k >= s.size {
		return nil
	}
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil && cur.next[level].distance <= k {
			k -= cur.next[level].distance
			nxt := cur.next[level].pointer
			k -= nxt.nodeWeight
			if k < 0 {
				return nxt
			}
			cur = nxt
		}
	}
	return cur
}

func (s *SkipListMultiset[T]) LowerBound(key T) *msNode[T] {
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil {
			if s.comp(cur.next[level].pointer.key, key) {
				cur = cur.next[level].pointer
			} else {
				break
			}
		}
	}
	candidate := cur.next[0].pointer
	if candidate != nil && !s.comp(candidate.key, key) {
		return candidate
	}
	return nil
}

func (s *SkipListMultiset[T]) UpperBound(key T) *msNode[T] {
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil && !s.comp(key, cur.next[level].pointer.key) {
			cur = cur.next[level].pointer
		}
	}
	return cur.next[0].pointer
}

func (s *SkipListMultiset[T]) SmallerBound(key T) *msNode[T] {
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil && s.comp(cur.next[level].pointer.key, key) {
			cur = cur.next[level].pointer
		}
	}
	return cur
}

func (s *SkipListMultiset[T]) Clear() {
	s.head = nil
	s.size = 0
	s.height = 0
}
