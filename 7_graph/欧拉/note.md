```JS
function getEulerPath(
  adjMap: Map<number, Set<number>>,
  start: number,
  isDirected: boolean
): number[] {
  let cur = start
  const stack: number[] = [start]
  const res: number[] = []

  while (stack.length > 0) {
    if (adjMap.has(cur) && adjMap.get(cur)!.size > 0) {
      stack.push(cur)
      const next = adjMap.get(cur)!.keys().next().value!
      // 无向图 要删两条边
      if (!isDirected) adjMap.get(next)!.delete(cur)
      cur = next
    } else {
      res.push(cur)
      cur = stack.pop()!
    }
  }

  // 有向图需要反转
  return res.reverse()
}
```

```JS

// 题目给定的图一定满足以下二者之一：
// 所有点入度等于出度；
// 恰有一个点出度 = 入度 + 1（欧拉路径的起点），且恰有一个点入度 = 出度 + 1（欧拉路径的终点），其他点入度等于出度。
function getStart(pairs: number[][]): number {
  const outdegree = new Map<number, number>()
  const indegree = new Map<number, number>()

  for (const [u, v] of pairs) {
    outdegree.set(u, (outdegree.get(u) || 0) + 1)
    indegree.set(v, (indegree.get(v) || 0) + 1)
  }

  // 有向图欧拉路径的起点出度比入度多 1
  const oddStartPoint = [...new Set(pairs.flat())].filter(
    key => (outdegree.get(key) || 0) - (indegree.get(key) || 0) === 1
  )

  if (oddStartPoint.length > 0) return oddStartPoint[0]
  else return pairs[0][0]
}
```
