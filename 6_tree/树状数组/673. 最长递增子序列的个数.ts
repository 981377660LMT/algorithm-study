import { BIT } from './BIT'

function findNumberOfLIS(nums: number[]): number {
  const set = new Set(nums)

  const map = new Map<number, number>()
  for (const [key, realValue] of [...set].sort((a, b) => a - b).entries()) {
    map.set(realValue, key + 1)
  }

  let res = 0
  const bit = new BIT(map.size)
  // 树状数组维护 (len, cnt) 信息

  return res
}
