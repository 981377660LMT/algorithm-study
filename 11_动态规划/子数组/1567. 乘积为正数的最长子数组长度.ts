/**
 * @param {number[]} nums
 * @return {number}
 * 请你求出乘积为正数的最长子数组的长度。
 * @description 正负一般使用二维的dp
 *
 */
function getMaxLen(nums: number[]): number {
  const n = nums.length
  const positive = Array<number>(n).fill(0)
  const negative = Array<number>(n).fill(0)
  let max = -Infinity

  if (nums[0] < 0) {
    negative[0] = 1
  } else if (nums[0] > 0) {
    positive[0] = 1
  }

  for (let i = 1; i < n; i++) {
    const cur = nums[i]
    if (cur > 0) {
      positive[i] = positive[i - 1] + 1
      // 需要判断之前负数组是否存在 nums[i] 自己没法成为一个乘积为负的数组。
      if (negative[i - 1] !== 0) negative[i] = negative[i - 1] + 1
    } else if (cur < 0) {
      negative[i] = positive[i - 1] + 1
      if (negative[i - 1] !== 0) positive[i] = negative[i - 1] + 1
    }
  }
  console.log(max)
  return Math.max.apply(null, positive)
}

console.log(getMaxLen([1, -2, -3, 4]))
console.log(getMaxLen([2, 3, -2, 4]))
export {}
