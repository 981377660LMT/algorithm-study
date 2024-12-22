function subsequencesWithMiddleMode(nums: number[]): number {
  const MOD = 1_000_000_007
  const n = nums.length
  if (n < 5) return 0

  const uniqueNums = Array.from(new Set(nums))
  const numToIndex: { [key: number]: number } = {}
  uniqueNums.forEach((num, idx) => {
    numToIndex[num] = idx
  })
  const m = uniqueNums.length

  const fact = new Array(n + 1).fill(1)
  for (let i = 1; i <= n; i++) {
    fact[i] = (fact[i - 1] * i) % MOD
  }

  const invFact = new Array(n + 1).fill(1)
  invFact[n] = powMod(fact[n], MOD - 2, MOD)
  for (let i = n - 1; i >= 0; i--) {
    invFact[i] = (invFact[i + 1] * (i + 1)) % MOD
  }

  function C(a: number, b: number): number {
    if (a < b || b < 0) return 0
    return (((fact[a] * invFact[b]) % MOD) * invFact[a - b]) % MOD
  }

  function powMod(x: number, power: number, mod: number): number {
    let res = 1
    x = x % mod
    while (power > 0) {
      if (power % 2 === 1) {
        res = (res * x) % mod
      }
      x = (x * x) % mod
      power = Math.floor(power / 2)
    }
    return res
  }

  const freqLeft = Array.from({ length: m }, () => new Array(n).fill(0))
  const freqRight = Array.from({ length: m }, () => new Array(n).fill(0))

  const currentLeft = new Array(m).fill(0)
  for (let i = 0; i < n; i++) {
    const num = nums[i]
    const idx = numToIndex[num]
    currentLeft[idx] += 1
    for (let j = 0; j < m; j++) {
      freqLeft[j][i] = currentLeft[j]
    }
  }

  const currentRight = new Array(m).fill(0)
  for (let i = n - 1; i >= 0; i--) {
    const num = nums[i]
    const idx = numToIndex[num]
    currentRight[idx] += 1
    for (let j = 0; j < m; j++) {
      freqRight[j][i] = currentRight[j]
    }
  }

  const leftNonCount = new Array(n).fill(0)
  const rightNonCount = new Array(n).fill(0)
  const commonNonCount = new Array(n).fill(0)

  for (let i = 0; i < n; i++) {
    const num = nums[i]
    const numIdx = numToIndex[num]
    let count = 0
    for (let j = 0; j < m; j++) {
      if (j === numIdx) continue
      if (freqLeft[j][i] > 0) count++
    }
    leftNonCount[i] = count
  }

  for (let i = 0; i < n; i++) {
    const num = nums[i]
    const numIdx = numToIndex[num]
    let count = 0
    for (let j = 0; j < m; j++) {
      if (j === numIdx) continue
      if (freqRight[j][i] > 0) count++
    }
    rightNonCount[i] = count
  }

  for (let i = 0; i < n; i++) {
    const num = nums[i]
    const numIdx = numToIndex[num]
    let count = 0
    for (let j = 0; j < m; j++) {
      if (j === numIdx) continue
      if (freqLeft[j][i] > 0 && freqRight[j][i] > 0) count++
    }
    commonNonCount[i] = count
  }

  let result = 0

  for (let i = 0; i < n; i++) {
    const num = nums[i]
    const numIdx = numToIndex[num]
    const lCount = freqLeft[numIdx][i]
    const rCount = freqRight[numIdx][i]
    const lnNon = leftNonCount[i]
    const rnNon = rightNonCount[i]
    const common = commonNonCount[i]

    if (lCount >= 2 && rCount >= 2) {
      const ways = (C(lCount, 2) * C(rCount, 2)) % MOD
      result = (result + ways) % MOD
    }

    if (lCount >= 2 && rCount >= 1 && rnNon >= 1) {
      const ways = (((C(lCount, 2) * C(rCount, 1)) % MOD) * C(rnNon, 1)) % MOD
      result = (result + ways) % MOD
    }
    if (lCount >= 1 && rCount >= 2 && lnNon >= 1) {
      const ways = (((C(lCount, 1) * C(rCount, 2)) % MOD) * C(lnNon, 1)) % MOD
      result = (result + ways) % MOD
    }

    if (lCount >= 0 && rCount >= 2 && lnNon >= 2 && rnNon >= 0) {
      const ways =
        (((C(lCount, 0) * C(lnNon, 2)) % MOD) * ((C(rCount, 2) * C(rnNon, 0)) % MOD)) % MOD
      result = (result + ways) % MOD
    }
    if (lCount >= 1 && rCount >= 1 && lnNon >= 1 && rnNon >= 1) {
      const ways =
        (((C(lCount, 1) * C(lnNon, 1)) % MOD) * ((C(rCount, 1) * C(rnNon, 1)) % MOD)) % MOD
      result = (result + ways) % MOD
    }
    if (lCount >= 2 && rCount >= 0 && lnNon >= 0 && rnNon >= 2) {
      const ways =
        (((C(lCount, 2) * C(lnNon, 0)) % MOD) * ((C(rCount, 0) * C(rnNon, 2)) % MOD)) % MOD
      result = (result + ways) % MOD
    }

    if (lCount >= 0 && rCount >= 1 && lnNon >= 2 && rnNon >= 1 && common >= 0) {
      const term1 = (C(lnNon, 2) * C(rnNon, 1)) % MOD
      const term2 = (C(common, 1) * C(lnNon - 1, 1)) % MOD
      const ways_f2_a = (term1 - term2 + MOD) % MOD
      const ways1 = ways_f2_a
      const term3 = (C(lnNon, 1) * C(rnNon, 2)) % MOD
      const term4 = (C(common, 1) * C(rnNon - 1, 1)) % MOD
      const ways_f2_b = (term3 - term4 + MOD) % MOD
      const ways_f2 = (ways1 + ways_f2_b) % MOD
      result = (result + ways_f2) % MOD
    }
  }

  return result
}
