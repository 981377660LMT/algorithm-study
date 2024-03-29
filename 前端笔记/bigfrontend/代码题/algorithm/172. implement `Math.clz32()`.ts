// https://bigfrontend.dev/zh/problem/clz32/discuss
// 实现 Math.clz32() 函数，返回一个 32 位无符号整数的前导 0 的个数

const len8tab =
  '' +
  '\x00\x01\x02\x02\x03\x03\x03\x03\x04\x04\x04\x04\x04\x04\x04\x04' +
  '\x05\x05\x05\x05\x05\x05\x05\x05\x05\x05\x05\x05\x05\x05\x05\x05' +
  '\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06' +
  '\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06\x06' +
  '\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07' +
  '\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07' +
  '\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07' +
  '\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07\x07' +
  '\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08' +
  '\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08' +
  '\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08' +
  '\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08' +
  '\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08' +
  '\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08' +
  '\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08' +
  '\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08'

/**
 * Len32 returns the minimum number of bits required to represent x;
 * the result is 0 for x == 0.
 */
function len32(x: number): number {
  x >>>= 0 // change to uint32
  let res = 0
  if (x >= 1 << 16) {
    x >>>= 16
    res = 16
  }
  if (x >= 1 << 8) {
    x >>>= 8
    res += 8
  }
  return res + len8tab.charCodeAt(x)
}

function clz32(num: number): number {
  return 32 - len32(num)
}
