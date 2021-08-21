/**
 * @param {number[][]} edges
 * @return {number[]}
 * @description
 * 该图由一个有着 n 个节点（节点值不重复，从 1 到 n）的树及一条附加的有向边构成。
 * 返回一条能删除的边，使得剩下的图是有 n 个节点的有根树。若有多个答案，返回最后出现在给定二维数组的答案。
 * @summary
 * 两种情况:入度为2的节点不合理 :We can keep a dict or a list parent to record each node’s parent.
 * / 形成环的边不合理 如果v,w已经联通，那么v,w不合理
 */
const findRedundantDirectedConnection = function (edges: number[][]): number[] {
  const n = edges.length

  // 并查集简化版
  // 为了简化将数组多加一位 第i位代表i
  const uf = () => {
    const parent = Array.from<number, number>({ length: n + 1 }, (_, i) => i)
    const find = (val: number) => {
      while (parent[val] !== val) {
        val = parent[val]
      }
      return val
    }
    const union = (pre: number, next: number) => {
      parent[next] = pre
    }
    return { union, find }
  }

  // 第一种情况:入度为2的节点不合理;这是肯定有两条侯选边，删一条边之后判断是不是树
  const isTreeAfterRemoveEdge = (edges: number[][], deleteEdge: number) => {
    const { union, find } = uf()
    for (let i = 0; i < n; i++) {
      if (i === deleteEdge) continue
      const [v, w] = edges[i]
      if (find(v) === w) return false
      union(v, w)
    }
    return true
  }

  // 第二种情况:形成环的边不合理
  const removeEdgeInCycle = (edges: number[][]) => {
    const { union, find } = uf()
    for (const [v, w] of edges) {
      if (find(v) === w) return [v, w]
      union(v, w)
    }
  }

  const main = () => {
    // 记录节点入度
    const indgrees = Array<number>(n + 1).fill(0)
    for (const [_, w] of edges) {
      indgrees[w]++
    }

    // 找入度为2的节点所对应的边，（如果有的话就两条边）注意要倒序
    const candidates: number[] = []
    for (let i = n - 1; i >= 0; i--) {
      if (indgrees[edges[i][1]] == 2) {
        candidates.push(i)
      }
    }

    // 第一种情况：入度2
    if (candidates.length) {
      if (isTreeAfterRemoveEdge(edges, candidates[0])) return edges[candidates[0]]
      else return edges[candidates[1]]
    }

    // 第二种情况：成环
    return removeEdgeInCycle(edges)!
  }

  return main()
}

// console.log(
//   findRedundantDirectedConnection([
//     [1, 2],
//     [2, 3],
//     [3, 4],
//     [4, 1],
//     [1, 5],
//   ])
// )
// // [4,1]
// console.log(
//   findRedundantDirectedConnection([
//     [2, 1],
//     [3, 1],
//     [4, 2],
//     [1, 4],
//   ])
// )

console.log(
  findRedundantDirectedConnection([
    [1, 2],
    [1, 3],
    [2, 3],
  ])
)

// [2,1]
export {}
