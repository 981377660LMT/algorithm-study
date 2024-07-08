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
// !每个长度对应一个counter，长度种类不超过sqrtn.

package main

const INF int = 1e18

func minimumCost(target string, words []string, costs []int) int {
	H := NewSimpleHash(BASE, MOD)
	table := H.Build(target)
	hashGroupByLen := make(map[int32]map[uint]int) // !每个长度对应一个counter
	add := func(word string, cost int) {
		hash := GetHash(word, BASE, MOD)
		len_ := int32(len(word))
		if _, ok := hashGroupByLen[len_]; !ok {
			hashGroupByLen[len_] = make(map[uint]int)
		}
		inner := hashGroupByLen[len_]
		if _, ok := inner[hash]; !ok {
			inner[hash] = cost
		} else {
			inner[hash] = min(inner[hash], cost)
		}
	}
	for i, word := range words {
		add(word, costs[i])
	}

	n := int32(len(target))
	dp := make([]int, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for i := int32(0); i < n; i++ {
		for len_, inner := range hashGroupByLen {
			if i+1-len_ >= 0 {
				hash := H.Query(table, i+1-len_, i+1)
				if cost, ok := inner[hash]; ok {
					dp[i+1] = min(dp[i+1], dp[i+1-len_]+cost)
				}
			}
		}
	}

	if dp[n] == INF {
		return -1
	}
	return dp[n]
}

const BASE uint = 13331
const MOD uint = 1e9 + 7

type S = string

func GetHash(s S, base uint, mod uint) uint {
	if len(s) == 0 {
		return 0
	}
	res := uint(0)
	for i := 0; i < len(s); i++ {
		res = (res*base + uint(s[i])) % mod
	}
	return res
}

type SimpleHash struct {
	base  uint
	mod   uint
	power []uint
}

// 131/13331/1713302033171(回文素数)
func NewSimpleHash(base uint, mod uint) *SimpleHash {
	return &SimpleHash{
		base:  base,
		mod:   mod,
		power: []uint{1},
	}
}

func (r *SimpleHash) Build(s S) (hashTable []uint) {
	sz := len(s)
	hashTable = make([]uint, sz+1)
	for i := 0; i < sz; i++ {
		hashTable[i+1] = (hashTable[i]*r.base + uint(s[i])) % r.mod
	}
	return hashTable
}

func (r *SimpleHash) Query(sTable []uint, start, end int32) uint {
	r.expand(end - start)
	return (r.mod + sTable[end] - sTable[start]*r.power[end-start]%r.mod) % r.mod
}

func (r *SimpleHash) expand(sz int32) {
	if int32(len(r.power)) < sz+1 {
		preSz := int32(len(r.power))
		r.power = append(r.power, make([]uint, sz+1-preSz)...)
		for i := preSz - 1; i < sz; i++ {
			r.power[i+1] = (r.power[i] * r.base) % r.mod
		}
	}
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
