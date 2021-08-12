/**
 * @param {number[]} nums
 * @return {number}
 * @description 除了一个数字出现一次，其他都出现了两次，让我们找到出现一次的数
 * @summary 出现两次的数异或全为0 出现一次的数与0异或等于自身
 */
var singleNumber = function (nums) {
  return nums.reduce((pre, cur) => pre ^ cur)
}

console.log(singleNumber([1, 2, 3, 2, 1]))
