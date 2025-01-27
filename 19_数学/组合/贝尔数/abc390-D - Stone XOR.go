// D - Stone XOR
// https://atcoder.jp/contests/abc390/tasks/abc390_d
// 枚举所有的集合划分数

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	subsum := make([]int, 1<<n)
	for i := 0; i < n; i++ {
		for s := 0; s < 1<<i; s++ {
			subsum[s|1<<i] = subsum[s] + nums[i]
		}
	}

	var res []int
	var dfs func(int, int)
	dfs = func(remain int, sum int) {
		if remain == 0 {
			res = append(res, sum)
			return
		}

		lb := remain & -remain
		next := remain ^ lb
		for s := next; ; s = (s - 1) & next {
			subset := s | lb
			dfs(remain^subset, sum^subsum[subset]) // op
			if s == 0 {
				break
			}
		}
	}

	dfs((1<<n)-1, 0)
	sort.Ints(res)
	res = Compact(res)
	fmt.Fprintln(out, len(res))
}

// !Like Rust's Vec::dedup.
func Compact[S ~[]E, E comparable](s S) S {
	if len(s) < 2 {
		return s
	}
	i := 1
	for k := 1; k < len(s); k++ {
		if s[k] != s[k-1] {
			if i != k {
				s[i] = s[k]
			}
			i++
		}
	}
	// clear(s[i:])
	return s[:i]
}
