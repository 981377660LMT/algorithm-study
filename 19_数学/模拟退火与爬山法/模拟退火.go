// https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/odt_bst.go#L15
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// https://leetcode-cn.com/contest/weekly-contest-197/problems/best-position-for-a-service-centre/ http://poj.org/problem?id=2420 UVa 10228 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=14&page=show_problem&problem=1169
func simulatedAnnealing(f func(x float64) float64) float64 {
	// 例：最小值
	x := .0
	ans := f(x)
	for t := 1e5; t > 1e-8; t *= 0.99 {
		y := x + (2*rand.Float64()-1)*t
		v := f(y)
		if v < ans || math.Exp((ans-v)/t) > rand.Float64() { // 最小直接取，或者以一定概率接受较大的值
			ans = v
			x = y
		}
	}
	return ans
}

// 另一种写法（利用时限）
// 此时 alpha 可以设大点，例如 0.999
func simulatedAnnealingWithinTimeLimit(f func(x float64) float64) float64 {
	const timeLimit = 2 - 0.1
	t0 := time.Now()
	// 例：最小值
	x := .0
	ans := f(x)
	for t := 1e5; time.Since(t0).Seconds() < timeLimit; {
		y := x + (2*rand.Float64()-1)*t
		v := f(y)
		if v < ans || math.Exp((ans-v)/t) > rand.Float64() { // 最小直接取，或者以一定概率接受较大的值
			ans = v
			x = y
		}
		t *= 0.999 // 置于末尾，方便在 roll 到不合适的数据时直接 continue，同时也保证不会因为 roll 不到合适的数据而超时
	}
	return ans
}

func main() {
	res := simulatedAnnealingWithinTimeLimit(
		func(x float64) float64 {
			return x*x*x*x + x
		},
	)

	fmt.Println(res)
}
