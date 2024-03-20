// https://www.cnblogs.com/TianMeng-hyl/p/14989441.html
// https://www.cnblogs.com/Dfkuaid-210/p/bit_divide.html

package main

import (
	"bufio"
	"fmt"
	"os"
)

// String Set Queries
// https://www.luogu.com.cn/problem/CF710F
// 1 s : 在数据结构中插入 s
// 2 s : 在数据结构中删除 s
// 3 s : 查询集合中的所有字符串在给出的模板串中出现的次数
//
// !短的串丢进 trie 维护，长的串用 kmp 维护
// !每次查询时，对文本串的所有后缀扔到 trie 里查询其前缀匹配个数，求总和。再把文本串与所有较大模式串单独 KMP。
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	type info = struct {
		str    string
		nexts  []int
		weight int
	}

	THRESHOLD := 1000
	longPattern := []info{}
	shortPattern := NewSimpleTrie(26, 'a')

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op int
		var s string
		fmt.Fscan(in, &op, &s)

		// 增加、删除
		if op == 1 || op == 2 {
			var weight int
			if op == 1 {
				weight = 1
			} else {
				weight = -1
			}

			len_ := len(s)

			// 模式串长度小于 THRESHOLD 时，用 trie 维护
			if len_ <= THRESHOLD {
				shortPattern.Add(len_, func(i int) int { return int(s[i]) }, weight)
			} else {
				// 模式串长度大于 THRESHOLD 时，用 kmp 维护
				longPattern = append(longPattern, info{str: s, nexts: GetNext(s), weight: weight})
			}
		} else {
			// 查询
			res := 0
			len_ := len(s)
			for start := 0; start < len_; start++ {
				res += shortPattern.Query(len_-start, func(i int) int { return int(s[start+i]) })
			}

			for _, p := range longPattern {
				res += p.weight * CountIndexOfAll(s, p.str, 0, p.nexts)
			}

			fmt.Fprintln(out, res)
			out.Flush() // 强制在线，需要刷新缓冲区
		}
	}
}

type SimpleTrie struct {
	sigma    int32     // 字符集大小.
	offset   int32     // 字符集的偏移量.
	weight   []int     // weight[v] 表示节点v对应的字符串出现的次数.
	children [][]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
}

func NewSimpleTrie(sigma, offset int) *SimpleTrie {
	res := &SimpleTrie{sigma: int32(sigma), offset: int32(offset)}
	res.newNode()
	return res
}

// 添加一个字符串.
func (trie *SimpleTrie) Add(n int, getAt func(int) int, delta int) {
	pos := int32(0)
	for i := 0; i < n; i++ {
		ord := int32(getAt(i)) - trie.offset
		if trie.children[pos][ord] == -1 {
			trie.children[pos][ord] = trie.newNode()
		}
		pos = trie.children[pos][ord]
	}

	// !afterInsert
	trie.weight[pos] += delta
}

// 查询一个字符串所有前缀出现的次数之和.
func (trie *SimpleTrie) Query(n int, getAt func(int) int) int {
	res := 0
	pos := int32(0)
	for i := 0; i < n; i++ {
		ord := int32(getAt(i)) - trie.offset
		if trie.children[pos][ord] == -1 {
			break
		}
		pos = trie.children[pos][ord]
		res += trie.weight[pos]
	}
	return res
}

func (trie *SimpleTrie) newNode() int32 {
	nexts := make([]int32, trie.sigma)
	for i := range nexts {
		nexts[i] = -1
	}
	trie.children = append(trie.children, nexts)
	trie.weight = append(trie.weight, 0)
	return int32(len(trie.children) - 1)
}

type Str = string

func GetNext(pattern Str) []int {
	next := make([]int, len(pattern))
	j := 0
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

// `O(n+m)` 寻找 `shorter` 在 `longer` 中匹配次数.
// nexts 数组为nil时, 会调用GetNext(shorter)求nexts数组.
func CountIndexOfAll(longer Str, shorter Str, position int, nexts []int) int {
	if len(shorter) == 0 {
		return 0
	}
	if len(longer) < len(shorter) {
		return 0
	}
	res := 0
	if nexts == nil {
		nexts = GetNext(shorter)
	}
	hitJ := 0
	for i := position; i < len(longer); i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = nexts[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == len(shorter) {
			res++
			hitJ = nexts[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
