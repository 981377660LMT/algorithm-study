// 3562. 折扣价交易股票的最大利润
// https://leetcode.cn/problems/maximum-profit-from-trading-stocks-with-discounts/description/
// 给你一个整数 n，表示公司中员工的数量。每位员工都分配了一个从 1 到 n 的唯一 ID ，其中员工 1 是 CEO。另给你两个下标从 1 开始的整数数组 present 和 future，两个数组的长度均为 n，具体定义如下：
// present[i] 表示第 i 位员工今天可以购买股票的 当前价格 。
// future[i] 表示第 i 位员工明天可以卖出股票的 预期价格 。
// 公司的层级关系由二维整数数组 hierarchy 表示，其中 hierarchy[i] = [ui, vi] 表示员工 ui 是员工 vi 的直属上司。
//
// 此外，再给你一个整数 budget，表示可用于投资的总预算。
//
// 公司有一项折扣政策：如果某位员工的直属上司购买了自己的股票，那么该员工可以以 半价 购买自己的股票（即 floor(present[v] / 2)）。
//
// 请返回在不超过给定预算的情况下可以获得的 最大利润 。
//
// 注意：
// 每只股票最多只能购买一次。
// 不能使用股票未来的收益来增加投资预算，购买只能依赖于 budget。
//
// 1 <= n <= 160
// present.length, future.length == n
// 1 <= present[i], future[i] <= 50
// hierarchy.length == n - 1
// hierarchy[i] == [ui, vi]
// 1 <= ui, vi <= n
// ui != vi
// 1 <= budget <= 160
// 没有重复的边。
// 员工 1 是所有员工的直接或间接上司。
// 输入的图 hierarchy 保证 无环 。
//
// !对每个顶点，dp[j][k] 表示预算<=j，能否半价购买时(k=0/1)，子树的最大利润之和.

package main

func maxProfit(n int, present []int, future []int, hierarchy [][]int, budget int) int {
	tree := make([][]int, n)
	for _, edge := range hierarchy {
		x, y := edge[0]-1, edge[1]-1
		tree[x] = append(tree[x], y)
	}

	var dfs func(int) [][2]int
	dfs = func(node int) [][2]int {
		subDp := make([][2]int, budget+1)
		for _, next := range tree[node] {
			childDp := dfs(next)
			for j := budget; j >= 0; j-- {
				for childJ, p := range childDp[:j+1] {
					for k, childRes := range p {
						subDp[j][k] = max(subDp[j][k], subDp[j-childJ][k]+childRes)
					}
				}
			}
		}

		dp := make([][2]int, budget+1)
		for j, p := range subDp {
			for k := 0; k < 2; k++ {
				cost := present[node] / (1 + k)
				if j >= cost {
					dp[j][k] = max(p[0], subDp[j-cost][1]+future[node]-cost)
				} else {
					dp[j][k] = p[0]
				}
			}
		}

		return dp
	}

	res := dfs(0)
	return res[budget][0]
}
