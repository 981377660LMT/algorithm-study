/**
 * @param {number[]} nums
 * @return {number}
 * @description 请你 原地 删除重复出现的元素，使每个元素 最多出现两次 ，返回删除后数组的新长度。
 */
var removeDuplicates = function (nums) {
  // 出现了几个不同的数
  let count = 0
  let duplicateCount = 0

  for (let i = 1; i < nums.length; i++) {
    if (nums[i] === nums[count] && duplicateCount === 1) {
      continue
    } else if (nums[i] === nums[count] && duplicateCount === 0) {
      count++
      nums[count] = nums[i]
      duplicateCount = 1
    } else {
      count++
      nums[count] = nums[i]
      duplicateCount = 0
    }
  }

  console.log(nums)
  return count + 1
}

console.log(removeDuplicates([1, 1, 1, 2, 2, 3]))
console.log(removeDuplicates([1, 1, 1, 1, 1, 1]))

// 返回5，[1,1,2,2,3]
