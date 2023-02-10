// st表

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type S = int

// DisjointSparseTable 维护不包含更新操作的半群的演算
//  op:只需要满足结合律 op(op(a,b),c) = op(a,op(b,c))
//  例如:乘积取模
func NewDisjointSparseTable(nums []S, op func(S, S) S) (query func(int, int) S) {
	n := len(nums)
	b := bits.Len(uint(n))
	lookup := make([]int, 1<<b)
	st := make([][]S, b)
	for i := range st {
		st[i] = make([]S, n)
	}
	for i := 0; i < n; i++ {
		st[0][i] = nums[i]
	}

	for i := 1; i < b; i++ {
		shift := 1 << i
		for j := 0; j < n; j += shift << 1 {
			t := min(j+shift, n)
			st[i][t-1] = nums[t-1]
			for k := t - 2; k >= j; k-- {
				st[i][k] = op(nums[k], st[i][k+1])
				if n <= t {
					break
				}
				st[i][t] = nums[t]
				r := min(t+shift, n)
				for k := t + 1; k < r; k++ {
					st[i][k] = op(st[i][k-1], nums[k])
				}
			}
		}
	}

	for i := 2; i < len(lookup); i++ {
		lookup[i] = lookup[i>>1] + 1
	}

	query = func(l, r int) S {
		if l >= r {
			return st[0][l]
		}
		p := lookup[l^r]
		return op(st[p][l], st[p][r])
	}

	return
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

func main() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	query := NewDisjointSparseTable(nums, min)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		r--
		fmt.Fprintln(out, query(l, r))
	}
}
