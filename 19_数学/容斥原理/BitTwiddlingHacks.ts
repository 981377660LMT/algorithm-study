export {}

// https://zhuanlan.zhihu.com/p/37014715
// https://hackmd.io/@0xff07/BTS
// http://graphics.stanford.edu/~seander/bithacks.html#SelectPosFromMSBRank

// x的符号(0, 1, -1)
function sign(x: number): number {
  return +(x > 0) - +(x < 0)
}

// x和y是否异号
function notSameSignInt32(x: number, y: number): boolean {
  return (x ^ y) < 0
}

// 不用分支计算绝对值
function absInt32(int32: number): number {
  const mask = int32 >>> 31
  return (int32 + mask) ^ mask
}

// -1 = 0b11111...111
function maxInt32(x: number, y: number): number {
  return x ^ ((x ^ y) & -(x < y))
}

function minInt32(x: number, y: number): number {
  return y ^ ((x ^ y) & -(x < y))
}

function isPow2Int32(x: number): boolean {
  return !!x && !!(x & (x - 1))
}

// 没看懂
function signExtendInt32(x: number, n: number): number {
  return (x << (32 - n)) >> (32 - n)
}

// 不用分支，根据条件设置/清除比特位.
// 当 flag 为true时，按照 mask 设置 x 的比特位, 否则清除.
// return flag ? x|(1<<bit) : x&~(1<<bit)
function setBitInt32(x: number, bit: number, flag: boolean): number {
  return x ^ ((-flag ^ x) & (1 << bit))
}

// 不用分支，根据条件求相反数
// return flag ? -x : x
function negateInt32(x: number, flag: boolean): number {
  return (x ^ -flag) + +flag
}

// 根据掩码来合并两个数
// 掩码，如果对应位为 1，则取 a 的值，否则取 b 的值
function mergeInt32(a: number, b: number, mask: number): number {
  return a ^ ((a ^ b) & mask) // 比下面少一个操作符
  return (a & mask) | (b & ~mask)
}

// 计算指定位的阶(Rank), 即从最高有效位到指定位中 1 的个数.
// k从低位开始,0<=k<bitlen(x).
// rank(0b101, 0) = 2.
function rankInt32(x: number, k: number): number {
  x >>>= k // 将指定位置向右移至最低位
  x = (x & 0x55555555) + ((x >> 1) & 0x55555555)
  x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
  x = (x & 0x0f0f0f0f) + ((x >> 4) & 0x0f0f0f0f)
  return x & 0xff
}

// 使用异或交换两个数字
// (a) ^ (b) 表达式可以被复用，所以这样写可能会更快
function swapInt32(): void {
  let a = 2
  let b = 3
  console.log(a, b)
  // eslint-disable-next-line no-multi-assign
  a ^ b && ((b ^= a ^= b), (a ^= b))
  console.log(a, b)
}

// 交换指定位置与长度的比特序列.
// pos1, pos2 从低位开始, 1<=pos1, pos2<=bitlen(x).
// 交换pos1-pos1+len-1 和 pos2-pos2+len-1 之间的比特位.
// swapRangeInt32(0b0010111, 1, 5, 3) => 0b1110001
function swapRangeInt32(x: number, pos1: number, pos2: number, len: number): number {
  const tmp = ((x >>> pos1) ^ (x >>> pos2)) & ((1 << len) - 1)
  return x ^ ((tmp << pos1) | (tmp << pos2))
}

const M1 = 0x55555555 // 01010101010101010101010101010101
const M2 = 0x33333333 // 00110011001100110011001100110011
const M4 = 0x0f0f0f0f // 00001111000011110000111100001111
const M8 = 0x00ff00ff // 00000000111111110000000011111111

/**
 * 比特位翻转.
 * 使用 5 × lg(N) 个操作符反向 N 个比特.
 */
function reverseBitUint32(uint32: number): number {
  uint32 = ((uint32 >>> 1) & M1) | ((uint32 & M1) << 1)
  uint32 = ((uint32 >>> 2) & M2) | ((uint32 & M2) << 2)
  uint32 = ((uint32 >>> 4) & M4) | ((uint32 & M4) << 4)
  uint32 = ((uint32 >>> 8) & M8) | ((uint32 & M8) << 8)
  return ((uint32 >>> 16) | (uint32 << 16)) >>> 0
}

// 取模,模为2的幂
function modPow2Int32(x: number, mod: number): number {
  return x & (mod - 1)
}

// 计算一个整数的以 2 为底的对数
function log2Uint32(uint32: number): number {
  let res = 0
  // eslint-disable-next-line no-cond-assign
  while ((uint32 >>>= 1)) {
    res++
  }
  return res
}
