function tictactoe(board: string[]): string {
  const m = board.length
  const n = board[0].length
  // 提取正对角线：
  const mainDiagonal = Array.from(board, (row, index) => row[index])
  // 提取副对角线
  const subDiagonal = Array.from(board, (row, index) => row[m - index - 1])
  // 提取每一列(转置)
  const col = Array.from<unknown, string[]>({ length: m }, () => Array(n).fill(''))
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      col[i][j] = board[j][i]
    }
  }

  for (const char of ['X', 'O']) {
    if (board.some(row => row === char.repeat(m))) return char
    if (col.some(col => col.join('') === char.repeat(m))) return char
    if (mainDiagonal.join('') === char.repeat(m)) return char
    if (subDiagonal.join('') === char.repeat(m)) return char
  }

  return board.some(row => row.includes(' ')) ? 'Pending' : 'Draw'
}

console.log(tictactoe(['O X', ' XO', 'X O']))
console.log(tictactoe(['OOX', 'XXO', 'OXO']))
console.log(tictactoe(['OOX', 'XXO', 'OX ']))
