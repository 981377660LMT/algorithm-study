/**
 * @param {number[]} nums   在一个长度为 n 的数组 nums 里的所有数字都在 0～n-1 的范围内
 * @return {number}
 * 时间复杂度O(n)，空间复杂度O(1)
 * @summary 让 位置i 的地方放元素i (鸽巢原理)
 * 请找出数组中任意一个重复的数字
 */
const findRepeatNumber = function (nums: number[]): number {
  const n = nums.length

  for (let i = 0; i < n; i++) {
    while (nums[i] !== i) {
      if (nums[i] === nums[nums[i]]) return nums[i]
      swap(nums[i], i)
    }
  }

  return -1

  function swap(i: number, j: number) {
    ;[nums[i], nums[j]] = [nums[j], nums[i]]
  }
}

console.log(findRepeatNumber([2, 3, 1, 0, 2, 5, 3]))
