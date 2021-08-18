/**
 * @param {character[][]} board
 * @return {void} Do not return anything, modify board in-place instead.
 */
const solveSudoku = function (board: string[][]): void {
  console.log(board.length)
  const isValidPosition = (board: string[][], row: number, col: number, char: string) => {
    // 判断行里是否重复
    for (let j = 0; j < 9; j++) {
      if (board[row][j] === char) return false
    }

    // 判断列里是否重复
    for (let i = 0; i < 9; i++) {
      if (board[i][col] === char) return false
    }

    // 判断9方格里是否重复
    const startRow = Math.floor(row / 3) * 3
    const startCol = Math.floor(col / 3) * 3
    for (let i = startRow; i < startRow + 3; i++) {
      for (let j = startCol; j < startCol + 3; j++) {
        if (board[i][j] === char) return false
      }
    }

    return true
  }

  /**
   *
   * @param board
   * @param n
   * @returns 表示点是否排好了
   */
  const bt = (board: string[][]): boolean => {
    for (let row = 0; row < 9; row++) {
      for (let col = 0; col < 9; col++) {
        if (board[row][col] !== '.') continue
        for (let i = 1; i <= 9; i++) {
          const char = i.toString()
          if (isValidPosition(board, row, col, char)) {
            board[row][col] = char
            if (bt(board)) return true
            board[row][col] = '.'
          }
        }
        // 回溯，还没有排好
        return false
      }
    }
    // 所有点都continue了，即所有点都已经填好了
    return true
  }
  bt(board)
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
