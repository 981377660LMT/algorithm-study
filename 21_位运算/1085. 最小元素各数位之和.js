/**
 * @param {number[]} nums
 * @return {number}
 * 然后计算 S，使其等于数组 A 当中最小的那个元素各个数位上数字之和。
 * 假如 S 所得计算结果是 奇数 ，返回 0 ；否则请返回 1。
 */
var sumOfDigits = function (nums) {
  return (
    (Math.min(...nums)
      .toString()
      .split('')
      .map(Number)
      .reduce((pre, cur) => pre + cur, 0) &
      1) ^
    1
  )
}
