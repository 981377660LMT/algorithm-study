import { useUnionFindArray } from './useUnionFind'

const DIR4 = [
  [1, 0],
  [0, 1],
  [-1, 0],
  [0, -1],
]

/**
 * @param {character[][]} board
 * @return {void} Do not return anything, modify board in-place instead.
 * @description 任何不在边界上，或不与边界上的 'O' 相连的 'O' 最终都会被填充为 'X'
 * @summary 将非边界的点填成第三种颜色
 */

function solve(board: string[][]): void {
  const row = board.length
  const col = board[0].length

  const dummy = row * col
  const uf = useUnionFindArray(dummy + 10)

  for (let r = 0; r < row; r++) {
    for (let c = 0; c < col; c++) {
      if (board[r][c] == 'O') {
        // 这种写法很好，减少循环二维数组的次数
        if (r == 0 || c == 0 || r == row - 1 || c == col - 1) {
          uf.union(r * col + c, dummy)
        } else {
          for (const [dr, dc] of DIR4) {
            const nr = r + dr
            const nc = c + dc
            if (board[nr][nc] == 'O') uf.union(nr * col + nc, r * col + c)
          }
        }
      }
    }
  }

  for (let i = 1; i < row - 1; i++) {
    for (let j = 1; j < col - 1; j++) {
      if (!uf.isConnected(i * col + j, dummy)) board[i][j] = 'X'
    }
  }
}

const res = [
  ['X', 'X', 'X', 'X'],
  ['X', 'O', 'O', 'X'],
  ['X', 'X', 'O', 'X'],
  ['X', 'O', 'X', 'X'],
]

solve(res)

console.table(res)
export {}
// 把那些不需要被替换的 O 看成一个拥有独门绝技的门派，
// 它们有一个共同祖师爷叫 dummy，这些 O 和 dummy 互相连通，
// 而那些需要被替换的 O 与 dummy 不连通。
