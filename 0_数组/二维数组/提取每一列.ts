const board = [
  ['a1', 'b1', 'c1'],
  ['a2', 'b2', 'c2'],
  ['a3', 'b3', 'c3'],
]

// 提取正对角线：
const mainDiagonal = Array.from(board, (row, index) => row[index])
// 提取副对角线
const subDiagonal = Array.from(board, (row, index) => row[board.length - index - 1])
// 提取每一列(转置)
const col = Array.from({ length: board.length }, () => Array(board[0].length).fill(0))
for (let i = 0; i < board.length; i++) {
  for (let j = 0; j < board[0].length; j++) {
    col[i][j] = board[j][i]
  }
}

console.log(mainDiagonal, subDiagonal, col)
