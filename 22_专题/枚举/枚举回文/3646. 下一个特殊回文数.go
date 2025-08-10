// 3646. 下一个特殊回文数(打表)
// https://leetcode.cn/problems/next-special-palindrome-number/solutions/3748548/bao-li-mei-ju-he-fa-pai-lie-by-endlessch-b5gw/
// 给你一个整数 n。
//
// 如果一个数满足以下条件，那么它被称为 特殊数 ：
//
// 它是一个 回文数 。
// 数字中每个数字 k 出现 恰好 k 次。
// 返回 严格 大于 n 的 最小 特殊数。
//
// 如果一个整数正向读和反向读都相同，则它是 回文数 。例如，121 是回文数，而 123 不是。

package main

import (
	"fmt"
	"sort"
	"strconv"
)

func main() {
	fmt.Println(specialPalindrome(23670))
}

var cands []int64

func init() {
	var halfCounter [10]int

	var dfs func(pos, targetLen, mid int, path []byte)
	dfs = func(pos, targetLen, mid int, path []byte) {
		if pos == targetLen {
			s := append([]byte(nil), path...)
			if mid != 0 {
				s = append(s, byte('0'+mid))
			}
			for i := targetLen - 1; i >= 0; i-- {
				s = append(s, path[i])
			}
			v, _ := strconv.ParseInt(string(s), 10, 64)
			cands = append(cands, v)
			return
		}

		for d := 1; d <= 9; d++ {
			if halfCounter[d] == 0 {
				continue
			}
			halfCounter[d]--
			path[pos] = byte('0' + d)
			dfs(pos+1, targetLen, mid, path)
			halfCounter[d]++
		}
	}

	for state := 1; state < 1<<9; state++ {
		allCount, mid := 0, 0
		for d := 1; d <= 9; d++ {
			if state&(1<<(d-1)) != 0 {
				allCount += d
				if d&1 == 1 {
					mid = d
				}
				halfCounter[d] = d / 2
			} else {
				halfCounter[d] = 0
			}
		}
		if allCount > 18 {
			continue
		}
		dfs(0, allCount/2, mid, make([]byte, allCount/2))
	}

	sort.Slice(cands, func(i, j int) bool { return cands[i] < cands[j] })
}

func specialPalindrome(n int64) int64 {
	pos := sort.Search(len(cands), func(i int) bool { return cands[i] > n })
	return cands[pos]
}
