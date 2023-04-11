/**
 * @param {number} n
 * @param {number[][]} edges
 * @return {number[]}
 * @description 第 i 条边连接节点 edges[i][0] 和 edges[i][1]
 * 返回一个表示节点 i 与其他所有节点距离之和的列表 ans。
 * @description 思路:无向图邻接表+pre+自底向上初始化数据+自顶向下求答案
 * @description 在不对其实施暴力的情况下，怎么样高效求出每个点作为的root的情况
 * @summary 整体来说，这道题算是相当有难度的一道题，
 * 同时考察了邻接链表的建立，无向图的遍历，树的dfs自顶向下先序和自底向上后序遍历，以及对复杂度的拆分能力
 */
const sumOfDistancesInTree = (n: number, edges: number[][]): number[] => {
  const adjList: number[][] = Array(n).fill(0)
  for (let i = 0; i < n; i++) adjList[i] = []

  // 用于计算根的值
  const depth = Array<number>(n).fill(0)
  // 用于计算每个节点的子节点数(包括自己)
  const count = Array<number>(n).fill(1)
  const res = Array<number>(n).fill(0)

  edges.forEach(([from, to]) => {
    adjList[from].push(to)
    adjList[to].push(from)
  })

  // 得到了所有点的层数和子节点数量（包含自己）
  const postOrderDFSForDepthAndCount = (cur: number, pre: number) => {
    for (const next of adjList[cur]) {
      if (next === pre) continue
      depth[next] = depth[cur] + 1
      // 自底向上后序
      postOrderDFSForDepthAndCount(next, cur)
      count[cur] += count[next]
    }
  }
  postOrderDFSForDepthAndCount(0, -1)

  // 初始值
  res[0] = depth.reduce((pre, cur) => pre + cur, 0)
  const preDFSForAnswer = (cur: number, pre: number) => {
    for (const next of adjList[cur]) {
      if (next === pre) continue
      // 自顶向下先序
      // 注意这个递推的含义 减是靠近子节点，加是远离非子节点
      res[next] = res[cur] - count[next] + (n - count[next])
      preDFSForAnswer(next, cur)
    }
  }
  preDFSForAnswer(0, -1)

  console.log(adjList, depth, count)
  return res
}

console.log(
  sumOfDistancesInTree(6, [
    [0, 1],
    [0, 2],
    [2, 3],
    [2, 4],
    [2, 5]
  ])
)

export {}
