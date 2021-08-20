/**
 * @param {string[]} strs
 * @param {number} m
 * @param {number} n
 * @return {number}
 * @description 找出并返回 strs 的最大子集的大小，该子集中 最多 有 m 个 0 和 n 个 1
 * @summary 两个维度的01背包问题 不是多重背包
 */
const findMaxForm = function (strs: string[], m: number, n: number) {
  // dp[i][j]：最多有i个0和j个1的strs的最大子集的大小为dp[i][j]
  const dp = Array.from({ length: m + 1 }, () => Array(n + 1).fill(0))
  // dp[0][0] = 1

  for (let stuffIndex = 0; stuffIndex < strs.length; stuffIndex++) {
    let str = strs[stuffIndex]
    let zeros = 0
    let ones = 0
    for (const letter of str) {
      letter === '0' && zeros++
      letter === '1' && ones++
    }
    // 记得01背包问题递推容量是倒序
    for (let i = m; i >= zeros; i--) {
      for (let j = n; j >= ones; j--) {
        dp[i][j] = Math.max(dp[i][j], dp[i - zeros][j - ones] + 1)
      }
    }
  }

  console.table(dp)
  return dp[m][n]
}

console.log(findMaxForm(['10', '0001', '111001', '1', '0'], 5, 3))
// 输出：4
// 最多有 5 个 0 和 3 个 1 的最大子集是 {"10","0001","1","0"}
export {}
