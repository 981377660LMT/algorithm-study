// 求有向图传递闭包模板题
// 1462. 课程表 IV
// https://leetcode.cn/problems/course-schedule-iv/

import { TransitiveClosure } from './TransitiveClosure'

function checkIfPrerequisite(n: number, prerequisites: number[][], queries: number[][]): boolean[] {
  const T = new TransitiveClosure(n)
  prerequisites.forEach(([u, v]) => T.addDirectedEdge(u, v))
  const res: boolean[] = Array(queries.length)
  queries.forEach(([pre, cur], i) => {
    res[i] = T.canReach(pre, cur)
  })
  return res
}

export {}

if (require.main === module) {
  // test 5000*5000
  const n = 5000
  const pre = Array.from({ length: n }, (_, i) => [i, i + 1])
  const queries = Array.from({ length: n }, (_, i) => [i, i + 1])
  console.time('checkIfPrerequisite')
  checkIfPrerequisite(n, pre, queries)
  console.timeEnd('checkIfPrerequisite')
}
