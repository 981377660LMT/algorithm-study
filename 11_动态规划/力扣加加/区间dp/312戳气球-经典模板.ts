/**
 * @param {number[]} nums
 * @return {number}
 * @description
 * 戳破第 i 个气球，你可以获得 nums[i - 1] * nums[i] * nums[i + 1] 枚硬币。
 * 超出了数组的边界，那么就当它是一个数字为 1 的气球。
 * 求所能获得硬币的最大数量。
 * @summary
 * k是开区间(i,j)   最后一个   被戳爆的气球！！！！！
 * dp[i][j] 表示开区间 (i,j) 内你能拿到的最多金币
 * dp[i][j]=dp[i][k]+nums[i]*nums[k]*nums[j]+dp[k][j]
 */
const maxCoins = function (nums: number[]): number {
  // 首尾添加1，方便处理边界情况
  nums.unshift(1)
  nums.push(1)

  const n = nums.length
  const dp = Array.from({ length: n + 1 }, () => Array(n).fill(0))
  for (let l = 0; l < n; l++) {
    for (let i = 0; i < n - l; i++) {
      const j = i + l
      for (let k = i + 1; k < j; k++) {
        // 我们假设最后戳破的气球是 k
        dp[i][j] = Math.max(dp[i][j], dp[i][k] + dp[k][j] + nums[i] * nums[k] * nums[j])
      }
    }
  }

  return dp[0][dp[0].length - 1]
}

console.log(maxCoins([3, 1, 5, 8]))
