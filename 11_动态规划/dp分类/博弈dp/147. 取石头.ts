// 假定有n (n > 0) 颗石子。
// 玩家A和B轮流取石子，一次可以取1个或两个，A先取。
// 谁取到最后一颗石子即为输。
// 请问 A或者B是否有必赢策略? 如果有请返回能必赢的玩家名。没有则返回null。

function canWinStonePicking(n: number): 'A' | 'B' | null {
  return n % 3 === 1 ? 'B' : 'A'
}

function canWinStonePicking2(n: number): 'A' | 'B' | null {
  // 第n步先手赢
  const dp = Array<boolean>(n + 1).fill(true)
  dp[1] = false

  for (let i = 2; i <= n; i++) {
    dp[i] = !dp[i - 1] || !dp[i - 2]
  }

  return dp[n] ? 'A' : 'B'
}

console.log(canWinStonePicking2(2))
console.log(canWinStonePicking2(3))
console.log(canWinStonePicking2(4))
console.log(canWinStonePicking2(50))
export {}
