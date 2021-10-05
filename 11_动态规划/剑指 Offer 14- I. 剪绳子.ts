// 贪心
/**
 * @param {number} n
 * @return {number}
 *给你一根长度为 n 的绳子，请把绳子剪成整数长度的 m 段
 请问 k[0]*k[1]*...*k[m-1] 可能的最大乘积是多少
 @summary 
 尽可能把绳子分成长度为3的小段，这样乘积最大
 如果n小于4 
 如果 n == 4，返回4
 如果n大于4，尽可能把绳子分成长度为3的小段，这样乘积最大
 */
const cuttingRope = function (n: number): number {
  if (n < 4) return n - 1
  let res = 1

  while (n > 4) {
    res *= 3
    n -= 3
  }

  return res * n
}

// dp
const cuttingRope2 = function (n: number): number {
  const dp = Array(n + 1).fill(0)
  dp[2] = 1
  for (let i = 3; i <= n; i++) {
    for (let j = 2; j < i; j++) {
      dp[i] = Math.max(dp[i], Math.max(j * dp[i - j], j * (i - j)))
    }
  }
  return dp[n]
}
// 剪了第一段后，剩下(i - j)长度可以剪也可以不剪。
// 如果不剪的话长度乘积即为j * (i - j)；
// 如果剪的话长度乘积即为j * dp[i - j]。取两者最大值max(j * (i - j), j * dp[i - j])
