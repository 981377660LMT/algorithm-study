package main

import (
	"fmt"
	"math/bits"
)

func main() {
	undirectedGraph := [][]int{{1, 2}, {0, 2}, {0, 1, 3}, {2}}
	cycle := CountCycle(undirectedGraph)
	fmt.Println(cycle)

	undirectedGraph = [][]int{{1, 2}, {0, 2}, {0, 1}}
	cycle = CountCycle(undirectedGraph) // [0 0 0 0 0 0 0 1] -> 0<<1 | 1<<1 | 2<<1
	fmt.Println(cycle)
}

// 给定一个无重边、自环的无向图，求图中`长度大于等于3的环`的集合.
// O(n^2 * 2^n)
func CountCycle(undirectedGraph [][]int) (cycleGroup []int) {
	n := uint32(len(undirectedGraph))
	nexts := make([]uint32, n)
	for i := uint32(0); i < n; i++ {
		for _, j := range undirectedGraph[i] {
			nexts[i] |= 1 << j
		}
	}

	cycleGroup = make([]int, 1<<n)
	for v := uint32(0); v < n; v++ {
		dp := make([]int, v<<v)
		for w := uint32(0); w < v; w++ {
			if nexts[v]>>w&1 == 1 {
				dp[(v<<w)+w] = 1
			}
		}

		mask := (uint32(1) << v) - 1
		for s := uint32(0); s < 1<<v; s++ {
			EnumerateBitsUint32(s, func(a uint32) {
				EnumerateBitsUint32(nexts[a]&mask&(^s), func(b uint32) {
					dp[v*(s|1<<b)+b] += dp[v*s+a]
				})
				if bits.OnesCount32(s) >= 2 && nexts[a]>>v&1 == 1 {
					cycleGroup[s|1<<v] += dp[v*s+a]
				}
			})
		}
	}

	for i := range cycleGroup {
		cycleGroup[i] /= 2
	}

	return
}

// 遍历每个为1的比特位
func EnumerateBitsUint32(s uint32, f func(bit uint32)) {
	for s != 0 {
		i := bits.TrailingZeros32(s)
		f(uint32(i))
		s ^= 1 << i
	}
}
