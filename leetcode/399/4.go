package main

// 给你一个整数数组 nums 和一个二维数组 queries，其中 queries[i] = [posi, xi]。

// 对于每个查询 i，首先将 nums[posi] 设置为 xi，然后计算查询 i 的答案，该答案为 nums 中 不包含相邻元素 的子序列的 最大 和。

// 返回所有查询的答案之和。

// 由于最终答案可能非常大，返回其对 109 + 7 取余 的结果。

// 子序列 是指从另一个数组中删除一些或不删除元素而不改变剩余元素顺序得到的数组。

// in := bufio.NewReader(os.Stdin)
// 	out := bufio.NewWriter(os.Stdout)
// 	defer out.Flush()

// 	const INF int = 1e18

// 	var n, q int32
// 	fmt.Fscan(in, &n, &q)
// 	weights := make([]int, n)
// 	for i := range weights {
// 		fmt.Fscan(in, &weights[i])
// 	}
// 	tree := make([][]int32, n)
// 	for i := int32(0); i < n-1; i++ {
// 		var u, v int32
// 		fmt.Fscan(in, &u, &v)
// 		u, v = u-1, v-1
// 		tree[u] = append(tree[u], v)
// 		tree[v] = append(tree[v], u)
// 	}

// 	dp0, dp1 := make([]int, n), make([]int, n)
// 	parent := make([]int32, n)
// 	for i := int32(0); i < n; i++ {
// 		parent[i] = -1
// 	}

// 	var dfs func(cur, pre int32)
// 	dfs = func(cur, pre int32) {
// 		dp0[cur], dp1[cur] = 0, weights[cur]
// 		parent[cur] = pre
// 		for _, next := range tree[cur] {
// 			if next == pre {
// 				continue
// 			}
// 			dfs(next, cur)
// 			dp0[cur] += max(dp1[next], dp0[next])
// 			dp1[cur] += dp0[next]
// 		}
// 	}
// 	dfs(0, -1)

// 	update := func(v int32) {
// 		// 重算所有祖先
// 		for ; v != -1; v = parent[v] {
// 			dp0[v], dp1[v] = 0, weights[v]
// 			for _, next := range tree[v] {
// 				if next == parent[v] {
// 					continue
// 				}
// 				dp0[v] += max(dp1[next], dp0[next])
// 				dp1[v] += dp0[next]
// 			}
// 		}
// 	}

//	for i := int32(0); i < q; i++ {
//		var x, y int
//		fmt.Fscan(in, &x, &y)
//		x--
//		weights[x] = y
//		update(int32(x))
//		fmt.Fprintln(out, max(dp0[0], dp1[0]))
//	}

const MOD int = 1e9 + 7

func maximumSumSubsequence(nums []int, queries [][]int) int {
	dp0, dp1 := make([]int, len(nums)), make([]int, len(nums))
	{
		dp0[0], dp1[0] = 0, nums[0]
		for i := 1; i < len(nums); i++ {
			dp0[i] += max(dp0[i-1], dp1[i-1])
			dp1[i] += nums[i] + dp0[i-1]
		}
	}

	update := func(pos, v int) {
		nums[pos] = v
		for i := pos; i < len(nums); i++ {
			preDp0, preDp1 := dp0[i], dp1[i]
			dp0[i], dp1[i] = 0, nums[i]
			if i > 0 {
				dp0[i] += max(dp0[i-1], dp1[i-1])
				dp1[i] += dp0[i-1]
			}
			if dp0[i] == preDp0 && dp1[i] == preDp1 {
				break
			}
		}
	}

	res := 0
	for _, query := range queries {
		pos, v := query[0], query[1]
		update(pos, v)
		res += max(dp0[len(nums)-1], dp1[len(nums)-1])
	}
	return res % MOD
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
