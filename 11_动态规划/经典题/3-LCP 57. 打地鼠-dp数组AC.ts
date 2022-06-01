const POS: number[][] = []
for (const r of [0, 1, 2]) {
  for (const c of [0, 1, 2]) {
    POS.push([r, c])
  }
}

// 注意js的cache很容易TLE
// 要优化就要把dfs优化成dp 这样常数会小一些
function getMaximumNumber(moles: number[][]): number {
  const times = new Set<number>([0])
  const indexMap = new Map<number, Set<number>>()
  for (const [t, r, c] of moles) {
    times.add(t)
    !indexMap.has(t) && indexMap.set(t, new Set())
    indexMap.get(t)!.add(r * 3 + c)
  }

  const allTimes = [...times].sort((a, b) => a - b)
  const n = allTimes.length

  let dp = Array.from({ length: 3 }, () => new Int32Array(3).fill(-(1 << 31)))

  // i=0 情况
  dp[1][1] = Number(indexMap.has(allTimes[0]) && indexMap.get(allTimes[0])!.has(4))

  for (let i = 1; i < n; i++) {
    const ndp = Array.from({ length: 3 }, () => new Int32Array(3).fill(-(1 << 31)))
    const curTime = allTimes[i]
    const preTime = allTimes[i - 1]
    const diff = curTime - preTime
    for (const [r, c] of POS) {
      for (const [preR, preC] of POS) {
        const score = Number(indexMap.has(curTime) && indexMap.get(curTime)!.has(r * 3 + c))
        if (Math.abs(r - preR) + Math.abs(c - preC) <= diff) {
          ndp[r][c] = Math.max(ndp[r][c], dp[preR][preC] + score)
        }
      }
    }

    dp = ndp
  }

  return Math.max(...dp.map(row => Math.max(...row)))
}

// 1
console.log(
  getMaximumNumber([
    [0, 0, 0],
    [1, 1, 0],
    [0, 2, 0],
    [1, 0, 1],
    [1, 2, 1],
  ])
)

// 0
console.log(getMaximumNumber([[1, 0, 0]]))

export {}
