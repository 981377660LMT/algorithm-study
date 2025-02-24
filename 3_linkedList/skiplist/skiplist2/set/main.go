package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	sl := NewSkipListSet[int](func(a, b int) bool { return a < b })

	sl.Insert(1)
	sl.Insert(3)
	sl.Insert(2)

	fmt.Println(sl.LowerBound(2).key)
	fmt.Println(sl.LowerBound(3).key)
	fmt.Println(sl.LowerBound(-1).key)

	fmt.Println(sl.Find(2))
}

type node[T any] struct {
	key  T
	next []nodeNext[T]
}

type nodeNext[T any] struct {
	pointer  *node[T]
	distance int
}

type SkipListSet[T any] struct {
	head   *node[T]
	comp   func(a, b T) bool
	height int
	size   int
}

func NewSkipListSet[T any](comp func(a, b T) bool) *SkipListSet[T] {
	s := &SkipListSet[T]{
		comp:   comp,
		height: 1,
		size:   0,
	}
	s.head = &node[T]{}
	s.head.next = []nodeNext[T]{{pointer: nil, distance: 0}}
	return s
}

func testJump() bool {
	return rand.Intn(64) < 30
}

func (s *SkipListSet[T]) insertNode(cur *node[T], level int, res **node[T], key T) int {
	distance := 0
	for {
		nxt := cur.next[level].pointer
		dis := cur.next[level].distance
		if nxt == nil || s.comp(key, nxt.key) {
			break
		} else if s.comp(nxt.key, key) {
			distance += dis + 1
			cur = nxt
		} else {
			*res = nil
			return -1
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
		cur.next[level] = nodeNext[T]{
			pointer:  *res,
			distance: below,
		}
		return distance + below
	}

	oldNext := cur.next[level].pointer
	newNode := &node[T]{
		key: key,
	}
	newNode.next = []nodeNext[T]{{pointer: oldNext, distance: 0}}
	*res = newNode
	cur.next[level].pointer = newNode
	return distance
}

func cloneLevel[T any](n *node[T], oldPointer *node[T], newDistance int) *node[T] {
	n.next = append(n.next, nodeNext[T]{pointer: oldPointer, distance: newDistance})
	return n
}

func (s *SkipListSet[T]) Insert(key T) {
	var res *node[T]
	dist := s.insertNode(s.head, s.height-1, &res, key)
	if res == nil {
		return
	}
	s.size++
	if dist < 0 || !testJump() {
		return
	}
	s.height++
	s.head.next = append(s.head.next, nodeNext[T]{pointer: res, distance: dist})
	res.next = append(res.next, nodeNext[T]{pointer: nil, distance: s.size - dist - 1})
}

func (s *SkipListSet[T]) Erase(key T) bool {
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
				cur.next[level] = nodeNext[T]{
					pointer:  nxt.next[level].pointer,
					distance: cur.next[level].distance + nxt.next[level].distance,
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

func (s *SkipListSet[T]) Find(key T) bool {
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil {
			nxt := cur.next[level].pointer
			if s.comp(nxt.key, key) {
				cur = nxt
			} else if !s.comp(key, nxt.key) {
				return true
			} else {
				break
			}
		}
	}
	return false
}

func (s *SkipListSet[T]) Rank(key T) int {
	cur := s.head
	ord := 0
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil && s.comp(cur.next[level].pointer.key, key) {
			ord += cur.next[level].distance + 1
			cur = cur.next[level].pointer
		}
	}
	return ord
}

func (s *SkipListSet[T]) Kth(k int) *node[T] {
	if k < 0 || k >= s.size {
		return nil
	}
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil && cur.next[level].distance <= k {
			k -= cur.next[level].distance
			k--
			if k < 0 {
				return cur.next[level].pointer
			}
			cur = cur.next[level].pointer
		}
	}
	return cur
}

func (s *SkipListSet[T]) SmallerBound(key T) *node[T] {
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil && s.comp(cur.next[level].pointer.key, key) {
			cur = cur.next[level].pointer
		}
	}
	return cur
}

func (s *SkipListSet[T]) LowerBound(key T) *node[T] {
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
	ret := cur.next[0].pointer
	if ret != nil && !s.comp(ret.key, key) {
		return ret
	}
	return nil
}

func (s *SkipListSet[T]) UpperBound(key T) *node[T] {
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil && !s.comp(key, cur.next[level].pointer.key) {
			cur = cur.next[level].pointer
		}
	}
	return cur.next[0].pointer
}

func (s *SkipListSet[T]) Size() int {
	return s.size
}

func (s *SkipListSet[T]) Empty() bool {
	return s.size == 0
}

func (s *SkipListSet[T]) Count(key T) int {
	if s.Find(key) {
		return 1
	}
	return 0
}
