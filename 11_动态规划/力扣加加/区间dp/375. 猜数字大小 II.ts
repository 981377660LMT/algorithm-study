// 1.为什么使用动态规划，不使用二分？
// 使用动态规划的好处在于我可以穷举所有的情况
// 对于这个题来说，就是指动态规划的方法可以把每一个数字都当作分割点，
// 而二分只能把中间的数字当作分割点。

// 当你猜了数字 x 并且猜错了的时候，你需要支付金额为 x 的现金
// 给定 n ≥ 1，计算你至少需要拥有多少现金才能确保你能赢得这个游戏。
// function getMoneyAmount(n: number): number {
//   // 区间[i,j]时需要的现金数
//   const dp = Array.from({ length: n + 1 }, () => Array(n + 1).fill(Infinity))
//   for (let i = 1; i <= n; i++) {
//     dp[i][i] = 0
//   }
//   dp[1][0] = 0

//   for (let l = 1; l < n; l++) {
//     for (let i = 0; i < n - l; i++) {
//       const j = i + l
//       for (let k = i; k < j; k++) {
//         dp[i][j] = Math.min(dp[i][j], k + Math.max(dp[i][k - 1], dp[k + 1][j]))
//       }
//     }
//   }

//   console.table(dp)
//   return dp[1][n]
// }

console.log(getMoneyAmount(10))
