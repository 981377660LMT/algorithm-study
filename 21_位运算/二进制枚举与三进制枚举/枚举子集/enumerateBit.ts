/* eslint-disable no-console */
/* eslint-disable no-param-reassign */

// 遍历bits(非常快)

/**
 * 遍历每个为1的比特位
 * @param int32 int32 < 2**31
 */
function enumerateBits(int32: number, callback: (bit: number) => void): void {
  while (int32 > 0) {
    const i = 31 - Math.clz32(int32 & -int32)
    callback(i)
    int32 ^= 1 << i
  }
}

if (require.main === module) {
  console.time('enumerateBits')
  for (let i = 0; i < 1e8; i++) {
    enumerateBits(i, bit => {})
  }
  console.timeEnd('enumerateBits')
}

export { enumerateBits }
