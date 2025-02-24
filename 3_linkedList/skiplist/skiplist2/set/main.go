// 跳表实现的Set，支持以下操作：
//
//  - NewSkipList[T](less func(a, b T) bool) *SkipList[T]：构造一个新的跳表，less 为比较函数
//  - NewSkipListFrom[T](less func(a, b T) bool, elements ...T) *SkipList[T]：从元素列表构造一个新的跳表
//  - sl.Add(key T) bool：向跳表中添加元素 key
//  - sl.Find(key T) *Node[T]：查找元素 key，返回对应的结点
//  - sl.Has(key T) bool：判断元素 key 是否存在
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
	"bufio"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"sort"
	"strings"
)

func main() {
	yosupo()
	// demo()
	// testRandomLevel()
}

const INF int = 1e18

func yosupo() {
	const eof = 0
	in := os.Stdin
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	_i, _n, buf := 0, 0, make([]byte, 1<<12)

	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return eof
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}

	NextByte := func() byte {
		b := rc()
		for ; '0' > b; b = rc() {
		}
		return b
	}
	_ = NextByte

	NextInt := func() (x int) {
		neg := false
		b := rc()
		for ; '0' > b || b > '9'; b = rc() {
			if b == eof {
				return
			}
			if b == '-' {
				neg = true
			}
		}
		for ; '0' <= b && b <= '9'; b = rc() {
			x = x*10 + int(b&15)
		}
		if neg {
			return -x
		}
		return
	}
	_ = NextInt

	n, q := NextInt(), NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = NextInt()
	}

	sl := NewSkipListSetFrom(func(a, b int) bool { return a < b }, nums...)
	for i := 0; i < q; i++ {
		op, x := NextInt(), NextInt()
		switch op {
		case 0:
			sl.Add(x)
		case 1:
			sl.Discard(x)
		case 2:
			if x > sl.Size() {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, sl.Kth(x-1).Key)
			}
		case 3:
			fmt.Fprintln(out, sl.BisectRight(x))
		case 4:
			if res := sl.Floor(x); res != nil {
				fmt.Fprintln(out, res.Key)
			} else {
				fmt.Fprintln(out, -1)
			}
		case 5:
			if res := sl.Ceil(x); res != nil {
				fmt.Fprintln(out, res.Key)
			} else {
				fmt.Fprintln(out, -1)
			}
		}
	}
}

func demo() {
	sl := NewSkipListSet[int](func(a, b int) bool { return a < b })
	sl.Add(1)
	sl.Add(3)
	sl.Add(2)
	sl.Add(3)

	fmt.Println(sl)
	sl.Range(3, 3, func(node *Node[int]) bool {
		fmt.Println(node.Key)
		return false
	})

	fmt.Println("---------")
	sl.Enumerate(0, 3, func(node *Node[int]) bool {
		fmt.Println(node.Key)
		return false
	})

	sl = NewSkipListSetFrom[int](func(a, b int) bool { return a < b }, 1, 3, 2, 3)
	fmt.Println(sl)
}

func testRandomLevel() {
	counter := make(map[int]int)
	for i := 0; i < 200000; i++ {
		level := randomLevel(20)
		counter[level]++
	}
	fmt.Println(counter)
}

const (
	maxLevel = 20
)

type Node[T comparable] struct {
	Key     T
	forward []*Node[T]
	span    []int // 第 i 层从当前结点到 forward[i] 之间跳过的元素个数
}

type SkipListSet[T comparable] struct {
	head  *Node[T]
	level int
	size  int
	less  func(a, b T) bool
}

func NewSkipListSet[T comparable](less func(a, b T) bool) *SkipListSet[T] {
	head := &Node[T]{
		forward: make([]*Node[T], maxLevel),
		span:    make([]int, maxLevel),
	}
	return &SkipListSet[T]{
		head:  head,
		level: 1,
		less:  less,
	}
}

// Type declarations inside generic functions are not currently supported.
type info[T any] struct {
	key   T
	level int
}

