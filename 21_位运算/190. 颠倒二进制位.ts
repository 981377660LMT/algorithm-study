/**
 * @param {number} n - a positive integer
 * @return {number} - a positive integer
 * 任何数字和 1 进行位运算的结果都取决于该数字最后一位
 * 颠倒给定的 32 位无符号整数的二进制位。
 */
const reverseBits = function (n: number): number {
  let res = 0

  // 最后一位肯定不用移 所以先移再加
  for (let i = 0; i < 32; i++) {
    res <<= 1
    res |= n & 1
    n >>= 1
  }

  // res转化为无符号,否则最高位会被解析为符号。
  return res >>> 0
}

console.log(reverseBits(0b00000010100101000001111010011100))
// console.log(0b0101)

// js 中的无符号右移 0 位。

// -1>>>0
// >>>0 实际上并没有发生数位变化，但是 js 却会把符号位替换成 0，
// java
// 11111111111111111111111111111111
// js
// 01111111111111111111111111111111
