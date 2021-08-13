/**
 * @param {number[]} nums
 * @param {number} target
 * @return {boolean}
 * 按照升序排序的数组在预先未知的某个点上进行了旋转。
 * 编写一个函数来判断给定的目标值是否存在于数组中。若存在返回 true，否则返回 false。
 * 这是 搜索旋转排序数组 的延伸题目，本题中的 nums  可能包含重复元素。
 * @summary 如果存在重复数字，就可能会发生 nums[mid] == nums[l] 了，比如 30333 。
 */
const search = function (nums: number[], target: number): boolean {
  let l = 0
  let r = nums.length - 1

  while (l <= r) {
    const mid = (l + r) >> 1
    if (nums[mid] === target) return true

    // 注意多的这句
    if (nums[l] === nums[mid]) {
      l++
      continue
    }

    // 左半部分有序
    if (nums[l] < nums[mid]) {
      if (nums[mid] > target && target >= nums[l]) {
        r = mid - 1
      } else {
        // target不在左半部分
        l = mid + 1
      }
    } else {
      // 右半部分有序
      if (nums[r] >= target && target > nums[mid]) {
        l = mid + 1
      } else {
        // target不在右半部分
        r = mid - 1
      }
    }
  }

  return false
}

console.log(search([1, 0, 1, 1, 1], 0))

export {}
