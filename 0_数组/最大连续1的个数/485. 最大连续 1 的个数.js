/**
 * @param {number[]} nums
 * @return {number}
 * 给定一个二进制数组， 计算其中最大连续 1 的个数。
 */
var findMaxConsecutiveOnes = function (nums) {
  return nums
    .join('')
    .split(0)
    .reduce((pre, cur) => Math.max(pre, cur.length), 0)
}

console.log(findMaxConsecutiveOnes([1, 1, 0, 1, 1, 1]))
