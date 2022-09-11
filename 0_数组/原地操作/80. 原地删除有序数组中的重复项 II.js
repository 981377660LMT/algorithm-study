/**
 * @param {number[]} nums
 * @return {number}
 */
let removeDuplicates = function (nums) {
  const K = 2
  let slow = 0
  let fast = 0

  // fast将值丢给前面的slow指针 slow指针检查是否有K个
  while (fast < nums.length) {
    if (nums[slow - K] !== nums[fast]) {
      nums[slow] = nums[fast]
      slow++
    }
    fast++
  }

  return slow
}
