import { UnionFind } from '../0_并查集'

/**
 * @param {number[][]} edges
 * @return {number[]}
 * @description  树可以看成是一个连通且 无环 的 无向 图。
 * 请找出一条可以删去的边，删除后可使得剩余部分是一个有着 n 个节点的树。如果有多个答案，则返回数组 edges 中最后出现的边。
 */
const findRedundantConnection = function (edges: number[][]): number[] {
  const uf = new UnionFind(1000)
  const res: number[][] = []
  for (const [v, w] of edges) {
    if (uf.isConnected(v, w)) res.push([v, w])
    else uf.union(v, w)
  }
  return res.pop()!
}

console.log(
  findRedundantConnection([
    [1, 2],
    [2, 3],
    [3, 4],
    [1, 4],
    [1, 5],
  ])
)
// [1,4]

export {}
