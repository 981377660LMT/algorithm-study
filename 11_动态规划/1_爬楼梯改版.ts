// 一步一个台阶，两个台阶，三个台阶，.......，直到 m个台阶。问有多少种不同的方法可以爬到楼顶呢？
// dp[i]：爬到有i个台阶的楼顶，有dp[i]种方法。

// 1阶，2阶，.... m阶就是物品，楼顶就是背包。
// 每一阶可以重复使用，例如跳了1阶，还可以继续跳1阶。
// 问跳到楼顶有几种方法其实就是问装满背包有几种方法。
// 完全背包的排列问题
const climbStairs = (n: number, m: number): number => {
  const dp = Array<number>(n + 1).fill(0)
  dp[0] = 1
  for (let i = 1; i <= n; i++) {
    for (let j = 1; j <= m; j++) {
      if (i - j >= 0) dp[i] += dp[i - j]
    }
  }
  return dp[n]
}

console.log(climbStairs(6, 2))

export {}
