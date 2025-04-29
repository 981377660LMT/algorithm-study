// !这个模板牺牲了性能，增加了通用性.

package main

import "fmt"

const INF int = 1e18

func main() {
	words := []string{"abc", "ab", "bc", "abcde"}
	acam := NewACAutoMatonArray(26, 'a')
	for _, word := range words {
		acam.AddString(word)
	}
	acam.BuildSuffixLink(true)

	indexes := make([][]int32, acam.Size())
	for i, pos := range acam.WordPos {
		indexes[pos] = append(indexes[pos], int32(i))
	}

	text := "abcdabc"
	SqrtSolverArray(
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
func SqrtSolverArray(
	acam *ACAutoMatonArray,
	text string,
	f func(textPrefixEnd int32 /** textPrefixEnd > 0 **/, wordPos int32),
) {
	hasWord := make([]bool, acam.Size())
	for _, pos := range acam.WordPos {
		hasWord[pos] = true
	}
	pos := int32(0)
	for i, c := range text {
		pos = acam.Move(pos, c)
		for cur := pos; cur != 0; cur = acam.LinkWord(cur) {
			if hasWord[cur] {
				f(int32(i+1), cur)
			}
		}
	}
}

// 3213. 最小代价构造字符串
// https://leetcode.cn/problems/construct-string-with-minimum-cost/description/
func minimumCost(target string, words []string, costs []int) int {
	acam := NewACAutoMatonArray(26, 'a')
	for _, word := range words {
		acam.AddString(word)
	}
	acam.BuildSuffixLink(true)

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

	SqrtSolverArray(
		acam, target,
		func(ti int32, wi int32) {
			dp[ti] = min(dp[ti], dp[ti-int32(depth[wi])]+nodeMinCost[wi])
		},
	)
	if dp[len(target)] == INF {
		return -1
	}
	return dp[len(target)]
}

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个前缀.
type ACAutoMatonArray struct {
	WordPos            []int32   // WordPos[i] 表示加入的第i个模式串对应的节点编号(单词结点).
	Parent             []int32   // parent[v] 表示节点v的父节点.
	Children           [][]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
	BfsOrder           []int32   // 结点的拓扑序,0表示虚拟节点.
	Depth              []int32   // !每个节点的深度.也就是对应的模式串前缀的长度.
	link               []int32   // 又叫fail.指向当前trie节点(对应一个前缀)的最长真后缀对应结点，例如"bc"是"abc"的最长真后缀.
	linkWord           []int32
	sigma              int32 // 字符集大小.
	offset             int32 // 字符集的偏移量.
	needUpdateChildren bool  // 是否需要更新children数组.
}

func NewACAutoMatonArray(sigma, offset int32) *ACAutoMatonArray {
	res := &ACAutoMatonArray{sigma: sigma, offset: offset}
	res.Clear()
	return res
}

// 添加一个字符串，返回最后一个字符对应的节点编号.
func (trie *ACAutoMatonArray) AddString(str string) int32 {
	if len(str) == 0 {
		return 0
	}
	pos := int32(0)
	for _, s := range str {
		ord := s - trie.offset
		if trie.Children[pos][ord] == -1 {
			trie.Children[pos][ord] = trie.newNode()
			trie.Parent[len(trie.Parent)-1] = pos
			trie.Depth[len(trie.Depth)-1] = trie.Depth[pos] + 1
		}
		pos = trie.Children[pos][ord]
	}
	trie.WordPos = append(trie.WordPos, pos)
	return pos
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *ACAutoMatonArray) Move(pos, ord int32) int32 {
	ord -= trie.offset
	if trie.needUpdateChildren {
		return trie.Children[pos][ord]
	}
	for {
		nexts := trie.Children[pos]
		if nexts[ord] != -1 {
			return nexts[ord]
		}
		if pos == 0 {
			return 0
		}
		pos = trie.link[pos]
	}
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *ACAutoMatonArray) Size() int32 {
	return int32(len(trie.Children))
}

func (trie *ACAutoMatonArray) Empty() bool {
	return len(trie.Children) == 1
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组(连接trie图).
//
// !move调用较少时，设置为false更快.
func (trie *ACAutoMatonArray) BuildSuffixLink(needUpdateChildren bool) {
	trie.needUpdateChildren = needUpdateChildren
	trie.link = make([]int32, len(trie.Children))
	for i := range trie.link {
		trie.link[i] = -1
	}
	trie.BfsOrder = make([]int32, len(trie.Children))
	head, tail := 0, 0
	trie.BfsOrder[tail] = 0
	tail++
	for head < tail {
		v := trie.BfsOrder[head]
		head++
		for i, next := range trie.Children[v] {
			if next == -1 {
				continue
			}
			trie.BfsOrder[tail] = next
			tail++
			f := trie.link[v]
			for f != -1 && trie.Children[f][i] == -1 {
				f = trie.link[f]
			}
			trie.link[next] = f
			if f == -1 {
				trie.link[next] = 0
			} else {
				trie.link[next] = trie.Children[f][i]
			}
		}
	}
	if !needUpdateChildren {
		return
	}
	for _, v := range trie.BfsOrder {
		for i, next := range trie.Children[v] {
			if next == -1 {
				f := trie.link[v]
				if f == -1 {
					trie.Children[v][i] = 0
				} else {
					trie.Children[v][i] = trie.Children[f][i]
				}
			}
		}
	}
}

// `linkWord`指向当前节点的最长后缀对应的节点.
// 区别于`_link`,`linkWord`指向的节点对应的单词不会重复.
// 即不会出现`_link`指向某个长串局部的恶化情况.
func (trie *ACAutoMatonArray) LinkWord(pos int32) int32 {
	if len(trie.linkWord) == 0 {
		hasWord := make([]bool, len(trie.Children))
		for _, p := range trie.WordPos {
			hasWord[p] = true
		}
		trie.linkWord = make([]int32, len(trie.Children))
		link, linkWord := trie.link, trie.linkWord
		for _, v := range trie.BfsOrder {
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
	return trie.linkWord[pos]
}

// 获取每个状态包含的模式串的索引.(模式串长度和较小时使用)
// fail指针每次命中，都至少有一个比指针深度更长的单词出现，因此每个位置最坏情况下不超过O(sqrt(n))次命中
// O(n*sqrt(n))
// TODO: roaring bitmaps 优化空间复杂度.
func (trie *ACAutoMatonArray) GetIndexes() [][]int32 {
	res := make([][]int32, len(trie.Children))
	for i, pos := range trie.WordPos {
		res[pos] = append(res[pos], int32(i))
	}
	for _, v := range trie.BfsOrder {
		if v != 0 {
			from, to := trie.link[v], v
			arr1, arr2 := res[from], res[to]
			arr3 := make([]int32, len(arr1)+len(arr2))
			i, j, k := 0, 0, 0
			for i < len(arr1) && j < len(arr2) {
				if arr1[i] < arr2[j] {
					arr3[k] = arr1[i]
					i++
				} else if arr1[i] > arr2[j] {
					arr3[k] = arr2[j]
					j++
				} else {
					arr3[k] = arr1[i]
					i++
					j++
				}
				k++
			}
			copy(arr3[k:], arr1[i:])
			k += len(arr1) - i
			copy(arr3[k:], arr2[j:])
			k += len(arr2) - j
			arr3 = arr3[:k:k]
			res[to] = arr3
		}
	}
	return res
}

func (trie *ACAutoMatonArray) Clear() {
	trie.WordPos = trie.WordPos[:0]
	trie.Parent = trie.Parent[:0]
	trie.Depth = trie.Depth[:0]
	trie.Children = trie.Children[:0]
	trie.link = trie.link[:0]
	trie.linkWord = trie.linkWord[:0]
	trie.BfsOrder = trie.BfsOrder[:0]
	trie.newNode()
}

func (trie *ACAutoMatonArray) newNode() int32 {
	trie.Parent = append(trie.Parent, -1)
	trie.Depth = append(trie.Depth, 0)
	nexts := make([]int32, trie.sigma)
	for i := range nexts {
		nexts[i] = -1
	}
	trie.Children = append(trie.Children, nexts)
	return int32(len(trie.Children) - 1)
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
