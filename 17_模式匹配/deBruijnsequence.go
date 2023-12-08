// de Bruijn sequence
// 德布鲁因序列(De Bruijn sequence)
// B(k,n)
// 一个长度为 k^n 的循环字符串，该字符串包含了所有长度为 n 的子串，每个字符的元素有 k 种可能。
// https://halfrost.com/go_s2_de_bruijn/
// 在图论中的应用：欧拉回路 和 汉密尔顿回路
// 德布鲁因序列用的比较广泛的一点应用就是 位扫描器。在 Google S2 中也是这个作用。
// 参考golang 的 bits.TrailingZeros() 函数

package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(deBruijn(2, 3))
	// [0 0 0 1 0 1 1 1]
	// -> 000 001 010 101 011 111 110 100

}

func deBruijn(k int, n int) []int {
	upper := int(math.Pow(float64(k), float64(n)))
	res := make([]int, upper)
	aux := make([]int, n*k)
	if k == 1 {
		res[0] = 0
		return res
	}
	for i := 0; i < k*n; i++ {
		aux[i] = 0
	}
	size := 0
	var dfs func(int, int)
	dfs = func(t int, p int) {
		if t > n {
			if n%p == 0 {
				for i := 1; i <= p; i++ {
					res[size] = aux[i]
					size++
				}
			}
		} else {
			aux[t] = aux[t-p]
			dfs(t+1, p)
			for i := aux[t-p] + 1; i < k; i++ {
				aux[t] = i
				dfs(t+1, t)
			}
		}
	}
	dfs(1, 1)
	return res
}
