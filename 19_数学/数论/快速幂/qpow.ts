/* eslint-disable no-param-reassign */

// 快速幂

/**
 * quick pow.
 *
 * @returns base**exp%mod
 * @warning base, exp, mod must be bigint.
 * Because js number precision is 2^53-1, which is not enough for mod**2.
 */
function qpowBigInt(base: bigint, exp: bigint, mod: bigint): bigint {
  if (!exp) return 1n % mod
  let res = 1n
  while (exp) {
    if (exp & 1n) res = (res * base) % mod
    base = (base * base) % mod
    exp >>= 1n
  }
  return res
}

function qpow(base: number, exp: number, mod = 1e9 + 7): number {
  if (!exp) return 1 % mod
  const bigMod = BigInt(mod)
  let bigBase = BigInt(base)
  let bigExp = BigInt(exp)
  let res = 1n
  while (bigExp) {
    if (bigExp & 1n) res = (res * bigBase) % bigMod
    bigBase = (bigBase * bigBase) % bigMod
    bigExp >>= 1n
  }
  return Number(res)
}

export { qpowBigInt, qpow }
