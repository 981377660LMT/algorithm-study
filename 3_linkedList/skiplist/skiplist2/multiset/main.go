// 跳表实现的MultiSet，支持以下操作：
//
//  - NewSkipList[T](less func(a, b T) bool) *SkipList[T]：构造一个新的跳表，less 为比较函数
//  - NewSkipListFrom[T](less func(a, b T) bool, elements ...T) *SkipList[T]：从元素列表构造一个新的跳表
//  - sl.Add(key T)：向跳表中添加元素 key
//  - sl.Find(key T) *Node[T]：查找元素 key，返回对应的结点
//  - sl.Count(key T) int：返回元素 key 的个数
//  - sl.Discard(key T) bool：删除元素 key，返回是否删除成功
//  - sl.Pop(index int) *Node[T]：删除排序后下标 index 处的元素（按 0-indexed 计）
//  - sl.BisectLeft(key T) int：返回 key 在排序后的下标（左侧边界）
//  - sl.BisectRight(key T) int：返回 key 在排序后的下标（右侧边界）
//  - sl.Kth(k int) *Node[T]：返回排序后下标为 k 的元素
//  - sl.Lower(key T) *Node[T]：返回小于 key 的最大元素
//  - sl.Floor(key T) *Node[T]：返回小于等于 key 的最大元素
//  - sl.Ceil(key T) *Node[T]：返回大于等于 key 的最小元素
//  - sl.Upper(key T) *Node[T]：返回大于 key 的最小元素
//  - sl.Range(low, high T, f func(node *Node[T]) (shouldBreak bool))：遍历 key 值在 [low, high] 范围内的结点
//  - sl.Enumerate(start, end int, f func(node *Node[T]) bool)：遍历排序后下标在 [start, end) 范围内的所有结点
//  - sl.Size() int：返回跳表的元素个数
//  - sl.Empty() bool：返回跳表是否为空
//  - sl.ForEach(f func(node *Node[T]) (shouldBreak bool))：遍历跳表中的所有结点
//  - sl.String() string：返回跳表的字符串表示

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"sort"
	"strings"
)

// https://leetcode.cn/problems/design-skiplist/description/
type Skiplist struct {
	list *SkipListMultiSet[int]
}

func Constructor() Skiplist {
	return Skiplist{NewSkipListMultiSet[int](func(a, b int) bool { return a < b })}
}

func (this *Skiplist) Search(target int) bool {
	return this.list.Find(target) != nil
}

func (this *Skiplist) Add(num int) {
	this.list.Add(num)
}

func (this *Skiplist) Erase(num int) bool {
	return this.list.Discard(num)
}

func main() {
	sl := NewSkipListMultiSet[int](func(a, b int) bool { return a < b })
	sl.Add(1)
	sl.Add(3)
	sl.Add(2)
	sl.Add(3)

	fmt.Println(sl)
	sl.Range(3, 3, func(node *Node[int]) bool {
		fmt.Println(node.key, node.Weight)
		return false
	})

	fmt.Println("---------")
	sl.Enumerate(3, 4, func(node *Node[int]) bool {
		fmt.Println(node.key, node.Weight)
		return false
	})

	sl = NewSkipListMultiSetFrom[int](func(a, b int) bool { return a < b }, 1, 3, 2, 3)
	fmt.Println(sl)
}

const (
	maxLevel = 20
)

type Node[T comparable] struct {
	key     T
	Weight  int
	forward []*Node[T]
	span    []int // 第 i 层从当前结点到 forward[i] 之间跳过的元素个数
}

type SkipListMultiSet[T comparable] struct {
	head  *Node[T]
	level int
	size  int
	less  func(a, b T) bool
}

func NewSkipListMultiSet[T comparable](less func(a, b T) bool) *SkipListMultiSet[T] {
	head := &Node[T]{
		forward: make([]*Node[T], maxLevel),
		span:    make([]int, maxLevel),
	}
	return &SkipListMultiSet[T]{
		head:  head,
		level: 1,
		less:  less,
	}
}

