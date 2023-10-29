function maximumPoints(edges: number[][], coins: number[], k: number): number {
  const n = coins.length
  const adjList: number[][] = Array(n)
  for (let i = 0; i < n; i++) adjList[i] = []
  edges.forEach(([u, v]) => {
    adjList[u].push(v)
    adjList[v].push(u)
  })

  const memo = new Map<number, number>()
  const max_ = Math.max(...coins)
  const maxLog = max_.toString(2).length

  const dfs = (cur: number, pre: number, log: number): number => {
    const curCoin = coins[cur] >> log
    const hash = cur * n * n + pre * n + log
    if (memo.has(hash)) return memo.get(hash)!
    let res1 = curCoin >> 1
    let res2 = curCoin - k
    const nexts = adjList[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next_ = nexts[i]
      if (next_ === pre) continue
      res1 += dfs(next_, cur, Math.min(log + 1, maxLog))
      res2 += dfs(next_, cur, log)
    }
    const res = Math.max(res1, res2)
    memo.set(hash, res)
    return res
  }

  return dfs(0, -1, 0)
}

export {}
