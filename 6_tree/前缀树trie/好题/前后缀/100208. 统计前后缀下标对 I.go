// 3045. 统计前后缀下标对 I
// https://leetcode.cn/problems/count-prefix-and-suffix-pairs-ii/
// 当 str1 同时是 str2 的前缀和后缀时，isPrefixAndSuffix(str1, str2) 返回 true，否则返回 false。
// 以整数形式，返回满足 i < j 且 isPrefixAndSuffix(words[i], words[j]) 为 true 的下标对 (i, j) 的 数量 。
//
// 把字符串 s 视作一个 pair 列表：
// [(s[0],s[n−1]),(s[1],s[n−2]),(s[2],s[n−3]),⋯,(s[n−1],s[0])]
// 只要这个 pair 列表是另一个字符串 t 的 pair 列表的前缀，那么 s 就是 t 的前后缀。

package main

func HashPair(first, second byte) int32 {
	return (int32(first)-'a')*26 + int32(second-'a')
}

type TrieNode struct {
	children map[int32]*TrieNode
	endCount int32
}

func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[int32]*TrieNode)}
}

func countPrefixSuffixPairs(words []string) int64 {
	res := 0
	root := NewTrieNode()
	for _, w := range words {
		cur := root
		for j := range w {
			p := HashPair(w[j], w[len(w)-1-j]) // (w[i], w[~i])
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
