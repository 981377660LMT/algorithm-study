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
	sl := NewSkipListMap[int, int](func(a, b int) bool { return a < b })

	sl.Insert(1, 1)
	sl.Insert(3, 3)
	sl.Insert(2, 2)

	fmt.Println(sl.LowerBound(2))
	fmt.Println(sl.LowerBound(3))
	fmt.Println(sl.LowerBound(1))

	fmt.Println(sl.Find(2).value)
}

type mapNode[K any, V any] struct {
	key   K
	value V
	next  []mapNodeNext[K, V]
}

type mapNodeNext[K any, V any] struct {
	pointer  *mapNode[K, V]
	distance int
}

type SkipListMap[K any, V any] struct {
	head   *mapNode[K, V]
	comp   func(a, b K) bool
	height int
	size   int
}

func NewSkipListMap[K any, V any](comp func(a, b K) bool) *SkipListMap[K, V] {
	s := &SkipListMap[K, V]{
		comp:   comp,
		height: 1,
		size:   0,
	}
	s.head = &mapNode[K, V]{}
	s.head.next = []mapNodeNext[K, V]{{pointer: nil, distance: 0}}
	return s
}

func testJump() bool {
	return rand.Intn(64) < 30
}

func cloneLevel[K any, V any](n *mapNode[K, V], oldPointer *mapNode[K, V], newDistance int) *mapNode[K, V] {
	n.next = append(n.next, mapNodeNext[K, V]{pointer: oldPointer, distance: newDistance})
	return n
}

func (s *SkipListMap[K, V]) insertNode(cur *mapNode[K, V], level int, res **mapNode[K, V], key K, value V) int {
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
			nxt.value = value
			*res = nil
			return -1
		}
	}

	if level > 0 {
		below := s.insertNode(cur, level-1, res, key, value)
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
		cur.next[level] = mapNodeNext[K, V]{
			pointer:  *res,
			distance: below,
		}
		return distance + below
	}

	oldNext := cur.next[level].pointer
	newNode := &mapNode[K, V]{
		key:   key,
		value: value,
	}
	newNode.next = []mapNodeNext[K, V]{{pointer: oldNext, distance: 0}}
	*res = newNode
	cur.next[level].pointer = newNode
	return distance
}

func (s *SkipListMap[K, V]) Insert(key K, value V) {
	var res *mapNode[K, V]
	dist := s.insertNode(s.head, s.height-1, &res, key, value)
	if res == nil {
		return
	}
	s.size++
	if dist < 0 || !testJump() {
		return
	}
	s.height++
	s.head.next = append(s.head.next, mapNodeNext[K, V]{pointer: res, distance: dist})
	res.next = append(res.next, mapNodeNext[K, V]{pointer: nil, distance: s.size - dist - 1})
}

func (s *SkipListMap[K, V]) Update(key K, value V) {
	if node := s.Find(key); node != nil {
		node.value = value
	} else {
		s.Insert(key, value)
	}
}

func (s *SkipListMap[K, V]) _erase(cur *mapNode[K, V], level int, key K) *mapNode[K, V] {
	for {
		nxt := cur.next[level].pointer
		if nxt == nil || s.comp(key, nxt.key) {
			break
		}
		cur = nxt
	}
	if level > 0 {
		res := s._erase(cur, level-1, key)
		if res != nil {
			if cur.next[level].pointer != res {
				cur.next[level].distance--
			} else {
				cur.next[level] = mapNodeNext[K, V]{
					pointer:  res.next[level].pointer,
					distance: cur.next[level].distance + res.next[level].distance,
				}
			}
			return res
		} else {
			return nil
		}
	} else {
		nxt := cur.next[level].pointer
		if nxt == nil || s.comp(key, nxt.key) {
			return nil
		} else {
			cur.next[level].pointer = nxt.next[level].pointer
			return nxt
		}
	}
}

func (s *SkipListMap[K, V]) Erase(key K) bool {
	res := s._erase(s.head, s.height-1, key)
	if res != nil {
		s.size--
		return true
	}
	return false
}

func (s *SkipListMap[K, V]) Rank(key K) int {
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

func (s *SkipListMap[K, V]) Kth(k int) *mapNode[K, V] {
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

func (s *SkipListMap[K, V]) Find(key K) *mapNode[K, V] {
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil {
			nxt := cur.next[level].pointer
			if s.comp(nxt.key, key) {
				cur = nxt
			} else if !s.comp(key, nxt.key) {
				return nxt
			} else {
				break
			}
		}
	}
	return nil
}

func (s *SkipListMap[K, V]) SmallerBound(key K) *mapNode[K, V] {
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil && s.comp(cur.next[level].pointer.key, key) {
			cur = cur.next[level].pointer
		}
	}
	return cur
}

func (s *SkipListMap[K, V]) LowerBound(key K) *mapNode[K, V] {
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
	ret := s.head.next[0].pointer
	if ret != nil && !s.comp(ret.key, key) {
		return ret
	}
	return nil
}

func (s *SkipListMap[K, V]) UpperBound(key K) *mapNode[K, V] {
	cur := s.head
	for level := s.height - 1; level >= 0; level-- {
		for cur.next[level].pointer != nil && !s.comp(key, cur.next[level].pointer.key) {
			cur = cur.next[level].pointer
		}
	}
	return cur.next[0].pointer
}

func (s *SkipListMap[K, V]) Size() int {
	return s.size
}

func (s *SkipListMap[K, V]) Empty() bool {
	return s.size == 0
}

func (s *SkipListMap[K, V]) Count(key K) int {
	if s.Find(key) != nil {
		return 1
	}
	return 0
}

func (s *SkipListMap[K, V]) Clear() {
	s.head = nil
	s.size = 0
	s.height = 0
}
