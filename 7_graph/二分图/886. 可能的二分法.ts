/**
 *
 * @param n
 * @param dislikes
 * @description
 * 给定一组 N 人（编号为 1, 2, ..., N）， 我们想把每个人分进任意大小的两组。
 * 每个人都可能不喜欢其他人，那么他们不应该属于同一组。
 * 当可以用这种方法将每个人分进两组时，返回 true；否则返回 false。
 * @summary
 * @link https://leetcode-cn.com/problems/possible-bipartition/solution/dfs-jin-xing-er-fen-tu-ran-se-wo-lai-gei-l2p3/
 * 👆节省空间的做法
 * 考虑由给定的 “不喜欢” 边缘形成的 N 人的图表。我们要检查这个图的每个连通分支是否为二分的。
 */
function possibleBipartition(n: number, dislikes: number[][]): boolean {
  // 邻接表
  const adjList = Array.from<number, number[]>({ length: n + 1 }, () => [])
  for (const [v, w] of dislikes) {
    adjList[v].push(w)
    adjList[w].push(v)
  }
  console.log(adjList)
  const colors = Array<number>(n + 1).fill(-1) // -1 0 1

  const dfs = (cur: number, curColor: number, colors: number[]): boolean => {
    colors[cur] = curColor
    for (const next of adjList[cur]) {
      if (colors[next] !== -1 && colors[next] === colors[cur]) return false
      if (colors[next] === -1 && !dfs(next, curColor ^ 1, colors)) return false
    }
    return true
  }

  for (let i = 0; i < n; i++) {
    if (colors[i] === -1 && !dfs(i, 0, colors)) return false
  }

  return true
}

console.log(
  possibleBipartition(4, [
    [1, 2],
    [1, 3],
    [2, 4],
  ])
)

export {}
