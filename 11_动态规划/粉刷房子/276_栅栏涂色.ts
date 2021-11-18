// 有 k 种颜色的涂料和一个包含 n 个栅栏柱的栅栏
// 每个栅栏柱可以用其中 一种 颜色进行上色。
// 相邻的栅栏柱 最多连续两个 颜色相同。

//  拆成2种情况
// （1）i和i-1不同，则i有k-1种情况
// （2）i和i-1相同，则这2根与i-2不同，
function numWays(n: number, k: number): number {
  const dp = Array<number>(n + 1).fill(0)
  dp[1] = k
  dp[2] = k ** 2

  for (let i = 3; i < n + 1; i++) {
    dp[i] = dp[i - 2] * (k - 1) + dp[i - 1] * (k - 1)
  }

  return dp[n]
}

console.log(numWays(3, 2))

export {}
