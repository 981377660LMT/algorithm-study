package main

import (
	"bufio"
	"fmt"
	"os"
)

const BASE uint = 131 // 自然溢出哈希

// 505. Prefixes and suffixes
// https://codeforces.com/problemsets/acmsguru/problem/99999/505
//
// 给定一个只包含小写字母的字符串集合words.
// 在线处理q个查询：words中有多少个字符串s满足prefix是s的一个前缀, suffix是s的一个后缀.
// n,q,|words[i]|<=1e5
//
// 分治+哈希：
// !查询的prefix和suffix长度很短(<=20)时，一共有n*20*20种(取出每个串的前后缀)，可以预处理出.
// !否则，直接查询遍历每个长度>=20的串查询哈希，时间复杂度O(∑*q/400)
// !这种做法求出现次数可以换成求单词的最大索引，同力扣那道题
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	const THRESHOLD int32 = 20

	prefixHash := make([]map[int32]uint, n) // 每个串每个前缀的哈希值
	suffixHash := make([]map[int32]uint, n) // 每个串每个后缀的哈希值
	for i := int32(0); i < n; i++ {
		prefixHash[i] = make(map[int32]uint)
		h1 := uint(0)
		for j, c := range words[i] {
			h1 = (h1*BASE + uint(c))
			prefixHash[i][int32(j+1)] = h1
		}
		suffixHash[i] = make(map[int32]uint)
		h2 := uint(0)
		for j := len(words[i]) - 1; j >= 0; j-- {
			h2 = (h2*BASE + uint(words[i][j]))
			suffixHash[i][int32(len(words[i])-j)] = h2
		}
	}

	shortHash := make(map[[2]uint]int32) // 长度<=20的前后缀哈希值
	for i, w := range words {
		for preLen := int32(1); preLen <= min32(THRESHOLD, int32(len(w))); preLen++ {
			for sufLen := int32(1); sufLen <= min32(THRESHOLD, int32(len(w))); sufLen++ {
				h1, h2 := prefixHash[i][preLen], suffixHash[i][sufLen]
				shortHash[[2]uint{h1, h2}]++
			}
		}
	}

	longWordIndex := []int32{} // 长度>20的串的下标
	for i, w := range words {
		if int32(len(w)) > THRESHOLD {
			longWordIndex = append(longWordIndex, int32(i))
		}
	}

	// 查询words中有多少个字符串s满足prefix是s的一个前缀, suffix是s的一个后缀.
	query := func(prefix, suffix string) int32 {
		n1, n2 := int32(len(prefix)), int32(len(suffix))
		h1 := uint(0)
		for j := int32(0); j < n1; j++ {
			h1 = (h1*BASE + uint(prefix[j]))
		}
		h2 := uint(0)
		for j := n2 - 1; j >= 0; j-- {
			h2 = (h2*BASE + uint(suffix[j]))
		}

		if n1 <= THRESHOLD && n2 <= THRESHOLD {
			return shortHash[[2]uint{h1, h2}]
		} else {
			res := int32(0)
			for _, index := range longWordIndex {
				n3 := int32(len(words[index]))
				if n1 <= n3 && n2 <= n3 && prefixHash[index][n1] == h1 && suffixHash[index][n2] == h2 {
					res++
				}
			}
			return res
		}
	}

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var prefix, suffix string
		fmt.Fscan(in, &prefix, &suffix)
		fmt.Fprintln(out, query(prefix, suffix))
	}
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
