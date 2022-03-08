/**
 * @description
 * 单调不减的数组nums1和nums2分别找到两个数，其和与target的差最小 返回这个最小差值
 */
function twoSum(nums1: number[], nums2: number[], target: number) {
  let l = 0
  let r = nums2.length - 1
  let res = Infinity

  while (l < nums1.length && r > -1) {
    const sum = nums1[l] + nums2[r]
    res = Math.min(res, Math.abs(target - sum))
    if (sum === target) return 0
    else if (sum > target) r--
    else l++
  }

  return res
}

export { twoSum }
