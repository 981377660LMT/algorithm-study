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

	isDebug := false

	var n, k int
	var nums []int
	if isDebug {
		n, k = 1, 2
		nums = []int{1, 2, 3, 4}
	} else {
		fmt.Fscan(in, &n, &k)
		nums = make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &nums[i])
		}
	}

	fmt.Fprintln(out, "res")
}
