/**
 * @param {number[]} stoneValue
 * @return {number}
 * Alice 会将这行石子分成两个 非空行（即，左侧行和右侧行）；
 * Bob 负责计算每一行的值，即此行中所有石子的值的总和。
 * Bob 会丢弃值最大的行，Alice 的得分为剩下那行的值（每轮累加）。
 * 如果两行的值相等，Bob 让 Alice 决定丢弃哪一行。下一轮从剩下的那一行开始。
   只 剩下一块石子 时，游戏结束。Alice 的分数最初为 0 。
   返回 Alice 能够获得的最大分数 。
 */
const stoneGameV = function (stoneValue: number[]): number {
  const len = stoneValue.length
  const pre = Array<number>(len + 1).fill(0)
  for (let i = 1; i <= len; i++) {
    pre[i] = pre[i - 1] + stoneValue[i - 1]
  }
  const memo = new Map<string, number>()

  // 剩下remain时先手是否胜利
  const dfs = (left: number, right: number): number => {
    if (right - left === 1) return 0
    const key = `${left}#${right}`
    if (memo.has(key)) return memo.get(key)!

    let res = -Infinity
    for (let i = left + 1; i < right; i++) {
      const leftSum = pre[i] - pre[left]
      const rightSum = pre[right] - pre[i]

      if (leftSum > rightSum) {
        res = Math.max(res, rightSum + dfs(i, right))
      } else if (leftSum < rightSum) {
        res = Math.max(res, leftSum + dfs(left, i))
      } else {
        res = leftSum + Math.max(dfs(left, i), dfs(i, right))
      }
    }

    memo.set(key, res)
    return res
  }

  return dfs(0, pre.length - 1)
}

console.log(stoneGameV([6, 2, 3, 4, 5, 5]))
console.log(stoneGameV([2, 1, 1]))
