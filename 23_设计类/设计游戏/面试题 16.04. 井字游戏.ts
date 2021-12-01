function judgeTictactoe(board: string[]): string {
  const m = board.length
  const n = board[0].length

  // 提取正对角线：
  const mainDiagonal = Array.from(board, (row, index) => row[index])
  // 提取副对角线
  const subDiagonal = Array.from(board, (row, index) => row[m - 1 - index])
  // 提取每一列(转置)
  const col = Array.from<unknown, string[]>({ length: m }, () => Array(n).fill(''))
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      col[i][j] = board[j][i]
    }
  }

  const lines = [...board, ...col.map(row => row.join('')), ...mainDiagonal, ...subDiagonal]
  console.log(lines)

  for (const char of ['A', 'B']) {
    const winLine = char.repeat(m)
    if (lines.includes(winLine)) return char
  }

  return board.some(row => row.includes(' ')) ? 'Pending' : 'Draw'
}

export { judgeTictactoe }

console.log(judgeTictactoe(['O X', ' XO', 'X O']))
// console.log(tictactoe(['OOX', 'XXO', 'OXO']))
// console.log(tictactoe(['OOX', 'XXO', 'OX ']))
