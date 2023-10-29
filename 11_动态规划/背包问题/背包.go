package copypasta

func _(min, max func(int, int) int, abs func(int) int) {

	/* 背包问题

	0-1 背包 0-1 Knapsack
	完全背包  Unbounded Knapsack
	多重背包  Bounded Knapsack

	这类问题可以从物品选择次序的无后效性入手
	子区间 -> 前缀和
	子序列 -> 背包
	https://en.wikipedia.org/wiki/Knapsack_problem
	https://codeforces.com/blog/entry/59606
	浅谈 ZKP 问题 https://www.luogu.com.cn/blog/xww666/qian-tan-zkp-wen-ti-gai-post
	另见 math_ntt.go 中的生成函数

	O(n√nw) 的 shuffle 做法（这里 w=max(wi)）
	https://codeforces.com/blog/entry/50036 https://codeforces.com/contest/739/problem/E
	https://arxiv.org/pdf/2308.11307.pdf
	http://acm.hdu.edu.cn/showproblem.php?pid=6804

	NOTE: 若求能否凑成 1,2,3,...,M，只需判断 dp[i] 是否为正 LC1049 https://leetcode-cn.com/problems/last-stone-weight-ii/
	套题 https://www.acwing.com/problem/
	混合背包 https://www.luogu.com.cn/problem/P1833
	*/

	// 0-1 背包 (n 个物品，背包容量为 maxW)   0-1 Knapsack
	// 状态：从前 i 个物品中选择若干个，当容量限制为 j 时能获得的最大价值和  i∈[0,n-1], j∈[0,maxW]
	// 初始值：f(0,j)=0  j∈[0,maxW]
	// 除开初始状态，每个状态有两个来源，决策为 max：
	// - 不选第 i 个物品：f(i-1,j) -> f(i,j)
	// - 选第 i 个物品：f(i-1,j-wi)+vi -> f(i,j)   j≥wi
	// 最优解为 f(n-1,maxW)
	// https://oi-wiki.org/dp/knapsack/
	// 模板题 https://www.luogu.com.cn/problem/P1048 https://atcoder.jp/contests/dp/tasks/dp_d
	// LC2291 https://leetcode.cn/problems/maximum-profit-from-trading-stocks/
	// 转换 LC494 https://leetcode.cn/problems/target-sum/
	//            https://atcoder.jp/contests/abc274/tasks/abc274_d
	// 转换 LC1049 https://leetcode-cn.com/problems/last-stone-weight-ii/
	// 转换 https://codeforces.com/problemset/problem/1381/B
	// 转换 https://codeforces.com/problemset/problem/1516/C
	// 转换 https://atcoder.jp/contests/dp/tasks/dp_x
	// 转换 https://leetcode.com/discuss/interview-question/2677093/Snowflake-oror-Tough-OA-question-oror-How-to-solve
	// 排序+转换 https://codeforces.com/problemset/problem/1203/F2
	// 正难则反 https://atcoder.jp/contests/tenka1-2019/tasks/tenka1_2019_d
	// 状压 LC1125 https://leetcode.cn/problems/smallest-sufficient-team/
	// 恰好组成 k 的数中能恰好组成哪些数 https://codeforces.com/problemset/problem/687/C
	// 转移对象是下标 https://codeforces.com/edu/course/2/lesson/9/3/practice/contest/307094/problem/I
	// - dp[i][j] 表示前 i 个数，凑成 j 的所有方案中，最小下标的最大值
	// 转移对象是下标 https://codeforces.com/problemset/problem/981/E
	// 打印方案 https://codeforces.com/problemset/problem/864/E
	// 变形，需要多加一个维度 https://atcoder.jp/contests/abc275/tasks/abc275_f
	// 贡献 https://atcoder.jp/contests/abc159/tasks/abc159_f
	// NOIP06·提高 金明的预算方案（也可以用树上背包做）https://www.luogu.com.cn/problem/P1064
	// EXTRA: 恰好装满（相当于方案数不为 0）LC416 https://leetcode-cn.com/problems/partition-equal-subset-sum/
	//        必须定义成恰好装满（紫书例题 9-5，UVa 12563）https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=441&page=show_problem&problem=4008
	// EXTRA: 背包容量为 0 https://codeforces.com/problemset/problem/366/C
	// EXTRA: 二维费用 https://www.acwing.com/problem/content/8/ https://www.luogu.com.cn/problem/P1507 LC474 https://leetcode-cn.com/problems/ones-and-zeroes/
	// EXTRA: 把一个维度转换成 DP 的定义 https://codeforces.com/problemset/problem/837/D
	// EXTRA: 离散化背包 https://codeforces.com/contest/366/submission/61452111
	zeroOneKnapsack := func(values, weights []int, maxW int) int {
		f := make([]int, maxW+1)
		for i, w := range weights {
			v := values[i]
			// 这里 j 的初始值可以优化成前 i 个物品的重量之和（但不能超过 maxW）
			for j := maxW; j >= w; j-- {
				f[j] = max(f[j], f[j-w]+v)
			}
		}
		return f[maxW]
	}

	// 0-1 背包 EXTRA: 恰好装满
	// https://leetcode.cn/contest/sf-tech/problems/cINqyA/
	// 转换 二维费用 https://codeforces.com/problemset/problem/730/J
	zeroOneKnapsackExactlyFull := func(values, weights []int, maxW int) {
		f := make([]int, maxW+1)
		for i := range f {
			f[i] = -1e9 // -1e18
		}
		f[0] = 0
		for i, w := range weights {
			v := values[i]
			for j := maxW; j >= w; j-- {
				f[j] = max(f[j], f[j-w]+v)
			}
		}
		for i := maxW; i >= 0; i-- {
			if f[i] >= 0 { // 能恰好装满 i，此时背包物品价值和的最大值是 dp[i]
				// ...
			}
		}
	}

	// 0-1 背包 EXTRA: 至少装入重量和为 maxW 的物品，求价值和的最小值 https://www.luogu.com.cn/problem/P4377
	// f[0] 表示至少为 0 的情况，也表示没有任何约束的情况
	// 比如选第 i 个物品后容量 <=0 了，那就表示前面的 i-1 个物品可以不受约束地随意选或不选了
	// 转换 https://codeforces.com/problemset/problem/19/B LC2742 https://leetcode.cn/problems/painting-the-walls/
	// 二维费用的情况+价值最小 https://ac.nowcoder.com/acm/contest/6218/C
	zeroOneKnapsackAtLeastFillUp := func(values, weights []int, maxW int) int {
		f := make([]int, maxW+1)
		for i := range f {
			f[i] = 1e9 // 1e18
		}
		f[0] = 0
		for i, v := range values {
			w := weights[i]
			for j := maxW; j >= 0; j-- {
				f[j] = min(f[j], f[max(j-w, 0)]+v) // max(j-w, 0) 蕴含了「至少」
			}
		}

		{
			// 另一种写法
			for i, v := range values {
				w := weights[i]
				for j := maxW; j >= 0; j-- {
					k := min(j+w, maxW)
					f[k] = min(f[k], f[j]+v)
				}
			}
		}

		return f[maxW]
	}

	// 0-1 背包 EXTRA: 从序列 a 中选若干个数，使其总和为 sum 的方案数
	// 常见题目是算严格分拆（必须用不同数字）  https://leetcode.cn/problems/ways-to-express-an-integer-as-sum-of-powers/
	// - https://oeis.org/A000009
	// NOTE: 1,1,1,...1(32个1),s-32,s-31,...,s 可以让方案数恰好为 2^32
	// 二维+上限+下限 LC879 https://leetcode-cn.com/problems/profitable-schemes/
	// https://atcoder.jp/contests/arc060/tasks/arc060_a
	// https://codeforces.com/problemset/problem/1673/C
	// 转换 https://atcoder.jp/contests/abc169/tasks/abc169_f
	// 转换 https://codeforces.com/problemset/problem/478/D
	// 转换 LC494 https://leetcode-cn.com/problems/target-sum/
	// 转换 LC1434 https://leetcode-cn.com/problems/number-of-ways-to-wear-different-hats-to-each-other/
	// 由于顺序不同也算方案，所以这题需要正序递推 LC377 https://leetcode-cn.com/problems/combination-sum-iv/
	zeroOneWaysToSum := func(a []int, sum int) int {
		f := make([]int, sum+1)
		f[0] = 1
		for _, v := range a {
			for j := sum; j >= v; j-- {
				f[j] += f[j-v] // % mod
			}
		}
		return f[sum]
	}

	// 0-1 背包 EXTRA: 打印字典序最小的方案
	// 倒序遍历物品，同时用 fa 数组记录转移来源，这样跑完 DP 后，从第一个物品开始即可得到字典序最小的方案
	// https://www.acwing.com/problem/content/description/12/
	zeroOneKnapsackLexicographicallySmallestResult := func(values, weights []int, maxW int) (ans []int) {
		n := len(values)
		f := make([]int, maxW+1) // fill
		//f[0] = 0
		fa := make([][]int, n)
		for i := n - 1; i >= 0; i-- {
			fa[i] = make([]int, maxW+1)
			for j := range fa[i] {
				fa[i][j] = j // 注意：<w 的转移来源也要标上！
			}
			v, w := values[i], weights[i]
			for j := maxW; j >= w; j-- {
				if f[j-w]+v >= f[j] { // 注意这里要取等号，从而保证尽可能地从字典序最小的方案转移过来
					f[j] = f[j-w] + v
					fa[i][j] = j - w
				}
			}
		}
		for i, j := 0, maxW; i < n; {
			if fa[i][j] == j { // &&  weights[i] > 0      考虑重量为 0 的情况，必须都选上
				i++
			} else {
				ans = append(ans, i+1) // 下标从 1 开始
				j = fa[i][j]
				i++ // 完全背包的情况，这行去掉
			}
		}
		return
	}

	// 0-1 背包 EXTRA: 价值主导的 0-1 背包
	// 适用于背包容量很大，但是物品价值不高的情况
	// 把重量看成价值，价值看成重量，求同等价值下能得到的最小重量，若该最小重量不超过背包容量，则该价值合法。所有合法价值的最大值即为答案
	// 时间复杂度 O(n * sum(values)) 或 O(n^2 * maxV)
	// https://atcoder.jp/contests/dp/tasks/dp_e
	// https://codeforces.com/contest/1650/problem/F
	zeroOneKnapsackByValue := func(values, weights []int, maxW int) int {
		totValue := 0
		for _, v := range values {
			totValue += v
		}
		f := make([]int, totValue+1)
		for i := range f {
			f[i] = 1e18
		}
		f[0] = 0
		totValue = 0
		for i, v := range values {
			w := weights[i]
			totValue += v
			for j := totValue; j >= v; j-- {
				f[j] = min(f[j], f[j-v]+w)
			}
		}
		for i := totValue; ; i-- {
			if f[i] <= maxW {
				return i
			}
		}
	}

	// todo 回退背包

	// 完全背包   Unbounded Knapsack
	// 更快的做法 https://www.zhihu.com/question/26547156/answer/1181239468
	// https://github.com/hqztrue/shared_materials/blob/master/codeforces/101064%20L.%20The%20Knapsack%20problem%20156ms_short.cpp
	// https://www.luogu.com.cn/problem/P1616
	// 至少 https://www.luogu.com.cn/problem/P2918
	// 恰好装满 LC322 https://leetcode-cn.com/problems/coin-change/
	// EXTRA: 恰好装满+打印方案 LC1449 https://leetcode-cn.com/problems/form-largest-integer-with-digits-that-add-up-to-target/
	// 【脑洞】求极限：lim_{maxW->∞} dp[maxW]/maxW
	unboundedKnapsack := func(values, weights []int, maxW int) int {
		f := make([]int, maxW+1) // fill
		//f[0] = 0
		for i, v := range values {
			w := weights[i]
			for j := w; j <= maxW; j++ {
				f[j] = max(f[j], f[j-w]+v)
			}
		}
		return f[maxW]
	}

	// 完全背包 EXTRA: 方案数
	// LC518 https://leetcode-cn.com/problems/coin-change-ii/
	// https://codeforces.com/problemset/problem/1673/C
	// https://www.luogu.com.cn/problem/P1832
	// https://www.luogu.com.cn/problem/P6205（需要高精）
	// 类似完全背包但是枚举的思路不一样 LC377 https://leetcode-cn.com/problems/combination-sum-iv/
	unboundedWaysToSum := func(a []int, total int) int {
		f := make([]int, total+1)
		f[0] = 1
		for _, v := range a {
			for j := v; j <= total; j++ {
				f[j] += f[j-v] // % mod
			}
		}
		return f[total]
	}

	// 完全背包 EXTRA: 二维费用方案数
	// 注意：「恰好使用 m 个物品」这个条件要当成一种费用来看待
	// https://codeforces.com/problemset/problem/543/A

	// 多重背包   Bounded Knapsack
	// 模板题 https://codeforces.com/problemset/problem/106/C
	//       https://www.luogu.com.cn/problem/P1776
	// todo 多重背包+完全背包 https://www.luogu.com.cn/problem/P1782 https://www.luogu.com.cn/problem/P1833 https://www.luogu.com.cn/problem/P2851
	// http://acm.hdu.edu.cn/showproblem.php?pid=2844 http://poj.org/problem?id=1742
	// https://www.luogu.com.cn/problem/P6771 http://poj.org/problem?id=2392
	// https://codeforces.com/contest/999/problem/F
	// todo 打印方案

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

	// 分组背包·每组至多选一个（恰好选一个见后面）
	// https://www.acwing.com/problem/content/9/
	// https://www.luogu.com.cn/problem/P1757
	// LC2218 https://leetcode.cn/problems/maximum-value-of-k-coins-from-piles/
	// https://codeforces.com/problemset/problem/148/E
	// todo 进一步优化 https://codeforces.com/problemset/problem/1442/D
	// 方案数（可以用前缀和优化）https://www.luogu.com.cn/problem/P1077
	// 方案数 LC2585 https://leetcode.cn/problems/number-of-ways-to-earn-points/
	type item struct{ v, w int }
	groupKnapsack := func(groups [][]item, maxW int) int {
		dp := make([]int, maxW+1)
		for _, g := range groups {
			// 这里 j 的初始值可以优化成前 i 个组的每组最大重量之和（但不能超过 maxW）
			for j := maxW; j >= 0; j-- {
				for _, it := range g {
					if v, w := it.v, it.w; w <= j {
						dp[j] = max(dp[j], dp[j-w]+v) // 如果 it.w 可能为 0 则需要用 dp[2][] 来滚动（或者保证每组至多一个 0 且 0 在该组最前面）
					}
				}
			}
		}
		return dp[maxW]
	}

	// todo 撤销计数
	//  https://leetcode.cn/circle/article/YnZBve/

	// 分组背包·每组恰好选一个
	// 允许物品重量为 0
	// https://atcoder.jp/contests/abc240/tasks/abc240_c
	// LC1981 https://leetcode-cn.com/problems/minimize-the-difference-between-target-and-chosen-elements/
	// 与二分图染色结合 https://codeforces.com/problemset/problem/1354/E
	// 转换 https://codeforces.com/problemset/problem/1637/D
	groupKnapsackFill := func(groups [][]int, maxW int) []bool {
		dp := make([]bool, maxW+1) // dp[i][j] 表示能否从前 i 组物品中选出重量恰好为 j 的，且每组都恰好选一个物品
		dp[0] = true
		for _, g := range groups {
		next:
			for j := maxW; j >= 0; j-- { // 这里 j 的初始值可以优化至前 i 组的最大元素值之和
				for _, w := range g {
					if w <= j && dp[j-w] {
						dp[j] = true
						continue next
					}
				}
				dp[j] = false // 由于我们是滚动数组的写法，dp[i][j] 无法满足时要标记成 false
			}
		}
		return dp // dp[j] 表示从每组恰好选一个，能否凑成重量 j
	}

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

	// 无向图简单环数量
	// https://blog.csdn.net/fangzhenpeng/article/details/49078233
	// https://codeforces.com/problemset/problem/11/D
	countCycle := func(g [][]int, n, m int) int {
		ans := 0
		// 取集合 s 的最小值作为起点
		dp := make([][]int, 1<<n)
		for i := range dp {
			dp[i] = make([]int, n)
		}
		for i := 0; i < n; i++ {
			dp[1<<i][i] = 1
		}
		for s := range dp {
			for v, dv := range dp[s] {
				if dv == 0 {
					continue
				}
				for _, w := range g[v] {
					if 1<<w < s&-s {
						continue
					}
					if 1<<w&s == 0 {
						dp[s|1<<w][w] += dv
					} else if 1<<w == s&-s {
						ans += dv
					}
				}
			}
		}
		return ans - m/2
	}

	// 树上最大独立集
	// 返回最大点权和（最大独立集的情形即所有点权均为一）
	// 每个点有选和不选两种决策，接受子树转移时，选的决策只能加上不选子树，而不选的决策可以加上 max{不选子树, 选子树}
	// https://brooksj.com/2019/06/20/%E6%A0%91%E7%9A%84%E6%9C%80%E5%B0%8F%E6%94%AF%E9%85%8D%E9%9B%86%EF%BC%8C%E6%9C%80%E5%B0%8F%E7%82%B9%E8%A6%86%E7%9B%96%E9%9B%86%EF%BC%8C%E6%9C%80%E5%A4%A7%E7%82%B9%E7%8B%AC%E7%AB%8B%E9%9B%86/
	// https://stackoverflow.com/questions/13544240/algorithm-to-find-max-independent-set-in-a-tree
	// 经典题：没有上司的舞会 LC337 https://leetcode.cn/problems/house-robber-iii/ https://www.luogu.com.cn/problem/P1352 https://ac.nowcoder.com/acm/problem/51178
	// 变形 LC2646 https://leetcode.cn/problems/minimize-the-total-price-of-the-trips/
	// 边权独立集 https://leetcode.cn/problems/choose-edges-to-maximize-score-in-a-tree/description/
	// 方案是否唯一 Tehran06，紫书例题 9-13，UVa 1220 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=247&page=show_problem&problem=3661
	maxIndependentSetOfTree := func(n int, g [][]int, a []int) int { // 无根树
		var f func(int, int) (notChosen, chosen int)
		f = func(v, fa int) (notChosen, chosen int) {
			chosen = a[v] // 1
			for _, w := range g[v] {
				if w != fa {
					nc, c := f(w, v)
					notChosen += max(nc, c)
					chosen += nc
				}
			}
			return
		}
		nc, c := f(0, -1)
		return max(nc, c)
	}

	// 树上最小顶点覆盖（每条边至少有一个端点被覆盖，也就是在已选的点集中）
	// 代码和树上最大独立集类似
	// https://brooksj.com/2019/06/20/%E6%A0%91%E7%9A%84%E6%9C%80%E5%B0%8F%E6%94%AF%E9%85%8D%E9%9B%86%EF%BC%8C%E6%9C%80%E5%B0%8F%E7%82%B9%E8%A6%86%E7%9B%96%E9%9B%86%EF%BC%8C%E6%9C%80%E5%A4%A7%E7%82%B9%E7%8B%AC%E7%AB%8B%E9%9B%86/
	// 经典题：战略游戏 https://www.luogu.com.cn/problem/P2016
	// 训练指南第一章例题 30，UVa10859 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=20&page=show_problem&problem=1800
	// - 求最小顶点覆盖，以及所有最小顶点覆盖中，两端点都被覆盖的边的最大个数
	// 构造 https://codeforces.com/problemset/problem/959/C
	minVertexCoverOfTree := func(n int, g [][]int, a []int) int { // 无根树
		var f func(int, int) (notChosen, chosen int)
		f = func(v, fa int) (notChosen, chosen int) {
			chosen = a[v] // 1
			for _, w := range g[v] {
				if w != fa {
					nc, c := f(w, v)
					notChosen += c
					chosen += min(nc, c)
				}
			}
			return
		}
		nc, c := f(0, -1)
		return min(nc, c)
	}

	// 树上最小支配集
	// 返回最小点权和（最小支配集的情形即所有点权均为一）
	// 下面的定义省去了（……时的最小支配集的元素个数）   w 为 i 的儿子
	// 视频讲解：https://www.bilibili.com/video/BV1oF411U7qL/
	// dp[i][0]：i 属于支配集 = a[i]+∑min(dp[w][0],dp[w][1],dp[w][2])
	// dp[i][1]：i 不属于支配集，且被儿子支配 = ∑min(dp[w][0],dp[w][1]) + 如果全选 dp[w][1] 则补上 min{dp[w][0]-dp[w][1]}
	// dp[i][2]：i 不属于支配集，且被父亲支配 = ∑min(dp[w][0],dp[w][1])
	// https://brooksj.com/2019/06/20/%E6%A0%91%E7%9A%84%E6%9C%80%E5%B0%8F%E6%94%AF%E9%85%8D%E9%9B%86%EF%BC%8C%E6%9C%80%E5%B0%8F%E7%82%B9%E8%A6%86%E7%9B%96%E9%9B%86%EF%BC%8C%E6%9C%80%E5%A4%A7%E7%82%B9%E7%8B%AC%E7%AB%8B%E9%9B%86/
	//
	// 监控二叉树 LC968 https://leetcode-cn.com/problems/binary-tree-cameras/
	// - https://codeforces.com/problemset/problem/1029/E
	// 保安站岗 https://www.luogu.com.cn/problem/P2458
	// 手机网络 https://www.luogu.com.cn/problem/P2899
	// https://ac.nowcoder.com/acm/problem/24953
	// todo EXTRA: 消防局的设立（支配距离为 2） https://www.luogu.com.cn/problem/P2279
	// todo EXTRA: 将军令（支配距离为 k） https://www.luogu.com.cn/problem/P3942
	//                                https://atcoder.jp/contests/arc116/tasks/arc116_e
	minDominatingSetOfTree := func(n int, g [][]int, a []int) int { // 无根树
		const inf int = 1e18
		var f func(int, int) (chosen, bySon, byFa int)
		f = func(v, fa int) (chosen, bySon, byFa int) {
			chosen = a[v] // 1
			extra := inf
			for _, w := range g[v] {
				if w != fa {
					c, bs, bf := f(w, v)
					m := min(c, bs)
					chosen += min(m, bf)
					bySon += m
					byFa += m
					extra = min(extra, c-bs)
				}
			}
			bySon += max(extra, 0)
			return
		}
		chosen, bySon, _ := f(0, -1)
		return min(chosen, bySon)
	}

	// EXTRA: 每个被支配的点，仅被一个点支配
	// Kaoshiung06，紫书例题 9-14，UVa 1218 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=247&page=show_problem&problem=3659

	// 树上最大匹配
	// g[v] = ∑{max(f[son],g[son])}
	// f[v] = max{1+g[son]+g[v]−max(f[son],g[son])}
	// https://codeforces.com/blog/entry/2059
	// https://blog.csdn.net/lycheng1215/article/details/78368002
	// https://vijos.org/p/1892
	maxMatchingOfTree := func(n int, g [][]int) int { // 无根树
		cover, nonCover := make([]int, n), make([]int, n)
		var f func(int, int)
		f = func(v, fa int) {
			for _, w := range g[v] {
				if w != fa {
					f(w, v)
					nonCover[v] += max(cover[w], nonCover[w])
				}
			}
			for _, w := range g[v] {
				cover[v] = max(cover[v], 1+nonCover[w]+nonCover[v]-max(cover[w], nonCover[w]))
			}
		}
		f(0, -1)
		return max(cover[0], nonCover[0])
	}

	// 树上所有路径的位运算与(&)的和
	// 单个点也算路径
	// 解法：对每一位，统计仅含 1 的路径个数
	// a[i] <= 2^20
	// https://ac.nowcoder.com/acm/contest/10167/C
	andPathSum := func(g [][]int, a []int) int {
		const mx = 21
		ans := 0
		for i := 0; i < mx; i++ {
			cntOnePath := 0
			var f func(v, fa int) int
			f = func(v, fa int) int {
				one := a[v] >> i & 1
				cntOnePath += one
				for _, w := range g[v] {
					if w != fa {
						o := f(w, v)
						if one > 0 {
							cntOnePath += one * o
							one += o
						}
					}
				}
				return one
			}
			f(0, -1)
			ans += 1 << i * cntOnePath
		}

		{
			// 另一种做法是对每一位，用并查集求出 1 组成的连通分量，每个连通分量对答案的贡献是 sz*(sz+1)/2
			n := len(a)
			fa := make([]int, n)
			var find func(int) int
			find = func(x int) int {
				if fa[x] != x {
					fa[x] = find(fa[x])
				}
				return fa[x]
			}
			merge := func(from, to int) { fa[find(from)] = find(to) }

			ans := 0
			for i := 0; i < mx; i++ {
				for j := range fa {
					fa[j] = j
				}
				sz := make([]int, n)
				for v, vs := range g {
					for _, w := range vs {
						if a[v]>>i&1 > 0 && a[w]>>i&1 > 0 {
							merge(v, w)
						}
					}
				}
				for j := 0; j < n; j++ {
					sz[find(j)]++
				}
				for j, f := range fa {
					if j == f && a[j]>>i&1 > 0 {
						ans += 1 << i * sz[j] * (sz[j] + 1) / 2
					}
				}
			}
		}
		return ans
	}

	// 树上所有路径的位运算或(|)的和
	// 单个点也算路径
	// 做法和上面类似，求出仅含 0 的路径的个数，然后用路径总数 n*(n+1) 减去该个数就得到了包含至少一个 1 的路径个数
	// 也可以用并查集求出 0 组成的连通分量

	// 树上所有路径的位运算异或(^)的和
	// 原题失效了，只找到几个题解可以参考 https://www.cnblogs.com/kuronekonano/p/11135742.html https://blog.csdn.net/qq_36876305/article/details/80060491
	// 上面链接是边权，这里改成点权，且路径至少有两个点
	// 解法：由于任意路径异或和可以用从根节点出发的路径异或和表示
	// 对每一位，统计从根节点出发的路径异或和在该位上的 0 的个数和 1 的个数，
	// 只有当 0 与 1 异或时才对答案有贡献，所以贡献即为这两个个数之积
	xorPathSum := func(g [][]int, a []int) int {
		n := len(a)
		const mx = 30
		cnt := [mx]int{}
		xor := 0
		var f func(v, fa int)
		f = func(v, fa int) {
			xor ^= a[v]
			for _, w := range g[v] {
				if w != fa {
					f(w, v)
				}
			}
			for i := 0; i < mx; i++ {
				cnt[i] += xor >> i & 1
			}
			xor ^= a[v]
		}
		f(0, -1)
		ans := 0
		for i, c := range cnt {
			ans += 1 << i * c * (n - c)
		}
		return ans
	}

	// 树上所有路径的位运算异或(^)的异或和
	// 这里的路径至少有两个点
	// 方法是考虑每个点出现在多少条路径上，若数目为奇数则对答案有贡献
	// 路径分两种情况，一种是没有父节点参与的，树形 DP 一下就行了；另一种是父节点参与的，个数就是 子树*(n-子树)
	// https://ac.nowcoder.com/acm/contest/272/B
	xorPathXorSum := func(g [][]int, a []int) int {
		n := len(a)
		ans := 0
		var f func(v, fa int) int
		f = func(v, fa int) int {
			cnt := 0
			sz := 1
			for _, w := range g[v] {
				if w != fa {
					s := f(w, v)
					cnt += sz * s
					sz += s
				}
			}
			cnt += sz * (n - sz)
			// 若一个点也算路径，那就再加一。或者在递归结束后把 ans^=a[0]^...^a[n-1]
			if cnt&1 > 0 {
				ans ^= a[v]
			}
			return sz
		}
		f(0, -1)
		return ans
	}

	_ = []any{
		prefixSumDP, mapDP,
		maxSubarraySum, maxSubarraySumWithRange, maxTwoSubarraySum,
		maxAlternatingSumDP, maxAlternatingSumGreedy,
		minCostSorted,
		lcs, lcsPath, longestPalindromeSubsequence,
		lisSlow, lis, lisAll, cntLis, lcis, lcisPath, countLIS,
		distinctSubsequence,
		lcp,
		palindromeO1Space, isPalindrome, minPalindromeCut,

		zeroOneKnapsack, zeroOneKnapsackExactlyFull, zeroOneKnapsackAtLeastFillUp, zeroOneWaysToSum, zeroOneKnapsackLexicographicallySmallestResult, zeroOneKnapsackByValue,
		unboundedKnapsack, unboundedWaysToSum,
		boundedKnapsack, boundedKnapsackBinary, boundedKnapsackMonotoneQueue, boundedKnapsackWays,
		groupKnapsack, groupKnapsackFill,
		treeKnapsack,

		mergeStones, countPalindromes,

		permDP, permDP2, tsp, countCycle, subsubDP, subsubDPMemo, sosDP, plugDP,

		digitDP, digitDP2, calcSum, kth666,

		binaryLifting,

		cht,

		diameter, countDiameter, countPath, countVerticesOnDiameter, maxPathSum,
		maxIndependentSetOfTree, minVertexCoverOfTree, minDominatingSetOfTree, maxMatchingOfTree,
		sumOfDistancesInTree, rerootDP,
		andPathSum, xorPathSum, xorPathXorSum,
	}
}
