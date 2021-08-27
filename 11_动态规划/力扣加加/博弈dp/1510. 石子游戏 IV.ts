/**
 * @param {number} n
 * @return {boolean}
 * 每个人轮流操作，正在操作的玩家可以从石子堆里拿走 任意 非零 平方数 个石子。
 * 如果石子堆里没有石子了，则无法操作的玩家输掉游戏。
 * 如果 Alice 会赢得比赛，那么返回 True ，否则返回 False 。
 */
const winnerSquareGame = function (n: number): boolean {
  const memo = new Map<number, boolean>()

  // 剩下remain时先手是否胜利
  const dfs = (remain: number): boolean => {
    if (remain === 0) return false
    if (memo.has(remain)) return memo.get(remain)!

    const end = ~~Math.sqrt(remain)
    for (let i = 1; i <= end; i++) {
      if (!dfs(remain - i ** 2)) {
        memo.set(remain, true)
        return true
      }
    }

    memo.set(remain, false)
    return false
  }

  return dfs(n)
}

console.log(winnerSquareGame(2))
