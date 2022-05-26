/**
 * @param {number} poured
 * @param {number} query_row
 * @param {number} query_glass  query_glass 和 query_row 的范围 [0, 99]
 * @return {number}
 * 现在当倾倒了非负整数杯香槟后，返回第 i 行 j 个玻璃杯所盛放的香槟占玻璃杯容积的比例（i 和 j都从0开始）。
 */
function champagneTower(poured: number, query_row: number, query_glass: number): number {
  // 这种初始化方式有一半空间是浪费的
  const dp = Array.from({ length: 101 }, () => Array(100).fill(0))
  dp[0][0] = poured
  for (let row = 0; row < 100; row++) {
    for (let col = 0; col <= row; col++) {
      if (dp[row][col] > 1) {
        const overflow = dp[row][col] - 1
        dp[row][col] = 1
        dp[row + 1][col] += overflow / 2
        dp[row + 1][col + 1] += overflow / 2
      }
    }
  }

  return dp[query_row][query_glass]
}

console.log(champagneTower(2, 1, 1))
console.log(champagneTower(100000009, 33, 33))

export {}
