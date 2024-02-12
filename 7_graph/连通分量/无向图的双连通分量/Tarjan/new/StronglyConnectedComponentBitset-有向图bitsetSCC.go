// StronglyConnectedComponentBitset-有向图bitsetSCC
// O(n^2/64)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// https://www.luogu.com.cn/problem/CF1268D
// 给定一个 n 个点的竞赛图，定义一次翻转操作为翻转一个结点相连的所有边方向，求至少要翻转多少个结点才能让图强联通，并求出最少操作的方案数
// 对 998244353 取模.
// 3 <= n <= 2000.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n int
	fmt.Fscan(in, &n)
	adjMatrix := make([]Bitset, n)
	rAdjMatrix := make([]Bitset, n)
	for i := 0; i < n; i++ {
		adjMatrix[i] = NewBitset(n)
		rAdjMatrix[i] = NewBitset(n)
	}
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		for j := 0; j < n; j++ {
			if s[j] == '1' {
				adjMatrix[i].Set(j)
				rAdjMatrix[j].Set(i)
			}
		}
	}

	count, belong := StronglyConnectedComponentBitset(n, adjMatrix, rAdjMatrix)
	if count <= 1 {
		fmt.Fprintln(out, "0 1")
		return
	}

	res := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				adjMatrix[i].Flip(j)
				adjMatrix[j].Flip(i)
				rAdjMatrix[j].Flip(i)
				rAdjMatrix[i].Flip(j)
			}
		}
		count, _ := StronglyConnectedComponentBitset(n, adjMatrix, rAdjMatrix)
		if count == 1 {
			res++
		}
		for j := 0; j < n; j++ {
			if i != j {
				adjMatrix[i].Flip(j)
				adjMatrix[j].Flip(i)
				rAdjMatrix[j].Flip(i)
				rAdjMatrix[i].Flip(j)
			}
		}
	}

	if res > 0 {
		fmt.Fprintln(out, "1", res)
		return
	}

	if count == 2 {
		n0, n1 := 0, 0
		for i := 0; i < n; i++ {
			if belong[i] == 0 {
				n0++
			} else {
				n1++
			}
		}
		if n == 4 && n0*n1 == 3 {
			fmt.Fprintln(out, "-1")
			return
		}
		fmt.Fprintln(out, "2", n0*n1*2%MOD)
	}

}

// kosaraju + bitset.
func StronglyConnectedComponentBitset(n int, adjMatrix []Bitset, rAdjMatrix []Bitset) (count int, belong []int) {
	belong = make([]int, n)
	unVisited := NewBitset(n)
	order := make([]int, 0, n)

	dfs := func(head int) {
		stack := []int{head}
		unVisited.Reset(head)
		for len(stack) > 0 {
			now := stack[len(stack)-1]
			unVisited.Reset(now)
			next := FindFirstOfAnd(adjMatrix[now], unVisited)
			if next != -1 {
				unVisited.Reset(next)
				stack = append(stack, next)
			} else {
				stack = stack[:len(stack)-1]
				order = append(order, now)
			}
		}
	}
	rdfs := func(head, k int) {
		stack := []int{head}
		unVisited.Reset(head)
		for len(stack) > 0 {
			now := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			belong[now] = k
			for {
				next := FindFirstOfAnd(rAdjMatrix[now], unVisited)
				if next == -1 {
					break
				}
				stack = append(stack, next)
				unVisited.Reset(next)
			}
		}
	}

	for i := 0; i < n; i++ {
		unVisited.Set(i)
	}
	for i := 0; i < n; i++ {
		if unVisited.Has(i) {
			dfs(i)
		}
	}
	for i := 0; i < n; i++ {
		unVisited.Set(i)
	}
	for i := len(order) - 1; i >= 0; i-- {
		if unVisited.Has(order[i]) {
			rdfs(order[i], count)
			count++
		}
	}
	return
}

type Bitset []uint

func NewBitset(n int) Bitset { return make(Bitset, n>>6+1) } // (n+64-1)>>6, 注意 n=0 的情况，n>>6+1的写法更好

func (b Bitset) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 } // get
func (b Bitset) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b Bitset) Set(p int)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b Bitset) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) } // 置 0

// 返回第一个 1 的下标，若不存在则返回-1.
func FindFirstOfAnd(b1, b2 Bitset) int {
	minLen := min(len(b1), len(b2))
	for i := 0; i < minLen; i++ {
		and := b1[i] & b2[i]
		if and != 0 {
			return i<<6 | bits.TrailingZeros(and)
		}
	}
	return -1
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
