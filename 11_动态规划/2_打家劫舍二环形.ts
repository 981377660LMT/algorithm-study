// 第一个房屋和最后一个房屋是紧挨着的,
// 分两种情况:偷第一个和不偷第一个
const maxMoney = (nums: number[]) => {
  const len = nums.length
  if (len === 0) return 0
  if (len === 1) return nums[0]
  const maxMoneyRange = (nums: number[], start: number, end: number): number => {
    if (start === end) return nums[start]
    const dp = Array(nums.length).fill(0)
    dp[start] = nums[start]
    dp[start + 1] = Math.max(nums[start], nums[start + 1])
    for (let i = start + 2; i <= end; i++) {
      dp[i] = Math.max(dp[i - 2] + nums[i], dp[i - 1])
    }
    return dp[end]
  }
  const money1 = maxMoneyRange(nums, 0, len - 2)
  const money2 = maxMoneyRange(nums, 1, len - 1)
  return Math.max(money1, money2)
}

console.log(maxMoney([1, 2, 3, 1]))

export {}
