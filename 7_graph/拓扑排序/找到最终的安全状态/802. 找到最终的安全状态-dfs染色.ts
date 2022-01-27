/**
 * @param {number[][]} graph 邻接表
 * @return {number[]}
 * 无论每一步选择沿哪条有向边行走，最后必然在有限步内到达终点，则将该起始节点称作是 安全 的。
 * 返回一个由图中所有安全的起始节点组成的数组作为答案。答案数组中的元素应当按 升序 排列。
 */
function eventualSafeNodes(graph: number[][]): number[] {
  const n = graph.length
  // -1:未访问 0:不安全 1:安全
  const colors = Array<-1 | 0 | 1>(n).fill(-1)

  for (let i = 0; i < n; i++) {
    if (colors[i] === -1) dfs(i)
  }

  const res: number[] = []
  for (const [index, color] of colors.entries()) {
    if (color === 1) res.push(index)
  }
  return res

  function dfs(cur: number): boolean {
    if (colors[cur] !== -1) return colors[cur] === 1

    colors[cur] = 0
    for (const next of graph[cur]) {
      if (!dfs(next)) return false
    }

    colors[cur] = 1
    return true
  }
}

console.log(eventualSafeNodes([[1, 2], [2, 3], [5], [0], [5], [], []]))

export {}
