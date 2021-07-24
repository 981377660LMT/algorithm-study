/**
 * @param {number} n
 * @return {string[][]}
 * @description 皇后彼此不能相互攻击，也就是说：任何两个皇后都不能处于同一条横行、纵行或斜线上。
 * @description 如何快速剪枝:// 将index行的皇后摆放在第i列，在每一行的0到n-1列上尝试 判断列和两条对角线是否存在元素
 */
const solveNQueens = function (n: number): string[][] {
  const res: string[][] = []

  // path表示每行存第几列，记录列数
  const bt = (path: number[], curRow: number) => {
    if (curRow === n) {
      res.push(path.map(col => '.'.repeat(col) + 'Q' + '.'.repeat(n - col - 1)))
      return
    }

    for (let i = 0; i < n; i++) {
      if (
        path.some((col, row) => col === i || col + row === i + curRow || col - row === i - curRow)
      )
        continue
      path.push(i)
      bt(path, curRow + 1)
      path.pop()
    }
  }
  bt([], 0)

  return res
}

console.dir(solveNQueens(9), { depth: null })
// 输出：[[".Q..","...Q","Q...","..Q."],["..Q.","Q...","...Q",".Q.."]]
export {}
