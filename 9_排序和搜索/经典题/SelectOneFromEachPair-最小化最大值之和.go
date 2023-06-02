// # SelectOneFromEachPair-最小化最大值之和
// # https://zhuanlan.zhihu.com/p/268630329

// # 给定n个对(ai,bi),对每个对,需要选择ai加入集合A,或者选择bi加入集合B
// # !最小化`max(A)+max(B)`

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

	var n, m int
	fmt.Fscan(in, &n, &m)
	robber := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &robber[i][0], &robber[i][1])
	}
	searchLight := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &searchLight[i][0], &searchLight[i][1])
	}

	pairs := make([][2]int, 0, n*m)
	for i := 0; i < n; i++ {
		rx, ry := robber[i][0], robber[i][1]
		for j := 0; j < m; j++ {
			sx, sy := searchLight[j][0], searchLight[j][1]
			if rx <= sx && ry <= sy {
				pairs = append(pairs, [2]int{sx - rx + 1, sy - ry + 1})
			}
		}
	}

	fmt.Fprintln(out, selectOneFromEachPairMinimizeMaxSum(pairs))
}

func selectOneFromEachPairMinimizeMaxSum(pairs [][2]int) int {
	if len(pairs) == 0 {
		return 0
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][0] < pairs[j][0]
	})
	sufMax := make([]int, len(pairs)+1)
	for i := len(pairs) - 1; i >= 0; i-- {
		sufMax[i] = max(sufMax[i+1], pairs[i][1])
	}
	res := sufMax[0]
	for i := 0; i < len(pairs); i++ {
		res = min(res, pairs[i][0]+sufMax[i+1])
	}
	return res

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
