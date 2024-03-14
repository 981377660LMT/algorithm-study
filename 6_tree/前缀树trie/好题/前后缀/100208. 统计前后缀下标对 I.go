// 3045. 统计前后缀下标对 I
// https://leetcode.cn/problems/count-prefix-and-suffix-pairs-ii/
// 当 str1 同时是 str2 的前缀和后缀时，isPrefixAndSuffix(str1, str2) 返回 true，否则返回 false。
// 以整数形式，返回满足 i < j 且 isPrefixAndSuffix(words[i], words[j]) 为 true 的下标对 (i, j) 的 数量 。
//
// 把字符串 s 视作一个 pair 列表：
// [(s[0],s[n−1]),(s[1],s[n−2]),(s[2],s[n−3]),⋯,(s[n−1],s[0])]
// 只要这个 pair 列表是另一个字符串 t 的 pair 列表的前缀，那么 s 就是 t 的前后缀。

package main

type Pair struct{ x, y byte }
type TrieNode struct {
	children map[Pair]*TrieNode
	endCount int32
}

func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[Pair]*TrieNode)}
}

func countPrefixSuffixPairs(words []string) int64 {
	res := 0
	root := NewTrieNode()
	for _, w := range words {
		cur := root
		for i := range w {
			p := Pair{w[i], w[len(w)-1-i]}
			if cur.children[p] == nil {
				cur.children[p] = NewTrieNode()
			}
			cur = cur.children[p]
			res += int(cur.endCount)
		}
		cur.endCount++
	}

	return int64(res)
}
