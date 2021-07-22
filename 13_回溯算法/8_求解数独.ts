/**
 * @param {character[][]} board
 * @return {void} Do not return anything, modify board in-place instead.
 */
const solveSudoku = function (board: string[][]): void {
  const n = board.length

  const isValidPosition = (
    board: string[][],
    row: number,
    col: number,
    n: number,
    char: string
  ) => {
    const blockRow = Math.floor(row / 3) * 3
    const blockCol = Math.floor(col / 3) * 3
    for (let i = 0; i < n; i++) {
      // 验证行列没有且九宫格内没有
      if (board[row][i] === char || board[i][col] === char) return false
      const curRow = blockRow + Math.floor(i / 3)
      const curCol = blockCol + Math.floor(i % 3)
      if (board[curRow][curCol] === char) return false
    }
    return true
  }

  const bt = (board: string[][], n: number): boolean => {
    for (let row = 0; row < n; row++) {
      for (let col = 0; col < n; col++) {
        if (board[row][col] !== '.') continue

        for (let i = 1; i <= 9; i++) {
          const char = i.toString()
          if (isValidPosition(board, row, col, n, char)) {
            board[row][col] = char
            if (bt(board, n)) return true
          }
        }

        board[row][col] = '.'
        return false
      }
    }
    return true
  }
  bt(board, n)
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

solveSudoku(sudoku)

console.table(sudoku)

export {}
