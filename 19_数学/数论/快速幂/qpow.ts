/* eslint-disable no-param-reassign */

/**
 * quick pow.
 *
 * @returns base**exp%mod
 * @warning base, exp, mod must be bigint.
 * Because js number precision is 2^53-1, which is not enough for mod**2.
 */
function qpow(base: bigint, exp: bigint, mod: bigint): bigint {
  if (exp === 0n) return 1n

  let res = 1n

  while (exp) {
    if (exp & 1n) {
      res *= base
      res %= mod
    }

    base *= base
    base %= mod
    exp >>= 1n
  }

  return res
}

export { qpow }
