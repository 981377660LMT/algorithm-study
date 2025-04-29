package main

import (
	"strings"
)

const INF int = 1e18

const INF32 int32 = 1e9 + 10

// 遍历连续key相同元素的分组.
func enumerateGroupByKey[K comparable](n int, key func(index int) K, f func(start, end int, curKey K)) {
	end := 0
	for end < n {
		start := end
		leader := key(end)
		end++
		for end < n && key(end) == leader {
			end++
		}
		f(start, end, leader)
	}
}

// 616. 给字符串添加加粗标签
func addBoldTag(s string, words []string) string {
	acm := NewACAutoMatonMap()
	for _, word := range words {
		acm.AddString([]byte(word))
	}
	acm.BuildSuffixLink()

	depth := acm.Depth
	boldDiff := make([]int, len(s)+1)

	pos := int32(0)
	for i := int32(0); i < int32(len(s)); i++ {
		pos = acm.Move(pos, s[i])
		end := i + 1
		start := end - depth[pos]
		boldDiff[start]++
		boldDiff[end]--

	}
	for i := 0; i < len(s); i++ {
		boldDiff[i+1] += boldDiff[i]
	}

	var sb strings.Builder
	enumerateGroupByKey(
		len(s), func(i int) bool { return boldDiff[i] > 0 },
		func(start, end int, b bool) {
			if b {
				sb.WriteString("<b>")
				sb.WriteString(s[start:end])
				sb.WriteString("</b>")
			} else {
				sb.WriteString(s[start:end])
			}
		},
	)

	return sb.String()
}

type T = byte

type ACAutoMatonMap struct {
	WordPos  []int32       // WordPos[i] 表示加入的第i个模式串对应的节点编号.
	Parent   []int32       // Parent[i] 表示第i个节点的父节点.
	Depth    []int32       // !Depth[i] 表示第i个节点的深度.也就是对应的模式串前缀的长度.
	children []map[T]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
	link     []int32       // 又叫fail.指向当前节点最长真后缀对应结点.
	linkWord []int32
	bfsOrder []int32 // 结点的拓扑序,0表示虚拟节点.
}

func NewACAutoMatonMap() *ACAutoMatonMap {
	res := &ACAutoMatonMap{}
	res.Clear()
	return res
}

func (ac *ACAutoMatonMap) AddString(s []T) int32 {
	if len(s) == 0 {
		return 0
	}
	pos := int32(0)
	for i := 0; i < len(s); i++ {
		ord := s[i]
		nexts := ac.children[pos]
		if next, ok := nexts[ord]; ok {
			pos = next
		} else {
			pos = ac.newNode2(pos, ord)
		}
	}
	ac.WordPos = append(ac.WordPos, pos)
	return pos
}

// 当前文本后缀能匹配的最长模式的前缀.
func (ac *ACAutoMatonMap) Move(pos int32, ord T) int32 {
	for {
		nexts := ac.children[pos]
		if next, ok := nexts[ord]; ok {
			return next
		}
		if pos == 0 {
			return 0
		}
		pos = ac.link[pos]
	}
}

func (ac *ACAutoMatonMap) BuildSuffixLink() {
	ac.link = make([]int32, len(ac.children))
	for i := range ac.link {
		ac.link[i] = -1
	}
	ac.bfsOrder = make([]int32, len(ac.children))
	head, tail := 0, 1
	for head < tail {
		v := ac.bfsOrder[head]
		head++
		for char, next := range ac.children[v] {
			ac.bfsOrder[tail] = next
			tail++
			f := ac.link[v]
			for f != -1 {
				if _, ok := ac.children[f][char]; ok {
					break
				}
				f = ac.link[f]
			}
			if f == -1 {
				ac.link[next] = 0
			} else {
				ac.link[next] = ac.children[f][char]
			}
		}
	}
}

// !对当前文本串后缀，找到每个模式串单词匹配的最长前缀.
func (ac *ACAutoMatonMap) LinkWord(pos int32) int32 {
	if len(ac.linkWord) == 0 {
		hasWord := make([]bool, len(ac.children))
		for _, p := range ac.WordPos {
			hasWord[p] = true
		}
		ac.linkWord = make([]int32, len(ac.children))
		link, linkWord := ac.link, ac.linkWord
		for _, v := range ac.bfsOrder {
			if v != 0 {
				p := link[v]
				if hasWord[p] {
					linkWord[v] = p
				} else {
					linkWord[v] = linkWord[p]
				}
			}
		}
	}
	return ac.linkWord[pos]
}

func (ac *ACAutoMatonMap) Empty() bool {
	return len(ac.children) == 1
}

func (ac *ACAutoMatonMap) Clear() {
	ac.WordPos = ac.WordPos[:0]
	ac.Parent = ac.Parent[:0]
	ac.Depth = ac.Depth[:0]
	ac.children = ac.children[:0]
	ac.link = ac.link[:0]
	ac.linkWord = ac.linkWord[:0]
	ac.bfsOrder = ac.bfsOrder[:0]
	ac.newNode()
}

func (ac *ACAutoMatonMap) Size() int32 {
	return int32(len(ac.children))
}

func (ac *ACAutoMatonMap) newNode() int32 {
	ac.children = append(ac.children, map[T]int32{})
	cur := int32(len(ac.children) - 1)
	ac.Parent = append(ac.Parent, -1)
	ac.Depth = append(ac.Depth, 0)
	return cur
}

func (ac *ACAutoMatonMap) newNode2(parent int32, ord T) int32 {
	ac.children = append(ac.children, map[T]int32{})
	cur := int32(len(ac.children) - 1)
	ac.Parent = append(ac.Parent, parent)
	ac.Depth = append(ac.Depth, ac.Depth[parent]+1)
	ac.children[parent][ord] = cur
	return cur
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
