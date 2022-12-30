const M1 = 0x55555555 // 01010101010101010101010101010101
const M2 = 0x33333333 // 00110011001100110011001100110011
const M4 = 0x0f0f0f0f // 00001111000011110000111100001111
const M8 = 0x00ff00ff // 00000000111111110000000011111111

/**
 * 颠倒给定的 32 位无符号整数的二进制位。
 */
function reverseBits(n: number): number {
  n = ((n >>> 1) & M1) | ((n & M1) << 1)
  n = ((n >>> 2) & M2) | ((n & M2) << 2)
  n = ((n >>> 4) & M4) | ((n & M4) << 4)
  n = ((n >>> 8) & M8) | ((n & M8) << 8)
  return ((n >>> 16) | (n << 16)) >>> 0
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
