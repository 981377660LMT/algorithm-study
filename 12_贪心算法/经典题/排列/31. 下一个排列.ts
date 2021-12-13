/**
 * @param {number[]} nums
 * @return {void} Do not return anything, modify nums in-place instead.
 * 算法需要将给定数字序列重新排列成字典序中下一个更大的排列。
 * 如果不存在下一个更大的排列，则将数字重新排列成最小的排列（即升序排列）。
 * @summary 倒序遍历  尽早交换  发现大的就叫唤
 */
const nextPermutation = function (nums: number[]) {
  const n = nums.length
  let mono = true

  loop: for (let left = n - 1; left > -1; left--) {
    for (let right = n - 1; right > left; right--) {
      // 找到了第一对后面大于前面
      if (nums[right] > nums[left]) {
        // 交换玩排序
        ;[nums[left], nums[right]] = [nums[right], nums[left]]
        reverseRange(nums, left + 1, n - 1)
        mono = false
        break loop
      }
    }
  }

  mono && nums.reverse()

  return nums

  function reverseRange(nums: number[], i: number, j: number) {
    while (i < j) {
      ;[nums[i], nums[j]] = [nums[j], nums[i]]
      i++
      j--
    }
  }
}

// console.log(nextPermutation([1, 2, 3]))
// // 输出：[1,3,2]
// console.log(nextPermutation([3, 2, 1]))
// // 输出：[1,2,3]
// console.log(nextPermutation([1, 2, 4, 3]))
// 1324
console.log(nextPermutation([2, 3, 1]))
export default 1

// 先找出最大的索引 k 满足 nums[k] < nums[k+1]，如果不存在，就翻转整个数组；
// 再找出另一个最大索引 l 满足 nums[l] > nums[k]；
// 交换 nums[l] 和 nums[k]；
// 最后翻转 nums[k+1:]。
