function findMinimumTime(strength: number[], K: number): number {
  const n = strength.length
  const mask = (1 << n) - 1
  const dp = Array(1 << n).fill(Infinity)
  dp[0] = 0
  for (let m = 0; m <= mask; m++) {
    const ones = onesCount32(m)
    const power = 1 + ones * K
    for (let i = 0; i < n; i++) {
      if ((m & (1 << i)) === 0) {
        const newMask = m | (1 << i)
        const need = strength[i]
        const cost = Math.floor((need + power - 1) / power)
        dp[newMask] = Math.min(dp[newMask], dp[m] + cost)
      }
    }
  }
  return dp[mask]
}

function onesCount32(uint32: number): number {
  uint32 -= (uint32 >>> 1) & 0x55555555
  uint32 = (uint32 & 0x33333333) + ((uint32 >>> 2) & 0x33333333)
  return (((uint32 + (uint32 >>> 4)) & 0x0f0f0f0f) * 0x01010101) >>> 24
}
