/**
 * @param {number} maxChoosableInteger maxChoosableInteger  不会大于 20
 * @param {number} desiredTotal
 * @return {boolean}
 * @description 两个玩家可以轮流从公共整数池中抽取从 1 到 maxChoosableInteger 的整数（不放回），直到累计整数和 >= desiredTotal
 * 数字不允许重复选:一个直观的思路是使用 set 记录已经被取的数字
 * 自底向上的dp
 */
const canIWin = function (maxChoosableInteger: number, desiredTotal: number): boolean {
  if (desiredTotal <= maxChoosableInteger) return true

  const dfs = (cur: number, visited: Set<number>): boolean => {
    // 自己还没选时
    if (cur >= desiredTotal) return false
    // 没选了
    if (visited.size === maxChoosableInteger) return false

    for (let i = 1; i <= maxChoosableInteger; i++) {
      // 自己选时,选这个数字就能赢了，返回 true
      if (!visited.has(i)) {
        visited.add(i)
        if (!dfs(cur + i, visited)) {
          visited.delete(i)
          return true
        }
        visited.delete(i)
      }
    }

    return false
  }

  return dfs(0, new Set())
}

console.log(canIWin(10, 11))
// false
export {}
