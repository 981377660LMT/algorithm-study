function hasHamiltonLoop(n: number, adjList: number[][], start = 0): boolean {
  const target = (1 << n) - 1
  const memo = new Map<string, boolean>()

  // 必须是这两个参数
  const dfs = (cur: number, visited: number): boolean => {
    if (visited === target && cur === start) return true
    const key = `${cur}#${visited}`
    if (memo.has(key)) return memo.get(key)!

    for (const next of adjList[cur]) {
      if (visited & (1 << next)) continue
      if (dfs(next, visited | (1 << next))) {
        memo.set(key, true)
        return true
      }
    }

    memo.set(key, false)
    return false
  }

  return dfs(0, 0)
}

// 回溯暴力法:O(n!)

console.log(
  hasHamiltonLoop(4, [
    [1, 2, 3],
    [0, 2, 3],
    [0, 1],
    [0, 1],
  ])
)
