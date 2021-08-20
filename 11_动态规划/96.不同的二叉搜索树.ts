/**
 * @param {number} n
 * @return {number}
 * 求恰由 n 个节点组成且节点值从 1 到 n 互不相同的 二叉搜索树 有多少种？
 */
const numTrees = function (n: number): number {
  const dp = Array<number>(n + 1).fill(0)
  dp[0] = 1
  for (let i = 1; i <= n; i++) {
    for (let j = 1; j <= i; j++) {
      // 例如i=3 讨论左右节点 左0右2 左1右1 左0右2 (关注形状而不关注值)
      dp[i] += dp[i - j] * dp[j - 1]
    }
  }

  return dp[n]
}

console.log(numTrees(3))

export default 1
// # Catalan Number  (2n)!/((n+1)!*n!)