func NewSkipListSetFrom[T comparable](less func(a, b T) bool, elements ...T) *SkipListSet[T] {
	if len(elements) == 0 {
		return NewSkipListSet(less)
	}
	elements = append(elements[:0:0], elements...)
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })

	var unique []info[T]
	maxLevelUsed := 1
	{
		current := elements[0]
		for i := 1; i < len(elements); i++ {
			if elements[i] != current {
				level := randomLevel(maxLevel)
				maxLevelUsed = max(maxLevelUsed, level)
				unique = append(unique, info[T]{key: current, level: level})
				current = elements[i]
			}
		}
		level := randomLevel(maxLevel)
		maxLevelUsed = max(maxLevelUsed, level)
		unique = append(unique, info[T]{key: current, level: level})
	}

	n := len(unique)
	total := n
	nodes := make([]*Node[T], n)
	for i := 0; i < n; i++ {
		level := unique[i].level
		nodes[i] = &Node[T]{
			Key:     unique[i].key,
			forward: make([]*Node[T], level),
			span:    make([]int, level),
		}
	}

	head := &Node[T]{
		forward: make([]*Node[T], maxLevel),
		span:    make([]int, maxLevel),
	}
	sl := &SkipListSet[T]{
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
			head.span[L] = firstIdx + 1
			for j := 0; j < len(chain)-1; j++ {
				iIdx := chain[j]
				jIdx := chain[j+1]
				nodes[iIdx].forward[L] = nodes[jIdx]
				gap := jIdx - (iIdx + 1)
				nodes[iIdx].span[L] = gap + 1
			}
			lastIdx := chain[len(chain)-1]
			nodes[lastIdx].forward[L] = nil
			nodes[lastIdx].span[L] = total - (lastIdx + 1)
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

func (sl *SkipListSet[T]) Add(key T) bool {
	levelPrevs := make([]*Node[T], maxLevel)
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
	newNode := &Node[T]{
		Key:     key,
		forward: make([]*Node[T], newLevel),
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

func (sl *SkipListSet[T]) Find(key T) *Node[T] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].Key, key) {
			cur = cur.forward[i]
		}
	}
	candidate := cur.forward[0]
	if candidate != nil && candidate.Key == key {
		return candidate
	}
	return nil
}

func (sl *SkipListSet[T]) Has(key T) bool {
	return sl.Find(key) != nil
}

func (sl *SkipListSet[T]) Discard(key T) bool {
	levelPrevs := make([]*Node[T], sl.level)
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

func (sl *SkipListSet[T]) Pop(index int) *Node[T] {
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

func (sl *SkipListSet[T]) BisectLeft(key T) int {
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

func (sl *SkipListSet[T]) BisectRight(key T) int {
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

func (sl *SkipListSet[T]) Kth(k int) *Node[T] {
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

func (sl *SkipListSet[T]) Lower(key T) *Node[T] {
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

func (sl *SkipListSet[T]) Floor(key T) *Node[T] {
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

func (sl *SkipListSet[T]) Ceil(key T) *Node[T] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && sl.less(cur.forward[i].Key, key) {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

func (sl *SkipListSet[T]) Upper(key T) *Node[T] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && !sl.less(key, cur.forward[i].Key) {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

// Range 遍历 key 值在 [low, high] 范围内的结点.
func (sl *SkipListSet[T]) Range(low, high T, f func(node *Node[T]) (shouldBreak bool)) {
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

// Enumerate 遍历排序后下标在 [start, end) 范围内的所有结点
func (sl *SkipListSet[T]) Enumerate(start, end int, f func(node *Node[T]) bool) {
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

func (sl *SkipListSet[T]) Size() int {
	return sl.size
}

func (sl *SkipListSet[T]) Empty() bool {
	return sl.size == 0
}

func (sl *SkipListSet[T]) ForEach(f func(node *Node[T]) (shouldBreak bool)) {
	cur := sl.head.forward[0]
	for cur != nil {
		if f(cur) {
			break
		}
		cur = cur.forward[0]
	}
}

func (sl *SkipListSet[T]) String() string {
	var sb strings.Builder
	sb.WriteString("SkipList[")
	cur := sl.head.forward[0]
	for cur != nil {
		sb.WriteString(fmt.Sprintf("%v", cur.Key))
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
