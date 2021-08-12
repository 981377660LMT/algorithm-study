import { BSet } from './BSet'

/**
 * @param {number} maxChoosableInteger maxChoosableInteger  不会大于 20(可以用状态压缩)
 * @param {number} desiredTotal
 * @return {boolean}
 * @description 两个玩家可以轮流从公共整数池中抽取从 1 到 maxChoosableInteger 的整数（不放回），直到累计整数和 >= desiredTotal
 * 数字不允许重复选:一个直观的思路是使用 set 记录已经被取的数字
 * @summary
 * 复杂度O(2^n) 2^n个状态
 * 自顶向下的记忆化搜索
 */
const canIWin = function (maxChoosableInteger: number, desiredTotal: number): boolean {
  if (desiredTotal <= maxChoosableInteger) return true
  // 全部拿完也无法到达
  const sum = (maxChoosableInteger * (maxChoosableInteger + 1)) / 2
  if (desiredTotal > sum) return false

  const memo = new Map<number, boolean>()
  // 这里使用数字,每次都是一个新的visited,不存在引用那种需要回溯的情况
  const helper = (remain: number, state: number): boolean => {
    if (memo.has(state)) return memo.get(state)!

    for (let i = 1; i <= maxChoosableInteger; i++) {
      const cur = 1 << i
      // 已经看过
      if (state & cur) continue

      // 直接获胜
      if (i >= remain) {
        memo.set(state, true)
        return true
      }

      // 或操作添加到visited中
      if (!helper(remain - i, state | cur)) {
        memo.set(state, true)
        return true
      }
    }

    memo.set(state, false)
    return false
  }

  return helper(desiredTotal, 0)
}

console.log(canIWin(10, 11))
// false
export {}
