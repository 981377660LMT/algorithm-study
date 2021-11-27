// 请你找到所有的 最小高度树 并按 任意顺序 返回它们的根节点标签列表。
// 思路：不断删除叶子节点
function findMinHeightTrees(n: number, edges: number[][]): number[] {
  if (n === 1) return [0]
  const adjList = Array.from<unknown, Set<number>>({ length: n }, () => new Set())
  for (const [u, v] of edges) {
    adjList[u].add(v)
    adjList[v].add(u)
  }

  // 加入叶子
  let leaves: number[] = []
  for (const [index, arr] of adjList.entries()) {
    if (arr.size === 1) leaves.push(index)
  }

  while (n >= 3) {
    const newLeaves: number[] = []
    // 进入减
    n -= leaves.length

    for (const cur of leaves) {
      const next = adjList[cur].keys().next().value
      adjList[next].delete(cur)
      if (adjList[next].size === 1) newLeaves.push(next)
    }

    leaves = newLeaves
  }

  return leaves
}

// 如图所示，当根是标签为 1 的节点时，树的高度是 1 ，这是唯一的最小高度树。
console.log(
  findMinHeightTrees(4, [
    [1, 0],
    [1, 2],
    [1, 3],
  ])
)
console.log(
  findMinHeightTrees(3, [
    [0, 1],
    [0, 2],
  ])
)
