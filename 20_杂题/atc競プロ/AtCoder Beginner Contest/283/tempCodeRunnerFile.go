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
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	for i := 0; i < n; i++ {
		fmt.Fprint(out, nums[i], " ")
	}
}
