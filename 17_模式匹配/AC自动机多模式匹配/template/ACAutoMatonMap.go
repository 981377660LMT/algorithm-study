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

func main() {
	// 示例用法
	s := "abcxyz123"
	words := []string{"abc", "123"}
	result := addBoldTag(s, words)
	println(result) // 输出: "<b>abc</b><b>abc</b>"
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
		// end := i + 1
		// start := end - depth[pos]
		// boldDiff[start]++
		// boldDiff[end]--
		// https://leetcode.cn/problems/add-bold-tag-in-string/submissions/626352379/
		// TODO: BUG
		for cur := pos; cur != 0; cur = acm.LinkWord(cur) {
			end := i + 1
			start := end - depth[cur]
			boldDiff[start]++
			boldDiff[end]--
		}
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

// 2781. 最长合法子字符串的长度
// https://leetcode.cn/problems/length-of-the-longest-valid-substring/description/
func longestValidSubstring(word string, forbidden []string) int {
	acm := NewACAutoMatonMap()
	for _, w := range forbidden {
		acm.AddString([]byte(w))
	}
	acm.BuildSuffixLink()

	minLen := make([]int32, acm.Size()) // 每个状态匹配到的模式串的最小长度
	for i := range minLen {
		minLen[i] = INF32
	}
	for i, pos := range acm.WordPos {
		minLen[pos] = min32(minLen[pos], int32(len(forbidden[i])))
	}
	acm.Dp(func(from, to int32) { minLen[to] = min32(minLen[to], minLen[from]) })

	res, left, pos := int32(0), int32(0), int32(0)
	for right := int32(0); right < int32(len(word)); right++ {
		char := word[right]
		pos = acm.Move(pos, byte(char))
		left = max32(left, right-minLen[pos]+2)
		res = max32(res, right-left+1)
	}
	return int(res)
}

// 3213. 最小代价构造字符串
// https://leetcode.cn/problems/construct-string-with-minimum-cost/description/
func minimumCost(target string, words []string, costs []int) int {
	trie := NewACAutoMatonMap()
	for _, word := range words {
		trie.AddString([]byte(word))
	}
	trie.BuildSuffixLink()

	nodeCosts := make([]int, trie.Size())
	for i := range nodeCosts {
		nodeCosts[i] = INF
	}
	for i, pos := range trie.WordPos {
		nodeCosts[pos] = min(nodeCosts[pos], costs[i])
	}
	nodeDepth := trie.Depth
	hasWord := make([]bool, trie.Size())
	for _, p := range trie.WordPos {
		hasWord[p] = true
	}

	dp := make([]int, len(target)+1)
	for i := 1; i <= len(target); i++ {
		dp[i] = INF
	}
	pos := int32(0)
	for i := int32(0); i < int32(len(target)); i++ {
		char := target[i]
		pos = trie.Move(pos, byte(char))
		for cur := pos; cur != 0; cur = trie.LinkWord(cur) {
			if hasWord[cur] {
				dp[i+1] = min(dp[i+1], dp[i+1-nodeDepth[cur]]+nodeCosts[cur])
			}
		}
	}
	if dp[len(target)] == INF {
		return -1
	}
	return dp[len(target)]
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

// 功能与 AddString 相同.
func (ac *ACAutoMatonMap) AddFrom(n int32, getOrd func(i int32) T) int32 {
	if n == 0 {
		return 0
	}
	pos := int32(0)
	for i := int32(0); i < n; i++ {
		ord := getOrd(i)
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

func (ac *ACAutoMatonMap) AddChar(pos int32, ord T) int32 {
	nexts := ac.children[pos]
	if next, ok := nexts[ord]; ok {
		return next
	}
	return ac.newNode2(pos, ord)
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

func (ac *ACAutoMatonMap) GetCounter() []int32 {
	counter := make([]int32, len(ac.children))
	for _, pos := range ac.WordPos {
		counter[pos]++
	}
	for _, v := range ac.bfsOrder {
		if v != 0 {
			counter[v] += counter[ac.link[v]]
		}
	}
	return counter
}

func (ac *ACAutoMatonMap) GetIndexes() [][]int32 {
	res := make([][]int32, len(ac.children))
	for i, pos := range ac.WordPos {
		res[pos] = append(res[pos], int32(i))
	}

	for _, v := range ac.bfsOrder {
		if v != 0 {
			from, to := ac.link[v], v
			arr1, arr2 := res[from], res[to]
			n1, n2 := int32(len(arr1)), int32(len(arr2))
			arr3 := make([]int32, 0, n1+n2)
			i, j := int32(0), int32(0)
			for i < n1 && j < n2 {
				if arr1[i] < arr2[j] {
					arr3 = append(arr3, arr1[i])
					i++
				} else if arr1[i] > arr2[j] {
					arr3 = append(arr3, arr2[j])
					j++
				} else {
					arr3 = append(arr3, arr1[i])
					i++
					j++
				}
			}
			for i < n1 {
				arr3 = append(arr3, arr1[i])
				i++
			}
			for j < n2 {
				arr3 = append(arr3, arr2[j])
				j++
			}
			res[to] = arr3
		}
	}
	return res
}

func (ac *ACAutoMatonMap) Dp(f func(from, to int32)) {
	for _, v := range ac.bfsOrder {
		if v != 0 {
			f(ac.link[v], v)
		}
	}
}

func (trie *ACAutoMatonMap) BuildFailTree() [][]int32 {
	res := make([][]int32, trie.Size())
	trie.Dp(func(pre, cur int32) {
		res[pre] = append(res[pre], cur)
	})
	return res
}

func (ac *ACAutoMatonMap) BuildTrieTree() [][]int32 {
	res := make([][]int32, ac.Size())
	var dfs func(int32)
	dfs = func(cur int32) {
		for _, next := range ac.children[cur] {
			res[cur] = append(res[cur], next)
			dfs(next)
		}
	}
	dfs(0)
	return res
}

// 返回str在trie树上的节点位置.如果不存在，返回0.
func (trie *ACAutoMatonMap) Search(str string) int32 {
	if len(str) == 0 {
		return 0
	}
	pos := int32(0)
	for i := 0; i < len(str); i++ {
		if pos >= int32(len(trie.children)) || pos < 0 {
			return 0
		}
		nexts := trie.children[pos]
		if next, ok := nexts[str[i]]; ok {
			pos = next
		} else {
			return 0
		}
	}
	return pos
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
