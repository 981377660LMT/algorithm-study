// 使用trie计算的编辑距离
// https://stevehanov.ca/blog/index.php?id=114
// 一旦某行 DP 的最小值大于 maxCost，便可在对应分支上剪枝，跳过进一步的深层搜索，显著提升在大规模词典下的匹配效率。
// 若仅仅是比较两个字符串的编辑距离，直接用二维 DP 即可；但当要对大量单词做近似匹配（模糊搜索）时，Trie + DP 的方法通常更具优势
//
// !这里的trie可以换成 DAWG（Directed Acyclic Word Graph）来进一步优化空间.
// 见 https://stevehanov.ca/blog/?id=115

package main

import (
	"fmt"
	"strings"
	"time"
)

const INF int = 1e18

func main() {
	// 5e4 len
	s1 := "a" + strings.Repeat("b", 1e5)
	s2 := "a" + strings.Repeat("c", 1e2)

	time1 := time.Now()
	fmt.Println(minDistance(s1, s2))
	fmt.Println(time.Since(time1))
}

// https://leetcode.cn/problems/edit-distance/description/
func minDistance(word1 string, word2 string) int {
	if len(word1) == 0 {
		return len(word2)
	}
	if len(word2) == 0 {
		return len(word1)
	}
	if len(word1) > len(word2) {
		word1, word2 = word2, word1
	}

	root := NewTrieNode()
	root.Insert(word1)

	results := Search(root, word2, INF)
	return results[0].distance
}

type TrieNode struct {
	word     string
	children map[rune]*TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
	}
}

func (node *TrieNode) Insert(word string) {
	curNode := node
	for _, w := range word {
		child, ok := curNode.children[w]
		if !ok {
			child = NewTrieNode()
			curNode.children[w] = child
		}
		curNode = child
	}
	curNode.word = word
}

type SearchResult struct {
	word     string
	distance int
}

// LevenshteinDistanceSearch.
// Search 返回所有与给定 `word` 的编辑距离 <= maxCost 的 (单词, 距离).
func Search(root *TrieNode, word string, maxCost int) []SearchResult {
	n := len(word)
	dp := make([]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = i
	}

	var res []SearchResult
	for letter, child := range root.children {
		searchRecursive(child, letter, word, dp, &res, maxCost)
	}
	return res
}

// Search 的辅助递归函数
//
//	node: 当前 Trie 节点
//	letter: 从父节点走过来的字符
//	word: 目标单词
//	previousRow: 上一行 DP(编辑距离) 数组
//	results: 存储 (单词, 距离) 的结果切片
//	maxCost: 最大可接受编辑距离
func searchRecursive(
	node *TrieNode, letter rune, word string,
	preDp []int, results *[]SearchResult, maxCost int,
) {
	columns := len(word) + 1
	curDp := make([]int, columns)
	curDp[0] = preDp[0] + 1
	dpMin := curDp[0]

	for col := 1; col < columns; col++ {
		insertCost := curDp[col-1] + 1
		deleteCost := preDp[col] + 1

		var replaceCost int
		if rune(word[col-1]) != letter {
			replaceCost = preDp[col-1] + 1
		} else {
			replaceCost = preDp[col-1]
		}

		curDp[col] = min3(insertCost, deleteCost, replaceCost)
		if curDp[col] < dpMin {
			dpMin = curDp[col]
		}
	}

	// 若该节点是一个完整单词，且编辑距离在 maxCost 范围内，则加入结果
	if node.word != "" && curDp[len(word)] <= maxCost {
		*results = append(*results, SearchResult{node.word, curDp[len(word)]})
	}

	// 如果本行最小值 <= maxCost，才有必要往下搜索
	if dpMin <= maxCost {
		for nextLetter, child := range node.children {
			searchRecursive(child, nextLetter, word, curDp, results, maxCost)
		}
	}
}

func mins(arr []int) int {
	minVal := arr[0]
	for _, v := range arr {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	} else {
		if b < c {
			return b
		}
		return c
	}
}
