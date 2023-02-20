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

	fmt.Fprint(out, "res")
}
