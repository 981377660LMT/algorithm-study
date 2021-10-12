/**
 * @param {number[]} nums - ascending array with duplicates
 * @param {number} target
 * @return {number}
 */
function firstIndex(nums: number[], target: number): number {
  let res = -1
  let l = 0
  let r = nums.length - 1

  while (l <= r) {
    const mid = (l + r) >> 1
    const midElement = nums[mid]
    if (midElement === target) {
      res = mid
      l++
    } else if (midElement < target) l = mid + 1
    else r = mid - 1
  }

  return res
}

export {}
