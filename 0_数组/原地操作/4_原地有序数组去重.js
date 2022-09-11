/* eslint-disable func-names */
/**
 * @param {number[]} nums
 * @return {number}
 * @summary 快指针是读指针， 慢指针是写指针
 */

function removeDuplicates(nums) {
  // 双指针 没见过的就搬过来
  let slow = 0
  for (let i = 0; i < nums.length; i++) {
    if (nums[i] !== nums[slow]) {
      // 先移后写
      slow++
      nums[slow] = nums[i]
    }
  }

  // 原地移除 类似于链表的slow.next = null
  nums.length = slow + 1
  return slow + 1
}

const a = [0, 0, 1, 1, 1, 2, 2, 3, 3, 4]
console.log(removeDuplicates(a))
console.log(a)
