// 629. K个逆序对数组

const MOD = 1e9 + 7
const N = 1e3 + 10
const stirling1 = Array.from<unknown, number[]>({ length: N }, () => Array(N).fill(0))
stirling1[0][0] = 1
for (let i = 1; i < N; i++) {
  for (let j = 1; j < N; j++) {
    stirling1[i][j] = (stirling1[i - 1][j - 1] + stirling1[i - 1][j] * (i - 1)) % MOD
  }
}

/**
 * @param {number} n  1 <= n <= 1000
 * @param {number} k k <= n
 * @return {number}
 * 请你将这些木棍排成一排，并满足从左侧 可以看到 恰好 k 根木棍。
 * @description 划分为k个部分，每个部分排列种数为圆排列种树=>第一类斯特林数
 */
function rearrangeSticks(n: number, k: number): number {
  return stirling1[n][k]
}

console.log(rearrangeSticks(3, 2))
// 输出：3
// 解释：[1,3,2], [2,3,1] 和 [2,1,3] 是仅有的能满足恰好 2 根木棍可以看到的排列。
// 可以看到的木棍已经用粗体+斜体标识。

export {}
