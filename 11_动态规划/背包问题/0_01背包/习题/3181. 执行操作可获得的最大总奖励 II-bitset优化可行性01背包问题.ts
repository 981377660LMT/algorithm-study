// BigInt 大数模拟
const maxTotalReward = function (rewardValues: number[]): number {
  rewardValues = [...new Set(rewardValues)].sort((a, b) => a - b)
  const Big1 = BigInt(1)
  let res = Big1
  rewardValues.forEach(v => {
    const BigV = BigInt(v)
    const low = ((Big1 << BigV) - Big1) & res
    res |= low << BigV
  })
  return res.toString(2).length - 1
}

export {}
