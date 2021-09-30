/**
 * @param {string} s
 * @return {number}
 * 打印机每次只能打印由 同一个字符 组成的序列
 * 每次可以在任意起始和结束位置打印新字符，并且会覆盖掉原来已有的字符
 * https://leetcode-cn.com/problems/strange-printer/solution/yi-wen-tuan-mie-qu-jian-dp-by-bnrzzvnepe-24jz/
 */
const strangePrinter = function (s: string): number {
  if (!s) return 0
  const n = s.length

  const dp = Array.from({ length: n + 1 }, () => Array(n).fill(0))
  for (let i = 0; i < n; i++) {
    dp[i][i] = 1
  }
  // ：从前向后
  // 区间长度 区间起点 区间终点 初始化 枚举分割点 首位一样可减少一次
  for (let l = 1; l < n; l++) {
    for (let i = 0; i < n - l; i++) {
      const j = i + l
      // dp[i][j]的上限：比前面多一次
      dp[i][j] = dp[i + 1][j] + 1
      for (let k = i + 1; k <= j; k++) {
        // 与首部相同 少打一次
        if (s[k] === s[i]) dp[i][j] = Math.min(dp[i][j], dp[i][k - 1] + dp[k + 1][j])
      }
    }
  }
  console.table(dp)
  return dp[0][dp[0].length - 1]
}

console.log(strangePrinter('aaabbb'))
// 输出：2
// 解释：首先打印 "aaa" 然后打印 "bbb"。
console.log(strangePrinter('aba'))
// 输出：2
// 解释：首先打印 "aaa" 然后在第二个位置打印 "b" 覆盖掉原来的字符 'a'。

// 与 312思想基本一致，枚举 k 时加一层限制
// 这道题每次可以在任意起始和结束位置打印新字符。 因此我们需要暴力枚举所有的起始位置和结束位置的笛卡尔积
