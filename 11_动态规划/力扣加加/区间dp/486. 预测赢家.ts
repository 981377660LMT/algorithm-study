/**
 * @param {number[]} nums
 * @return {boolean}  玩家 1 从数组任意一端拿取一个分数，随后玩家 2 继续从剩余数组任意一端拿取分数，然后玩家 1 拿
 * 最终获得分数总和最多的玩家获胜
 */
const PredictTheWinner = function (nums: number[]): boolean {
  const n = nums.length
  // 先手的得分
  const dp = Array.from({ length: n + 1 }, () => Array(n).fill(0))
  for (let i = 0; i < n; i++) {
    dp[i][i] = nums[i]
  }

  for (let l = 1; l < n; l++) {
    for (let i = 0; i < n - l; i++) {
      const j = i + l
      dp[i][j] = Math.max(nums[i] - dp[i + 1][j], nums[j] - dp[i][j - 1])
    }
  }

  return dp[0][dp[0].length - 1] >= 0
}

console.log(PredictTheWinner([1, 5, 2]))
console.log(PredictTheWinner([1, 5, 233, 7]))
