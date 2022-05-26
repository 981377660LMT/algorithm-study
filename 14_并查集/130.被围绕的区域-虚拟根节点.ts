/**
 * @param {character[][]} board
 * @return {void} Do not return anything, modify board in-place instead.
 * @description 任何不在边界上，或不与边界上的 'O' 相连的 'O' 最终都会被填充为 'X'
 * @summary 将非边界的点填成第三种颜色
 */

import { useUnionFindArray } from './useUnionFind'

/**
 * @param {character[][]} board
 * @return {void} Do not return anything, modify board in-place instead.
 * @description 任何不在边界上，或不与边界上的 'O' 相连的 'O' 最终都会被填充为 'X'
 * @summary 将非边界的点填成第三种颜色
 */
const solve = (board: string[][]): void => {
  const row = board.length
  if (row == 0) return
  const col = board[0].length
  const dummy = row * col
  const uf = useUnionFindArray(dummy + 1)
  const arr = [
    [1, 0],
    [0, 1],
    [-1, 0],
    [0, -1],
  ]

  for (let i = 0; i < row; i++) {
    for (let j = 0; j < col; j++) {
      if (board[i][j] == 'O') {
        // 这种写法很好，减少循环二维数组的次数
        if (i == 0 || j == 0 || i == row - 1 || j == col - 1) {
          uf.union(i * col + j, dummy)
        } else {
          //考察四个方向
          for (let k = 0; k < 4; k++) {
            let x = i + arr[k][0]
            let y = j + arr[k][1]
            if (board[x][y] == 'O') uf.union(x * col + y, i * col + j)
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
