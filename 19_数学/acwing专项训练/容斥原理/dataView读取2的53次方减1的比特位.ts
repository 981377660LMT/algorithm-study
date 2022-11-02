/* eslint-disable no-param-reassign */
// https://github.com/981377660LMT/ts/issues/135
// Buffer ArrayBuffer DataView 的区别/使用 以及浮点数float32与float64的存储

const float64View = new DataView(new ArrayBuffer(8))

function _bitCount53ByDataView(num: number): number {
  if (num === 0) return 0
  float64View.setFloat64(0, num, false) // 大端序存储
  const low32 = float64View.getUint32(4, false)
  const low32BitCount = _bitCount32(low32)
  const high32 = float64View.getUint32(0, false) & 0x000fffff // sign 1 + exponent 11 = 12 移除前12位
  const high32BitCount = _bitCount32(high32) + 1 // 1 表示科学计数法整数部分的1 (num不为0时)
  return high32BitCount + low32BitCount
}

function _bitCount32(uint32: number): number {
  uint32 -= (uint32 >>> 1) & 0x55555555
  uint32 = (uint32 & 0x33333333) + ((uint32 >>> 2) & 0x33333333)
  return (((uint32 + (uint32 >>> 4)) & 0x0f0f0f0f) * 0x01010101) >>> 24
}

if (require.main === module) {
  console.log(_bitCount53ByDataView(2 ** 53 - 1))
  console.log(_bitCount53ByDataView(0))
}

export {}
