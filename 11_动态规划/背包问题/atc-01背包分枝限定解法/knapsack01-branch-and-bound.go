// 剪枝 n<=2000,vi<=1e9,wi<=1e9,limit<=1e9
// O(2^n)

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

	var n, limit int
	fmt.Fscan(in, &n, &limit)
	values := make([]int, n)
	weights := make([]int, n)
	for i := range values {
		fmt.Fscan(in, &values[i], &weights[i])
	}
	fmt.Fprintln(out, knapsackBranchBound(values, weights, limit))
}

func knapsackBranchBound(values []int, weights []int, limit int) int {
	n := len(values)
	goods := make([][]int, n)
	for i := range goods {
		goods[i] = []int{values[i], weights[i]}
	}
	sort.Slice(goods, func(i, j int) bool {
		return goods[i][0]*goods[j][1] > goods[j][0]*goods[i][1]
	})

	best := 0
	relax := func(i, v, w int) (int, bool) {
		res := v
		flag := true
		for i < n {
			if w == 0 {
				break
			}
			if w >= goods[i][1] {
				w -= goods[i][1]
				res += goods[i][0]
				i++
				continue
			}
			flag = false
			res += goods[i][0] * w / goods[i][1]
			break
		}
		return res, flag
	}

	var dfs func(i int, v int, w int)
	dfs = func(i int, v int, w int) {
		if i == n {
			if v > best {
				best = v
			}
			return
		}

		rel, flag := relax(i, v, w)
		if flag {
			if rel > best {
				best = rel
			}
			return
		}

		if rel < best {
			return
		}

		if w >= goods[i][1] {
			dfs(i+1, v+goods[i][0], w-goods[i][1])
		}

		dfs(i+1, v, w)
	}

	dfs(0, 0, limit)
	return best
}
