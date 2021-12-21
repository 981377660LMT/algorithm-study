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
  // 抹去最右边的 1
  while (left < right) {
    right = right & (right - 1)
  }

  return right
}

console.log(rangeBitwiseAnd(5, 7))

export {}
