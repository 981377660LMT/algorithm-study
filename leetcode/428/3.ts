export {}

class StringHasher1 {
  private readonly _base: number
  private readonly _mod: number
  private readonly _power: number[]
  constructor(base = 13131, mod = 1e7 + 169) {
    this._base = base
    this._mod = mod
    this._power = [1]
  }
  build(s: ArrayLike<number> | string): number[] {
    const n = s.length
    const hashTable = Array(n + 1).fill(0)
    const mod = this._mod
    const base = this._base
    if (typeof s === 'string') {
      for (let i = 0; i < n; i++) {
        hashTable[i + 1] = (hashTable[i] * base + s.charCodeAt(i)) % mod
      }
    } else {
      for (let i = 0; i < n; i++) {
        hashTable[i + 1] = (hashTable[i] * base + s[i]) % mod
      }
    }
    return hashTable
  }
  query(sTable: number[], start: number, end: number): number {
    this._expand(end - start)
    const res = (sTable[end] - sTable[start] * this._power[end - start]) % this._mod
    return res < 0 ? res + this._mod : res
  }
  combine(h1: number, h2: number, h2len: number): number {
    this._expand(h2len)
    return (h1 * this._power[h2len] + h2) % this._mod
  }
  addChar(hash: number, c: string | number): number {
    if (typeof c === 'string') {
      return (hash * this._base + c.charCodeAt(0)) % this._mod
    }
    return (hash * this._base + c) % this._mod
  }
  lcp(
    sTable: number[],
    start1: number,
    end1: number,
    tTable: number[],
    start2: number,
    end2: number
  ): number {
    const len1 = end1 - start1
    const len2 = end2 - start2
    const len = Math.min(len1, len2)
    let low = 0
    let high = len + 1
    while (high - low > 1) {
      const mid = (low + high) >>> 1
      if (this.query(sTable, start1, start1 + mid) === this.query(tTable, start2, start2 + mid)) {
        low = mid
      } else {
        high = mid
      }
    }
    return low
  }
  private _expand(size: number): void {
    if (this._power.length < size + 1) {
      const preSize = this._power.length
      const diff = size + 1 - preSize
      const power = this._power
      const base = this._base
      const mod = this._mod
      for (let i = 0; i < diff; i++) {
        power.push(0)
      }
      for (let i = preSize - 1; i < size; i++) {
        power[i + 1] = (power[i] * base) % mod
      }
    }
  }
}

function beautifulSplits(nums: number[]): number {
  const n = nums.length
  const hasher = new StringHasher1()
  const table = hasher.build(nums)
  let res = 0
  for (let i = 1; i <= n - 2; i++) {
    for (let j = i + 1; j <= n - 1; j++) {
      let ok = false
      const len2 = j - i
      const len3 = n - j
      if (j >= 2 * i) {
        const h1 = hasher.query(table, 0, i)
        const h2 = hasher.query(table, i, i + i)
        if (h1 === h2) ok = true
      }
      if (!ok && len2 <= len3) {
        const h2 = hasher.query(table, i, j)
        const h3 = hasher.query(table, j, j + len2)
        if (h2 === h3) ok = true
      }
      if (ok) res++
    }
  }
  return res
}
