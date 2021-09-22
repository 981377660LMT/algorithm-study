/**
 * 判断有向图中不存在环，此处使用dfs
 */

// dfs: 先构建出图和isVisted集合,再dfs
const canFinish = (numCourses: number, prerequisites: [number, number][]) => {
  if (!prerequisites.length) return true
  const hasCycle = (n: number, prerequisites: [number, number][]) => {
    const visited = Array<boolean>(n).fill(false)
    const path = new Set<number>()
    const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
    for (const [cur, pre] of prerequisites) {
      adjList[pre].push(cur)
    }

    /**
     *
     * @param cur
     * @param visited
     * @param path 用于检测有向图的环,回溯需要删除
     * @returns
     */
    const dfs = (cur: number, visited: boolean[], path: Set<number>): boolean => {
      visited[cur] = true
      path.add(cur)

      for (const next of adjList[cur]) {
        if (!visited[next]) {
          if (dfs(next, visited, path)) return true
        } else {
          if (path.has(next)) {
            return true
          }
        }
      }

      path.delete(cur)
      return false
    }

    for (let i = 0; i < n; i++) {
      if (visited[i]) continue
      if (dfs(i, visited, path)) return true
    }

    return false
  }

  return !hasCycle(numCourses, prerequisites)
}

console.log(
  canFinish(2, [
    [1, 0],
    [0, 1],
  ])
)

// console.log(
//   canFinish(3, [
//     [1, 0],
//     [0, 1],
//     [1, 2],
//   ])
// )
// console.log(
//   canFinish(20, [
//     [0, 10],
//     [3, 18],
//     [5, 5],
//     [6, 11],
//     [11, 14],
//     [13, 1],
//     [15, 1],
//     [17, 4],
//   ])
// )
// console.log(
//   canFinish(4, [
//     [1, 0],
//     [2, 0],
//     [3, 1],
//     [3, 2],
//   ])
// )
export {}
