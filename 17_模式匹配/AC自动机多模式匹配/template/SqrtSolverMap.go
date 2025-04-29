// !这个模板牺牲了性能，增加了通用性.

package main

import "fmt"

const INF int = 1e18

func main() {
	words := []string{"abc", "ab", "bc", "abcde"}
	acam := NewACAutoMatonMap()
	for _, word := range words {
		acam.AddString(word)
	}
	acam.BuildSuffixLink()

	indexes := make([][]int32, acam.Size())
	for i, pos := range acam.WordPos {
		indexes[pos] = append(indexes[pos], int32(i))
	}

	text := "abcdabc"
	SqrtSolverMap(
		acam, text,
		func(textPrefixEnd int32, wordPos int32) {
			for _, wid := range indexes[wordPos] {
				prefix := text[:textPrefixEnd]
				fmt.Printf("text 的前缀 %s 的后缀匹配模式串 %s\n", prefix, words[wid])
			}
		},
	)
}

// 给定一些模式串和一个文本串.
// 对文本串的每个前缀 `text[0:textPrefixEnd)`，其后缀匹配模式串对应的树节点编号为 `wordPos`.
// f 被调用的次数为 `O(len(text) * sqrt(∑len(patterns)))`.
func SqrtSolverMap(
	acm *ACAutoMatonMap,
	text string,
	f func(textPrefixEnd int32 /** textPrefixEnd > 0 **/, wordPos int32),
) {
	hasWord := make([]bool, acm.Size())
	for _, pos := range acm.WordPos {
		hasWord[pos] = true
	}
	pos := int32(0)
	for i, c := range text {
		pos = acm.Move(pos, c)
		for cur := pos; cur != 0; cur = acm.LinkWord(cur) {
			if hasWord[cur] {
				f(int32(i+1), cur)
			}
		}
	}
}

// 3213. 最小代价构造字符串
// https://leetcode.cn/problems/construct-string-with-minimum-cost/description/
func minimumCost(target string, words []string, costs []int) int {
	acam := NewACAutoMatonMap()
	for _, word := range words {
		acam.AddString(word)
	}
	acam.BuildSuffixLink()

	depth := acam.Depth
	nodeMinCost := make([]int, acam.Size())
	for i := range nodeMinCost {
		nodeMinCost[i] = INF
	}
	for i, pos := range acam.WordPos {
		nodeMinCost[pos] = min(nodeMinCost[pos], costs[i])
	}
	dp := make([]int, len(target)+1)
	for i := 1; i <= len(target); i++ {
		dp[i] = INF
	}
	dp[0] = 0

	SqrtSolverMap(acam, target, func(ti int32, pi int32) {
		dp[ti] = min(dp[ti], dp[ti-depth[pi]]+nodeMinCost[pi])
	})
	if dp[len(target)] == INF {
		return -1
	}
	return dp[len(target)]
}

type T = rune

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

func (ac *ACAutoMatonMap) AddString(s string) int32 {
	if len(s) == 0 {
		return 0
	}
	pos := int32(0)
	for _, ord := range s {
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

// !对当前文本串后缀，找到每个模式串单词匹配的最大前缀.
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
