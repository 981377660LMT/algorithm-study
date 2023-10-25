/* eslint-disable no-param-reassign */

import { modMul } from '../../卷积/modInt'

// 快速幂

/**
 * quick pow.
 *
 * @returns base**exp%mod
 * @warning base, exp, mod must be bigint.
 * Because js number precision is 2^53-1, which is not enough for mod**2.
 */
function qpowBigInt(base: bigint, exp: bigint, mod: bigint): bigint {
  base %= mod
  let res = 1n % mod
  while (exp) {
    if (exp & 1n) res = (res * base) % mod
    base = (base * base) % mod
    exp >>= 1n
  }
  return res
}

function qpow(base: number, exp: number, mod = 1e9 + 7): number {
  base %= mod
  let res = 1 % mod
  while (exp > 0) {
    if (exp & 1) res = modMul(res, base, mod)
    base = modMul(base, base, mod)
    exp = Math.floor(exp / 2)
  }
  return res
}

export { qpowBigInt, qpow, qpow as pow }
