package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	// 检查长度为3的子串
	for i := 0; i < n-2; i++ {
		A := []byte{}
		set := make(map[byte]bool)
		for j := 0; j < 3; j++ {
			if s[i+j] != '?' {
				A = append(A, s[i+j])
				set[s[i+j]] = true
			}
		}

		if len(A) != len(set) {
			fmt.Fprintln(out, 0)
			return
		}
	}

	// dp
	var dfs func(index int, pre1 int, pre2 int) int
	memo := [5050][30][30]int{}
	for i := 0; i < 5050; i++ {
		for j := 0; j < 30; j++ {
			for k := 0; k < 30; k++ {
				memo[i][j][k] = -1
			}
		}
	}
	dfs = func(index int, pre1 int, pre2 int) int {
		if index == n {
			return 1
		}
		if memo[index][pre1][pre2] != -1 {
			return memo[index][pre1][pre2]
		}
		if s[index] != '?' {
			if pre1 == int(s[index]-'a') || pre2 == int(s[index]-'a') {
				return 0
			}
			return dfs(index+1, pre2, int(s[index]-'a'))
		}
		res := 0
		for cand := 0; cand < 26; cand++ {
			if cand != pre1 && cand != pre2 {
				res += dfs(index+1, pre2, cand)
				res %= MOD
			}
		}
		memo[index][pre1][pre2] = res
		return res
	}

	res := dfs(0, 27, 28)
	fmt.Fprintln(out, res)
}
