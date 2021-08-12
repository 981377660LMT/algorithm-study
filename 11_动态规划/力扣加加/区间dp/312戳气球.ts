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
  const len = nums.length
  console.log(len)
  const dp = Array.from({ length: len }, () => Array(len).fill(0))
  const rangeBest = (i: number, j: number) => {
    let max = 0
    for (let k = i + 1; k < j; k++) {
      const left = dp[i][k]
      const right = dp[k][j]
      const sum = left + nums[i] * nums[k] * nums[j] + right
      max = Math.max(max, sum)
    }
    dp[i][j] = max
  }

  // l为区间长度,从2开始
  for (let l = 2; l < len; l++) {
    for (let i = 0; l + i < len; i++) {
      rangeBest(i, i + l)
    }
  }
  console.table(dp)

  return dp[0][len - 1]
}

console.log(maxCoins([3, 1, 5, 8]))
