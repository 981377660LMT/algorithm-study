/**
 * @param {number} left
 * @param {number} right
 * @return {number}
 * 只要有一个0，那么无论有多少个 1都是 0
 * 11000
 * ...
 * 11101
 * 向右移3位后相等
 */
const rangeBitwiseAnd = function (left: number, right: number): number {
  // 向右移动到最高低相同时为止
  let i = 0
  while (left !== right) {
    left = left >> 1
    right = right >> 1
    i++
  }

  return left << i
}

console.log(rangeBitwiseAnd(5, 7))

export {}
