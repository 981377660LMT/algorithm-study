/* eslint-disable consistent-return */
/* eslint-disable no-param-reassign */

import { UnionFindArray } from '../../../14_并查集/UnionFind'

/**
 * @param n 顶点数
 * @param edges 无向边
 * @returns 最小生成树的权值和，如果不存在则返回 -1
 */
function kruskal1(n: number, edges: [u: number, v: number, w: number][]): number {
  const uf = new UnionFindArray(n)
  let res = 0
  let count = 0

  edges = edges.slice().sort((a, b) => a[2] - b[2])
  edges.forEach(([u, v, w]) => {
    const root1 = uf.find(u)
    const root2 = uf.find(v)
    if (root1 !== root2) {
      uf.union(root1, root2)
      res += w
      count++
      if (count === n - 1) return res
    }
  })

  return -1
}

// !给定无向图的边，求出一个最小生成树(如果不存在,则求出的是森林中的多个最小生成树)
//  这个树也叫Kruskal重构树.

function kruskal2(
  n: number,
  edges: [u: number, v: number, w: number][] | number[][]
): readonly [forestEdges: [u: number, v: number, w: number][], ok: boolean] {
  const uf = new UnionFindArray(n)
  const forestEdges: [u: number, v: number, w: number][] = []
  let edgeCount = 0

  edges = edges.slice().sort((a, b) => a[2] - b[2])
  edges.forEach(([u, v, w]) => {
    const root1 = uf.find(u)
    const root2 = uf.find(v)
    if (root1 !== root2) {
      uf.union(root1, root2)
      forestEdges.push([u, v, w])
      edgeCount++
      if (edgeCount === n - 1) return [forestEdges, true]
    }
  })

  return [forestEdges, false]
}

export { kruskal1, kruskal2 }