func NewSkipListMultiSetFrom[T comparable](less func(a, b T) bool, elements ...T) *SkipListMultiSet[T] {
	if len(elements) == 0 {
		return NewSkipListMultiSet(less)
	}
	elements = append(elements[:0:0], elements...)
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })

	type info[T any] struct {
		key    T
		weight int
		level  int
	}
	var unique []info[T]
	maxLevelUsed := 1
	{
		current := elements[0]
		count := 1
		for i := 1; i < len(elements); i++ {
			if elements[i] == current {
				count++
			} else {
				level := randomLevel(maxLevel)
				maxLevelUsed = max(maxLevelUsed, level)
				unique = append(unique, info[T]{key: current, weight: count, level: level})
				current = elements[i]
				count = 1
			}
		}

		level := randomLevel(maxLevel)
		maxLevelUsed = max(maxLevelUsed, level)
		unique = append(unique, info[T]{key: current, weight: count, level: level})
	}

	n := len(unique)
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
			Weight:  unique[i].weight,
			forward: make([]*Node[T], level),
			span:    make([]int, level),
		}
	}

	head := &Node[T]{
		forward: make([]*Node[T], maxLevel),
		span:    make([]int, maxLevel),
	}
	sl := &SkipListMultiSet[T]{
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
			head.span[L] = prefix[firstIdx] + nodes[firstIdx].Weight
			for j := 0; j < len(chain)-1; j++ {
				iIdx := chain[j]
				jIdx := chain[j+1]
				nodes[iIdx].forward[L] = nodes[jIdx]
				gap := prefix[jIdx] - (prefix[iIdx] + nodes[iIdx].Weight)
				nodes[iIdx].span[L] = gap + nodes[jIdx].Weight
			}
			lastIdx := chain[len(chain)-1]
			nodes[lastIdx].forward[L] = nil
			nodes[lastIdx].span[L] = total - (prefix[lastIdx] + nodes[lastIdx].Weight)
		}
	}

	return sl
}

func randomLevel(maxLevel int) int {
	var x uint64 = rand.Uint64() & ((1 << uint(maxLevel-1)) - 1)
	zeroes := bits.TrailingZeros64(x)
	if zeroes >= maxLevel {
		return maxLevel - 1
	}
	if zeroes == 0 {
		return 1
	}
	return zeroes
}

func (sl *SkipListMultiSet[T]) Add(key T) {
	levelPrevs := make([]*Node[T], maxLevel)
	suffixSum := make([]int, maxLevel)
	cur := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			suffixSum[i] = 0
		} else {
			suffixSum[i] = suffixSum[i+1]
		}
		for cur.forward[i] != nil && sl.less(cur.forward[i].key, key) {
			suffixSum[i] += cur.span[i]
			cur = cur.forward[i]
		}
		levelPrevs[i] = cur
	}

	candidate := cur.forward[0]
	if candidate != nil && candidate.key == key {
		candidate.Weight++
		for i := 0; i < sl.level; i++ {
			if levelPrevs[i].forward[i] == candidate {
				levelPrevs[i].span[i]++
			}
		}
		sl.size++
		return
	}

	newLevel := randomLevel(maxLevel)
	if newLevel > sl.level {
		for i := sl.level; i < newLevel; i++ {
			suffixSum[i] = 0
			levelPrevs[i] = sl.head
			levelPrevs[i].span[i] = sl.size
		}
		sl.level = newLevel
	}
	newNode := &Node[T]{
		key:     key,
		Weight:  1,
		forward: make([]*Node[T], newLevel),
		span:    make([]int, newLevel),
	}
	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = levelPrevs[i].forward[i]
		newNode.span[i] = levelPrevs[i].span[i] - (suffixSum[0] - suffixSum[i])
		levelPrevs[i].forward[i] = newNode
		levelPrevs[i].span[i] = (suffixSum[0] - suffixSum[i]) + newNode.Weight
	}
	for i := newLevel; i < sl.level; i++ {
		levelPrevs[i].span[i] += newNode.Weight
	}
	sl.size += newNode.Weight
}

