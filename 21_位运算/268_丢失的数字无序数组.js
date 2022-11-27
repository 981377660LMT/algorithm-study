/**
 * @param {number[]} nums
 * @return {number}
 题中数组为 0~n 中缺失一个数字。那么加上数组的下标。用异或刚好可以得出缺失的数字。
 如： [3, 0, 1]
 可得：3^0^1^0^1^2^3 = 2
 */
const missingNumber = nums => nums.reduce((pre, cur, index) => pre ^ cur ^ index, nums.length)

console.log(missingNumber([3, 0, 1]))
