package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/enumerate_quotients
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	res := FloorRange(n)
	fmt.Fprintln(out, len(res))
	for i := len(res) - 1; i >= 0; i-- {
		fmt.Fprint(out, res[i][2], " ")
	}
}

// 将 [1,n] 内的数分成O(2*sqrt(n))段, 每段内的 n//i 相同.
// 每个段为(left,right,div)，表示 left <= i <= right 内的 n//i == div.
func FloorRange(n int) [][3]int {
	if n <= 0 {
		return nil
	}
	res := [][3]int{}
	m := 1
	for m*m <= n {
		res = append(res, [3]int{m, m, n / m})
		m++
	}
	for i := m; i > 0; i-- {
		left := n/(i+1) + 1
		right := n / i
		if left <= right && len(res) > 0 && res[len(res)-1][1] < left {
			res = append(res, [3]int{left, right, n / left})
		}
	}
	return res
}
