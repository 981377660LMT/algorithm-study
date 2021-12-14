const enum State {
  Unvisited,
  Oncycle,
  Safe,
}

/**
 * @param {number[][]} graph 邻接表
 * @return {number[]}
 * 无论每一步选择沿哪条有向边行走，最后必然在有限步内到达终点，则将该起始节点称作是 安全 的。
 * 返回一个由图中所有安全的起始节点组成的数组作为答案。答案数组中的元素应当按 升序 排列。
 * @summary
   若起始节点位于一个环内，或者能到达一个环，则该节点不是安全的。
   否则，该节点是安全的
 */
var eventualSafeNodes = function (graph: number[][]): number[] {
  const n = graph.length
  const color = Array<number>(n).fill(0)
  const res: number[] = []
  // 是否有环经过cur
  const dfs = (cur: number): boolean => {
    if (color[cur] !== State.Unvisited) return color[cur] === State.Safe
    color[cur] = State.Oncycle
    for (const next of graph[cur]) {
      if (!dfs(next)) return false
    }

    color[cur] = State.Safe
    return true
  }

  for (let i = 0; i < n; i++) {
    if (dfs(i)) res.push(i)
  }

  return res
}

console.log(eventualSafeNodes([[1, 2], [2, 3], [5], [0], [5], [], []]))
// 三色标记法
// 白色（用 0 表示）：该节点尚未被访问；
// 灰色（用 1 表示）：该节点位于递归栈中，或者在某个环上；
// 黑色（用 2 表示）：该节点搜索完毕，是一个安全节点。
// 当我们首次访问一个节点时，将其标记为灰色，
// 并继续搜索与其相连的节点。
// 如果在搜索过程中遇到了一个灰色节点，则说明找到了一个环，
// 此时退出搜索，栈中的节点仍保持为灰色
export {}
