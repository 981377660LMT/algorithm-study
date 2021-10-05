// 答案需要取模 1e9+7（1000000007），如计算初始结果为：1000000008，请返回 1。
// 这一题已经不能用动态规划了，取余之后max函数就不能用来比大小了。
const cuttingRope3 = function (n: number): number {
  const MOD = 10 ** 9 + 7
  if (n < 4) return n - 1
  let res = 1

  while (n > 4) {
    res = (res * 3) % MOD
    n -= 3
  }

  return (res * n) % MOD
}
