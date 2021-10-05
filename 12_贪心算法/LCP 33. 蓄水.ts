/**
 * @param {number[]} bucket
 * @param {number[]} vat
 * @return {number}
 * 枚举蓄水的次数
 */
const storeWater = function (bucket: number[], vat: number[]): number {
  const n = bucket.length
  const maxVat = Math.max(...vat)
  if (maxVat === 0) return 0
  let res = maxVat + n

  // 枚举蓄水次数
  for (let i = 0; i < maxVat + n + 1; i++) {
    let cur = i
    for (let j = 0; j < n; j++) {
      if (i * bucket[j] >= vat[j]) continue
      else cur += Math.ceil(vat[j] / i) - bucket[j] // 需要的升级次数
    }
    res = Math.min(res, cur)
  }

  return res
}

console.log(storeWater([1, 3], [6, 8]))

// 实际蓄水量 达到或超过 最低蓄水量，即完成蓄水要求
// 每个水缸对应最低蓄水量记作 vat[i]，返回小扣至少需要多少次操作可以完成所有水缸蓄水要求