func (sl *SkipListMultiSet[T]) Find(key T) *Node[T] {
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

func (sl *SkipListMultiSet[T]) Count(key T) int {
	node := sl.Find(key)
	if node != nil {
		return node.Weight
	}
	return 0
}

func (sl *SkipListMultiSet[T]) Discard(key T) bool {
	levelPrevs := make([]*Node[T], sl.level)
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].key, key) {
			cur = cur.forward[i]
		}
		levelPrevs[i] = cur
	}
	target := cur.forward[0]
	if target == nil || target.key != key {
		return false
	}

	if target.Weight > 1 {
		target.Weight--
		for i := 0; i < sl.level; i++ {
			if levelPrevs[i].forward[i] == target {
				levelPrevs[i].span[i]--
			}
		}
		sl.size--
		return true
	}

	for i := 0; i < sl.level; i++ {
		if levelPrevs[i].forward[i] == target {
			levelPrevs[i].span[i] += target.span[i] - target.Weight
			levelPrevs[i].forward[i] = target.forward[i]
		} else {
			levelPrevs[i].span[i] -= target.Weight
		}
	}
	sl.size--
	for sl.level > 1 && sl.head.forward[sl.level-1] == nil {
		sl.level--
	}
	return true
}

// Pop 删除排序后下标 index 处的元素（按 0-indexed 计），
func (sl *SkipListMultiSet[T]) Pop(index int) *Node[T] {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	levelPrevs := make([]*Node[T], sl.level)
	cur := sl.head
	traversed := 0
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && traversed+cur.span[i] <= index {
			traversed += cur.span[i]
			cur = cur.forward[i]
		}
		levelPrevs[i] = cur
	}

	target := cur.forward[0]
	if target.Weight > 1 {
		target.Weight--
		for i := 0; i < sl.level; i++ {
			if levelPrevs[i].forward[i] == target {
				levelPrevs[i].span[i]--
			}
		}
		sl.size--
		return target
	}
	for i := 0; i < sl.level; i++ {
		if levelPrevs[i].forward[i] == target {
			levelPrevs[i].span[i] += target.span[i] - target.Weight
			levelPrevs[i].forward[i] = target.forward[i]
		} else {
			levelPrevs[i].span[i] -= target.Weight
		}
	}
	sl.size--
	for sl.level > 1 && sl.head.forward[sl.level-1] == nil {
		sl.level--
	}
	return target
}

func (sl *SkipListMultiSet[T]) BisectLeft(key T) int {
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

func (sl *SkipListMultiSet[T]) BisectRight(key T) int {
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

func (sl *SkipListMultiSet[T]) Kth(k int) *Node[T] {
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

func (sl *SkipListMultiSet[T]) Lower(key T) *Node[T] {
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

func (sl *SkipListMultiSet[T]) Floor(key T) *Node[T] {
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

func (sl *SkipListMultiSet[T]) Ceil(key T) *Node[T] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].key, key) {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

func (sl *SkipListMultiSet[T]) Upper(key T) *Node[T] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && !sl.less(key, cur.forward[i].key) {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

// Range 遍历 key 值在 [low, high] 范围内的结点，
func (sl *SkipListMultiSet[T]) Range(low, high T, f func(node *Node[T]) (shouldBreak bool)) {
	cur := sl.head
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

// Enumerate 遍历排序后下标在 [start, end) 范围内的所有结点，
func (sl *SkipListMultiSet[T]) Enumerate(start, end int, f func(node *Node[T]) bool) {
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
	for i := offset; i < node.Weight && count < end; i++ {
		if f(node) {
			return
		}
		count++
	}

	node = node.forward[0]
	for node != nil && count < end {
		for i := 0; i < node.Weight && count < end; i++ {
			if f(node) {
				return
			}
			count++
		}
		node = node.forward[0]
	}
}

func (sl *SkipListMultiSet[T]) Size() int {
	return sl.size
}

func (sl *SkipListMultiSet[T]) Empty() bool {
	return sl.size == 0
}

func (sl *SkipListMultiSet[T]) ForEach(f func(node *Node[T]) (shouldBreak bool)) {
	cur := sl.head.forward[0]
	for cur != nil {
		if f(cur) {
			break
		}
		cur = cur.forward[0]
	}
}

func (sl *SkipListMultiSet[T]) String() string {
	var sb strings.Builder
	sb.WriteString("SkipList[")
	cur := sl.head.forward[0]
	for cur != nil {
		sb.WriteString(fmt.Sprintf("%v:%v", cur.key, cur.Weight))
		if cur.forward[0] != nil {
			sb.WriteString("->")
		}
		cur = cur.forward[0]
	}
	sb.WriteString("]")
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
