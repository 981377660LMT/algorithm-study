function makeStringGood(s: string): number {
  const n = s.length
  const aCode = 'a'.charCodeAt(0)

  const count = new Array(26).fill(0)
  const chars: number[] = []
  for (let i = 0; i < n; i++) {
    const c = s.charCodeAt(i) - aCode
    count[c]++
    chars.push(c)
  }

  const transformCosts: number[][] = Array.from({ length: 26 }, () => [])
  for (let l = 0; l < 26; l++) {
    const costs = new Array(n)
    for (let i = 0; i < n; i++) {
      const c = chars[i]
      if (l >= c) {
        costs[i] = l - c
      } else {
        costs[i] = 2
      }
    }
    costs.sort((a, b) => a - b)
    transformCosts[l] = costs
  }

  function costToProduce(l: number, f: number, x: number): number {
    let transformSum = 0
    for (let i = 0; i < x; i++) {
      transformSum += transformCosts[l][i]
    }
    const insertCost = Math.max(0, f - x)
    return transformSum + insertCost
  }

  function solveKF(k: number, f: number): number {
    if (k > 26) return Number.MAX_SAFE_INTEGER

    const letterOptions = new Array(26)
    for (let l = 0; l < 26; l++) {
      const opts = []
      for (let x = 0; x <= Math.min(f, n); x++) {
        const cst = costToProduce(l, f, x)
        opts.push({ x, val: cst - x })
      }
      letterOptions[l] = opts
    }

    const dp = Array.from({ length: k + 1 }, () => Array(n + 1).fill(Number.MAX_SAFE_INTEGER))
    dp[0][0] = 0

    for (let l = 0; l < 26; l++) {
      const opts = letterOptions[l]
      for (let chosen = k; chosen >= 0; chosen--) {
        for (let usedX = 0; usedX <= n; usedX++) {
          if (dp[chosen][usedX] === Number.MAX_SAFE_INTEGER) continue
          if (chosen < k) {
            for (let { x, val } of opts) {
              if (usedX + x <= n) {
                const newVal = dp[chosen][usedX] + val
                if (newVal < dp[chosen + 1][usedX + x]) {
                  dp[chosen + 1][usedX + x] = newVal
                }
              }
            }
          }
        }
      }
    }

    let ans = Number.MAX_SAFE_INTEGER
    for (let X = 0; X <= n; X++) {
      const val = dp[k][X]
      if (val < ans) {
        ans = val
      }
    }
    if (ans === Number.MAX_SAFE_INTEGER) return ans

    return n + ans
  }

  let ans = Number.MAX_SAFE_INTEGER

  const maxL = n + 26
  for (let L = 1; L <= maxL; L++) {
    for (let f = 1; f <= L; f++) {
      if (L % f !== 0) continue
      const k = L / f
      if (k > 26) continue
      const cost = solveKF(k, f)
      if (cost < ans) ans = cost
    }
  }

  return ans
}
