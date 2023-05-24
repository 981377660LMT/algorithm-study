// 你也得到一个数组 queries ，其中 queries[j] = [uj, vj]。
// 对于第 j 个查询，您应该回答课程 uj 是否是课程 vj 的先决条件。
// 返回一个布尔数组 answer ，其中 answer[j] 是第 j 个查询的答案。

import { Floyd } from './Floyd'

function checkIfPrerequisite(n: number, prerequisites: number[][], queries: number[][]): boolean[] {
  const F = new Floyd(n)
  prerequisites.forEach(([u, v]) => F.addDirectedEdge(u, v, 1))
  return queries.map(([u, v]) => F.dist(u, v) !== -1)
}

export {}
