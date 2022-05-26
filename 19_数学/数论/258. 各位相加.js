/**
 * @param {number} num
 * @return {number}
 * @summary 一个数能被9整除，则这个数的各个位数相加也能被9整除
 */
var addDigits = function (num) {
  if (num < 10) return num
  return num % 9 === 0 ? 9 : num % 9
}

console.log(addDigits(38))
// 解释: 各位相加的过程为：3 + 8 = 11, 1 + 1 = 2。 由于 2 是一位数，所以返回 2。
