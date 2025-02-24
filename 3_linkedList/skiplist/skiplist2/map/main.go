// 跳表实现的有序 Map, ，支持以下操作：
//
// - NewSkipListMap(less func(a, b K) bool) *SkipListMap[K, V]：创建一个空的 SkipListMap，less 函数用来确定键的顺序.
// - Put(key K, value V) bool：插入一个键值对；如果键已存在则更新对应的值，返回是否为新插入.
// - Get(key K) (V, bool)：查找 key 对应的值，存在时返回 (value, true)，否则返回零值和 false.
// - Has(key K) bool：判断 key 是否存在.
// - Delete(key K) bool：移除 key 对应的节点，成功时返回 true.
// - Pop(index int) *Node[K, V]：根据下标弹出节点，支持负数索引（-1 表示最后一个）。
// - BisectLeft(key K) int：返回 key 的左侧下标（所有小于 key 的元素个数）.
// - BisectRight(key K) int：返回 key 的右侧下标（所有小于等于 key 的元素个数）.
// - Kth(k int) *Node[K, V]：返回排序后下标为 k 的节点.
// - Lower(key K) *Node[K, V]：返回 key 的严格前驱节点.
// - Floor(key K) *Node[K, V]：返回 key 的前驱节点（允许等于 key）.
// - Ceil(key K) *Node[K, V]：返回 key 的后继节点（允许等于 key）.
// - Upper(key K) *Node[K, V]：返回 key 的严格后继节点.
// - Range(low, high K, f func(node *Node[K, V]) (shouldBreak bool))：遍历键值在 [low, high] 范围内的节点.
// - Enumerate(start, end int, f func(node *Node[K, V]) bool)：遍历排序后下标在 [start, end) 范围内的所有节点.
// - Size() int：返回跳表中节点的个数.
// - Empty() bool：判断跳表是否为空.
// - ForEach(f func(node *Node[K, V]) (shouldBreak bool))：遍历所有节点.
// - String() string：返回跳表的字符串表示.

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"strings"
)

func main() {
	// 定义键的顺序，比如 int 类型按升序排列
	sl := NewSkipListMap[int, string](func(a, b int) bool { return a < b })

	// 插入键值对
	sl.Put(5, "five")
	sl.Put(2, "two")
	sl.Put(8, "eight")
	sl.Put(3, "three")

	// 更新已有键
	sl.Put(3, "THREE")

	// 输出跳表内容
	fmt.Println(sl.String())

	// 查找
	if v, ok := sl.Get(3); ok {
		fmt.Println("Key 3:", v)
	}

	// 删除
	sl.Delete(2)
	fmt.Println("After deleting key 2:", sl.String())
}

const (
	maxLevel = 20
)

type Node[K comparable, V any] struct {
	Key     K
	Value   V
	forward []*Node[K, V]
	span    []int // 第 i 层从当前结点到 forward[i] 之间跳过的节点个数
}

type SkipListMap[K comparable, V any] struct {
	head  *Node[K, V]
	level int
	size  int
	less  func(a, b K) bool
}

func NewSkipListMap[K comparable, V any](less func(a, b K) bool) *SkipListMap[K, V] {
	head := &Node[K, V]{
		forward: make([]*Node[K, V], maxLevel),
		span:    make([]int, maxLevel),
	}
	return &SkipListMap[K, V]{
		head:  head,
		level: 1,
		less:  less,
	}
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

// Put 插入一个键值对；如果键已存在则更新对应的值，返回是否为新插入.
func (sl *SkipListMap[K, V]) Put(key K, value V) bool {
	levelPrevs := make([]*Node[K, V], maxLevel)
	suffixSum := make([]int, maxLevel)
	cur := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			suffixSum[i] = 0
		} else {
			suffixSum[i] = suffixSum[i+1]
		}
		for cur.forward[i] != nil && sl.less(cur.forward[i].Key, key) {
			suffixSum[i] += cur.span[i]
			cur = cur.forward[i]
		}
		levelPrevs[i] = cur
	}

	candidate := cur.forward[0]
	if candidate != nil && candidate.Key == key {
		candidate.Value = value
		return false
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
	newNode := &Node[K, V]{
		Key:     key,
		Value:   value,
		forward: make([]*Node[K, V], newLevel),
		span:    make([]int, newLevel),
	}

	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = levelPrevs[i].forward[i]
		newNode.span[i] = levelPrevs[i].span[i] - (suffixSum[0] - suffixSum[i])
		levelPrevs[i].forward[i] = newNode
		levelPrevs[i].span[i] = (suffixSum[0] - suffixSum[i]) + 1
	}
	for i := newLevel; i < sl.level; i++ {
		levelPrevs[i].span[i] += 1
	}
	sl.size++
	return true
}

