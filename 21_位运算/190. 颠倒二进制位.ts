/**
 * @param {number} n - a positive integer
 * @return {number} - a positive integer
 * 任何数字和 1 进行位运算的结果都取决于该数字最后一位
 */
const reverseBits = function (n: number): number {
  var result = 0
  var count = 32

  while (count--) {
    result *= 2
    result += n & 1
    n = n >> 1
  }
  return result
}

console.log(reverseBits(0b00000010100101000001111010011100))
// console.log(0b0101)
