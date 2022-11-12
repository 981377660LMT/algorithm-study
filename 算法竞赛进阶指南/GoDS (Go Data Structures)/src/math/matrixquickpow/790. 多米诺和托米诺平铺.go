package matrixquickpow

// 790. 多米诺和托米诺平铺
// # !一维DP状态转移方程：dp[i] = 2 * dp[i - 1] + dp[i - 3]

// # !即:
// # [ai  ]   =    [2,0,1,0]  * [ai-1]
// # [ai-1]        [1,0,0,0]    [ai-2]
// # [ai-2]        [0,1,0,0]    [ai-3]
// # [ai-3]        [0,0,1,0]    [ai-4]
// # n<=3时 取init
// # n>3 时，转移n-3次

func numTilings(n int) int {
	init := Matrix([][]int{{5}, {2}, {1}, {0}})
	if n <= 3 {
		return init[len(init)-1-n][0]
	}

	trans := Matrix([][]int{{2, 0, 1, 0}, {1, 0, 0, 0}, {0, 1, 0, 0}, {0, 0, 1, 0}})
	resTrans := trans.Pow(n-3, 1e9+7)
	res := resTrans.Mul(init, 1e9+7)
	return res[0][0]
}
