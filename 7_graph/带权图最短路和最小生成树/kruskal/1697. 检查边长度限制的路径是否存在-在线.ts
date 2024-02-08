// https://leetcode.cn/problems/checking-existence-of-edge-length-limited-paths/solution/zai-xian-zuo-fa-shu-shang-bei-zeng-lca-b-lzjq/

import { UnionFindArray } from '../../../14_并查集/UnionFind'
import { SparseTable } from '../../../22_专题/RMQ问题/SparseTable'
import { Tree } from '../../../6_tree/重链剖分/Tree'
import { kruskal2 } from './Kruskal'

// 1.求出最小生成树森林是最优的(Kruskal重构树)
// 2.求出最小生成树森林后，再求出两个点路径上的最大边权

function distanceLimitedPathsExist(n: number, edgeList: number[][], queries: number[][]): boolean[] {
  const [forestEdges] = kruskal2(n, edgeList)
  const uf = new UnionFindArray(n)
  const tree = new Tree(n)
  forestEdges.forEach(([u, v, w]) => {
    uf.union(u, v)
    tree.addEdge(u, v, w)
  })
  tree.build(-1)

  const leaves = Array(n).fill(0)
  forestEdges.forEach(([u, v, w]) => {
    const eid = tree.eid(u, v)
    leaves[eid] = w
  })
  const st = new SparseTable(leaves, () => 0, Math.max)

  const res: boolean[] = Array(queries.length).fill(false)
  for (let i = 0; i < queries.length; i++) {
    const [u, v, limit] = queries[i]
    if (!uf.isConnected(u, v)) continue
    let max = 0
    tree.enumeratePathDecomposition(u, v, false, (start, end) => {
      max = Math.max(max, st.query(start, end))
    })
    res[i] = max < limit
  }
  return res
}

if (require.main === module) {
  const n = 5
  const edgeList = [
    [0, 1, 10],
    [1, 2, 5],
    [2, 3, 9],
    [3, 4, 13]
  ]
  const queries = [
    [0, 4, 14],
    [1, 4, 13]
  ]

  console.log(distanceLimitedPathsExist(n, edgeList, queries))
}