// findNode 查找 key 应该所在的位置，返回第一个大于等于 key 的节点.
func (sl *SkipListMap[K, V]) findNode(key K) *Node[K, V] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].Key, key) {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

// Get 查找 key 对应的值，存在时返回 (value, true)，否则返回零值和 false.
func (sl *SkipListMap[K, V]) Get(key K) (V, bool) {
	candidate := sl.findNode(key)
	if candidate != nil && candidate.Key == key {
		return candidate.Value, true
	}
	var zero V
	return zero, false
}

func (sl *SkipListMap[K, V]) Has(key K) bool {
	candidate := sl.findNode(key)
	return candidate != nil && candidate.Key == key
}

// Delete 移除 key 对应的节点，成功时返回 true.
func (sl *SkipListMap[K, V]) Delete(key K) bool {
	levelPrevs := make([]*Node[K, V], sl.level)
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].Key, key) {
			cur = cur.forward[i]
		}
		levelPrevs[i] = cur
	}
	target := cur.forward[0]
	if target == nil || target.Key != key {
		return false
	}

	for i := 0; i < sl.level; i++ {
		if levelPrevs[i].forward[i] == target {
			levelPrevs[i].span[i] += target.span[i] - 1
			levelPrevs[i].forward[i] = target.forward[i]
		} else {
			levelPrevs[i].span[i] -= 1
		}
	}
	sl.size--
	for sl.level > 1 && sl.head.forward[sl.level-1] == nil {
		sl.level--
	}
	return true
}

// Pop 根据下标弹出节点，支持负数索引（-1 表示最后一个）。
func (sl *SkipListMap[K, V]) Pop(index int) *Node[K, V] {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	levelPrevs := make([]*Node[K, V], sl.level)
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
	for i := 0; i < sl.level; i++ {
		if levelPrevs[i].forward[i] == target {
			levelPrevs[i].span[i] += target.span[i] - 1
			levelPrevs[i].forward[i] = target.forward[i]
		} else {
			levelPrevs[i].span[i] -= 1
		}
	}
	sl.size--
	for sl.level > 1 && sl.head.forward[sl.level-1] == nil {
		sl.level--
	}
	return target
}

// BisectLeft 返回 key 的左侧下标（所有小于 key 的元素个数）.
func (sl *SkipListMap[K, V]) BisectLeft(key K) int {
	rank := 0
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].Key, key) {
			rank += cur.span[i]
			cur = cur.forward[i]
		}
	}
	return rank
}

// BisectRight 返回 key 的右侧下标（所有小于等于 key 的元素个数）.
func (sl *SkipListMap[K, V]) BisectRight(key K) int {
	rank := 0
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && !sl.less(key, cur.forward[i].Key) {
			rank += cur.span[i]
			cur = cur.forward[i]
		}
	}
	return rank
}

// Kth 返回排序后下标为 k 的节点.
func (sl *SkipListMap[K, V]) Kth(k int) *Node[K, V] {
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

// Lower 返回 key 的严格前驱节点.
func (sl *SkipListMap[K, V]) Lower(key K) *Node[K, V] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].Key, key) {
			cur = cur.forward[i]
		}
	}
	if cur == sl.head {
		return nil
	}
	return cur
}

// Floor 返回 key 的前驱节点（允许等于 key）.
func (sl *SkipListMap[K, V]) Floor(key K) *Node[K, V] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && !sl.less(key, cur.forward[i].Key) {
			cur = cur.forward[i]
		}
	}
	if cur == sl.head {
		return nil
	}
	return cur
}

// Ceil 返回 key 的后继节点（允许等于 key）.
func (sl *SkipListMap[K, V]) Ceil(key K) *Node[K, V] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].Key, key) {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

// Upper 返回 key 的严格后继节点.
func (sl *SkipListMap[K, V]) Upper(key K) *Node[K, V] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && !sl.less(key, cur.forward[i].Key) {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

// Range 遍历键值在 [low, high] 范围内的节点.
func (sl *SkipListMap[K, V]) Range(low, high K, f func(node *Node[K, V]) (shouldBreak bool)) {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].Key, low) {
			cur = cur.forward[i]
		}
	}
	cur = cur.forward[0]
	for cur != nil {
		if sl.less(high, cur.Key) {
			break
		}
		if f(cur) {
			break
		}
		cur = cur.forward[0]
	}
}

// Enumerate 遍历排序后下标在 [start, end) 范围内的所有节点.
func (sl *SkipListMap[K, V]) Enumerate(start, end int, f func(node *Node[K, V]) bool) {
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
	count := start
	for node != nil && count < end {
		if f(node) {
			return
		}
		count++
		node = node.forward[0]
	}
}

func (sl *SkipListMap[K, V]) Size() int {
	return sl.size
}

func (sl *SkipListMap[K, V]) Empty() bool {
	return sl.size == 0
}

func (sl *SkipListMap[K, V]) ForEach(f func(node *Node[K, V]) (shouldBreak bool)) {
	cur := sl.head.forward[0]
	for cur != nil {
		if f(cur) {
			break
		}
		cur = cur.forward[0]
	}
}

func (sl *SkipListMap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString("SkipListMap[")
	cur := sl.head.forward[0]
	for cur != nil {
		sb.WriteString(fmt.Sprintf("(%v:%v)", cur.Key, cur.Value))
		if cur.forward[0] != nil {
			sb.WriteString("->")
		}
		cur = cur.forward[0]
	}
	sb.WriteString("]")
	return sb.String()
}
