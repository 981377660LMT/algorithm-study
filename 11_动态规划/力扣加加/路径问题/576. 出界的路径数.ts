/**
 * @param {number} m
 * @param {number} n
 * @param {number} maxMove
 * @param {number} startRow
 * @param {number} startColumn
 * @return {number} 找出并返回可以将球移出边界的路径数量
 * 因为答案可能非常大，返回对 109 + 7 取余 后的结果。
 * Time = O(m * n * maxMove) = O(n^3) Space = O(m * n * maxMove) = O(n^3)
 * @summary 这道题一看就是记忆化搜索
 */
const findPaths = function (
  m: number,
  n: number,
  maxMove: number,
  startRow: number,
  startColumn: number
): number {
  const mod = 10 ** 9 + 7
  const memo = new Map<string, number>()
  const dfs = (x: number, y: number, remain: number): number => {
    if (remain < 0) return 0
    if (x < 0 || x >= m || y < 0 || y >= n) return 1
    const key = `${x}#${y}#${remain}`
    if (memo.has(key)) return memo.get(key)!
    let res = 0
    res += dfs(x - 1, y, remain - 1)
    res += dfs(x + 1, y, remain - 1)
    res += dfs(x, y - 1, remain - 1)
    res += dfs(x, y + 1, remain - 1)
    res %= mod
    memo.set(key, res)
    return res
  }
  return dfs(startRow, startColumn, maxMove)
}

console.log(findPaths(2, 2, 2, 0, 0))

export {}
