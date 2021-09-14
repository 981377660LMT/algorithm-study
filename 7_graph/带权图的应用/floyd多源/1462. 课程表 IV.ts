/**
 *
 * @param numCourses  2 <= n <= 100
 * @param prerequisites
 * @param queries
 * 对于每个查询对 queries[i] ，请判断 queries[i][0] 是否是 queries[i][1] 的先修课程
 * @summary
 * 判断两点间是否联通，等价于距离不为Infinity
 */
function checkIfPrerequisite(
  numCourses: number,
  prerequisites: number[][],
  queries: number[][]
): boolean[] {
  // 构建dist矩阵
  const dist = Array.from<number, number[]>({ length: numCourses }, () =>
    Array(numCourses).fill(Infinity)
  )

  for (const [i, j] of prerequisites) {
    dist[i][j] = 1 // 有通路 权重为1
  }

  for (let i = 0; i < numCourses; i++) {
    dist[i][i] = 0
  }

  for (let m = 0; m < numCourses; m++) {
    for (let i = 0; i < numCourses; i++) {
      for (let j = 0; j < numCourses; j++) {
        dist[i][j] = Math.min(dist[i][j], dist[i][m] + dist[m][j])
      }
    }
  }

  return queries.map(([v, w]) => dist[v][w] !== Infinity)
}

console.log(
  checkIfPrerequisite(
    2,
    [[1, 0]],
    [
      [0, 1],
      [1, 0],
    ]
  )
)
