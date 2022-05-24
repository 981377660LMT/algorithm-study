// 第一个房屋和最后一个房屋是紧挨着的,
// 分两种情况:偷第一个和不偷第一个
function rob(nums: number[]) {
  const n = nums.length
  if (n === 1) return nums[0]
  // 取不取第一块
  const money1 = maxMoneyRange(0, n - 2)
  const money2 = maxMoneyRange(1, n - 1)
  return Math.max(money1, money2)

  function maxMoneyRange(start: number, end: number): number {
    if (start === end) return nums[start]
    const dp = Array<number>(nums.length).fill(0)
    dp[start] = nums[start]
    dp[start + 1] = Math.max(nums[start], nums[start + 1])
    for (let i = start + 2; i <= end; i++) {
      dp[i] = Math.max(dp[i - 2] + nums[i], dp[i - 1])
    }
    return dp[end]
  }
}

console.log(rob([1, 2, 3, 1]))

export {}
