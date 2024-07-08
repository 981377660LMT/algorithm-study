// 单词拼接dp
// 100350. 最小代价构造字符串
// https://leetcode.cn/problems/construct-string-with-minimum-cost/description/
// 给你一个字符串 target、一个字符串数组 words 以及一个整数数组 costs，这两个数组长度相同。
// 设想一个空字符串 s。
// 你可以执行以下操作任意次数（包括零次）：
// 选择一个在范围  [0, words.length - 1] 的索引 i。
// 将 words[i] 追加到 s。
// 该操作的成本是 costs[i]。
// 返回使 s 等于 target 的 最小 成本。如果不可能，返回 -1。
//
// 1 <= target.length <= 5e4
// 所有 words[i].length 的总和小于或等于 5e4
//
// dp[i]=min(dp[i],dp[j+1]+cost) 当target[i:j+1]出现在字典里.
// !解法：短串插入到trie中匹配(哈希表也可以)，长串个数不多，使用kmp预处理.

package main

const INF32 int32 = 1e9 + 10
const INF int = 1e18

type Str = string

func GetNext(pattern Str) []int32 {
	next := make([]int32, len(pattern))
	j := int32(0)
	for i := 1; i < len(pattern); i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = next[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		next[i] = j
	}
	return next
}

// `O(n+m)` 寻找 `shorter` 在 `longer` 中的所有匹配位置.
// nexts 数组为nil时, 会调用GetNext(shorter)求nexts数组.
func IndexOfAll(longer Str, shorter Str, position int32, nexts []int32) []int32 {
	if len(shorter) == 0 {
		return []int32{}
	}
	if len(longer) < len(shorter) {
		return nil
	}
	res := []int32{}
	if nexts == nil {
		nexts = GetNext(shorter)
	}
	hitJ := int32(0)
	for i := position; i < int32(len(longer)); i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = nexts[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == int32(len(shorter)) {
			res = append(res, i-int32(len(shorter))+1)
			hitJ = nexts[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}

type trieNode struct {
	children map[int32]*trieNode
	cost     int
}

func newTrieNode() *trieNode {
	return &trieNode{children: make(map[int32]*trieNode), cost: INF}
}

func (t *trieNode) Insert(word string, cost int) {
	node := t
	for _, c := range word {
		if _, ok := node.children[c]; !ok {
			node.children[c] = newTrieNode()
		}
		node = node.children[c]
	}
	node.cost = min(node.cost, cost)
}

func minimumCost(target string, words []string, costs []int) int {
	const THRESHOLD int32 = 500

	// !1.1 小串维护trie
	smallTrieRoot := newTrieNode()
	for i, word := range words {
		if int32(len(word)) <= THRESHOLD {
			smallTrieRoot.Insert(word, costs[i])
		}
	}

	// !1.2 大串kmp预处理
	n := int32(len(target))
	type pair struct {
		to   int32
		cost int
	}
	bigTo := make([][]pair, n+1)
	nexts := GetNext(target)
	for i, word := range words {
		m := int32(len(word))
		if m > THRESHOLD {
			starts := IndexOfAll(target, word, 0, nexts)
			for _, start := range starts {
				bigTo[start] = append(bigTo[start], pair{start + m, costs[i]})
			}
		}
	}

	dp := make([]int, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[n] = 0
	for i := n - 1; i >= 0; i-- {
		// !2.1 小串匹配
		root := smallTrieRoot
		for j := i; j < min32(n, i+THRESHOLD); j++ {
			if _, ok := root.children[int32(target[j])]; !ok {
				break
			}
			root = root.children[int32(target[j])]
			dp[i] = min(dp[i], dp[j+1]+int(root.cost))
		}

		// !2.2 大串匹配
		for _, v := range bigTo[i] {
			to, cost := v.to, v.cost
			dp[i] = min(dp[i], dp[to]+cost)
		}
	}

	if dp[0] == INF {
		return -1
	}
	return dp[0]
}

func min(a, b int) int {
	if a < b {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
