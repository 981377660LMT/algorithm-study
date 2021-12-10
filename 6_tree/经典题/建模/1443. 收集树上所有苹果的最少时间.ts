// 通过树上的一条边，需要花费 1 秒钟
// 你从 节点 0 出发，请你返回最少需要多少秒，可以收集到所有苹果，并回到节点 0 。
// 1 <= n <= 10^5
// 无向边edges

// 总结：看走了多少条边
// 1.将需要收集的苹果的祖先节点的hasApple状态自下而上标记为True
// 2.遍历无向树edges中各个边，如果构成边的两点的hasApple状态都为True，说明需要经过这条边，res+=1，而收集的过程中每条边需要走两次，所以2*res即为所求
// https://leetcode-cn.com/problems/minimum-time-to-collect-all-apples-in-a-tree/solution/python3-zi-di-xiang-shang-dfs-by-yim-6-aub7/
function minTime(n: number, edges: number[][], hasApple: boolean[]): number {
  const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
  for (const [u, v] of edges) {
    adjList[u].push(v)
    adjList[v].push(u)
  }

  let res = 0
  dfs(0, new Set())
  for (const [u, v] of edges) {
    if (hasApple[u] && hasApple[v]) res++
  }

  return res * 2

  function dfs(cur: number, visited: Set<number>): void {
    visited.add(cur)

    for (const next of adjList[cur]) {
      if (visited.has(next)) continue

      dfs(next, visited)
      // 从下往上标记
      if (hasApple[next]) hasApple[cur] = true
    }
  }
}

console.log(
  minTime(
    7,
    [
      [0, 1],
      [0, 2],
      [1, 4],
      [1, 5],
      [2, 3],
      [2, 6],
    ],
    [false, false, true, false, true, true, false]
  )
)
// 输出：8
// 解释：上图展示了给定的树，其中红色节点表示有苹果。
// 一个能收集到所有苹果的最优方案由绿色箭头表示。
