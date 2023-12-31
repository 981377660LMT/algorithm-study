export function imul(a: number, b: number): number {
  a |= 0 // int
  b |= 0 // int
  const c = a & 0xffff
  const d = b & 0xffff
  // Shift by 0 fixes the sign on the high part.
  return (c * d + ((((a >>> 16) * d + c * (b >>> 16)) << 16) >>> 0)) | 0 // int
}

// v8 has an optimization for storing 31-bit signed numbers.
// Values which have either 00 or 11 as the high order bits qualify.
// This function drops the highest order bit in a signed number, maintaining
// the sign bit.
// 保留 i32 的符号位和低30位，丢弃最高位（第31位），并将第30位替换为 i32 右移一位后的第30位。
// 这样做的目的是为了优化V8引擎对31位有符号整数的存储。
export function smi(i32: number): number {
  return ((i32 >>> 1) & 0x40000000) | (i32 & 0xbfffffff)
}
