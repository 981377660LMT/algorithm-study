/**
 * 有向图环检测
 * 想要学习课程 0 ，你需要先完成课程 1 ，我们用一个匹配来表示：[0,1] 。
 */

// dfs: 先构建出图和isVisted集合,再dfs
const canFinish = (numCourses: number, prerequisites: [number, number][]) => {
  if (prerequisites.length === 0) return true

  const inDegrees = Array<number>(numCourses).fill(0)
  const adjList = Array.from<unknown, number[]>({ length: numCourses }, () => [])
  for (const [cur, pre] of prerequisites) {
    inDegrees[cur]++
    adjList[pre].push(cur)
  }

  const queue: number[] = []
  inDegrees.forEach((v, i) => !v && queue.push(i))

  let count = 0
  while (queue.length) {
    const cur = queue.shift()!
    count++
    for (const next of adjList[cur]) {
      inDegrees[next]--
      if (inDegrees[next] === 0) queue.push(next)
    }
  }

  return count === numCourses
}
console.log(
  canFinish(2, [
    [1, 0],
    [0, 1],
  ])
)
console.log(
  canFinish(3, [
    [1, 0],
    [0, 1],
    [1, 2],
  ])
)
console.log(
  canFinish(20, [
    [0, 10],
    [3, 18],
    [5, 5],
    [6, 11],
    [11, 14],
    [13, 1],
    [15, 1],
    [17, 4],
  ])
)
console.log(
  canFinish(4, [
    [1, 0],
    [2, 0],
    [3, 1],
    [3, 2],
  ])
)
export {}
