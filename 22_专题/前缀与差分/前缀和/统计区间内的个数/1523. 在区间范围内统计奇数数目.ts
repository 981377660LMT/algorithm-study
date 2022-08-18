// 请你返回 [low , high] 之间奇数的数目。
// 0 <= low <= high <= 10^9

function cal(upper: number): number {
  if (upper <= 0) return 0
  return (upper + 1) >> 1
}

function countOdds(low: number, high: number): number {
  return cal(high) - cal(low - 1)
}

export {}
