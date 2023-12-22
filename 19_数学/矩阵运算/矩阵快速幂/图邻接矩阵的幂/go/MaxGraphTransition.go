// !floyd 也是一种矩阵乘法.
// Floyd算法是一种解决所有节点对之间最短路径问题的动态规划算法。
// 在这个算法中，我们使用一个二维矩阵来存储每个节点对之间的最短距离。
// 这个过程可以看作是一种特殊的矩阵乘法。
// 在标准的矩阵乘法中，我们计算结果矩阵的每个元素，作为两个输入矩阵对应行和列元素的乘积之和。
// 在Floyd算法中，我们也在做类似的事情，但是乘法和加法被替换为了取最小值和加法。
// 具体来说，对于Floyd算法中的每一对节点(i, j)，
// 我们尝试通过所有其他节点k来寻找一条可能的更短路径。
// 这就相当于在计算结果矩阵的元素dist[i][j]，
// 作为输入矩阵dist[i][k]和dist[k][j]的和的最小值。
// 因此，Floyd算法可以被看作是一种特殊的矩阵乘法，
// 其中的乘法和加法操作被替换为了取最小值和加法操作。

package main

const INF int = 1e18

// 给定一个图的临接矩阵，求转移k次(每次转移可以从一个点到到任意一个点)后的最长路(矩阵).
// 转移1次就是floyd.
func MaxGraphTransition(adjMatrix [][]int, k int) [][]int {
	n := len(adjMatrix)
	adjMatrixCopy := make([][]int, n)
	for i := range adjMatrixCopy {
		adjMatrixCopy[i] = make([]int, n)
		copy(adjMatrixCopy[i], adjMatrix[i])
	}

	dist := newMatrix(n, n, -INF)
	for i := 0; i < n; i++ {
		dist[i][i] = 0
	}

	for k > 0 {
		if k&1 == 1 {
			dist = MatMul(dist, adjMatrixCopy)
		}
		k >>= 1
		adjMatrixCopy = MatMul(adjMatrixCopy, adjMatrixCopy)
	}

	return dist
}

// 转移的自定义函数.
// ed:内部的结合律为取max(Floyd).
func MatMul(m1, m2 [][]int) [][]int {
	res := newMatrix(len(m1), len(m2[0]), -INF)
	for i := 0; i < len(m1); i++ {
		for k := 0; k < len(m2); k++ {
			for j := 0; j < len(m2[0]); j++ {
				res[i][j] = max(res[i][j], m1[i][k]+m2[k][j])
			}
		}
	}
	return res
}

func newMatrix(row, col int, fill int) [][]int {
	res := make([][]int, row)
	for i := range res {
		row := make([]int, col)
		for j := range row {
			row[j] = fill
		}
		res[i] = row
	}
	return res
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
