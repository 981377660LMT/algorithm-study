// 给你一个整数 n ，返回和为 n 的完全数的 最少数量 。
// 这个方法太慢了复杂度O(n**3/2)
// 完全平方数就是物品（可以无限件使用），凑个正整数n就是背包，问凑满这个背包最少有多少物品？
const numSqures = (n: number) => {
  // 递推式:最接近的一个完全平方数加上剩下的走递推
  const dp = Array(n + 1).fill(Infinity)
  dp[0] = 0
  dp[1] = 1

  for (let i = 1; i <= n; i++) {
    // 使用j取代了不必要的数组
    for (let j = 1; j * j <= i; j++) {
      dp[i] = Math.min(dp[i], dp[i - j * j] + 1)
    }
  }

  return dp[n]
}

console.log(numSqures(12))

export {}
