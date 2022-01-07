// 629. K个逆序对数组
/**
 * @param {number} n  1 <= n <= 1000
 * @param {number} k
 * @return {number}
 * 请你将这些木棍排成一排，并满足从左侧 可以看到 恰好 k 根木棍。
 * @description
 * n 根木棍被划分成了 k 个部分，每个部分的第一根木棍即为可以看到的木棍。
 * 即长为 n 的排列划分成 k 个非空圆排列的方案数:第一类斯特林数
 * 假设1能被看到，也就是1必然放在最前面,那我们需要找2,3,4,5中的只能有3-1=2个被看到的解。
 * 假设1不被看到，也就是1放在除了最前面的位置都不会影响解，那我们需要找2,3,4,5中有3个被看到的解。
 */
function rearrangeSticks(n: number, k: number): number {
  const MOD = 10 ** 9 + 7
  const dp = Array.from<unknown, number[]>({ length: n + 1 }, () => Array(k + 1).fill(0))

  dp[0][0] = 1 // ???

  for (let i = 1; i < n + 1; i++) {
    for (let j = 1; j < k + 1; j++) {
      dp[i][j] = (dp[i - 1][j - 1] + dp[i - 1][j] * (i - 1)) % MOD
    }
  }

  return dp[n][k]
}

console.log(rearrangeSticks(3, 2))
// 输出：3
// 解释：[1,3,2], [2,3,1] 和 [2,1,3] 是仅有的能满足恰好 2 根木棍可以看到的排列。
// 可以看到的木棍已经用粗体+斜体标识。
