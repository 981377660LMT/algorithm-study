package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// abc201_e_1()
	abc201_e_2()
}

// E - Xor Distances (树上异或距离)
// https://atcoder.jp/contests/abc201/tasks/abc201_e
// 给定一棵树.
// 求树上所有点对的异或距离之和模 10^9+7.
//
// https://atcoder.jp/contests/abc201/editorial/1827
// 以任意一点root为根.
// !dist(i, j) = dist(lca, i) ^ dist(lca, j)
//
//	= dist(lca, i) ^ dist(lca, j) ^ dist(lca, root) ^ dist(lca, root)
//	= (dist(root, lca) ^ dist(lca, i)) ^ (dist(lca, j) ^ dist(root, lca))
//
// !           = dist(root, i) ^ dist(root, j)
// 等价于：
// !树上`两点的异或距离之和`等于`两点到根的异或距离之和`.
//
// 先求出所有点到根的异或距离.
// 然后按位统计每一位的贡献.
func abc201_e_1() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7

	var n int
	fmt.Fscan(in, &n)
	tree := make([][][2]int, n)
	for i := 1; i < n; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		a--
		b--
		tree[a] = append(tree[a], [2]int{b, c})
		tree[b] = append(tree[b], [2]int{a, c})
	}

	xorToRoot := make([]int, n)
	var dfs func(v, p, xor int)
	dfs = func(v, p, xor int) {
		xorToRoot[v] = xor
		for _, e := range tree[v] {
			if e[0] == p {
				continue
			}
			dfs(e[0], v, xor^e[1])
		}
	}
	dfs(0, -1, 0)

	res := 0
	for i := 0; i < 60; i++ {
		counter := [2]int{}
		for _, xor := range xorToRoot {
			counter[xor>>i&1]++
		}
		pair := counter[0] * counter[1] % MOD
		res += (1 << i) % MOD * pair % MOD
		res %= MOD
	}
	fmt.Fprintln(out, res)
}

// 按位统计每一位的贡献.
// => 树的边权为0或1.求点对距离为奇数的点对个数.
// !=> 可以将树用黑白染色,然后求黑*白点对的个数.
func abc201_e_2() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7

	var n int
	fmt.Fscan(in, &n)
	tree := make([][][2]int, n)
	for i := 1; i < n; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		a--
		b--
		tree[a] = append(tree[a], [2]int{b, c})
		tree[b] = append(tree[b], [2]int{a, c})
	}

	counter := [60]int{}
	var dfs func(v, p, xor int)
	dfs = func(v, p, xor int) {
		for i := 0; i < len(counter); i++ {
			counter[i] += (xor >> i) & 1
		}
		for _, e := range tree[v] {
			if e[0] == p {
				continue
			}
			dfs(e[0], v, xor^e[1])
		}
	}
	dfs(0, -1, 0)

	res := 0
	for i := 0; i < len(counter); i++ {
		a, b := counter[i], n-counter[i]
		pair := a * b % MOD
		res += (1 << i) % MOD * pair % MOD
		res %= MOD
	}
	fmt.Fprintln(out, res)
}
