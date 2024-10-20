package main

import (
	"fmt"
	"math/rand"
)

func demo() {
	s := "abcba"
	n := int32(len(s))
	R := NewRollingHash(0)
	table1 := R.Build(n, func(i int32) uint64 { return uint64(s[i]) })
	table2 := R.Build(n, func(i int32) uint64 { return uint64(s[n-i-1]) })
	isPalindrome := func(start, end int32) bool {
		return R.Query(table1, start, end) == R.Query(table2, n-end, n-start)
	}

	fmt.Println(isPalindrome(0, 5)) // true
}

const (
	mod61  uint64 = (1 << 61) - 1
	mask30 uint64 = (1 << 30) - 1
	mask31 uint64 = (1 << 31) - 1
	mask61 uint64 = mod61
)

// RollingHash61
type RollingHash struct {
	base  uint64
	power []uint64
}

// base: 0 表示随机生成.
func NewRollingHash(base uint64) *RollingHash {
	for base == 0 {
		base = rand.Uint64() % mod61 // rng61
	}
	return &RollingHash{base: base, power: []uint64{1}}
}

func (rh *RollingHash) Build(n int32, f func(i int32) uint64) (table []uint64) {
	table = make([]uint64, n+1)
	for i := int32(0); i < n; i++ {
		table[i+1] = rh.add(rh.mul(table[i], rh.base), rh.mod(f(i)))
	}
	return
}

func (rh *RollingHash) Eval(n int32, f func(i int32) uint64) (res uint64) {
	for i := int32(0); i < n; i++ {
		res = rh.add(rh.mul(res, rh.base), rh.mod(f(i)))
	}
	return
}

func (rh *RollingHash) Query(table []uint64, start, end int32) uint64 {
	if start < 0 {
		start = 0
	}
	if end > int32(len(table)) {
		end = int32(len(table))
	}
	if start >= end {
		return 0
	}
	rh.expand(end - start)
	return rh.sub(table[end], rh.mul(table[start], rh.power[end-start]))
}

func (rh *RollingHash) Combine(h1, h2 uint64, h2len int32) uint64 {
	rh.expand(h2len)
	return rh.add(rh.mul(h1, rh.power[h2len]), h2)
}

func (rh *RollingHash) AddChar(h uint64, x uint64) uint64 {
	return rh.add(rh.mul(h, rh.base), rh.mod(x))
}

// s1[start1:end1] 与 s2[start2:end2] 的最长公共前缀长度.
func (rh *RollingHash) LCP(table1 []uint64, start1, end1 int32, table2 []uint64, start2, end2 int32) int32 {
	n := min32(end1-start1, end2-start2)
	low, high := int32(0), n+1
	for high-low > 1 {
		mid := (low + high) >> 1
		if rh.Query(table1, start1, start1+mid) == rh.Query(table2, start2, start2+mid) {
			low = mid
		} else {
			high = mid
		}
	}
	return low
}

func (rh *RollingHash) expand(size int32) {
	if int32(len(rh.power)) < size+1 {
		preSize := int32(len(rh.power))
		for i := preSize - 1; i < size; i++ {
			rh.power = append(rh.power, rh.mul(rh.power[i], rh.base))
		}
	}
}

// x % (2^61-1)
func (rh *RollingHash) mod(x uint64) uint64 {
	xu := x >> 61
	xd := x & mask61
	res := xu + xd
	if res >= mod61 {
		res -= mod61
	}
	return res
}

// a*b % (2^61-1)
func (rh *RollingHash) mul(a, b uint64) uint64 {
	au := a >> 31
	ad := a & mask31
	bu := b >> 31
	bd := b & mask31
	mid := ad*bu + au*bd
	midu := mid >> 30
	midd := mid & mask30
	return rh.mod(au*bu<<1 + midu + (midd << 31) + ad*bd)
}

// a,b: modint61
func (rh *RollingHash) add(a, b uint64) uint64 {
	res := a + b
	if res >= mod61 {
		res -= mod61
	}
	return res
}

// a,b: modint61
func (rh *RollingHash) sub(a, b uint64) uint64 {
	res := a - b
	if res >= mod61 {
		res += mod61
	}
	return res
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

// Problems:
//
//
//
//
//
//

// 2223. 构造字符串的总得分和
// https://leetcode.cn/problems/sum-of-scores-of-built-strings/description/
func sumScores(s string) int64 {
	hasher := NewRollingHash(0)
	table := hasher.Build(int32(len(s)), func(i int32) uint64 { return uint64(s[i]) })
	countPre := func(curLen, start int32) int32 {
		left, right := int32(1), curLen
		for left <= right {
			mid := (left + right) >> 1
			hash0 := hasher.Query(table, start, start+mid)
			hash1 := hasher.Query(table, 0, mid)
			if hash0 == hash1 {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}

		return right
	}

	n := int32(len(s))
	res := 0
	for i := int32(1); i < n+1; i++ {
		count := countPre(i, n-i)
		res += int(count)
	}
	return int64(res)
}

// 3327. 判断 DFS 字符串是否是回文串
// https://leetcode.cn/problems/check-if-dfs-strings-are-palindromes/description/
func findAnswer(parent []int, s string) []bool {
	n := int32(len(s))
	tree := make([][]int32, n)
	for i := int32(1); i < n; i++ {
		tree[parent[i]] = append(tree[parent[i]], i)
	}

	H := NewRollingHash(0)
	hash1 := make([]uint64, n)
	hash2 := make([]uint64, n)
	size := make([]int32, n)

	var dfs func(cur, pre int32)
	dfs = func(cur, pre int32) {
		curH1, curH2 := uint64(0), uint64(s[cur])
		curSize := int32(1)
		nexts := tree[cur]

		for i := 0; i < len(nexts); i++ {
			if nexts[i] == pre {
				continue
			}
			dfs(nexts[i], cur)
			curH1 = H.Combine(curH1, hash1[nexts[i]], size[nexts[i]])
			curSize += size[nexts[i]]
		}

		for i := int32(len(nexts)) - 1; i >= 0; i-- {
			if nexts[i] == pre {
				continue
			}
			curH2 = H.Combine(curH2, hash2[nexts[i]], size[nexts[i]])
		}

		curH1 = H.AddChar(curH1, uint64(s[cur]))
		hash1[cur] = curH1
		hash2[cur] = curH2
		size[cur] = curSize
	}
	dfs(0, -1)

	res := make([]bool, n)
	for i := int32(0); i < n; i++ {
		res[i] = hash1[i] == hash2[i]
	}
	return res
}
