// 在16进制下，[1,n]中出现k种不同数的个数
// dp[pos][hasLeadingZero][isLimit][count]
// 位数<=2e5 k<=16
// !因为每个digit都是等价的 所以可以不用visited来表示状态 而是count来表示

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
)

const MOD int = 1e9 + 7

func cal(upper string, k int) int {
	n := len(upper)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		res, _ := strconv.ParseInt(string(upper[i]), 16, 0)
		nums[i] = int(res)
	}

	memo := [2e5 + 10][2][2][17]int{}
	for i1 := 0; i1 < 2e5+10; i1++ {
		for i2 := 0; i2 < 2; i2++ {
			for i3 := 0; i3 < 2; i3++ {
				for i4 := 0; i4 < 17; i4++ {
					memo[i1][i2][i3][i4] = -1
				}
			}
		}
	}
	var dfs func(pos int, hasLeadingZero int, isLimit int, visited int) int
	dfs = func(pos int, hasLeadingZero int, isLimit int, visited int) int {
		bitCount := bits.OnesCount(uint(visited))
		if bitCount > k {
			return 0
		}
		if pos == n {
			if bitCount == k {
				return 1
			}
			return 0
		}
		if cand := memo[pos][hasLeadingZero][isLimit][bitCount]; cand != -1 {
			return cand
		}

		res := 0
		up := 15
		if isLimit == 1 {
			up = nums[pos]
		}
		for cur := 0; cur <= up; cur++ {
			nextIsLimit := 0
			if isLimit == 1 && cur == up {
				nextIsLimit = 1
			}
			if hasLeadingZero == 1 && cur == 0 {
				res += dfs(pos+1, 1, nextIsLimit, visited)
				res %= MOD
			} else {
				res += dfs(pos+1, 0, nextIsLimit, (visited | (1 << cur)))
				res %= MOD
			}
		}

		memo[pos][hasLeadingZero][isLimit][bitCount] = res
		return res
	}

	return dfs(0, 1, 1, 0)
}

func digitsParadiseInHexadecimal(n string, k int) int {
	return cal(n, k) - cal("0", k)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n string
	var k int
	fmt.Fscan(in, &n, &k)
	fmt.Fprintln(out, digitsParadiseInHexadecimal(n, k))
}
