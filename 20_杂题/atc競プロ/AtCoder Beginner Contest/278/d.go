package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = int(1e18)
const MOD int = 998244353

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

	// !分块
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var x int
			fmt.Fscan(in, &x)
			tree.Update(1, n, x, 0) // 覆盖
		} else if op == 2 {
			var i, x int
			fmt.Fscan(in, &i, &x)
			tree.Update(i, i, x, 1) // 单点修改
		} else {
			var i int
			fmt.Fscan(in, &i)
			fmt.Fprintln(out, tree.Query(i, i))
		}
	}
}
