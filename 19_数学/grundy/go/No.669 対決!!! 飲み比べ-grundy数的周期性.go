// https://yukicoder.me/problems/no/669

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// !喝酒
	// 二人轮流喝酒,每轮喝 1-k 毫升
	// 不能继续喝的人输
	// 问先手是否必胜
	// !1<n<=1000; 1<=k<=1000; 1<=nums[i]<=1e6
	// !结论:每次限定从1个堆取[1,k]个石头,石头堆剩余i个石头时,grundy[i] = i%(k+1)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	xor := 0
	for i := range nums {
		xor ^= nums[i] % (k + 1)
	}

	if xor == 0 {
		fmt.Fprintln(out, "NO")
	} else {
		fmt.Fprintln(out, "YES")
	}
}
