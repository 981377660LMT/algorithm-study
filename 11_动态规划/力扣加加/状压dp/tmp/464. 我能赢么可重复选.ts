/**
 * @param {number} maxChoosableInteger
 * @param {number} desiredTotal
 * @return {boolean}
 * @description 两个玩家可以轮流从公共整数池中抽取从 1 到 maxChoosableInteger 的整数（不放回），直到累计整数和 >= desiredTotal
 *
 */
const canIWin = function (maxChoosableInteger: number, desiredTotal: number): boolean {
  if (desiredTotal <= maxChoosableInteger) return true

  const dfs = (cur: number): boolean => {
    // 自己还没选时
    if (cur >= desiredTotal) return false

    for (let i = 1; i <= maxChoosableInteger; i++) {
      // 自己选时,选这个数字就能赢了，返回 true
      if (!dfs(cur + i)) return true
    }

    return false
  }

  return dfs(0)
}

console.log(canIWin(10, 11))
// false

export {}
