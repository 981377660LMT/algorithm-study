import { UnionFind } from '../0_并查集'

/**
 * @param {number} n
 * @param {number[][]} edgeList
 * @param {number[][]} queries
 * @return {boolean[]}
 * @description 无向有权图 两个点之间可能有 超过一条边 。
 * 对于每个查询 queries[j] ，判断是否存在从 pj 到 qj 的路径，且这条路径上的每一条边都 严格小于 limitj 。
 * @summary 从最严格的限制开始存边 直到大于limt停止 存在路径等价于看是否联通
 */
const distanceLimitedPathsExist = function (
  n: number,
  edgeList: number[][],
  queries: number[][]
): boolean[] {
  const res = Array<boolean>(queries.length).fill(false)
  edgeList.sort((a, b) => a[2] - b[2])
  queries.sort((a, b) => a[2] - b[2])
  const uf = new UnionFind()

  let j = 0
  for (let i = 0; i < queries.length; i++) {
    const [from, to, limit] = queries[i]

    while (j < edgeList.length && edgeList[j][2] < limit) {
      const [u, v] = edgeList[j]
      uf.add(u).add(v).union(u, v)
      j++
    }

    res[i] = uf.isConnected(from, to)
  }

  return res
}

console.log(
  distanceLimitedPathsExist(
    3,
    [
      [0, 1, 2],
      [1, 2, 4],
      [2, 0, 8],
      [1, 0, 16],
    ],
    [
      [0, 1, 2],
      [0, 2, 5],
    ]
  )
)
// [false,true]
