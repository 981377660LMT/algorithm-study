/* eslint-disable no-console */
import { Rerooting } from './Rerooting'

// 求每个点作为根节点时，到叶子节点的最大点权和(不包括自身)
function maxOutput(n: number, edges: number[][], price: number[]): number {
  const R = new Rerooting(n)
  for (const [u, v] of edges) {
    R.addEdge(u, v)
  }

  const res = R.reRooting({
    e: () => 0,
    op: (childRes1, childRes2) => Math.max(childRes1, childRes2),
    composition: (fromRes, parent, cur, direction) => {
      if (direction === 0) return fromRes + price[cur]
      return fromRes + price[parent]
    }
  })

  return Math.max(...res)
}

if (require.main === module) {
  console.log(
    maxOutput(
      6,
      [
        [0, 1],
        [1, 2],
        [1, 3],
        [3, 4],
        [3, 5]
      ],
      [9, 8, 7, 6, 10, 5]
    )
  ) // 24
}

export {}
