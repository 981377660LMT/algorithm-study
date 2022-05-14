// 给你一个整数 n ，返回 和为 n 的完全平方数的最少数量 。

import { SimpleQueue } from '../../2_queue/Deque/Queue'

// bfs求无权图最短路
function numSquares(n: number): number {
  const queue = new SimpleQueue<[cur: number, cost: number]>([[n, 0]])
  const visited = new Set<number>()

  while (queue.length) {
    const [cur, steps] = queue.shift()!
    if (cur === 0) return steps

    for (let i = 1; cur - i ** 2 >= 0; i++) {
      const next = cur - i ** 2
      if (visited.has(next)) continue
      visited.add(next)
      queue.push([next, steps + 1])
    }
  }

  throw new Error('No Solution')
}

console.log(numSquares(12))

export {}
