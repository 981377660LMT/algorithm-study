/* eslint-disable no-console */
/* eslint-disable no-param-reassign */

// 遍历bits(非常快)

/**
 * 遍历每个为1的比特位
 * @param int32 int32 < 2**31
 */
function enumerateBits(int32: number, callback: (bit: number) => void): void {
  for (let i = 0; int32 > 0; i++) {
    if (int32 & 1) {
      callback(i)
    }
    int32 >>= 1
  }
}

function enumerateBits32(s: number, f: (bit: number) => void): void {
  while (s) {
    const i = 31 - Math.clz32(s & -s) // lowbit.bit_length() - 1
    f(i)
    s ^= 1 << i
  }
}

if (require.main === module) {
  console.time('enumerateBits')
  for (let i = 0; i < 1e6; i++) {
    enumerateBits(0b110011011, bit => {})
  }
  console.timeEnd('enumerateBits')

  console.time('enumerateBits32')
  for (let i = 0; i < 1e6; i++) {
    enumerateBits32(0b1100110111, bit => {})
  }
  console.timeEnd('enumerateBits32')
}

export { enumerateBits, enumerateBits32 }
