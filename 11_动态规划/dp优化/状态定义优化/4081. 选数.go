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

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	fmt.Fprintln(out, MaximizeTrailingZeros(nums, k))
}

const INF int = 1e18

// !给定一个正整数数组,恰好选k个数,要求选出的数的乘积的末尾0的个数最多。
func MaximizeTrailingZeros(nums []int, k int) int {
	n := len(nums)
	C2, C5 := make([]int, n), make([]int, n)

	for i, num := range nums {
		for num%2 == 0 {
			num /= 2
			C2[i]++
		}
		for num%5 == 0 {
			num /= 5
			C5[i]++
		}
	}

	sorted := append(C5[:0:0], C5...)
	sort.Ints(sorted)
	max5 := 0
	for i := n - 1; i > n-1-k; i-- {
		max5 += sorted[i]
	}

	cache := make([]int, (n+1)*(k+1)*(max5+1))
	for i := range cache {
		cache[i] = -1
	}

	var dfs func(index, remain, c5 int) int
	dfs = func(index, remain, c5 int) int {
		if index == n {
			if remain == 0 && c5 == 0 {
				return 0
			}
			return -INF
		}

		key := index*(k+1)*(max5+1) + remain*(max5+1) + c5
		if cache[key] != -1 {
			return cache[key]
		}
		res := dfs(index+1, remain, c5)
		a, b := C2[index], C5[index]
		if remain > 0 && c5-b >= 0 {
			cand := dfs(index+1, remain-1, c5-b) + a
			if cand > res {
				res = cand
			}
		}
		cache[key] = res
		return res
	}

	res := 0
	for c5 := 0; c5 <= max5; c5++ {
		c2 := dfs(0, k, c5)
		if c2 >= c5 {
			res = c5
		}
	}
	return res
}
