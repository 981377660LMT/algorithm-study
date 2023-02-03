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

func optimalBinaryTree(keys, freqs []int) int {
	n := len(freqs)
	sort.Slice(freqs, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + freqs[i-1]
	}

	return KnuthYao2(n, func(i, k, j int) int {
		return preSum[j+1] - preSum[i] - freqs[k]
	})
}

// !dp[i][j]=min(dp[i][k-1]+dp[k+1][j]+cost(i,k,j))
func KnuthYao2(n int, cost func(i, k, j int) int) int {
	dp := make([][]int, n+2)
	pos := make([][]int, n+2)
	for i := 0; i < n+2; i++ {
		dp[i] = make([]int, n+2)
		pos[i] = make([]int, n+2)
	}

	for i := 1; i <= n; i++ {
		for j := i; j <= n; j++ {
			dp[i][j] = INF
		}
	}

	for i := 1; i <= n; i++ {
		dp[i][i] = 0
		pos[i][i] = i
	}

	for len := 2; len <= n; len++ {
		for i := 1; i+len-1 <= n; i++ {
			j := i + len - 1
			for k := pos[i][j-1]; k <= pos[i+1][j]; k++ {
				res := dp[i][k-1] + dp[k+1][j] + cost(i-1, k-1, j-1)
				if res < dp[i][j] {
					dp[i][j] = res
					pos[i][j] = k
				}
			}
		}
	}

	return dp[1][n]
}

func main() {
	fmt.Println(optimalBinaryTree([]int{10, 12, 20}, []int{5, 10, 20})) // 20
}
