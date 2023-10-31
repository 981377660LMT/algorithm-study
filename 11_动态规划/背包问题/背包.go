package copypasta

func _(min, max func(int, int) int, abs func(int) int) {

	// 练习：
	// https://www.luogu.com.cn/training/5197
	// https://www.luogu.com.cn/training/8917
	// https://www.luogu.com.cn/training/231055

	// 树上背包/树形背包/依赖背包
	// todo 树上背包的上下界优化 https://ouuan.github.io/post/%E6%A0%91%E4%B8%8A%E8%83%8C%E5%8C%85%E7%9A%84%E4%B8%8A%E4%B8%8B%E7%95%8C%E4%BC%98%E5%8C%96/
	//   子树合并背包的复杂度证明 https://blog.csdn.net/lyd_7_29/article/details/79854245
	//   复杂度 https://leetcode.cn/circle/discuss/t7l62c/
	//   https://www.cnblogs.com/shaojia/p/15520224.html
	//   https://snuke.hatenablog.com/entry/2019/01/15/211812
	//   复杂度优化 https://loj.ac/d/3144
	//   https://zhuanlan.zhihu.com/p/103813542
	//
	// todo https://loj.ac/p/160
	//   https://www.luogu.com.cn/problem/P2014 https://www.acwing.com/problem/content/10/ https://www.acwing.com/problem/content/288/
	//   加强版 https://www.luogu.com.cn/problem/U53204
	//   https://www.luogu.com.cn/problem/P1272
	//   加强版 https://www.luogu.com.cn/problem/U53878
	//   https://www.luogu.com.cn/problem/P3177
	// NOIP06·提高 金明的预算方案 https://www.luogu.com.cn/problem/P1064
	treeKnapsack := func(g [][]int, items []item, root, maxW int) int {
		var f func(int) []int
		f = func(v int) []int {
			it := items[v]
			dp := make([]int, maxW+1)
			for i := it.w; i <= maxW; i++ {
				dp[i] = it.v // 根节点必须选
			}
			for _, to := range g[v] {
				dt := f(to)
				for j := maxW; j >= it.w; j-- {
					// 类似分组背包，枚举分给子树 to 的容量 w，对应的子树的最大价值为 dt[w]
					// w 不可超过 j-it.w，否则无法选择根节点
					for w := 0; w <= j-it.w; w++ {
						dp[j] = max(dp[j], dp[j-w]+dt[w])
					}
				}
			}
			return dp
		}
		return f(root)[maxW]
	}

}
