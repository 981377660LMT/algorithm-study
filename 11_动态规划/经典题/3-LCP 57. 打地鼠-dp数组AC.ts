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
  const record = new Map<number, Set<number>>()
  for (const [t, r, c] of moles) {
    times.add(t)
    !record.has(t) && record.set(t, new Set())
    record.get(t)!.add(r * 3 + c)
  }

  const allTimes = [...times].sort((a, b) => a - b)
  const n = allTimes.length

  // 也可以Int32Array
  const dp = Array.from({ length: n }, () => [
    [-Infinity, -Infinity, -Infinity],
    [-Infinity, -Infinity, -Infinity],
    [-Infinity, -Infinity, -Infinity],
  ])

  // i=0 情况
  dp[0][1][1] = Number(record.has(allTimes[0]) && record.get(allTimes[0])!.has(4))

  for (let i = 1; i < n; i++) {
    const curTime = allTimes[i]
    const preTime = allTimes[i - 1]
    const diff = curTime - preTime
    for (const [r, c] of POS) {
      for (const [preR, preC] of POS) {
        const cur = Number(record.has(curTime) && record.get(curTime)!.has(r * 3 + c))
        if (Math.abs(r - preR) + Math.abs(c - preC) <= diff) {
          dp[i][r][c] = Math.max(dp[i][r][c], dp[i - 1][preR][preC] + cur)
        }
      }
    }
  }

  return Math.max(...dp[n - 1].flat())
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
