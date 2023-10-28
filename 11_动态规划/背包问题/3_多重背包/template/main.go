package main

func main() {

	// 多重背包 - 未优化
	// 转换（价值主导）（由于要取 min 所以不能用二进制优化）https://codeforces.com/problemset/problem/922/E
	boundedKnapsack := func(stocks, values, weights []int, maxW int) int {
		n := len(stocks)
		f := make([][]int, n+1)
		for i := range f {
			f[i] = make([]int, maxW+1)
		}
		for i, num := range stocks {
			v, w := values[i], weights[i]
			for j := range f[i] {
				// 枚举选了 k=0,1,2,...num 个第 i 种物品
				for k := 0; k <= num && k*w <= j; k++ {
					f[i+1][j] = max(f[i+1][j], f[i][j-k*w]+k*v)
				}
			}
		}
		return f[n][maxW]
	}

	// 多重背包 - 优化 1 - 二进制优化
	boundedKnapsackBinary := func(stocks, values, weights []int, maxW int) int {
		f := make([]int, maxW+1)
		for i, num := range stocks {
			v, w := values[i], weights[i]
			for k1 := 1; num > 0; k1 <<= 1 {
				k := min(k1, num)
				for j := maxW; j >= k*w; j-- {
					f[j] = max(f[j], f[j-k*w]+k*v)
				}
				num -= k
			}
		}
		return f[maxW]
	}

	// 多重背包 - 优化 2 - 单调队列优化
	// 参考挑战 p.340
	// 时间复杂度 O(n*maxW)
	boundedKnapsackMonotoneQueue := func(stocks, values, weights []int, maxW int) int {
		f := make([]int, maxW+1)
		for i, num := range stocks {
			v, w := values[i], weights[i]
			for rem := 0; rem < w; rem++ { // 按照 j%w 的结果，分组转移，rem 表示 remainder
				type pair struct{ maxF, j int }
				q := []pair{}
				// 为什么压缩维度了还可以正着枚举？因为转移来源都存到单调队列里面了，正序倒序都可以
				// 并且这样相比倒着枚举，不需要先往队列里面塞 num 个数据，更加简洁
				for j := 0; j*w+rem <= maxW; j++ {
					t := f[j*w+rem] - j*v
					for len(q) > 0 && q[len(q)-1].maxF <= t {
						q = q[:len(q)-1] // 及时去掉无用数据
					}
					q = append(q, pair{t, j})
					// 本质是查表法，q[0].maxF 就表示 f[(j-1)*w+r]-(j-1)*v, f[(j-2)*w+r]-(j-2)*v, …… 这些转移来源的最大值
					f[j*w+rem] = q[0].maxF + j*v // 把物品个数视作两个 j 的差（前缀和思想）
					if j-q[0].j == num {         // 至多选 num 个物品
						q = q[1:] // 及时去掉无用数据
					}
				}
			}
		}
		return f[maxW]
	}

	// 多重背包 - 求方案数 - 同余前缀和优化
	// 讲解 https://leetcode.cn/problems/count-of-sub-multisets-with-bounded-sum/solution/duo-zhong-bei-bao-fang-an-shu-cong-po-su-f5ay/
	// LC2902 https://leetcode.cn/problems/count-of-sub-multisets-with-bounded-sum/
	// LC1155 https://leetcode.cn/problems/number-of-dice-rolls-with-target-sum/
	boundedKnapsackWays := func(a []int) []int {
		const mod = 1_000_000_007
		total := 0
		cnt := map[int]int{}
		for _, x := range a {
			total += x
			cnt[x]++
		}

		f := make([]int, total+1)
		f[0] = cnt[0] + 1
		delete(cnt, 0)

		maxJ := 0
		for x, c := range cnt {
			maxJ += x * c
			for j := x; j <= maxJ; j++ {
				f[j] = (f[j] + f[j-x]) % mod // 同余前缀和
			}
			for j := maxJ; j >= x*(c+1); j-- {
				f[j] = (f[j] - f[j-x*(c+1)] + mod) % mod
			}
		}
		return f
	}

}
