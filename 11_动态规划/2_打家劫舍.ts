const maxMoney = (money: number[]) => {
  if (money.length === 0) {
    return 0
  }

  // 前两项初始值
  const dp = [0, money[0]]
  for (let index = 2; index <= money.length; index++) {
    dp[index] = Math.max(dp[index - 2] + money[index - 1], dp[index - 1])
  }

  return dp
}

console.log(maxMoney([1, 2, 3, 4, 5]))
