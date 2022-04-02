// 你需要返回一个长度恰好为 N 的序列，第一个元素为长度为 1 的子数组的最小值，

import { getRange } from './4_每个元素作为最值的影响范围'

// 第二个元素为长度为 2 的子数组的最小值
function getMinimums(nums: number[]): number[] {
  const n = nums.length
  const res: number[] = Array(n).fill(Infinity)
  const ranges = getRange(nums, false, true)
  for (const [index, [left, right]] of ranges.entries()) {
    const length = right - left + 1
    res[length - 1] = Math.min(res[length - 1], nums[index])
  }

  for (let i = n - 2; ~i; i--) {
    res[i] = Math.min(res[i], res[i + 1])
  }

  return res
}

console.log(getMinimums([1, 3, 5, 2, 4, 6]))
