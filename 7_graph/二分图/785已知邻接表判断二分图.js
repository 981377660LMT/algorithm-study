/**
 * @param {number[][]} graph 邻接表
 * @return {boolean}
 * @description 如果能将一个图的节点集合分割成两个独立的子集 A 和 B ，
 * 并使图中的每一条边的两个节点一个来自 A 集合，
 * 一个来自 B 集合，就将这个图称为 二分图
 * @summary dfs染色 visited时相邻节点不同色 0/-1
 * 求二分图有两种思路，一个是着色法，另外一个是并查集
 */
const isBipartite = graph => {
  const visited = new Set()
  const colors = Array(graph.length).fill(-1)
  let res = true

  /**
   *
   * @param {number} cur
   * @param {number[]} colors
   * @param {Set<number>} visited
   * 优化：使用color===-1 代替visited集合
   */
  const dfs = (cur, color, colors, visited) => {
    visited.add(cur)
    colors[cur] = color
    for (const next of graph[cur]) {
      if (visited.has(next)) {
        console.log(colors, cur, next)
        if (colors[cur] === colors[next]) {
          return (res = false)
        }
      } else {
        dfs(next, 1 - color, colors, visited)
      }
    }
  }

  for (let v = 0; v < graph.length; v++) {
    !visited.has(v) && dfs(v, 0, colors, visited)
  }

  return res
}

console.log(
  isBipartite([
    [1, 2, 3],
    [0, 2],
    [0, 1, 3],
    [0, 2],
  ])
)
