/* eslint-disable no-param-reassign */

import assert from 'assert'

const N = 1e5
const MOD = BigInt(1e9 + 7)
const fac = new BigUint64Array(N + 1).fill(1n)
const ifac = new BigUint64Array(N + 1).fill(1n)

for (let i = 2; i <= N; i++) {
  fac[i] = (fac[i - 1] * BigInt(i)) % MOD
  ifac[i] = calInv(fac[i], MOD)
}

function comb(n: number, k: number): bigint {
  if (n < 0 || k < 0 || n < k) return 0n
  return (((fac[n] * ifac[k]) % MOD) * ifac[n - k]) % MOD
}

function calInv(a: bigint, mod: bigint): bigint {
  return qpow(a, mod - 2n, mod)
}

// 注意js精度丢失问题 需要用大数
function qpow(base: bigint, exp: bigint, mod: bigint): bigint {
  if (exp === 0n) return 1n

  let res = 1n

  while (exp) {
    if (exp & 1n) {
      res *= base
      res %= mod
    }

    // 此处超出2^53-1 需要bigint ((1e9+7*1e9+7)已经超出))
    base *= base
    base %= mod
    exp >>= 1n
  }

  return res
}

if (require.main === module) {
  assert.strictEqual(comb(12345, 123), 52132431n)
}

export { comb }
