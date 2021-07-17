'use strict'
/**
 * 构造入度矩阵
 */
Object.defineProperty(exports, '__esModule', { value: true })
// dfs: 先构建出图和isVisted集合,再dfs
const canFinish = (numCourses, prerequisites) => {
  if (prerequisites.length === 0) return []
  const res = []
  // 1. 计算入度
  const inDegrees = Array.from({ length: numCourses }, () => 0)
  prerequisites.forEach(tuple => inDegrees[tuple[0]]++)
  // 2. bfs queue里全部为为入度为0的点
  const bfsQueue = inDegrees.filter(inDegree => inDegree === 0)
  while (bfsQueue.length) {
    const head = bfsQueue.shift()
    res.push(head)
    // 3. 每次shift出一个就更新所有依赖的点的入度
    for (const [k, v] of prerequisites) {
      //   console.log(inDegrees)
      if (v === head) inDegrees[k]--
      // 4. 推入新的入度为0的点
      if (inDegrees[k] === 0) bfsQueue.push(k)
    }
  }
  return res
}

console.log(
  canFinish(4, [
    [1, 0],
    [2, 0],
    [3, 1],
    [3, 2],
  ])
)
