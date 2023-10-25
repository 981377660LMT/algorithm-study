import { modDiv, modMul } from '../../../../卷积/modInt'

/**
 * 求组合数.适用于 n 巨大但 k 或 n-k 较小的情况.
 */
function bigModSmallKComb(n: number, k: number, mod: number): number {
  if (k > n - k) {
    k = n - k
  }
  let a = 1
  let b = 1
  for (let i = 1; i <= k; i++) {
    a = modMul(a, n, mod)
    n--
    b = modMul(b, i, mod)
  }
  return modDiv(a, b, mod)
}

if (require.main === module) {
  console.log(bigModSmallKComb(100, 2, 1e9 + 7))
}
