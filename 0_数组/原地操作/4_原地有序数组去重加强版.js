/**
 * @param {number[]} nums
 * @return {number}
 * @description 请你 原地 删除重复出现的元素，使每个元素 最多出现两次 ，返回删除后数组的新长度。
 */
let removeDuplicates = function (nums) {
  // 出现了几个不同的数
  let slow = 0
  let dupCount = 0

  for (let i = 1; i < nums.length; i++) {
    if (nums[i] === nums[slow] && dupCount === 1) {
      continue
    } else if (nums[i] === nums[slow] && dupCount === 0) {
      // 先移后写
      slow++
      nums[slow] = nums[i]
      dupCount = 1
    } else {
      slow++
      nums[slow] = nums[i]
      dupCount = 0
    }
  }

  console.log(nums)
  return slow + 1
}

console.log(removeDuplicates([1, 1, 1, 2, 2, 3]))
console.log(removeDuplicates([1, 1, 1, 1, 1, 1]))

// 返回5，[1,1,2,2,3]
