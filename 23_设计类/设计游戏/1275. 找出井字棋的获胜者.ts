import { judgeTictactoe } from './面试题 16.04. 井字游戏'

/**
 * @param {number[][]} moves
 * @return {string}
 */
const tictactoe = function (moves: number[][]): string {
  const res = Array.from({ length: 3 }, () => Array(3).fill(' '))
  for (const [index, [x, y]] of moves.entries()) {
    res[x][y] = index & 1 ? 'B' : 'A'
  }

  return judgeTictactoe(res.map(row => row.join('')))
}

console.log(
  tictactoe([
    [0, 0],
    [2, 0],
    [1, 1],
    [2, 1],
    [2, 2],
  ])
)

export {}
