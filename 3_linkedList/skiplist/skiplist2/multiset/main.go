package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"sort"
	"strings"
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
	sl.Insert(3)

	fmt.Println(sl)
	sl.Range(3, 3, func(node *Node[int]) bool {
		fmt.Println(node.key, node.weight)
		return false
	})

	fmt.Println("---------")
	sl.Enumerate(3, 4, func(node *Node[int]) bool {
		fmt.Println(node.key, node.weight)
		return false
	})

	sl = NewSkipListFrom[int](func(a, b int) bool { return a < b }, 1, 3, 2, 3)
	fmt.Println(sl)
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

func NewSkipListFrom[T comparable](less func(a, b T) bool, elements ...T) *SkipList[T] {
	if len(elements) == 0 {
		return NewSkipList(less)
	}
	elements = append(elements[:0:0], elements...)
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })

	type info[T any] struct {
		key    T
		weight int
		level  int
	}
	var unique []info[T]
	{
		current := elements[0]
		count := 1
		for i := 1; i < len(elements); i++ {
			if elements[i] == current {
				count++
			} else {
				unique = append(unique, info[T]{key: current, weight: count})
				current = elements[i]
				count = 1
			}
		}
		unique = append(unique, info[T]{key: current, weight: count})
	}

	n := len(unique)
	maxLevelUsed := 1
	for i := 0; i < n; i++ {
		newLevel := randomLevel(maxLevel)
		unique[i].level = newLevel
		if newLevel > maxLevelUsed {
			maxLevelUsed = newLevel
		}
	}

	prefix := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + unique[i].weight
	}
	total := prefix[n]

	nodes := make([]*Node[T], n)
	for i := 0; i < n; i++ {
		level := unique[i].level
		nodes[i] = &Node[T]{
			key:     unique[i].key,
			weight:  unique[i].weight,
			forward: make([]*Node[T], level),
			span:    make([]int, level),
		}
	}

	head := &Node[T]{
		forward: make([]*Node[T], maxLevelUsed),
		span:    make([]int, maxLevelUsed),
	}
	sl := &SkipList[T]{
		head:  head,
		level: maxLevelUsed,
		size:  total,
		less:  less,
	}

	levelNodes := make([][]int, maxLevelUsed)
	for i := 0; i < n; i++ {
		level := unique[i].level
		for l := 0; l < level; l++ {
			levelNodes[l] = append(levelNodes[l], i)
		}
	}
	for L := 0; L < maxLevelUsed; L++ {
		chain := levelNodes[L]
		if len(chain) == 0 {
			head.forward[L] = nil
			head.span[L] = total
		} else {
			firstIdx := chain[0]
			head.forward[L] = nodes[firstIdx]
			head.span[L] = prefix[firstIdx] + nodes[firstIdx].weight
			for j := 0; j < len(chain)-1; j++ {
				iIdx := chain[j]
				jIdx := chain[j+1]
				nodes[iIdx].forward[L] = nodes[jIdx]
				gap := prefix[jIdx] - (prefix[iIdx] + nodes[iIdx].weight)
				nodes[iIdx].span[L] = gap + nodes[jIdx].weight
			}
			lastIdx := chain[len(chain)-1]
			nodes[lastIdx].forward[L] = nil
			nodes[lastIdx].span[L] = total - (prefix[lastIdx] + nodes[lastIdx].weight)
		}
	}

	return sl
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

// Pop 删除排序后下标 index 处的元素（按 0-indexed 计），
func (sl *SkipList[T]) Pop(index int) *Node[T] {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	update := make([]*Node[T], sl.level)
	cur := sl.head
	traversed := 0
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && traversed+cur.span[i] <= index {
			traversed += cur.span[i]
			cur = cur.forward[i]
		}
		update[i] = cur
	}

	target := cur.forward[0]
	if target.weight > 1 {
		target.weight--
		for i := 0; i < sl.level; i++ {
			if update[i].forward[i] == target {
				update[i].span[i]--
			}
		}
		sl.size--
		return target
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
	return target
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

// Range 遍历 key 值在 [low, high] 范围内的结点，
func (sl *SkipList[T]) Range(low, high T, f func(node *Node[T]) (shouldBreak bool)) {
	cur := sl.head
	// 找到第一个 key >= low 的结点
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].key, low) {
			cur = cur.forward[i]
		}
	}
	cur = cur.forward[0]
	for cur != nil {
		if sl.less(high, cur.key) {
			break
		}
		if f(cur) {
			break
		}
		cur = cur.forward[0]
	}
}

// Enumerate 遍历排序后下标在 [start, end) 范围内的所有出现。
func (sl *SkipList[T]) Enumerate(start, end int, f func(node *Node[T]) bool) {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return
	}

	cur := sl.head
	traversed := 0
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && traversed+cur.span[i] <= start {
			traversed += cur.span[i]
			cur = cur.forward[i]
		}
	}

	node := cur.forward[0]
	if node == nil {
		return
	}

	offset := start - traversed
	count := start

	for i := offset; i < node.weight && count < end; i++ {
		if f(node) {
			return
		}
		count++
	}

	node = node.forward[0]
	for node != nil && count < end {
		for i := 0; i < node.weight && count < end; i++ {
			if f(node) {
				return
			}
			count++
		}
		node = node.forward[0]
	}
}

func (sl *SkipList[T]) Size() int {
	return sl.size
}

func (sl *SkipList[T]) Empty() bool {
	return sl.size == 0
}

func (sl *SkipList[T]) ForEach(f func(node *Node[T]) (shouldBreak bool)) {
	cur := sl.head.forward[0]
	for cur != nil {
		if f(cur) {
			break
		}
		cur = cur.forward[0]
	}
}

func (sl *SkipList[T]) String() string {
	var sb strings.Builder
	sb.WriteString("SkipList[")
	cur := sl.head.forward[0]
	for cur != nil {
		sb.WriteString(fmt.Sprintf("%v:%v", cur.key, cur.weight))
		if cur.forward[0] != nil {
			sb.WriteString("->")
		}
		cur = cur.forward[0]
	}
	sb.WriteString("]")
	return sb.String()
}
