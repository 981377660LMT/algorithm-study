/**
 * @param {number[]} nums
 * @return {number}
 * 数组 中心下标 是数组的一个下标，其左侧所有元素相加的和等于右侧所有元素相加的和。
 * 如果数组有多个中心下标，应该返回 最靠近左边 的那一个。如果数组不存在中心下标，返回 -1 。
 */
var pivotIndex = function (nums) {
  const sum = nums.reduce((r, n) => r + n, 0)
  let left = 0
  for (let i = 0; i < nums.length; i++) {
    const right = sum - left - nums[i]
    if (left === right) return i
    left += nums[i]
  }
  return -1
}

console.log(pivotIndex([1, 7, 3, 6, 5, 6]))
