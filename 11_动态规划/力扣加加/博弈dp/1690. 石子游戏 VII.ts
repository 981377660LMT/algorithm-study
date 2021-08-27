/**
 * @param {number[]} stones 1 <= stones[i] <= 1000
 * @return {number}
 * 每个玩家的回合中，可以从行中 移除 最左边的石头或最右边的石头，并获得与该行中剩余石头值之 和 相等的得分。当没有石头可移除时，得分较高者获胜。
 * 当没有石头可移除时，得分较高者获胜。
 * 如果爱丽丝和鲍勃都 发挥出最佳水平 ，请返回他们 得分的差值 。
 * @summary 最终的得分之差其实就是鲍勃选的数字之和。
 */
const stoneGameVII = function (stones: number[]): number {
  const len = stones.length
  const sum = stones.reduce((pre, cur) => pre + cur)
  const memo = new Map<string, number>()

  // 当前先手时能拿到的相对最大值
  const dfs = (left: number, right: number, sum: number): number => {
    if (right === left) return 0
    const key = `${left}#${right}`
    if (memo.has(key)) return memo.get(key)!

    const chooseLeft = sum - stones[left] - dfs(left + 1, right, sum - stones[left])
    const chooseRight = sum - stones[right] - dfs(left, right - 1, sum - stones[right])
    const res = Math.max(chooseLeft, chooseRight)

    memo.set(key, res)
    return res
  }

  return dfs(0, len - 1, sum)
}

console.log(stoneGameVII([5, 3, 1, 4, 2]))
