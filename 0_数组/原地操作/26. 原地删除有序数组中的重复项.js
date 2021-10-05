/**
 * @param {number[]} nums
 * @return {number}
 * 直接比较相邻的两项 即可知道该数是否重复出现。
 */
var removeDuplicates = function (nums) {
  let slow = 0
  for (let fast = 1; fast < nums.length; fast++) {
    if (nums[slow] !== nums[fast]) {
      slow++
      nums[slow] = nums[fast]
    }
  }

  // 前面已经原地删除了 这里看一下结果
  return nums.slice(0, slow + 1)
}

console.log(removeDuplicates([1, 1, 2]))
