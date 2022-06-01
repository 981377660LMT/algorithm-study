// 给你一个整数 n ，返回和为 n 的完全数的 最少数量 。
// 完全平方数就是物品（可以无限件使用），凑个正整数n就是背包，问凑满这个背包最少有多少物品？
function numSquares(n: number) {
  // 递推式:最接近的一个完全平方数加上剩下的走递推
  const dp = new Uint32Array(n + 1).fill(-1)
  dp[0] = 0
  dp[1] = 1

  for (let i = 1; i <= n; i++) {
    for (let j = 1; j * j <= i; j++) {
      dp[i] = Math.min(dp[i], dp[i - j * j] + 1)
    }
  }

  return dp[n]
}

console.log(numSquares(12))

export {}
