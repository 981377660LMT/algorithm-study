// 1 <= cnt <= cards.length <= 10^5
// 若这 cnt 张卡牌数字总和为偶数，则选手成绩「有效」且得分为 cnt 张卡牌数字总和
// 请帮参赛选手计算最大的有效得分。若不存在获取有效得分的卡牌方案，则返回 0。
function maxmiumScore(cards: number[], cnt: number): number {
  cards.sort((a, b) => b - a)
  const oddPre: number[] = [0]
  const evenPre: number[] = [0]

  for (const card of cards) {
    if (card & 1) {
      oddPre.push(oddPre[oddPre.length - 1] + card)
    } else {
      evenPre.push(evenPre[evenPre.length - 1] + card)
    }
  }

  // 枚举所有组合中奇数的个数 k（k必须是偶数） 和 cnt - k（需判断是否足够）个偶数，它们都取最大则该轮组合结果最大
  let res = 0
  for (let oddCount = 0; oddCount < oddPre.length; oddCount += 2) {
    if (cnt - oddCount >= 0 && cnt - oddCount < evenPre.length)
      res = Math.max(res, oddPre[oddCount] + evenPre[cnt - oddCount])
  }

  return res
}

console.log(maxmiumScore([1, 2, 8, 9], 3))
console.log(maxmiumScore([3, 3, 1], 1))
