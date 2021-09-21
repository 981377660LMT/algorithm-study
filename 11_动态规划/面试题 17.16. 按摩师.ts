// 她不能接受相邻的预约
function massage(nums: number[]): number {
  // dp表示接收到第几个数的最大和
  const len = nums.length
  const dp = Array<number>(len).fill(0)
  dp[0] = nums[0]
  dp[1] = Math.max(nums[0], nums[1])
  for (let i = 2; i < len; i++) {
    dp[i] = Math.max(dp[i - 1], dp[i - 2] + nums[i])
  }

  return dp[len - 1]
}

console.log(massage([2, 7, 9, 3, 1]))
