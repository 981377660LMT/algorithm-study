// 把4个数求24点, 转化为3个数求24点, 再转化为2个数, 最后变1个数
// 一共有 9216 种可能性
function judgePoint24(cards: number[]): boolean {
  if (cards.length === 1) return Math.abs(cards[0] - 24) < 1e-6

  for (let i = 0; i < cards.length; i++) {
    for (let j = i + 1; j < cards.length; j++) {
      const rest = cards.filter((_, index) => index !== i && index !== j)
      if (
        judgePoint24([cards[i] + cards[j], ...rest]) ||
        judgePoint24([cards[i] * cards[j], ...rest]) ||
        judgePoint24([cards[i] - cards[j], ...rest]) ||
        judgePoint24([cards[j] - cards[i], ...rest]) ||
        judgePoint24([cards[i] / cards[j], ...rest]) ||
        judgePoint24([cards[j] / cards[i], ...rest])
      )
        return true
    }
  }

  return false
}

console.log(judgePoint24([4, 1, 8, 7]))
console.log(judgePoint24([1, 2, 1, 2]))
