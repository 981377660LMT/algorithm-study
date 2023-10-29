package main

// 树上最大独立集
// 返回最大点权和（最大独立集的情形即所有点权均为一）
// 每个点有选和不选两种决策，接受子树转移时，选的决策只能加上不选子树，而不选的决策可以加上 max{不选子树, 选子树}
// https://brooksj.com/2019/06/20/%E6%A0%91%E7%9A%84%E6%9C%80%E5%B0%8F%E6%94%AF%E9%85%8D%E9%9B%86%EF%BC%8C%E6%9C%80%E5%B0%8F%E7%82%B9%E8%A6%86%E7%9B%96%E9%9B%86%EF%BC%8C%E6%9C%80%E5%A4%A7%E7%82%B9%E7%8B%AC%E7%AB%8B%E9%9B%86/
// https://stackoverflow.com/questions/13544240/algorithm-to-find-max-independent-set-in-a-tree
// 经典题：没有上司的舞会 LC337 https://leetcode.cn/problems/house-robber-iii/ https://www.luogu.com.cn/problem/P1352 https://ac.nowcoder.com/acm/problem/51178
// 变形 LC2646 https://leetcode.cn/problems/minimize-the-total-price-of-the-trips/
// 边权独立集 https://leetcode.cn/problems/choose-edges-to-maximize-score-in-a-tree/description/
// 方案是否唯一 Tehran06，紫书例题 9-13，UVa 1220 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=247&page=show_problem&problem=3661
func MaxIndependentSetOfTree(tree [][]int, values []int) int { // 无根树
	var dfs func(int, int) (notChosen, chosen int)
	dfs = func(cur, pre int) (notChosen, chosen int) {
		chosen = values[cur] // 1
		for _, next := range tree[cur] {
			if next != pre {
				nc, c := dfs(next, cur)
				notChosen += max(nc, c)
				chosen += nc
			}
		}
		return
	}
	nc, c := dfs(0, -1)
	return max(nc, c)
}

// 树上最小顶点覆盖（每条边至少有一个端点被覆盖，也就是在已选的点集中）
// 代码和树上最大独立集类似
// https://brooksj.com/2019/06/20/%E6%A0%91%E7%9A%84%E6%9C%80%E5%B0%8F%E6%94%AF%E9%85%8D%E9%9B%86%EF%BC%8C%E6%9C%80%E5%B0%8F%E7%82%B9%E8%A6%86%E7%9B%96%E9%9B%86%EF%BC%8C%E6%9C%80%E5%A4%A7%E7%82%B9%E7%8B%AC%E7%AB%8B%E9%9B%86/
// 经典题：战略游戏 https://www.luogu.com.cn/problem/P2016
// 训练指南第一章例题 30，UVa10859 https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=20&page=show_problem&problem=1800
// - 求最小顶点覆盖，以及所有最小顶点覆盖中，两端点都被覆盖的边的最大个数
// 构造 https://codeforces.com/problemset/problem/959/C
func MinVertexCoverOfTree(tree [][]int, values []int) int { // 无根树
	var dfs func(int, int) (notChosen, chosen int)
	dfs = func(cur, pre int) (notChosen, chosen int) {
		chosen = values[cur] // 1
		for _, next := range tree[cur] {
			if next != pre {
				nc, c := dfs(next, cur)
				notChosen += c
				chosen += min(nc, c)
			}
		}
		return
	}
	nc, c := dfs(0, -1)
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
func MinDominatingSetOfTree(tree [][]int, values []int) int { // 无根树
	const INF int = 1e18
	var dfs func(int, int) (chosen, bySon, byFa int)
	dfs = func(cur, pre int) (chosen, bySon, byFa int) {
		chosen = values[cur] // 1
		extra := INF
		for _, next := range tree[cur] {
			if next != pre {
				c, bs, bf := dfs(next, cur)
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
	chosen, bySon, _ := dfs(0, -1)
	return min(chosen, bySon)
}

// 树上最大匹配
// g[v] = ∑{max(f[son],g[son])}
// f[v] = max{1+g[son]+g[v]−max(f[son],g[son])}
// https://codeforces.com/blog/entry/2059
// https://blog.csdn.net/lycheng1215/article/details/78368002
// https://vijos.org/p/1892
func MaxMatchingOfTree(n int, g [][]int) int { // 无根树
	cover, nonCover := make([]int, n), make([]int, n)
	var dfs func(int, int)
	dfs = func(cur, pre int) {
		for _, next := range g[cur] {
			if next != pre {
				dfs(next, cur)
				nonCover[cur] += max(cover[next], nonCover[next])
			}
		}
		for _, next := range g[cur] {
			cover[cur] = max(cover[cur], 1+nonCover[next]+nonCover[cur]-max(cover[next], nonCover[next]))
		}
	}
	dfs(0, -1)
	return max(cover[0], nonCover[0])
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
