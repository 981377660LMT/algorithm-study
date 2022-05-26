import { findLastLarge } from '../对每个数，寻找右侧最后一个比自己大的数'

// 最大宽度坡
function maxWidthRamp(nums: number[]): number {
  const lastLarge = findLastLarge(nums)
  let res = 0
  for (let i = 0; i < lastLarge.length; i++) {
    if (lastLarge[i] === -1) continue
    res = Math.max(res, lastLarge[i] - i)
  }

  return res
}

console.log(maxWidthRamp([9, 8, 1, 0, 1, 9, 4, 0, 4, 1]))
