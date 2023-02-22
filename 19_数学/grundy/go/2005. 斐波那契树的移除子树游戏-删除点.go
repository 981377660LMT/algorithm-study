// https://leetcode.cn/problems/subtree-removal-game-with-fibonacci-tree/
// 斐波那契树是一种按这种规则函数 order(n) 创建的二叉树：
// order(0) 是空树。
// order(1) 是一棵只有一个节点的二叉树。
// order(n) 是一棵根节点的左子树为 order(n - 2) 、右子树为 order(n - 1) 的二叉树。
// Alice 和 Bob 在玩一种关于斐波那契树的游戏，由 Alice 先手。
// !在每个回合中，每个玩家选择一个节点，然后移除该节点及其子树。只能删除根节点 root 的玩家输掉这场游戏。
// 给定一个整数 n，假定两名玩家都按最优策略进行游戏，若 Alice 赢得这场游戏，返回 true 。
// 若 Bob 赢得这场游戏，返回 false 。
// n<=100

// 克朗原理:
// !对于树上的某一个点，ta 的分支可以转化成以这个点为根的一根竹子，
// 这个竹子的长度就是 ta 各个分支的边的数量的异或和
// !这里删除的是树节点,可以看成每个结点有一个出点和入点,删点等价于删边

package main

func findGameWinner(n int) bool {
	memo := make([]int, n+1)
	for i := range memo {
		memo[i] = -1
	}

	var dfs func(int) int
	dfs = func(root int) int {
		if memo[root] != -1 {
			return memo[root]
		}
		if root == 1 {
			return 0
		}

		g := 0
		if root-1 >= 1 {
			g ^= dfs(root-1) + 1
		}
		if root-2 >= 1 {
			g ^= dfs(root-2) + 1
		}

		memo[root] = g
		return g
	}

	groundy := dfs(n)
	return groundy != 0
}
