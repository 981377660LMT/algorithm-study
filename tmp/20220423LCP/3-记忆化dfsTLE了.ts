const POS: number[][] = []
for (const r of [0, 1, 2]) {
  for (const c of [0, 1, 2]) {
    POS.push([r, c])
  }
}

// 注意js的cache很容易TLE
// 如果dfs超时 要优化就要把dfs优化成dp 这样常数会小一些
function getMaximumNumber(moles: number[][]): number {
  const times = new Set<number>([0])
  const record = new Map<number, Set<number>>()
  for (const [t, r, c] of moles) {
    times.add(t)
    !record.has(t) && record.set(t, new Set())
    record.get(t)!.add(r * 3 + c)
  }

  const allTimes = [...times].sort((a, b) => a - b)
  const n = allTimes.length

  const cache = Array.from({ length: n + 10 }, () => new Int32Array(9).fill(-1))
  const res = dfs(0, 1, 1)

  return res

  function dfs(index: number, row: number, col: number): number {
    if (index >= n) return 0
    // const key = `${index}_${row}_${col}`
    if (cache[index][row * 3 + col] !== -1) {
      return cache[index][row * 3 + col]
    }

    const curTime = allTimes[index]
    const nextTime = allTimes[index + 1] ?? curTime
    const diff = nextTime - curTime

    const cur = Number(record.has(curTime) && record.get(curTime)!.has(row * 3 + col))
    let nextMax = 0

    for (const [nr, nc] of POS) {
      if (Math.abs(row - nr) + Math.abs(col - nc) <= diff) {
        nextMax = Math.max(nextMax, dfs(index + 1, nr, nc))
      }
    }

    cache[index][row * 3 + col] = cur + nextMax
    return cur + nextMax
  }
}

console.log(
  getMaximumNumber([
    [1, 1, 0],
    [2, 0, 1],
    [4, 2, 2],
  ])
)

console.log(
  getMaximumNumber([
    [2, 0, 2],
    [6, 2, 0],
    [4, 1, 0],
    [2, 2, 2],
    [3, 0, 2],
  ])
)
export {}
