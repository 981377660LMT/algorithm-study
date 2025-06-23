package main

import (
	"fmt"
	"math/bits"
)

// 解答错误
// 507 / 600 个通过的测试用例
// 输入
// n =
// 4
// k =
// 2
// m =
// 4
// time =
// [51,57,22,73]
// mul =
// [0.98,1.68,0.55,1.71]
// 添加到测试用例
// 输出
// 223.37000
// 预期结果
// 215.96
func main() {
	fmt.Println(minTime(4, 2, 4, []int{51, 57, 22, 73}, []float64{0.98, 1.68, 0.55, 1.71}))
}

const inf float64 = 2e15

// 降序枚举子集(包含s自身).
func EnumerateSubset(s int, allowEmpty bool, f func(b int)) {
	for b := s; b > 0; b = (b - 1) & s {
		f(b)
	}
	if allowEmpty {
		f(0)
	}
}

func minTime(n int, k int, m int, timeArr []int, mul []float64) float64 {
	full := (1 << n) - 1
	dp := make([]float64, (1<<n)*m)
	for i := range dp {
		dp[i] = -1
	}

	var dfs func(visited, mulIndex int) float64
	dfs = func(visited, mulIndex int) float64 {
		if visited == full {
			return 0
		}
		hash := visited*m + mulIndex
		if dp[hash] != -1 {
			return dp[hash]
		}
		dp[hash] = inf

		res := inf
		EnumerateSubset((^visited)&full, false, func(g int) {
			size := bits.OnesCount(uint(g))
			if size > k {
				return
			}

			nextVisited := visited | g

			maxTime := 0
			for i := 0; i < n; i++ {
				if (g>>i)&1 == 1 && timeArr[i] > maxTime {
					maxTime = timeArr[i]
				}
			}
			crossTime := float64(maxTime) * mul[mulIndex]
			nextMulIndex := (mulIndex + int(crossTime)%m) % m
			if nextVisited == full {
				res = min64(res, crossTime)
			} else {
				for i := 0; i < n; i++ {
					if (g>>i)&1 == 1 {
						returnTime := float64(timeArr[i]) * mul[nextMulIndex]
						nextMulIndex2 := (nextMulIndex + int(returnTime)%m) % m
						remain := dfs(nextVisited^(1<<i), nextMulIndex2)
						if remain < inf {
							res = min64(res, crossTime+returnTime+remain)
						}
					}
				}
			}
		})

		dp[hash] = res
		return res
	}

	res := dfs(0, 0)
	if res >= inf {
		return -1
	}
	return res
}

func min64(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
