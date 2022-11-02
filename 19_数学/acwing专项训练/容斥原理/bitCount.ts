/* eslint-disable no-param-reassign */

import assert from 'assert'

/**
 * @param num num <= 2**53 - 1
 * @returns num 的二进制表示中1的个数
 */
function bitCount53(num: number): number {
  return bitCount32(num) + bitCount32(Math.floor(num / 0x100000000))
}

function bitCount32(uint32: number): number {
  uint32 -= (uint32 >>> 1) & 0x55555555
  uint32 = (uint32 & 0x33333333) + ((uint32 >>> 2) & 0x33333333)
  return (((uint32 + (uint32 >>> 4)) & 0x0f0f0f0f) * 0x01010101) >>> 24
}

/**
 * @param num num <= 2**53 - 1
 * @returns num 的二进制表示的长度
 */
function bitLength53(num: number): number {
  return bitLength32(num) + bitLength32(Math.floor(num / 0x100000000))
}

function bitLength32(uint32: number): number {
  return 32 - Math.clz32(uint32)
}

if (require.main === module) {
  assert.strictEqual(bitCount53(0), 0)
  assert.strictEqual(bitCount53(2 ** 53 - 1), 53)

  assert.strictEqual(bitLength53(0), 0)
  assert.strictEqual(bitLength53(2 ** 53 - 1), 53)

  console.time('bitCount53')
  for (let i = 0; i < 1e8; i++) {
    bitCount53(2 ** 53 - 1)
  }
  console.timeEnd('bitCount53') // bitCount53: 53.431ms

  console.time('bitLength53')
  for (let i = 0; i < 1e8; i++) {
    bitLength53(2 ** 53 - 1)
  }
  console.timeEnd('bitLength53') // bitLength53: 54.321ms
}

export { bitCount32, bitCount53, bitLength32, bitLength53 }

// 结论:
// !当n为32位整数时
// !32 - Math.clz32(n) == n 的二进制表示的长度 == log2(n) + 1

// 使用 DataView 读取 2**53 - 1 的比特位
const float64View = new DataView(new ArrayBuffer(8))
function _bitCount53ByDataView(num: number): number {
  if (num === 0) return 0
  float64View.setFloat64(0, num, false) // 大端序存储
  const low32 = float64View.getUint32(4, false)
  const low32BitCount = bitCount32(low32)
  const high32 = float64View.getUint32(0, false) & 0x000fffff // sign 1 + exponent 11 = 12 移除前12位
  const high32BitCount = bitCount32(high32) + 1 // 1 表示科学计数法整数部分的1 (num不为0时)
  return high32BitCount + low32BitCount
}
