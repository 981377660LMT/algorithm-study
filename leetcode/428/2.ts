export {}

// floyd 好题
function maxAmount(
  initialCurrency: string,
  pairs1: string[][],
  rates1: number[],
  pairs2: string[][],
  rates2: number[]
): number {
  const allVertex = new Set<string>()
  allVertex.add(initialCurrency)
  for (const [s, t] of pairs1) {
    allVertex.add(s)
    allVertex.add(t)
  }
  for (const [s, t] of pairs2) {
    allVertex.add(s)
    allVertex.add(t)
  }

  const keys = Array.from(allVertex)
  const mp = new Map<string, number>()
  for (let i = 0; i < keys.length; i++) mp.set(keys[i], i)
  const n = keys.length

  const day1 = Array.from({ length: n }, () => Array(n).fill(-Infinity))
  const day2 = Array.from({ length: n }, () => Array(n).fill(-Infinity))
  for (let i = 0; i < n; i++) {
    day1[i][i] = 1
    day2[i][i] = 1
  }

  for (let i = 0; i < pairs1.length; i++) {
    const [s, t] = pairs1[i]
    const r = rates1[i]
    const si = mp.get(s)!
    const ti = mp.get(t)!
    day1[si][ti] = Math.max(day1[si][ti], r)
    day1[ti][si] = Math.max(day1[ti][si], 1 / r)
  }

  for (let i = 0; i < pairs2.length; i++) {
    const [s, t] = pairs2[i]
    const r = rates2[i]
    const si = mp.get(s)!
    const ti = mp.get(t)!
    day2[si][ti] = Math.max(day2[si][ti], r)
    day2[ti][si] = Math.max(day2[ti][si], 1 / r)
  }

  for (let k = 0; k < n; k++) {
    for (let i = 0; i < n; i++) {
      for (let j = 0; j < n; j++) {
        if (day1[i][k] > 0 && day1[k][j] > 0) {
          const tmp = day1[i][k] * day1[k][j]
          if (tmp > day1[i][j]) day1[i][j] = tmp
        }
        if (day2[i][k] > 0 && day2[k][j] > 0) {
          const tmp = day2[i][k] * day2[k][j]
          if (tmp > day2[i][j]) day2[i][j] = tmp
        }
      }
    }
  }

  const cur = mp.get(initialCurrency)!
  let res = 1
  for (let i = 0; i < n; i++) {
    if (day1[cur][i] > 0 && day2[i][cur] > 0) {
      const tmp = day1[cur][i] * day2[i][cur]
      if (tmp > res) res = tmp
    }
  }
  return res
}
