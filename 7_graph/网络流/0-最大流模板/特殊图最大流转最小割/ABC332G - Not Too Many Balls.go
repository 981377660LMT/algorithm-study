// [ABC332G] Not Too Many Balls
// https://atcoder.jp/contests/abc332/tasks/abc332_g
//
// 现在有 n 种球和 m 个盒子，第 i 种球有 ai​  个，第 j 个箱子里最多放 bj​  个球。
// 除此之外，第 j 个箱子里最多放 i×j 个第 i 种球。
// 求最多能放多少个球。
// n<=500,m<=5e5
// TODO
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

	var n, m int
	fmt.Fscan(in, &n, &m)

	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	B := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &B[i])
	}

}
