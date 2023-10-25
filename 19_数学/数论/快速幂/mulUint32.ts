/* eslint-disable no-inner-declarations */
/* eslint-disable implicit-arrow-linebreak */

import { qpow } from './qpow'

/**
 * 两个uint32数的乘积模`mod`.
 */
const mulUint32 = (num1: number, num2: number, mod = 1e9 + 7): number =>
  (((Math.floor(num1 / 65536) * num2) % mod) * 65536 + (num1 & 65535) * num2) % mod

/**
 * 多个uint32数的乘积模`mod`.
 */
const mulUint32Array = (arr: ArrayLike<number>, mod = 1e9 + 7): number => {
  if (!arr.length) throw new Error('mul: no args')
  if (arr.length === 1) return arr[0] % mod
  let res = 1
  for (let i = 0; i < arr.length; i++) {
    res = mulUint32(res, arr[i], mod)
  }
  return res
}

/**
 * uint32数的快速幂.
 */
const powUint32 = (base: number, exp: number, mod = 1e9 + 7): number => {
  base %= mod
  let res = 1 % mod
  while (exp) {
    if (exp & 1) res = mulUint32(res, base, mod)
    base = mulUint32(base, base, mod)
    exp = Math.floor(exp / 2)
  }
  return res
}

export { mulUint32, mulUint32Array, powUint32 }

if (require.main === module) {
  // 2550. 猴子碰撞的方法数
  // https://leetcode.cn/problems/count-collisions-of-monkeys-on-a-polygon/
  function monkeyMove(n: number): number {
    const MOD = 1e9 + 7
    let res = (powUint32(2, n, MOD) - 2) % MOD
    if (res < 0) res += MOD
    return res
  }

  testPerf()

  //
  function testPerf(): void {
    const base = 1000
    const exp = 1e9
    const mod = 1e9 + 7

    console.time('powUint32')
    for (let i = 0; i < 1e6; i++) {
      powUint32(base, exp, mod)
    }
    console.timeEnd('powUint32')

    console.time('powBigInt')
    for (let i = 0; i < 1e6; i++) {
      qpow(base, exp, mod)
    }
    console.timeEnd('powBigInt')
  }
}
