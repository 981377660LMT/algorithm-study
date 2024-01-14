/**
 * 单调不减的数组 nums1 和 nums2 分别找到两个数，其和与target的差最小, 返回这个最小差值.
 */
function twoSum(nums1: number[], nums2: number[], target: number) {
  let left = 0
  let right = nums2.length - 1
  let res = Infinity

  while (left < nums1.length && right >= 0) {
    const sum = nums1[left] + nums2[right]
    res = Math.min(res, Math.abs(target - sum))
    if (sum === target) return 0
    if (sum > target) {
      right--
    } else {
      left++
    }
  }

  return res
}

export { twoSum }
