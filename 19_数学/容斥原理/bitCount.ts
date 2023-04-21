/* eslint-disable no-inner-declarations */
/* eslint-disable no-console */
/* eslint-disable no-param-reassign */

import assert from 'assert'

/**
 * @param num num <= 2**53 - 1
 * @returns num 的二进制表示中1的个数
 */
function bitCount53(num: number): number {
  return bitCount32(num >>> 0) + bitCount32(~~(num / 0x100000000))
}

// popCount32
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
  return bitLength32(num >>> 0) + bitLength32(~~(num / 0x100000000))
}

function bitLength32(uint32: number): number {
  return 32 - Math.clz32(uint32)
}

/**
 * @param num num <= 2**53 - 1
 * @returns num 的二进制表示中尾随零的个数
 *
 * !注意：如果 num === 0，返回 53
 */
function trailingZero53(num: number): number {
  if (num === 0) return 53
  const low32 = trailingZero32(num >>> 0)
  return low32 !== 1 << 5 ? low32 : (1 << 5) + trailingZero32(~~(num / 0x100000000))
}

//
/**
 * @returns `32位无符号整形数字` 的二进制表示中尾随零的个数
 *
 * !注意：如果 num === 0，返回 32
 * lowbit + bitLength
 */
function trailingZero32(uint32: number): number {
  if (uint32 === 0) return 32
  return 31 - Math.clz32(uint32 & -uint32)
}

// uin32 > 0
function floorLog2(uint32: number): number {
  return 31 - Math.clz32(uint32)
}
function ceilLog2(uint32: number): number {
  return floorLog2(uint32) + ((uint32 & -uint32) !== uint32 ? 1 : 0)
}

/**
 * 第k(k>=0)个1的位置,从低位向高位寻找.
 * 如果不存在，返回-1.
 * @link http://graphics.stanford.edu/~seander/bithacks.html#SelectPosFromMSBRank
 *       https://github.com/pranjalssh/CP_codes/blob/984d3a7155e9f474855114d8aa96936458eb2397/anta/!DepthDescendantsQuery.cpp#L12
 * @alias rank32
 */
function select32(uint32: number, k: number): number {
  if (!uint32 || bitCount32(uint32) <= k) return -1
  const a = (uint32 & 0x55555555) + ((uint32 >>> 1) & 0x55555555)
  const b = (a & 0x33333333) + ((a >>> 2) & 0x33333333)
  const c = (b & 0x0f0f0f0f) + ((b >>> 4) & 0x0f0f0f0f)
  let t = (c & 0xff) + ((c >>> 8) & 0xff)
  let s = 0
  s += ((t - k - 1) & 128) >>> 3
  k -= t & ((t - k - 1) >>> 8) // if(k >= t) s += 16, k -= t;
  t = (c >>> s) & 0xf
  s += ((t - k - 1) & 128) >>> 4
  k -= t & ((t - k - 1) >>> 8) // if(k >= t) s += 8, k -= t;
  t = (b >>> s) & 0x7
  s += ((t - k - 1) & 128) >>> 5
  k -= t & ((t - k - 1) >>> 8) // if(k >= t) s += 4, k -= t;
  t = (a >>> s) & 0x3
  s += ((t - k - 1) & 128) >>> 6
  k -= t & ((t - k - 1) >>> 8) // if(k >= t) s += 2, k -= t;
  t = (uint32 >>> s) & 0x1
  s += ((t - k - 1) & 128) >>> 7 // if(k >= t) s += 1;
  return s
}

export { bitCount32, bitCount53, bitLength32, bitLength53, trailingZero32, trailingZero53 }

// 结论:
// !当n为32位整数时
// !32 - Math.clz32(n) <=> n 的二进制表示的长度
// !                   <=> Math.ceil(Math.log2(n + 1)) <=> 1 + Math.floor(Math.log2(n))
// bit_len(n-1) <=> log(n)
// bit_len(n) <=> 1+log(n)

if (require.main === module) {
  assert.strictEqual(bitCount53(0), 0)
  assert.strictEqual(bitCount53(2 ** 53 - 1), 53)

  assert.strictEqual(bitLength53(0), 0)
  assert.strictEqual(bitLength53(2 ** 53 - 1), 53)

  assert.strictEqual(trailingZero53(0), 53)
  assert.strictEqual(trailingZero53(2 ** 53 - 1), 0)
  console.log(2 ** 52, (2 ** 52).toString(2))
  assert.strictEqual(trailingZero53(2 ** 52), 52)

  console.time('bitCount53')
  for (let i = 0; i < 1e8; i++) {
    bitCount53(i)
  }
  console.timeEnd('bitCount53') // 260.361ms

  console.time('bitCount32')
  for (let i = 0; i < 1e8; i++) {
    bitCount32(i)
  }
  console.timeEnd('bitCount32') // 53.436ms

  console.time('bitLength53')
  for (let i = 0; i < 1e8; i++) {
    bitLength53(i)
  }
  console.timeEnd('bitLength53') // 131.062ms

  console.time('bitLength32')
  for (let i = 0; i < 1e8; i++) {
    bitLength32(i)
  }
  console.timeEnd('bitLength32') // 55.096ms

  console.time('trailingZero53')
  for (let i = 0; i < 1e8; i++) {
    trailingZero53(i)
  }
  console.timeEnd('trailingZero53') // 234.997ms

  console.time('trailingZero32')
  for (let i = 0; i < 1e8; i++) {
    trailingZero32(i)
  }
  console.timeEnd('trailingZero32') // 129.897ms

  console.time('floorLog2/ceilLog2')
  for (let i = 1; i < 1e6; i++) {
    assert.strictEqual(floorLog2(i), Math.floor(Math.log2(i)))
    assert.strictEqual(ceilLog2(i), Math.ceil(Math.log2(i)))
  }
  console.timeEnd('floorLog2/ceilLog2')

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
  assert.strictEqual(_bitCount53ByDataView(12345), bitCount53(12345))

  function uint32(num: number): number {
    return num >>> 0
    // return num & 0xffffffff
  }
}
