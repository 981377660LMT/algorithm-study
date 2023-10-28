package main

import "sort"

func main() {

	// 长度为 m 的 LIS 个数
	// 赤壁之战 https://www.acwing.com/problem/content/299/
	// 定义 dp[i][j] 表示 a[:j+1] 的长度为 i 且以 a[j] 结尾的 LIS
	// 则有 dp[i][j] = ∑dp[i-1][k]  (k<j && a[k]<a[j])
	// 注意到当 j 增加 1 时，只多了 k=j 这一个新决策，这样可以用树状数组来维护
	// 复杂度 O(mnlogn)
	countLIS := func(a []int, m int) int {
		// 将 a 离散化成从 2 开始的序列
		b := append([]int(nil), a...)
		sort.Ints(b)
		for i, v := range a {
			a[i] = sort.SearchInts(b, v) + 2
		}

		n := len(a)
		const mod int = 1e9 + 7
		tree := make([]int, n+2)
		add := func(i, val int) {
			for ; i < n+2; i += i & -i {
				tree[i] = (tree[i] + val) % mod
			}
		}
		sum := func(i int) (res int) {
			for ; i > 0; i &= i - 1 {
				res = (res + tree[i]) % mod
			}
			return
		}

		dp := make([][]int, m+1)
		for i := range dp {
			dp[i] = make([]int, n)
		}
		for i := 1; i <= m; i++ {
			tree = make([]int, n+2)
			if i == 1 {
				add(1, 1)
			}
			for j, v := range a {
				dp[i][j] = sum(v - 1)
				add(v, dp[i-1][j])
			}
		}
		ans := 0
		for _, v := range dp[m] {
			ans = (ans + v) % mod
		}
		return ans
	}

	// LIS 方案数 O(nlogn)
	// 原理见下面这题官方题解的方法二
	// LC673 https://leetcode-cn.com/problems/number-of-longest-increasing-subsequence/
	cntLis := func(a []int) int {
		g := [][]int{}   // 保留所有历史信息
		cnt := [][]int{} // 个数前缀和
		for _, v := range a {
			p := sort.Search(len(g), func(i int) bool { return g[i][len(g[i])-1] >= v })
			c := 1
			if p > 0 {
				// 根据 g[p-1] 来计算 cnt
				i := sort.Search(len(g[p-1]), func(i int) bool { return g[p-1][i] < v })
				c = cnt[p-1][len(cnt[p-1])-1] - cnt[p-1][i]
			}
			if p == len(g) {
				g = append(g, []int{v})
				cnt = append(cnt, []int{0, c})
			} else {
				g[p] = append(g[p], v)
				cnt[p] = append(cnt[p], cnt[p][len(cnt[p])-1]+c)
			}
		}
		c := cnt[len(cnt)-1]
		return c[len(c)-1]
	}

}
