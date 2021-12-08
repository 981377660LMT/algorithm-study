```JS
  function getEulerLoop(adjMap: Map<number, number[]>, start: number): number[] {
    let cur = start
    const stack: number[] = [cur]
    const res: number[] = []

    while (stack.length > 0) {
      if (adjMap.has(cur) && adjMap.get(cur)!.length > 0) {
        stack.push(cur)
        const next = adjMap.get(cur)!.pop()!
        cur = next
      } else {
        res.push(cur)
        cur = stack.pop()!
      }
    }

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

  const oddStartPoint = [...new Set(pairs.flat())].filter(
    key => (outdegree.get(key) || 0) - (indegree.get(key) || 0) === 1
  )

  if (oddStartPoint.length > 0) return oddStartPoint[0]
  else return pairs[0][0]
}
```
