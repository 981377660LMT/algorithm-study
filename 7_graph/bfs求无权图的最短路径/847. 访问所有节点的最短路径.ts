// 基于状态压缩的广度优先搜索算法
/**
 * @param {number[][]} graph
 * @return {number}
 * @description
 * // 最短路径:bfs
   // 题目的特点: 每个节点都可以被遍历多次
   // 我们需要找到另外一种状态 state， state 能唯一确定节点 x 在遍历中的信息
   // 哈密尔顿路径的变形版
 */
const shortestPathLength = function (graph: number[][]): number {
  const len = graph.length
  // 每个节点都有一个set 记录以此节点为起点的不同state
  const visited = Array.from<unknown, Set<number>>({ length: len }, () => new Set())
  // 当前点,当前状态,走过的距离
  const queue: [number, number, number][] = []

  // 每个节点开始bfs
  for (let i = 0; i < len; i++) {
    queue.push([i, 1 << i, 0])
    visited[i].add(i << i)
  }

  while (queue.length) {
    const [cur, state, dis] = queue.shift()!
    for (const next of graph[cur]) {
      const newState = state | (1 << next)
      // 全为1 所有节点都走到了
      if (newState === (1 << len) - 1) {
        console.log(visited, cur)
        return dis + 1
      }
      if (!visited[next].has(newState)) {
        visited[next].add(newState)
        queue.push([next, newState, dis + 1])
      }
    }
  }

  return 0
}

console.log(shortestPathLength([[1, 2, 3], [0], [0], [0]]))

export {}
