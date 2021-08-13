/**
 * @param {number[]} nums
 * @param {number} k 给定一个数组，将数组中的元素向右移动 k 个位置，其中 k 是非负数。
 * @return {void} Do not return anything, modify nums in-place instead.
 * 空间复杂度为 O(1) 的 原地 算法解决
 */
const rotate = function (nums: number[], k: number): void {
  k = k % nums.length
  const reverse = (l: number, r: number) => {
    while (l < r) {
      ;[nums[l], nums[r]] = [nums[r], nums[l]]
      l++
      r--
    }
  }

  reverse(0, nums.length - 1)
  reverse(0, k - 1)
  reverse(k, nums.length - 1)

  console.log(nums)
}

console.log(rotate([1, 2, 3, 4, 5, 6, 7], 3))
// 输出: [5,6,7,1,2,3,4]
// 解释:
// 向右旋转 1 步: [7,1,2,3,4,5,6]
// 向右旋转 2 步: [6,7,1,2,3,4,5]
// 向右旋转 3 步: [5,6,7,1,2,3,4]
