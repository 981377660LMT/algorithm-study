import { memo } from '../../../5_map/memo'

/**
 * @param {number[]} stoneValue
 * @return {string}
 * Alice 和 Bob 轮流取石子，Alice 总是先开始。在每个玩家的回合中，该玩家可以拿走剩下石子中的的前 1、2 或 3 堆石子 。
 * 得分最高的选手将会赢得比赛，比赛也可能会出现平局。
 */
const stoneGameIII = function (stoneValue: number[]): string {
  const len = stoneValue.length

  // 从index拿时先手比对手多的数量
  let dfs = (index: number): number => {
    if (index >= len) return 0

    let res = -Infinity
    let myStone = 0

    for (let i = index; i < index + 3; i++) {
      if (i >= len) break
      myStone += stoneValue[i]
      // 减去对手最多拿的
      res = Math.max(res, myStone - dfs(i + 1))
    }

    return res
  }
  dfs = memo(dfs)

  const res = dfs(0)
  if (res > 0) return 'Alice'
  else if (res < 0) return 'Bob'
  else return 'Tie'
}

console.log(stoneGameIII([1, 2, 3, 5]))
console.log(stoneGameIII([1]))
console.log(stoneGameIII([1, 2, 3, 6]))
