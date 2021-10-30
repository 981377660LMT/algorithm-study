import { useUnionFindArray } from '../../14_并查集/推荐使用并查集精简版'

const enum Color {
  Red = 0,
  Black,
  Unvisited,
}

/**
 * @param {number[][]} adjList 邻接表
 * @return {boolean}
 * @description 如果能将一个图的节点集合分割成两个独立的子集 A 和 B ，
 * 并使图中的每一条边的两个节点一个来自 A 集合，
 * 一个来自 B 集合，就将这个图称为 二分图
 * @summary dfs染色 visited时相邻节点不同色 0/-1
 * 求二分图有两种思路，一个是着色法，另外一个是并查集
 */
const isBipartite = (adjList: number[][]): boolean => {
  const colors = Array<Color>(adjList.length).fill(Color.Unvisited)
  let res = true

  for (let i = 0; i < adjList.length; i++) {
    if (colors[i] === Color.Unvisited) dfs(i, Color.Red)
  }

  return res

  function dfs(cur: number, color: Color) {
    colors[cur] = color

    for (const next of adjList[cur]) {
      if (colors[next] === Color.Unvisited) {
        dfs(next, color ^ 1)
      } else {
        if (colors[cur] === colors[next]) {
          return (res = false)
        }
      }
    }
  }
}

// 并查集
// 遍历每个顶点，将`当前顶点`的所有邻接点进行合并
// 若某个邻接点与`当前顶点`已经在一个集合中了，说明不是二分图，返回 false。
const isBipartite2 = (adjList: number[][]): boolean => {
  const uf = useUnionFindArray(adjList.length)

  for (let cur = 0; cur < adjList.length; cur++) {
    const nextPoint = adjList[cur]
    for (const next of adjList[cur]) {
      if (uf.isConnected(cur, next)) return false
      uf.union(nextPoint[0], next)
    }
  }

  return true
}

console.log(
  isBipartite([
    [1, 2, 3],
    [0, 2],
    [0, 1, 3],
    [0, 2],
  ])
)

console.log(
  isBipartite([
    [1, 3],
    [0, 2],
    [1, 3],
    [0, 2],
  ])
)
