import { hasCycle } from './有向图环检测dfs'

function canFinish(numCourses: number, prerequisites: number[][]): boolean {
  const adjList = Array.from<any, number[]>({ length: numCourses }, () => [])
  for (const [u, v] of prerequisites) {
    adjList[u].push(v)
  }

  return !hasCycle(numCourses, adjList)
}

console.log(
  canFinish(2, [
    [1, 0],
    [0, 1]
  ])
)
