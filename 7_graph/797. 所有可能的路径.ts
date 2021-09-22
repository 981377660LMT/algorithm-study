// 给你一个有 n 个节点的 有向无环图（DAG），请你找出所有从节点 0 到节点 n-1 的路径并输出（不要求按特定顺序）
function allPathsSourceTarget(graph: number[][]): number[][] {
  const res: number[][] = []
  const n = graph.length

  const dfs = (cur: number, path: number[], target: number) => {
    if (cur === target) {
      return res.push(path.slice())
    }

    for (const next of graph[cur]) {
      path.push(next)
      dfs(next, path, target)
      path.pop()
    }
  }

  dfs(0, [0], n - 1)

  return res
}

console.log(allPathsSourceTarget([[1, 2], [3], [3], []]))
