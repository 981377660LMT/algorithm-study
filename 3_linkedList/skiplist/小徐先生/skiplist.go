// https://mp.weixin.qq.com/s?__biz=MzkxMjQzMjA0OQ==&mid=2247484204&idx=1&sn=54817591aa44359cde9b1b88d386b31b

package main

import "math/rand"

func main() {

}

type Skiplist struct {
	head *node
}

type node struct {
	nexts    []*node
	key, val int
}

func (s *Skiplist) Get(key int) (int, bool) {
	if node := s.search(key); node != nil {
		return node.val, true
	}
	return -1, false
}

func (s *Skiplist) Put(key, val int) {
	if node := s.search(key); node != nil {
		node.val = val
		return
	}

	// roll 出新节点的高度
	level := s.roll()

	// 新节点高度超出跳表最大高度，则需要对高度进行补齐
	for len(s.head.nexts)-1 < level {
		s.head.nexts = append(s.head.nexts, nil)
	}

	// 创建出新的节点
	newNode := node{
		key:   key,
		val:   val,
		nexts: make([]*node, level+1),
	}

	// 从头节点的最高层出发
	node := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		// 向右遍历，直到右侧节点不存在或者 key 值大于 key
		for node.nexts[level] != nil && node.nexts[level].key < key {
			node = node.nexts[level]
		}

		// 调整指针关系，完成新节点的插入
		newNode.nexts[level] = node.nexts[level]
		node.nexts[level] = &newNode
	}
}

func (s *Skiplist) Del(key int) {
	if n := s.search(key); n == nil {
		return
	}

	// 从头节点的最高层出发
	node := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for node.nexts[level] != nil && node.nexts[level].key < key {
			node = node.nexts[level]
		}

		if node.nexts[level] == nil || node.nexts[level].key > key {
			continue
		}

		// 走到此处意味着右侧节点的 key 值必然等于 key，则调整指针引用
		node.nexts[level] = node.nexts[level].nexts[level]
	}

	var diff int
	// 倘若某一层已经不存在数据节点，高度需要递减
	for level := len(s.head.nexts) - 1; level > 0 && s.head.nexts[level] == nil; level-- {
		diff++
	}
	s.head.nexts = s.head.nexts[:len(s.head.nexts)-diff]
}

// 找到 skiplist 当中 ≥ start，且 ≤ end 的 kv 对
func (s *Skiplist) Range(start, end int) [][2]int {
	// 首先通过 ceiling 方法，找到 skiplist 中 key 值大于等于 start 且最接近于 start 的节点 ceilNode
	ceilNode := s.ceiling(start)
	// 如果不存在，直接返回
	if ceilNode == nil {
		return [][2]int{}
	}

	// 从 ceilNode 首层出发向右遍历，把所有位于 [start,end] 区间内的节点统统返回
	var res [][2]int
	for move := ceilNode; move != nil && move.key <= end; move = move.nexts[0] {
		res = append(res, [2]int{move.key, move.val})
	}
	return res
}

// 找到 skiplist 中，key 值大于等于 target 且最接近于 target 的 key-value 对
func (s *Skiplist) Ceiling(target int) ([2]int, bool) {
	if ceilNode := s.ceiling(target); ceilNode != nil {
		return [2]int{ceilNode.key, ceilNode.val}, true
	}
	return [2]int{}, false
}

// 找到 skiplist 中，key 值小于等于 target 且最接近于 target 的 key-value 对
func (s *Skiplist) Floor(target int) ([2]int, bool) {
	if floorNode := s.floor(target); floorNode != nil {
		return [2]int{floorNode.key, floorNode.val}, true
	}
	return [2]int{}, false
}

// 找到 key 值大于等于 target 且 key 值最接近于 target 的节点
func (s *Skiplist) ceiling(target int) *node {
	move := s.head

	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}
		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}

	// 此时 move 已经对应于在首层 key 值小于 key 且最接近于 key 的节点，其右侧第一个节点即为所寻找的目标节点
	return move.nexts[0]
}

// 找到 key 值小于等于 target 且 key 值最接近于 target 的节点
func (s *Skiplist) floor(target int) *node {
	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}
		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}
	return move
}

func (s *Skiplist) search(key int) *node {
	node := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		// 在每一层中持续向右遍历，直到下一个节点不存在或者 key 值大于等于 key
		for node.nexts[level] != nil && node.nexts[level].key < key {
			node = node.nexts[level]
		}
		if node.nexts[level] != nil && node.nexts[level].key == key {
			return node.nexts[level]
		}
	}
	return nil
}

// TODO
func (s *Skiplist) roll() int {
	var level int
	// 每次投出 1，则层数加 1
	for rand.Intn(2) > 0 {
		level++
	}
	return level
}
