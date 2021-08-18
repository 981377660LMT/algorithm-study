/**
 * @param {character[][]} board
 * 等价于对每个九宫格没看过
 * 已知第i个元素，求在矩阵r行c列中的位置，即为~~(i / c) 行和 i % c 列。
 * @summary 思路：需要27个set记录 (9行，9列，9个九宫格) 遍历每一个位置时判断即可
 */
const isValidSudoku = function (board: string[][]): boolean {
  const rowSet = Array.from<unknown, Set<string>>({ length: 9 }, () => new Set())
  const colSet = Array.from<unknown, Set<string>>({ length: 9 }, () => new Set())
  const gridSet = Array.from<unknown, Set<string>>({ length: 9 }, () => new Set())

  for (let i = 0; i < 9; i++) {
    for (let j = 0; j < 9; j++) {
      const num = board[i][j]
      if (num === '.') continue

      if (rowSet[i].has(num)) {
        return false
      } else {
        rowSet[i].add(num)
      }

      if (colSet[j].has(num)) {
        return false
      } else {
        colSet[j].add(num)
      }

      const gridIndex = 3 * ~~(i / 3) + ~~(j / 3)
      if (gridSet[gridIndex].has(num)) {
        return false
      } else {
        gridSet[gridIndex].add(num)
      }
    }
  }

  return true
}

const sudoku = [
  ['5', '3', '.', '.', '7', '.', '.', '.', '.'],
  ['6', '.', '.', '1', '9', '5', '.', '.', '.'],
  ['.', '9', '8', '.', '.', '.', '.', '6', '.'],
  ['8', '.', '.', '.', '6', '.', '.', '.', '3'],
  ['4', '.', '.', '8', '.', '3', '.', '.', '1'],
  ['7', '.', '.', '.', '2', '.', '.', '.', '6'],
  ['.', '6', '.', '.', '.', '.', '2', '8', '.'],
  ['.', '.', '.', '4', '1', '9', '.', '.', '5'],
  ['.', '.', '.', '.', '8', '.', '.', '7', '9'],
]

console.log(isValidSudoku(sudoku))

export {}
