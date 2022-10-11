/**
 * @param {number[]} nums  nums 存在重复数字
 * @param {number} target
 * @return {boolean}
 * 按照升序排序的数组在预先未知的某个点上进行了旋转。
 * 编写一个函数来判断给定的目标值是否存在于数组中。若存在返回 true，否则返回 false。
 * !这是 搜索旋转排序数组 的延伸题目，本题中的 nums  可能包含重复元素。
 * @summary 如果存在重复数字，就可能会发生 nums[mid] == nums[l] 了，比如 30333 。
 */
function search(nums: number[], target: number): boolean {
  let left = 0
  let right = nums.length - 1

  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (nums[mid] === target) {
      return true
    }

    // 注意多的这句
    if (nums[left] === nums[mid]) {
      left++
      continue
    }

    // 左半部分有序
    if (nums[left] <= nums[mid]) {
      if (nums[mid] >= target && target >= nums[left]) {
        right = mid - 1
      } else {
        // target不在左半部分
        left = mid + 1
      }
    } else if (nums[mid] <= nums[right]) {
      // 右半部分有序
      if (nums[right] >= target && target >= nums[mid]) {
        left = mid + 1
      } else {
        // target不在右半部分
        right = mid - 1
      }
    }
  }

  return false
}

console.log(search([1, 0, 1, 1, 1], 0))

export {}

// 第一类
// 1011110111 和 1110111101 这种。此种情况下 nums[start] == nums[mid]，分不清到底是前面有序还是后面有序，
// 此时 start++ 即可。相当于去掉一个重复的干扰项。
// 第二类
// 22 33 44 55 66 77 11 这种，也就是 nums[start] < nums[mid]。此例子中就是 2 < 5；
// 这种情况下，前半部分有序。因此如果 nums[start] <=target<nums[mid]，则在前半部分找，否则去后半部分找。
// 第三类
// 66 77 11 22 33 44 55 这种，也就是 nums[start] > nums[mid]。此例子中就是 6 > 2；
// 这种情况下，后半部分有序。因此如果 nums[mid] <target<=nums[end]。则在后半部分找，否则去前半部分找。
