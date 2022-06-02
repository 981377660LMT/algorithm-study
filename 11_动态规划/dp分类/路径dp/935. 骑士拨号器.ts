const DIR: Record<number, number[]> = {
  0: [4, 6],
  1: [6, 8],
  2: [7, 9],
  3: [4, 8],
  4: [0, 3, 9],
  5: [],
  6: [0, 1, 7],
  7: [2, 6],
  8: [1, 3],
  9: [2, 4],
}
const MOD = 1e9 + 7

/**
 * @param {number} n
 * @return {number}
 * 🐎走日字
 * 我们将 “骑士” 放在电话拨号盘的任意数字键（如上图所示）上，接下来，骑士将会跳 N-1 步。
 */
function knightDialer(n: number): number {
  let res = 0

  const memo = new Map<string, number>()
  for (let i = 0; i < 10; i++) {
    res += dfs(i, n - 1)
    res %= MOD
  }

  return res

  function dfs(cur: number, remain: number): number {
    if (remain === 0) return 1
    const key = `${cur}#${remain}`
    if (memo.has(key)) return memo.get(key)!

    let res = 0
    for (const next of DIR[cur]) {
      res += dfs(next, remain - 1)
      res %= MOD
    }

    memo.set(key, res)
    return res % MOD
  }
}

console.log(knightDialer(2))
console.log(knightDialer(1))
console.log(knightDialer(3))
console.log(knightDialer(3131))
export {}
