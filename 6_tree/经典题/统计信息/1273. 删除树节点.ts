// 请你删除节点值之和为 0 的每一棵子树。
// 在完成所有删除之后，返回树中剩余节点的数目。

function deleteTreeNodes(nodes: number, parent: number[], value: number[]): number {
  const adjList = Array.from<unknown, number[]>({ length: nodes }, () => [])
  for (let i = 0; i < nodes; i++) {
    if (parent[i] == -1) continue
    adjList[parent[i]].push(i)
  }

  const subTreeSum = value.slice()
  const subTreeCount = Array<number>(nodes).fill(1)
  dfs(0)
  return subTreeCount[0]

  // 关键：返回前总体更新subTreeCount；对每个next局部更新subTreeSum/subTreeCount
  function dfs(cur: number): void {
    for (const next of adjList[cur]) {
      dfs(next)

      // 后序dfs，由这个分支更新
      subTreeSum[cur] += subTreeSum[next]
      subTreeCount[cur] += subTreeCount[next]
    }

    // 统计完了 准备回溯
    if (subTreeSum[cur] === 0) subTreeCount[cur] = 0
  }
}

console.log(deleteTreeNodes(7, [-1, 0, 0, 1, 2, 2, 2], [1, -2, 4, 0, -2, -1, -1]))
// 输出：2
