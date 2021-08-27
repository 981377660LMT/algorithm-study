/**
 * @param {number[]} aliceValues
 * @param {number[]} bobValues
 * @return {number}
 * 轮到某个玩家时，他可以 移出 一个石子并得到这个石子的价值。Alice 和 Bob 对石子价值有 不一样的的评判标准 。双方都知道对方的评判标准。
 * aliceValues[i] 和 bobValues[i] 分别表示 Alice 和 Bob 认为第 i 个石子的价值。
 * 如果 Alice 赢，返回 1 。 如果游戏平局，返回 0
 * @summary 选取石头的时候要么让自己得分最高，要么让对手失分最多，双方都基于这个贪心策略选择石头。
 * alice拿走石头能获得的价值是 石头本身的价值+恶心bob的价值。
 */
const stoneGameVI = function (aliceValues: number[], bobValues: number[]): number {
  const length = aliceValues.length
  const zipped = Array.from<number, [number, number]>({ length }, (_, i) => [
    aliceValues[i],
    bobValues[i],
  ])
  zipped.sort((a, b) => b[0] + b[1] - a[0] - a[1])
  let aliceSum = 0
  let bobSum = 0
  for (let i = 0; i < length; i += 2) {
    aliceSum += zipped[i][0]
  }
  for (let i = 1; i < length; i += 2) {
    bobSum += zipped[i][1]
  }
  if (aliceSum > bobSum) return 1
  else if (aliceSum < bobSum) return -1
  else return 0
}

console.log(stoneGameVI([1, 3], [2, 1]))
