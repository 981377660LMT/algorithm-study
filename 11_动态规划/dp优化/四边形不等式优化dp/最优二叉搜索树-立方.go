// 最优二叉搜索树
// https://blog.csdn.net/weixin_43914593/article/details/105150937
// https://beet-aizu.github.io/library/algorithm/optimalbinarytree.cpp
// 给定n个不同元素的集合S=(e1,e2,...,en) e1<e2<...<en
// 把S 的元素建一棵二叉搜索树，希望查询频率越高的元素离根越近。
// 访问树中元素ei的成本cost(ei)等于从根到该元素结点的路径边数
// 给定元素的查询频率f(ei) 求最小的成本∑f(ei)cost(ei)
// n<=250

// 解:
// 从小到大排序后 可以方便组成一颗BST (根节点k在中间滑动)
// dp[i][j]表示区间[i,j]组成的BST的最小成本
// 当把两棵左右子树连在根结点上时，本身的深度增加1，所以每个元素都多计算一次(除开根节点)
// dp[i][j]=min(dp[i][k-1]+dp[k+1][j]+freqSum(i,j)-freqs[k])

package main

import (
	"fmt"
	"sort"
)

const INF int = 1e18

// O(n^3)
func optimalBinaryTree(keys, freqs []int) int {
	n := len(freqs)
	sort.Slice(freqs, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + freqs[i-1]
	}

	memo := make([][]int, n)
	for i := 0; i < n; i++ {
		memo[i] = make([]int, n)
		for j := 0; j < n; j++ {
			memo[i][j] = -1
		}
	}
	var dfs func(left, right int) int
	dfs = func(left, right int) int {
		if left >= right {
			return 0
		}
		if memo[left][right] != -1 {
			return memo[left][right]
		}
		res := INF
		for k := left; k <= right; k++ {
			res = min(res, dfs(left, k-1)+dfs(k+1, right)+preSum[right+1]-preSum[left]-freqs[k])
		}

		memo[left][right] = res
		return res
	}

	return dfs(0, n-1)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println(optimalBinaryTree([]int{10, 12, 20}, []int{5, 10, 20})) // 20
}
