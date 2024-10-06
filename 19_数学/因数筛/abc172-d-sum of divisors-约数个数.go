package main

import (
	"bufio"
	"fmt"
	"os"
)

// ABC172-D Sum of Divisors
// https://atcoder.jp/contests/abc172/tasks/abc172_d
// 记f(x)为x的约数个数，求∑f(i)*i (1<=i<=n)
// !筛法求约数个数.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int
	fmt.Fscan(in, &N)

	c2 := make([]int, N+1)
	for f := 1; f <= N; f++ {
		for m := f; m <= N; m += f {
			c2[m]++
		}
	}

	res := 0
	for i := 1; i <= N; i++ {
		res += i * c2[i]
	}

	fmt.Fprintln(out, res)
}
