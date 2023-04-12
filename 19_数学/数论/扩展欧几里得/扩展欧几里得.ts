/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable no-shadow */
/* eslint-disable semi-style */

import assert from 'assert'

function gcd<T extends number | bigint>(a: T, b: T): T {
  if (Number.isNaN(a) || Number.isNaN(b)) {
    return NaN as T
  }

  return b === 0 ? a : gcd(b, (a % b) as T)
}

/**
 * 扩展欧几里得求 `a*x + b*y = gcd(a,b)` 的一组解
 */
function exgcd(a: number, b: number): [x: number, y: number, gcd: number] {
  if (b === 0) return [1, 0, a]
  let [x, y, gcd] = exgcd(b, a % b)

  // 根据(b，a%b)的解推出(a,b)的解
  // 辗转相除法反向推导每层a、b的因子使得gcd(a,b)=ax+by成立
  const tmp = x - y * Math.floor(a / b)
  x = y
  y = tmp
  return [x, y, gcd]
}

/**
 * 求出逆元 `inv` 满足 `a*inv ≡ 1 (mod m)`
 */
function modularInverse(a: number, mod: number): number {
  const [x, _, gcd] = exgcd(a, mod)
  if (gcd !== 1) throw new Error('No inverse')
  return ((x % mod) + mod) % mod
}

if (require.main === module) {
  assert.strictEqual(modularInverse(3, 10), 7)
  const MOD = 998244353
  const INV2 = (MOD + 1) / 2 // 结论:2的逆元(模998244353)为499122177
  assert.strictEqual(modularInverse(2, MOD), INV2)
}

export { exgcd, gcd, modularInverse }
