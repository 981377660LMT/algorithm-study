import { ArrayDeque } from '../../2_queue/Deque/ArrayDeque'

type Num = number
type Step = number

/**
 * @param {number[][]} board  大小为 n x n 的整数矩阵 board ，方格按从 1 到 n2 编号，
 * @return {number} 返回达到编号为 n2 的方格所需的最少移动次数，如果不可能，则返回 -1。
 */
const snakesAndLadders = function (board: number[][]): number {
  const n = board.length
  const numToPosition = (num: number, n: number) => {
    const [row, col] = [~~((num - 1) / n), (num - 1) % n]
    if (row & 1) return [n - 1 - row, n - 1 - col]
    else return [n - 1 - row, col]
  }

  const visited = new Set<number>()
  const queue = new ArrayDeque<[Num, Step]>(10 ** 4)
  queue.push([1, 0])
  while (queue.length) {
    let [num, step] = queue.shift()!
    const [row, col] = numToPosition(num, n)

    // 蛇或者梯子
    if (board[row][col] !== -1) num = board[row][col]
    if (num === n * n) return step
    // 掷色子
    for (let i = 1; i <= 6; i++) {
      const nextNum = num + i
      if (nextNum > n * n || visited.has(nextNum)) continue
      visited.add(nextNum)
      queue.push([nextNum, step + 1])
    }
  }

  return -1
}

console.log(
  snakesAndLadders([
    [-1, -1, -1, -1, -1, -1],
    [-1, -1, -1, -1, -1, -1],
    [-1, -1, -1, -1, -1, -1],
    [-1, 35, -1, -1, 13, -1],
    [-1, -1, -1, -1, -1, -1],
    [-1, 15, -1, -1, -1, -1],
  ])
)
