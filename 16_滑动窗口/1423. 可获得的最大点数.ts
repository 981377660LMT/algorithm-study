/**
 * @param {number[]} cardPoints
 * @param {number} k
 * @return {number}
 * @description
 * 每次行动，你可以从行的开头或者末尾拿一张卡牌，最终你必须正好拿 k 张卡牌。
 * 给你一个整数数组 cardPoints 和整数 k，请你返回可以获得的最大点数。
 * @summary
 * 注意先减后结算
 */
function maxScore(cardPoints: number[], k: number): number {
  const total = cardPoints.reduce((pre, cur) => pre + cur, 0)
  const needLength = cardPoints.length - k
  let [sum, minSum] = [0, Infinity]

  for (let index = 0; index < cardPoints.length; index++) {
    sum += cardPoints[index]
    if (index >= needLength) sum -= cardPoints[index - needLength]
    if (index >= needLength - 1) minSum = Math.min(minSum, sum)
  }

  return total - minSum
}

console.log(maxScore([1, 2, 3, 4, 5, 6, 1], 3))
