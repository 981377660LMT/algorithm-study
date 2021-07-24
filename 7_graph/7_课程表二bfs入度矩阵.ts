/**
 * 入度+bfs
 */

// dfs: 先构建出图和isVisted集合,再dfs
const canFinish = (numCourses: number, prerequisites: [number, number][]) => {
  const res: number[] = []

  // 1. 计算入度数组
  const inDegrees = Array.from({ length: numCourses }, () => 0)
  prerequisites.forEach(tuple => inDegrees[tuple[0]]++)

  // 2. bfs queue里全部为为入度为0的点
  const bfsQueue = inDegrees.reduce<number[]>(
    (pre, cur, index) => (cur === 0 ? pre.concat(index) : pre),
    []
  )

  while (bfsQueue.length) {
    const head = bfsQueue.shift()!
    res.push(head)
    // 3. 每次shift出一个就更新所有依赖的点的入度
    for (const [k, v] of prerequisites) {
      if (v === head) {
        inDegrees[k]--
        // 4. 推入新的入度为0的点
        if (inDegrees[k] === 0) bfsQueue.push(k)
      }
    }
  }

  return res.length === numCourses ? res : []
}

// console.log(
//   canFinish(3, [
//     [1, 0],
//     [0, 1],
//     [1, 2],
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
// console.log(canFinish(1, []))
console.log(canFinish(2, [[0, 1]]))

export {}
