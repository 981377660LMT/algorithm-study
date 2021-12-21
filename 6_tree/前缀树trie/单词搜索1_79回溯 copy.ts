/**
 * @param {character[][]} board
 * @param {string} word
 * @return {boolean}
 * @description 注意不要用forEach forEach里面return无效！
 */
const exist = (board: string[][], word: string): boolean => {
  if (board.length === 0) return false
  const r = board.length
  const c = board[0].length
  const next = [
    [-1, 0],
    [0, 1],
    [1, 0],
    [0, -1],
  ]

  // 一般题目返回布尔值，dfs就返回布尔值
  const dfs = (x: number, y: number, step: number): boolean => {
    console.log(board[x][y], word[step])
    // 1. 回溯终点
    if (board[x][y] !== word[step]) return false
    if (step === word.length - 1) return true
    // console.log(step, board, word[step])
    // 2.回溯处理
    // 标记为visited
    board[x][y] = '$'
    for (const [dx, dy] of next) {
      const nextRow = x + dx
      const nextColumn = y + dy
      // 在矩阵中
      if (nextRow >= 0 && nextRow < r && nextColumn >= 0 && nextColumn < c) {
        if (dfs(nextRow, nextColumn, step + 1)) {
          return true
        }
      }
    }

    // 3. 回溯重置
    board[x][y] = word[step]
    return false
  }

  // 4.每个点开始回溯
  for (let i = 0; i < r; i++) {
    for (let j = 0; j < c; j++) {
      if (dfs(i, j, 0)) return true
    }
  }

  return false
}

console.dir(
  exist(
    [
      ['A', 'B', 'C', 'E'],
      ['S', 'F', 'C', 'S'],
      ['A', 'D', 'E', 'E'],
    ],
    'ABCC'
  ),
  { depth: null }
)

export {}
