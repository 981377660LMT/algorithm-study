/**
 * @param {number[][]} box
 * @return {number}
 * 箱子宽 wi、深 di、高 hi。
 * 将箱子堆起来时，下面箱子的宽度、高度和深度必须大于上面的箱子
 */
const pileBox = function (box: number[][]): number {
  box.sort((a, b) => a[0] - b[0] || b[1] - a[1] || b[2] - a[2])
  // dp[i] 代表以第 i 个箱子放在最底下的最大高度。
  const dp = Array(box.length).fill(0)

  for (let i = 0; i < box.length; i++) {
    dp[i] = box[i][2]
    for (let j = 0; j < i; j++) {
      if (box[j][0] < box[i][0] && box[j][1] < box[i][1] && box[j][2] < box[i][2]) {
        dp[i] = Math.max(dp[i], dp[j] + box[i][2])
      }
    }
  }

  return Math.max.apply(null, dp)
}

console.log(
  pileBox([
    [1, 1, 1],
    [2, 2, 2],
    [3, 3, 3],
  ])
)

export default 1
