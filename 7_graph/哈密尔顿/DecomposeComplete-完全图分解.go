// DecomposeComplete-完全图分解

package main

import "fmt"

func main() {
	fmt.Println(ToPaths(4))             // [[0 1 3 2] [1 2 0 3]]
	fmt.Println(ToCycles(5))            // [[0 1 3 2 4] [1 2 0 3 4]]
	fmt.Println(ToMatchings(4))         // [[3 0 2 1] [3 1 0 2] [3 2 1 0]]
	fmt.Println(ToCyclesAndMatching(4)) // [[0 2 1 3]] [0 1 2 3]
}

// 将偶数个顶点的无向完全图分解为哈密尔顿路径.
func ToPaths(n int) (hamiltonPath [][]int) {
	if n < 0 || n%2 != 0 {
		return
	}
	k := n / 2
	hamiltonPath = make([][]int, k)
	for i := range hamiltonPath {
		hamiltonPath[i] = make([]int, n)
	}

	for i := 0; i < k; i++ {
		for j := 0; j < n; j++ {
			x := &hamiltonPath[i][j]
			if j&1 == 1 {
				*x = i + (j+1)/2
			} else {
				*x = i - j/2
			}
			if *x < 0 {
				*x += n
			}
			if *x >= n {
				*x -= n
			}
		}
	}

	return
}

// 将奇数个顶点的无向完全图分解为哈密尔顿回路.
func ToCycles(n int) (hamiltonCycle [][]int) {
	if n < 0 || n%2 != 1 {
		return
	}
	hamiltonCycle = ToPaths(n - 1)
	for i := range hamiltonCycle {
		hamiltonCycle[i] = append(hamiltonCycle[i], n-1)
	}
	return
}

// 将偶数个顶点的无向完全图分解为匹配.
//  从中任意选择两个匹配都可以构成哈密尔顿回路.
func ToMatchings(n int) (matching [][]int) {
	if n <= 0 || n%2 != 0 {
		return
	}
	matching = make([][]int, n-1)
	for i := range matching {
		matching[i] = make([]int, n)
	}
	mod := n - 1
	for a := 0; a < mod; a++ {
		matching[a][0] = n - 1
		matching[a][1] = a
		for k := 1; k < n/2; k++ {
			b, c := a-k, a+k
			if b < 0 {
				b += mod
			}
			if c >= mod {
				c -= mod
			}
			matching[a][2*k] = b
			matching[a][2*k+1] = c
		}
	}
	return
}

// 将偶数顶点的无向完全图分解为哈密尔顿回路和一个匹配.
func ToCyclesAndMatching(n int) (cycles [][]int, matching []int) {
	if n <= 0 || n%2 != 0 {
		return
	}
	mod := n - 1
	cycles = make([][]int, n/2-1)
	for i := range cycles {
		cycles[i] = make([]int, n)
	}
	for a := 0; a < mod-1; a++ {
		if a%2 == 0 {
			cycle := cycles[a/2]
			cycle[0] = a
			for i := 0; i < n-2; i++ {
				var nxt int
				if i%2 == 0 {
					nxt = 2*a + 2 - cycle[i]
				} else {
					nxt = 2*a - cycle[i]
				}
				if nxt < 0 {
					nxt += mod
				}
				if nxt >= mod {
					nxt -= mod
				}
				cycle[i+1] = nxt
			}
			cycle[n-1] = mod
		}
	}
	matching = make([]int, n)
	for a := 0; a < mod/2; a++ {
		b := mod - 2 - a
		matching[2*a] = a
		matching[2*a+1] = b
	}
	matching[n-2] = n - 2
	matching[n-1] = n - 1
	return
}
