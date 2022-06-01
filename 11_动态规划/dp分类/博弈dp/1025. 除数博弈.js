/**
 * @param {number} n
 * @return {boolean}

 */
var divisorGame = function (n) {
  const dp = Array(n + 1).fill(false) // 表示先手胜利
  dp[1] = false
  dp[2] = true
  for (let i = 3; i <= n; i++) {
    for (let j = 1; j < i; j++) {
      if (i % j === 0 && !dp[i - j]) {
        dp[i] = true
        break
      }
    }
  }

  return dp[n]
}

// * 最初，黑板上有一个数字 N
// 在每个玩家的回合，玩家需要执行以下操作：
// 选出任一 x，满足 0 < x < N 且 N % x == 0 。
// 用 N - x 替换黑板上的数字 N 。
// 如果玩家无法执行这些操作，就会输掉游戏。
