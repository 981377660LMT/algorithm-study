const maxMoney = (money: number[]) => {
  const len = money.length
  if (money.length === 0) {
    return 0
  }

  // dp[i]：考虑下标i（包括i）以内的房屋，最多可以偷窃的金额为dp[i]。
  const dp = Array<number>(money.length).fill(0)
  dp[0] = money[0]
  dp[1] = Math.max(money[0], money[1])

  for (let i = 2; i < len; i++) {
    dp[i] = Math.max(dp[i - 2] + money[i], dp[i - 1])
  }

  return dp[len - 1]
}

console.log(maxMoney([1, 2, 3, 4, 5]))
export {}
