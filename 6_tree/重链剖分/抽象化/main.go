package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc417G()
}

// G - Binary Cat
// https://atcoder.jp/contests/abc417/tasks/abc417_g
func abc417G() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	for i := 0; i < q; i++ {
		var l, r, x int
		fmt.Fscan(in, &l, &r, &x)
	}
}
