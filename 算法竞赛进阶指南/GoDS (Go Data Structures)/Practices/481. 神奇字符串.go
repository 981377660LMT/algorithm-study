package main

import (
	"fmt"
)

const N int = 1e5

var preSum [N + 1]int // !静态数组

func init() { // !init 里打表
	res := make([]byte, 0, N+10)
	res = append(res, 1, 2, 2)
	i := 2
	for len(res) < N {
		last := res[len(res)-1] ^ 3
		repeat := res[i]
		for j := 0; j < int(repeat); j++ {
			res = append(res, last)
		}
		i++
	}

	for i := 1; i <= N; i++ {
		preSum[i] = preSum[i-1] + int(2-res[i-1])
	}
}

func magicalString(n int) int {
	return preSum[n]
}

func main() {
	fmt.Println(magicalString(6))
}
