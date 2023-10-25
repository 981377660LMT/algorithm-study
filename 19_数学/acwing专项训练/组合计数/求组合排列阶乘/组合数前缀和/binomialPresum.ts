/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

import { MoAlgo } from '../../../../../22_专题/离线查询/莫队/Moalgo'
import { modAdd, modMul } from '../../../../卷积/modInt'
import { qpow } from '../../../../数论/快速幂/qpow'
import { Enumeration } from '../Enumeration'

interface IEnumeration {
  readonly mod: number
  inv(k: number): number
  comb(n: number, k: number): number
}

/**
 * 莫队求组合数前缀和.
 * @param queries 每项为 [n, k]，表示组合数 C(n, k).
 * @param enumeration 组合数计算类.
 * @returns 数组第i项为组合数前缀和 `C(ni,0) + C(ni,1) + ... + C(ni,ki)`.
 * @summary 利用递推公式 `S(n,k) = 2*S(n-1,k) - C(n-1,k)`求解.
 * @alias multiPointBinomimalSum
 */
function binomialPresum(queries: [n: number, k: number][], enumeration: IEnumeration): number[] {
  let maxN = 2
  queries.forEach(q => {
    maxN = Math.max(maxN, q[0])
  })

  const q = queries.length
  const mo = new MoAlgo(maxN + 1, q)
  queries.forEach(q => {
    mo.addQuery(q[1], q[0])
  })

  const res: number[] = Array(q).fill(0)
  const inv2 = enumeration.inv(2)
  const { mod } = enumeration
  let cur = 1
  let curN = 0
  let curK = 0
  const addLeft = () => {
    cur -= enumeration.comb(curN, curK--)
    cur %= mod
  }
  const addRight = () => {
    cur += cur - enumeration.comb(curN++, curK)
    cur %= mod
  }
  const removeLeft = () => {
    cur += enumeration.comb(curN, ++curK)
    cur %= mod
  }
  const removeRight = () => {
    cur = modMul(modAdd(cur, enumeration.comb(--curN, curK), mod), inv2, mod)
  }
  const query = (i: number) => {
    if (cur < 0) cur += mod
    res[i] = cur
  }
  mo.run(addLeft, addRight, removeLeft, removeRight, query)
  return res
}

export { binomialPresum }

if (require.main === module) {
  // 华科大-04. 美丽字符串
  // https://leetcode.cn/contest/hust_1024_2023/problems/yH1vqC/
  const MOD = 998244353
  const enumeration = new Enumeration(2e5, MOD)

  function beautifulString(s: string): number {
    let curOne = 0
    let curZero = 0
    const todo: [n: number, k: number][] = []
    let res = 0
    for (let i = 0; i < s.length; i++) {
      if (s[i] === '1') {
        curOne++
      } else {
        curZero++
      }

      const atLeastSelect = i + 1 - (s[i] === '1' ? curOne : curZero)
      res += qpow(2, i, MOD)
      res %= MOD
      todo.push([i, atLeastSelect - 1])
    }

    const preSum = binomialPresum(todo, enumeration)
    for (let i = 0; i < preSum.length; i++) {
      res -= preSum[i]
      res %= MOD
    }

    if (res < 0) res += MOD
    return res
  }
}
