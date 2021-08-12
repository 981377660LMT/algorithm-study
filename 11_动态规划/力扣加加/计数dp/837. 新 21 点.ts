/**
 * @param {number} n
 * @param {number} k
 * @param {number} maxPts
 * @return {number}
 * 爱丽丝以 0 分开始，并在她的得分少于 K 分时抽取数字。 抽取时，她从 [1, maxPts] 的范围中随机获得一个整数作为分数进行累计，
 * 当爱丽丝获得不少于 K 分时，她就停止抽取数字。 停止抽牌后，她的牌面小于等于 N 时，她就获胜了，求她获胜的概率。
 * @summary dp[x] 为她手上牌面为x时，能获胜的概率，
 * 当x>=K时，爱丽丝会停止抽牌，这个时候游戏已经结束了，她是赢是输也已经确定了，所以此时赢的概率要么1，要么0
 * 当x<K时 dp[x]=1/w * dp[x+1]+ 1/w * dp[x+2] + 1/w * dp[x+3]...+ 1/w * dp[x+w]
 */
const new21Game = function (n: number, k: number, maxPts: number): number {
  // 爱丽丝停止抽牌时，她可能达到的最大牌面是 K+W-1
  const dp = Array(k + maxPts - 1 + 1).fill(0)
  let sum = 0

  for (let i = k; i < k + maxPts; i++) {
    dp[i] = i <= n ? 1 : 0
    sum += dp[i]
  }
  for (let i = k - 1; i >= 0; i--) {
    dp[i] = sum / maxPts
    sum = sum - dp[i + maxPts] + dp[i]
  }

  return dp[0]
}

console.log(new21Game(6, 1, 10))
// 输出：0.60000
