package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"time"
)

// https://leetcode.cn/problems/design-skiplist/description/
type Skiplist struct {
	list *SkipList[int]
}

func Constructor() Skiplist {
	return Skiplist{NewSkipList[int](func(a, b int) bool { return a < b })}
}

func (this *Skiplist) Search(target int) bool {
	return this.list.Find(target) != nil
}

func (this *Skiplist) Add(num int) {
	this.list.Insert(num)
}

func (this *Skiplist) Erase(num int) bool {
	return this.list.Erase(num)
}

func main() {
	sl := NewSkipList[int](func(a, b int) bool { return a < b })

	sl.Insert(1)
	sl.Insert(3)
	sl.Insert(2)

	fmt.Println(sl.BisectLeft(1))
	fmt.Println(sl.BisectLeft(2))
	fmt.Println(sl.BisectRight(2))
	fmt.Println(sl.BisectLeft(3))
	fmt.Println(sl.BisectLeft(4))

	fmt.Println(sl.Kth(0))
	fmt.Println(sl.Kth(1))
	fmt.Println(sl.Kth(2))
	fmt.Println(sl.Kth(3))

	fmt.Println(sl.Find(2).key)

	sl.Erase(2)
	fmt.Println(sl.Lower(1))
	fmt.Println(sl.Floor(1))

	fmt.Println(sl.Ceil(3))
	fmt.Println(sl.Upper(3))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	maxLevel = 20
)

type Node[T comparable] struct {
	key     T
	weight  int
	forward []*Node[T]
	span    []int // 第 i 层从当前结点到 forward[i] 之间跳过的元素个数
}

type SkipList[T comparable] struct {
	head  *Node[T]
	level int
	size  int
	less  func(a, b T) bool
}

func NewSkipList[T comparable](less func(a, b T) bool) *SkipList[T] {
	head := &Node[T]{
		forward: make([]*Node[T], maxLevel),
		span:    make([]int, maxLevel),
	}
	return &SkipList[T]{
		head:  head,
		level: 1,
		less:  less,
	}
}

func randomLevel(maxLevel int) int {
	var x uint64 = rand.Uint64() & ((1 << uint(maxLevel-1)) - 1)
	zeroes := bits.TrailingZeros64(x)
	if zeroes > maxLevel {
		return maxLevel
	}
	if zeroes == 0 {
		return 1
	}
	return zeroes
}

func (sl *SkipList[T]) Insert(key T) {
	update := make([]*Node[T], maxLevel)
	rank := make([]int, maxLevel)
	cur := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for cur.forward[i] != nil && sl.less(cur.forward[i].key, key) {
			rank[i] += cur.span[i]
			cur = cur.forward[i]
		}
		update[i] = cur
	}

	candidate := cur.forward[0]
	if candidate != nil && candidate.key == key {
		candidate.weight++
		for i := 0; i < sl.level; i++ {
			if update[i].forward[i] == candidate {
				update[i].span[i]++
			}
		}
		sl.size++
		return
	}

	newLevel := randomLevel(maxLevel)
	if newLevel > sl.level {
		for i := sl.level; i < newLevel; i++ {
			rank[i] = 0
			update[i] = sl.head
			update[i].span[i] = sl.size
		}
		sl.level = newLevel
	}
	newNode := &Node[T]{
		key:     key,
		weight:  1,
		forward: make([]*Node[T], newLevel),
		span:    make([]int, newLevel),
	}
	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		newNode.span[i] = update[i].span[i] - (rank[0] - rank[i])
		update[i].forward[i] = newNode
		update[i].span[i] = (rank[0] - rank[i]) + newNode.weight
	}
	for i := newLevel; i < sl.level; i++ {
		update[i].span[i] += newNode.weight
	}
	sl.size += newNode.weight
}

func (sl *SkipList[T]) Find(key T) *Node[T] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].key, key) {
			cur = cur.forward[i]
		}
	}
	candidate := cur.forward[0]
	if candidate != nil && candidate.key == key {
		return candidate
	}
	return nil
}

func (sl *SkipList[T]) Count(key T) int {
	node := sl.Find(key)
	if node != nil {
		return node.weight
	}
	return 0
}

func (sl *SkipList[T]) Erase(key T) bool {
	update := make([]*Node[T], sl.level)
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].key, key) {
			cur = cur.forward[i]
		}
		update[i] = cur
	}
	target := cur.forward[0]
	if target == nil || target.key != key {
		return false
	}

	if target.weight > 1 {
		target.weight--
		for i := 0; i < sl.level; i++ {
			if update[i].forward[i] == target {
				update[i].span[i]--
			}
		}
		sl.size--
		return true
	}

	for i := 0; i < sl.level; i++ {
		if update[i].forward[i] == target {
			update[i].span[i] += target.span[i] - target.weight
			update[i].forward[i] = target.forward[i]
		} else {
			update[i].span[i] -= target.weight
		}
	}
	sl.size--
	for sl.level > 1 && sl.head.forward[sl.level-1] == nil {
		sl.level--
	}
	return true
}

func (sl *SkipList[T]) BisectLeft(key T) int {
	rank := 0
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].key, key) {
			rank += cur.span[i]
			cur = cur.forward[i]
		}
	}
	return rank
}

func (sl *SkipList[T]) BisectRight(key T) int {
	rank := 0
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && !sl.less(key, cur.forward[i].key) {
			rank += cur.span[i]
			cur = cur.forward[i]
		}
	}
	return rank
}

func (sl *SkipList[T]) Kth(k int) *Node[T] {
	if k < 0 || k >= sl.size {
		return nil
	}
	cur := sl.head
	traversed := 0
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && traversed+cur.span[i] <= k {
			traversed += cur.span[i]
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

func (sl *SkipList[T]) Lower(key T) *Node[T] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].key, key) {
			cur = cur.forward[i]
		}
	}
	if cur == sl.head {
		return nil
	}
	return cur
}

func (sl *SkipList[T]) Floor(key T) *Node[T] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && !sl.less(key, cur.forward[i].key) {
			cur = cur.forward[i]
		}
	}
	if cur == sl.head {
		return nil
	}
	return cur
}

func (sl *SkipList[T]) Ceil(key T) *Node[T] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].key, key) {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

func (sl *SkipList[T]) Upper(key T) *Node[T] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && !sl.less(key, cur.forward[i].key) {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

func (sl *SkipList[T]) Size() int {
	return sl.size
}

func (sl *SkipList[T]) Empty() bool {
	return sl.size == 0
}
