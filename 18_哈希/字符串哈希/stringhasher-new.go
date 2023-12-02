// API:
//  RollingHash(base) : 以base为底的哈希函数.
//  Build(s) : O(n)返回s的哈希值表.
//  Query(sTable, l, r) : 返回s[l, r)的哈希值.
//  Combine(h1, h2, h2len) : 哈希值h1和h2的拼接, h2的长度为h2len. (线段树的op操作)
//  AddChar(h,c) : 哈希值h和字符c拼接形成的哈希值.
//  Lcp(aTable, l1, r1, bTable, l2, r2) : O(logn)返回a[l1, r1)和b[l2, r2)的最长公共前缀.

package main

import (
	"bufio"
	"fmt"
	"os"
)

// 顺序遍历每个单词,问之前是否见过类似的单词.
//
//	`类似`: 最多交换一次相邻字符,可以得到相同的单词.
//	n<=2e5 sum(len(s))<=1e6
func ConditionalReflection(words []string) []bool {

	R := NewRollingHash(13331)
	res := make([]bool, len(words))
	visited := make(map[uint]struct{})

	for i := 0; i < len(words); i++ {
		m := len(words[i])
		word := words[i]
		table := R.Build(word)
		hashes := []uint{R.Query(table, 0, m)} // 不交换
		for j := 0; j < m-1; j++ {             // 交换
			newWord := string(word[j+1]) + string(word[j])
			mid := R.Query(R.Build(newWord), 0, 2)
			left := R.Query(table, 0, j)
			right := R.Query(table, j+2, m)
			left = R.Combine(left, mid, 2)
			left = R.Combine(left, right, m-j-2)
			hashes = append(hashes, left)
		}

		for _, h := range hashes {
			if _, ok := visited[h]; ok {
				res[i] = true
				break
			}
		}

		visited[hashes[0]] = struct{}{}
	}

	return res
}

func main() {
	// https://yukicoder.me/problems/no/2102
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	res := ConditionalReflection(words)
	for i := 0; i < n; i++ {
		if res[i] {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

type RollingHash struct {
	base  uint
	power []uint
}

// 131/13331/1713302033171(回文素数)
func NewRollingHash(base uint) *RollingHash {
	return &RollingHash{
		base:  base,
		power: []uint{1},
	}
}

func (r *RollingHash) Build(s string) (hashTable []uint) {
	sz := len(s)
	hashTable = make([]uint, sz+1)
	for i := 0; i < sz; i++ {
		hashTable[i+1] = hashTable[i]*r.base + uint(s[i])
	}
	return hashTable
}

func (r *RollingHash) Query(sTable []uint, start, end int) uint {
	r.expand(end - start)
	return sTable[end] - sTable[start]*r.power[end-start]
}

func (r *RollingHash) Combine(h1, h2 uint, h2len int) uint {
	r.expand(h2len)
	return h1*r.power[h2len] + h2
}

func (r *RollingHash) AddChar(hash uint, c byte) uint {
	return hash*r.base + uint(c)
}

// 两个字符串的最长公共前缀长度.
func (r *RollingHash) LCP(sTable []uint, start1, end1 int, tTable []uint, start2, end2 int) int {
	len1 := end1 - start1
	len2 := end2 - start2
	len := min(len1, len2)
	low := 0
	high := len + 1
	for high-low > 1 {
		mid := (low + high) / 2
		if r.Query(sTable, start1, start1+mid) == r.Query(tTable, start2, start2+mid) {
			low = mid
		} else {
			high = mid
		}
	}
	return low
}

func (r *RollingHash) expand(sz int) {
	if len(r.power) < sz+1 {
		preSz := len(r.power)
		r.power = append(r.power, make([]uint, sz+1-preSz)...)
		for i := preSz - 1; i < sz; i++ {
			r.power[i+1] = r.power[i] * r.base
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
