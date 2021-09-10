/**
 * @param {number} n
 * @return {number}
 * @description 不同的瓷砖一共有六种
 * 有两种形状的瓷砖：一种是 2x1 的多米诺形，另一种是形如 "L" 的托米诺形。两种形状都可以旋转。
 * 给定 N 的值，有多少种方法可以平铺 2 x N 的面板？返回值 mod 10^9 + 7。
 */
const numTilings = function (n: number): number {
  const mod = 10e9 + 7
  const dp = Array(n + 3).fill(0)
  dp[0] = 1
  dp[1] = 1
  dp[2] = 2
  dp[3] = 5
  for (let i = 4; i <= n; i++) {
    dp[i] = ((((2 * dp[i - 1]) % mod) % mod) + (dp[i - 3] % mod)) % mod
  }

  return dp[n] % mod
}
console.log(10e1)
console.log(10 ** 9)

console.log(numTilings(3))
console.log(numTilings(4))

export {}
