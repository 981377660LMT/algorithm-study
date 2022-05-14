/**
 * @param {number[]} nums
 * @param {number} k
 * @return {boolean}
 * @description 断数组中是否存在两个不同的索引 i 和 j，使得 nums [i] = nums [j]，并且 i 和 j 的差的 绝对值 至多为 k。
 * @summary 遍历 记录最后的位置 即可
 */
const containsNearbyDuplicate = function (nums, k) {
  const map = new Map()
  for (let i = 0; i < nums.length; i++) {
    if (i - map.get(nums[i]) <= k) return true
    map.set(nums[i], i)
  }
  return false
}

export {}
