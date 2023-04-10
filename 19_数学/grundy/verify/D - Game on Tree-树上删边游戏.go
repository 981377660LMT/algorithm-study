// https://blog.hamayanhamayan.com/entry/2017/07/10/022214
// Alice喝Bob在树上比赛
// 轮流选择一条还存在的边,删除这条边,同时删去所有在删除后不再与根相连的部分.
// 最后无法选择的人输.
// 问先手是否必胜
// https://blog.csdn.net/wu_tongtong/article/details/79311284

// 克朗原理:
// 对于树上的某一个点，ta 的分支可以转化成以这个点为根的一根竹子(nim游戏中的石子堆)，
// 这个竹子的长度就是 ta 各个分支的边的数量的异或和

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		tree[a] = append(tree[a], b)
		tree[b] = append(tree[b], a)
	}

	var dfs func(int, int) int
	dfs = func(cur, pre int) int {
		g := 0
		for _, next := range tree[cur] {
			if next != pre {
				g ^= dfs(next, cur) + 1 // 统计分支的边数(grundy数)
			}
		}
		return g
	}

	groundy := dfs(0, -1)
	if groundy == 0 {
		fmt.Fprintln(out, "Bob")
	} else {
		fmt.Fprintln(out, "Alice")
	}
}
