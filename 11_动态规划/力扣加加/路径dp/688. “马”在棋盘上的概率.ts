/**
 * @param {number} n  NxN 的国际象棋棋盘
 * @param {number} k  打算进行 K 次移动
 * @param {number} row  初始位置
 * @param {number} column  初始位置
 * @return {number}
 * @description 这题很像 576. 出界的路径数  935. 骑士拨号器
 * dfs(x,y,remain)
 */
var knightProbability = function (n: number, k: number, row: number, column: number): number {
  const memo = new Map<string, number>()
  const next: [number, number][] = []
  // `这段逻辑很巧妙`
  for (const dx of [-2, -1, 1, 2]) {
    for (const dy of [-2, -1, 1, 2]) {
      if (Math.abs(dx) + Math.abs(dy) === 3) next.push([dx, dy])
    }
  }
  const isInBoard = (x: number, y: number) => x >= 0 && x < n && y >= 0 && y < n

  const dfs = (x: number, y: number, remain: number): number => {
    if (remain === 0) return 1

    const key = `${x}#${y}#${remain}`
    if (memo.has(key)) return memo.get(key)!
    let res = 0
    for (const [dx, dy] of next) {
      const nextX = x + dx
      const nextY = y + dy
      isInBoard(nextX, nextY) && (res += dfs(nextX, nextY, remain - 1) / 8)
    }

    memo.set(key, res)
    return res
  }

  return dfs(row, column, k)
}

console.log(knightProbability(3, 2, 0, 0))

export {}
