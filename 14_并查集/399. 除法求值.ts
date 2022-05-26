/**
 * @param {string[][]} equations
 * @param {number[]} values
 * @param {string[][]} queries
 * @return {number[]}
 * @description 如果存在某个无法确定的答案，则用 -1.0 替代这个答案。如果问题中出现了给定的已知条件中没有出现的字符串，也需要用 -1.0 替代这个答案。
 * @summary 有向带权图的dfs
 */
const calcEquation = function (
  equations: string[][],
  values: number[],
  queries: string[][]
): number[] {
  const res: number[] = Array<number>(queries.length).fill(Infinity)
  const adjList = new Map<string, [string, number][]>()

  // 建图
  for (let i = 0; i < equations.length; i++) {
    const u = equations[i][0]
    const v = equations[i][1]
    const weight = values[i]
    if (!adjList.has(u)) adjList.set(u, [[v, weight]])
    else adjList.get(u)!.push([v, weight])
    if (!adjList.has(v)) adjList.set(v, [[u, 1 / weight]])
    else adjList.get(v)!.push([u, 1 / weight])
  }

  // 大海捞针般的dfs适合用生成器
  function* dfs(
    cur: string,
    from: string,
    to: string,
    visited: Set<string>,
    val: number
  ): Generator<number> {
    if (!adjList.has(from) || !adjList.has(to)) yield -1
    if (from === to) yield 1
    if (cur === to) yield val
    for (const [next, weight] of adjList.get(cur)!) {
      if (!visited.has(next)) {
        visited.add(next)
        yield* dfs(next, from, to, visited, val * weight)
      }
    }
  }

  for (let i = 0; i < queries.length; i++) {
    const [from, to] = queries[i]
    const val = dfs(from, from, to, new Set([from]), 1).next().value || -1
    res[i] = val
  }
  console.log(adjList)
  return res
}

console.log(
  calcEquation(
    [
      ['x1', 'x2'],
      ['x2', 'x3'],
      ['x3', 'x4'],
      ['x4', 'x5'],
    ],
    [3.0, 4.0, 5.0, 6.0],
    [
      ['x1', 'x5'],
      ['x5', 'x2'],
      ['x2', 'x4'],
      ['x2', 'x2'],
      ['x2', 'x9'],
      ['x9', 'x9'],
    ]
  )
)
// [ 360, 0.008333333333333333, 20, 1, -1, -1 ]
console.log(
  calcEquation(
    [
      ['a', 'b'],
      ['c', 'd'],
    ],
    [1.0, 1.0],
    [
      ['a', 'c'],
      ['b', 'd'],
      ['b', 'a'],
      ['d', 'c'],
    ]
  )
)
// [-1.0,-1.0,1.0,1.0]
