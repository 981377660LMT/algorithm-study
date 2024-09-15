package main

const INF int = 1e18
const INF32 int32 = 1e9 + 10

// 单词拆分
// 3213. 最小代价构造字符串
// https://leetcode.cn/problems/construct-string-with-minimum-cost/description/
//
// 给你一个字符串 target、一个字符串数组 words 以及一个整数数组 costs，这两个数组长度相同。
// 设想一个空字符串 s。
// 你可以执行以下操作任意次数（包括零次）：
// 选择一个在范围  [0, words.length - 1] 的索引 i。
// 将 words[i] 追加到 s。
// 该操作的成本是 costs[i]。
// 返回使 s 等于 target 的 最小 成本。如果不可能，返回 -1
//
// !最坏情况: 文本串全是a，模式串是[a, aa, aaa, aaaa, ...]，最多有O(nsqrt(n))个匹配项.
// O(nsqrtn)
// !跳fail解法.
func minimumCost(target string, words []string, costs []int) int {
	trie := NewACAutoMatonArray(26, 97)
	for _, word := range words {
		trie.AddString(word)
	}
	trie.BuildSuffixLink(true)

	nodeCosts := make([]int, trie.Size())
	for i := range nodeCosts {
		nodeCosts[i] = INF
	}
	for i, pos := range trie.WordPos {
		nodeCosts[pos] = min(nodeCosts[pos], costs[i])
	}
	nodeDepth := trie.Depth

	dp := make([]int, len(target)+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0

	pos := int32(0)
	for i, char := range target {
		pos = trie.Move(pos, char)
		// 对当前文本串后缀，找到每个匹配的模式串.
		for cur := pos; cur != 0; cur = trie.LinkWord(cur) {
			dp[i+1] = min(dp[i+1], dp[int32(i)+1-nodeDepth[cur]]+nodeCosts[cur])
		}
	}
	if dp[len(target)] == INF {
		return -1
	}
	return dp[len(target)]
}

// 3292. 形成目标字符串需要的最少字符串数 II
// https://leetcode.cn/problems/minimum-number-of-valid-strings-to-form-target-ii/description/
// 给你一个字符串数组 words 和一个字符串 target。
// 如果字符串 x 是 words 中 任意 字符串的 `前缀`，则认为 x 是一个 有效 字符串。
// 现计划通过 连接 有效字符串形成 target ，请你计算并返回需要连接的 最少 字符串数量。如果无法通过这种方式形成 target，则返回 -1。
//
// !AC自动机接受所有模式串的前缀.转移时判断即可.
func minValidStrings(words []string, target string) int {
	trie := NewACAutoMatonArray(26, 97)
	for _, s := range words {
		trie.AddString(s)
	}
	trie.BuildSuffixLink(true)

	nodeDepth := trie.Depth

	dp := make([]int, len(target)+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0

	pos := int32(0)
	for i := int32(0); i < int32(len(target)); i++ {
		c := int32(target[i])
		pos = trie.Move(pos, c)
		if pos == 0 { // 匹配失败.
			return -1
		}
		dp[i+1] = dp[i+1-nodeDepth[pos]] + 1
	}
	return dp[len(target)]
}

// 单词拆分
// 3213. 最小代价构造字符串
// https://leetcode.cn/problems/construct-string-with-minimum-cost/description/
//
// 给你一个字符串 target、一个字符串数组 words 以及一个整数数组 costs，这两个数组长度相同。
// 设想一个空字符串 s。
// 你可以执行以下操作任意次数（包括零次）：
// 选择一个在范围  [0, words.length - 1] 的索引 i。
// 将 words[i] 追加到 s。
// 该操作的成本是 costs[i]。
// 返回使 s 等于 target 的 最小 成本。如果不可能，返回 -1
//
// !最坏情况: 文本串全是a，模式串是[a, aa, aaa, aaaa, ...]，最多有O(nsqrt(n))个匹配项.
// !预处理dp转移解法.
// O(nsqrtn)
func minimumCost2(target string, words []string, costs []int) int {
	trie := NewACAutoMatonArray(26, 97)
	for _, word := range words {
		trie.AddString(word)
	}

	trie.BuildSuffixLink(true)
	dp := make([]int, len(target)+1)
	for i := 1; i <= len(target); i++ {
		dp[i] = INF
	}
	pos := int32(0)
	indexes := trie.GetIndexes() // 每个状态包含的模式串的索引
	for i, char := range target {
		pos = trie.Move(pos, char)
		for _, wordIndex := range indexes[pos] {
			wordLen := len(words[wordIndex])
			if i+1 >= wordLen {
				dp[i+1] = min(dp[i+1], dp[i+1-wordLen]+costs[wordIndex])
			}
		}
	}
	if dp[len(target)] == INF {
		return -1
	}
	return dp[len(target)]
}

// 2781. 最长合法子字符串的长度
// https://leetcode.cn/problems/length-of-the-longest-valid-substring/
// 给你一个字符串 word 和一个字符串数组 forbidden 。
// 如果一个字符串不包含 forbidden 中的任何字符串，我们称这个字符串是 合法 的。
// 请你返回字符串 word 的一个 最长合法子字符串 的长度。
// 子字符串 指的是一个字符串中一段连续的字符，它可以为空。
//
// 1 <= word.length <= 1e5
// word 只包含小写英文字母。
// 1 <= forbidden.length <= 1e5
// !1 <= forbidden[i].length <= 1e5
// !sum(len(forbidden)) <= 1e7
// forbidden[i] 只包含小写英文字母。
//
// 思路:
// !子串是前缀的后缀.
// 类似字符流, 需要处理出每个位置为结束字符的包含至少一个模式串的`最短后缀`.
// !那么此时左端点就对应这个位置+1
func longestValidSubstring(word string, forbidden []string) int {
	acm := NewACAutoMatonArray(26, 97)
	for _, w := range forbidden {
		acm.AddString(w)
	}
	acm.BuildSuffixLink(false)

	minLen := make([]int32, acm.Size()) // 每个状态(前缀)匹配的最短单词
	for i := range minLen {
		minLen[i] = INF32
	}
	for i, pos := range acm.WordPos {
		minLen[pos] = min32(minLen[pos], int32(len(forbidden[i])))
	}
	acm.Dp(func(from, to int32) { minLen[to] = min32(minLen[to], minLen[from]) })

	res, left, pos := int32(0), int32(0), int32(0)
	for right, char := range word {
		pos = acm.Move(pos, char)
		left = max32(left, int32(right)-minLen[pos]+2)
		res = max32(res, int32(right)-left+1)
	}
	return int(res)
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

// 在pos位置添加一个字符，返回新的节点编号.
func (trie *ACAutoMatonArray) AddChar(pos, ord int32) int32 {
	ord -= trie.offset
	if trie.Children[pos][ord] != -1 {
		return trie.Children[pos][ord]
	}
	trie.Children[pos][ord] = trie.newNode()
	trie.Parent[len(trie.Parent)-1] = pos
	trie.Depth[len(trie.Depth)-1] = trie.Depth[pos] + 1
	return trie.Children[pos][ord]
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

// 获取每个状态包含的模式串的个数.
func (trie *ACAutoMatonArray) GetCounter() []int32 {
	counter := make([]int32, len(trie.Children))
	for _, pos := range trie.WordPos {
		counter[pos]++
	}
	for _, v := range trie.BfsOrder {
		if v != 0 {
			counter[v] += counter[trie.link[v]]
		}
	}
	return counter
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

// 按照拓扑序自顶向下进行转移(EnumerateFail).
func (trie *ACAutoMatonArray) Dp(f func(link, cur int32)) {
	for _, v := range trie.BfsOrder {
		if v != 0 {
			f(trie.link[v], v)
		}
	}
}

// 按照拓扑序逆序自底向上进行转移(EnumerateFailReverse).
func (trie *ACAutoMatonArray) DpReverse(f func(link, cur int32)) {
	for i := len(trie.BfsOrder) - 1; i >= 0; i-- {
		v := trie.BfsOrder[i]
		if v != 0 {
			f(trie.link[v], v)
		}
	}
}

func (trie *ACAutoMatonArray) BuildFailTree() [][]int32 {
	res := make([][]int32, trie.Size())
	trie.Dp(func(pre, cur int32) {
		res[pre] = append(res[pre], cur)
	})
	return res
}

func (trie *ACAutoMatonArray) BuildTrieTree() [][]int32 {
	res := make([][]int32, trie.Size())
	for i := int32(1); i < trie.Size(); i++ {
		res[trie.Parent[i]] = append(res[trie.Parent[i]], i)
	}
	return res
}

// 返回str在trie树上的节点位置.如果不存在，返回0.
func (trie *ACAutoMatonArray) Search(str string) int32 {
	if len(str) == 0 {
		return 0
	}
	pos := int32(0)
	for _, char := range str {
		if pos >= int32(len(trie.Children)) || pos < 0 {
			return 0
		}
		ord := char - trie.offset
		if next := trie.Children[pos][ord]; next == -1 {
			return 0
		} else {
			pos = next
		}
	}
	return pos
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
