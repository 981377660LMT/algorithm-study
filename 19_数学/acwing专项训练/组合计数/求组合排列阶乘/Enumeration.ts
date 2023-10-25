/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

import { modMul } from '../../../卷积/modInt'
import { qpow } from '../../../数论/快速幂/qpow'

class Enumeration {
  readonly mod: number
  private readonly _fac: number[] = [1]
  private readonly _ifac: number[] = [1]
  private readonly _inv: number[] = [1]

  constructor(size: number, mod: number) {
    this.mod = mod
    this._expand(size)
  }

  /** 阶乘. */
  fac(k: number): number {
    this._expand(k)
    return this._fac[k]
  }

  /** 阶乘逆元. */
  ifac(k: number): number {
    this._expand(k)
    return this._ifac[k]
  }

  /** 模逆元. */
  inv(k: number): number {
    this._expand(k)
    return this._inv[k]
  }

  /** 组合数. */
  comb(n: number, k: number): number {
    if (n < 0 || k < 0 || n < k) {
      return 0
    }
    const mod = this.mod
    return modMul(modMul(this.fac(n), this.ifac(k), mod), this.ifac(n - k), mod)
  }

  /** 排列数. */
  perm(n: number, k: number): number {
    if (n < 0 || k < 0 || n < k) {
      return 0
    }
    const mod = this.mod
    return modMul(this.fac(n), this.ifac(n - k), mod)
  }

  /** 可重复选取元素的组合数. */
  combWithReplacement(n: number, k: number): number {
    if (n === 0) return k === 0 ? 1 : 0
    return this.comb(n + k - 1, k)
  }

  /** n个相同的球放入k个不同的盒子(盒子可放任意个球)的方案数. */
  put(n: number, k: number): number {
    return this.comb(n + k - 1, n)
  }

  /** 卡特兰数. */
  catalan(n: number): number {
    return modMul(this.comb(2 * n, n), this.inv(n + 1), this.mod)
  }

  lucas(n: number, k: number): number {
    if (k === 0) return 1
    const mod = this.mod
    return modMul(this.comb(n % mod, k % mod), this.lucas(Math.floor(n / mod), Math.floor(k / mod)), mod)
  }

  private _expand(size: number): void {
    size = Math.min(size, this.mod - 1)
    if (this._fac.length < size + 1) {
      const mod = this.mod
      const preSize = this._fac.length
      const diff = size + 1 - preSize
      this._fac.length += diff
      this._ifac.length += diff
      this._inv.length += diff
      for (let i = preSize; i <= size; ++i) {
        this._fac[i] = modMul(this._fac[i - 1], i, mod)
      }
      this._ifac[size] = qpow(this._fac[size], mod - 2, mod) // !modInv
      for (let i = size - 1; i >= preSize; --i) {
        this._ifac[i] = modMul(this._ifac[i + 1], i + 1, mod)
      }
      for (let i = preSize; i <= size; ++i) {
        this._inv[i] = modMul(this._ifac[i], this._fac[i - 1], mod)
      }
    }
  }
}

export { Enumeration }

if (require.main === module) {
  // 2842. 统计一个字符串的 k 子序列美丽值最大的数目
  // https://leetcode.cn/contest/biweekly-contest-112/problems/count-k-subsequences-of-a-string-with-maximum-beauty/
  const MOD = 1e9 + 7
  const C = new Enumeration(1e5, MOD)

  function countKSubsequencesWithMaxBeauty(s: string, k: number): number {
    const counter = new Map<string, number>()
    for (let i = 0; i < s.length; i++) {
      const c = s[i]
      counter.set(c, (counter.get(c) || 0) + 1)
    }
    if (counter.size < k) return 0

    const freqCounter = new Map<number, number>()
    for (const count of counter.values()) {
      freqCounter.set(count, (freqCounter.get(count) || 0) + 1)
    }
    const items = [...freqCounter.entries()].sort((a, b) => b[0] - a[0])
    let res = 1
    let remain = k
    for (let i = 0; i < items.length; i++) {
      if (remain <= 0) break
      const { 0: freq, 1: sameCount } = items[i]
      const min_ = Math.min(remain, sameCount)
      const tmp = modMul(C.comb(sameCount, min_), qpow(freq, min_, MOD), MOD)
      res = modMul(res, tmp, MOD)
      remain -= min_
    }
    return res
  }
}
