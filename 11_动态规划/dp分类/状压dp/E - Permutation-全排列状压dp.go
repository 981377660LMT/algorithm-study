/**
求 1-n的排列个数, 且满足 m 个限制条件 [i,valueUpper,countUpper]
!A[1],A[2],...A[i]中不超过valueUpper的数的个数不超过countUpper

n<=18 m<=100
!全排列状压dp O(m*2^n)
**/
package main

import (
	"bufio"
	"fmt"
	"os"
)

func countPermutation(n int, restrictions [][3]int) int {
	type pair struct {
		valueUpper, countUpper int
	}

	group := make([][]pair, n) // 每个index处需要满足的限制条件
	for i := 0; i < len(restrictions); i++ {
		group[restrictions[i][0]-1] = append(group[restrictions[i][0]-1], pair{restrictions[i][1], restrictions[i][2]})
	}

	check := func(index int, nums []int) bool {
		for i := 0; i < len(group[index]); i++ {
			count := 0
			for j := 0; j < len(nums); j++ {
				if nums[j] <= group[index][i].valueUpper {
					count++
				}
				if count > group[index][i].countUpper {
					return false
				}
			}
		}
		return true
	}

	memo := [20][1 << 20]int{}
	for i := 0; i < 20; i++ {
		for j := 0; j < 1<<20; j++ {
			memo[i][j] = -1
		}
	}

	var dfs func(index int, visited int) int
	dfs = func(index int, visited int) int {
		if index == n {
			return 1
		}
		if memo[index][visited] != -1 {
			return memo[index][visited]
		}

		res := 0
		selected := []int{}
		for i := 0; i < n; i++ {
			if (visited>>i)&1 == 1 {
				selected = append(selected, i+1)
			}
		}

		for next := 0; next < n; next++ {
			if (visited>>next)&1 == 1 {
				continue
			}
			selected = append(selected, next+1)
			if check(index, selected) {
				res += dfs(index+1, visited|1<<next)
			}
			selected = selected[:len(selected)-1]
		}

		memo[index][visited] = res
		return res
	}

	res := dfs(0, 0)
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	restrictions := make([][3]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &restrictions[i][0], &restrictions[i][1], &restrictions[i][2])
	}

	res := countPermutation(n, restrictions)
	fmt.Fprintln(out, res)
}
