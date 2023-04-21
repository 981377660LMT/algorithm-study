/**
 * 使用了双哈希的字符串哈希.
 * @important 比 {@link StringHasher2} 慢.
 */
class SafeStringHasher {
  private readonly _n: number
  private readonly _base1: number
  private readonly _base2: number
  private readonly _mod1: number
  private readonly _mod2: number
  private readonly _power1: Uint32Array
  private readonly _power2: Uint32Array
  private readonly _hash1: Uint32Array
  private readonly _hash2: Uint32Array

  /**
   * @param base
   * 131/13331/1333331/12821/12721/12421
   * @param mod
   * 1e7+19/1e7+79/1e7+103/1e7+121/1e7+139/1e7+141/1e7+169
   */
  constructor(
    s: ArrayLike<number> | string,
    base1 = 131,
    base2 = 13331,
    mod1 = 1e7 + 19,
    mod2 = 1e7 + 79
  ) {
    const n = s.length
    this._n = n
    this._base1 = base1
    this._base2 = base2
    this._mod1 = mod1
    this._mod2 = mod2
    const hash1 = new Uint32Array(n + 1)
    const hash2 = new Uint32Array(n + 1)
    const power1 = new Uint32Array(n + 1)
    const power2 = new Uint32Array(n + 1)
    power1[0] = 1
    power2[0] = 1
    if (typeof s === 'string') {
      for (let i = 0; i < n; i++) {
        const v = s.charCodeAt(i)
        hash1[i + 1] = (hash1[i] * base1 + v) % mod1
        hash2[i + 1] = (hash2[i] * base2 + v) % mod2
        power1[i + 1] = (power1[i] * base1) % mod1
        power2[i + 1] = (power2[i] * base2) % mod2
      }
    } else {
      for (let i = 0; i < n; i++) {
        hash1[i + 1] = (hash1[i] * base1 + s[i]) % mod1
        hash2[i + 1] = (hash2[i] * base2 + s[i]) % mod2
        power1[i + 1] = (power1[i] * base1) % mod1
        power2[i + 1] = (power2[i] * base2) % mod2
      }
    }
    this._hash1 = hash1
    this._hash2 = hash2
    this._power1 = power1
    this._power2 = power2
  }

  /**
   * 返回[start, end)的哈希值.
   */
  query(start = 0, end = this._n): number {
    if (start >= end) return 0
    const diff = end - start
    let res1 = this._hash1[end] - ((this._hash1[start] * this._power1[diff]) % this._mod1)
    let res2 = this._hash2[end] - ((this._hash2[start] * this._power2[diff]) % this._mod2)
    if (res1 < 0) res1 += this._mod1
    if (res2 < 0) res2 += this._mod2
    return res1 * this._mod2 + res2
  }

  /**
   * 哈希值h1和h2的拼接, h2的长度为h2len.
   */
  combine(h1: number, h2: number, h2len: number): number {
    const res1 = (h1 * this._power1[h2len] + h2) % this._mod1
    const res2 = (h1 * this._power2[h2len] + h2) % this._mod2
    return res1 * this._mod2 + res2
  }

  /**
   * 哈希值h和字符c拼接形成的哈希值.
   */
  addChar(h: number, c: number | string): number {
    const v = typeof c === 'string' ? c.charCodeAt(0) : c
    const res1 = (h * this._base1 + v) % this._mod1
    const res2 = (h * this._base2 + v) % this._mod2
    return res1 * this._mod2 + res2
  }
}

// --------------------------------------------
// API:
//  new StringHasher(base,mod)
//  build(s) : O(n)返回s的哈希值表.
//  query(sTable, l, r) : 返回s[l, r)的哈希值.
//  combine(h1, h2, h2len) : 哈希值h1和h2的拼接, h2的长度为h2len. (线段树的op操作)
//  addChar(h,c) : 哈希值h和字符c拼接形成的哈希值.
//  lcp(aTable, l1, r1, bTable, l2, r2) : O(logn)返回a[l1, r1)和b[l2, r2)的最长公共前缀.

type Hash = readonly [h1: number, h2: number]

/**
 * !使用双哈希判断,且mod和base的选取使得运算过程中不超过2^53-1,这种方式比BigInt64要快,
 * 比手动使用两个 {@link StringHasher1} 要慢(1000ms vs 1800ms).
 *
 * @example
 * ```ts
 * const H = new StringHasher2()
 * const s = 'abcde'
 * const table = H.build(s)
 * const h1 = H.query(table, 0, 3)
 * const h2 = H.query(table, 3, 6)
 * if (h1[0] === h2[0] && h1[1] === h2[1]) console.log('same')
 * ```
 *
 * 哈希值计算方法：
 * hash(s, p, m) = (val(s[0]) * p^k-1 + val(s[1]) * p^k-2 + ... + val(s[k-1]) * p^0) mod m.
 * 越靠左字符权重越大.
 *
 * @todo 更快的方法是不使用tuple, 而是将哈希表示为 `hash1*mod2+hahs2`.
 * @link https://github.com/drken1215/algorithm/blob/83d8e9bae9c46fb3b69ed8058d5a8a6db01f8300/String/rolling_hash.cpp#L18
 */
class StringHasher2 {
  private readonly _base0: number
  private readonly _base1: number
  private readonly _mod0: number
  private readonly _mod1: number
  private readonly _power0: number[]
  private readonly _power1: number[]

  /**
   * @param base
   * 131/13331/1333331/12821/12721/12421
   * @param mod
   * 1e7+19/1e7+79/1e7+103/1e7+121/1e7+139/1e7+141/1e7+169
   */
  constructor(base0 = 131, base1 = 13331, mod1 = 1e7 + 19, mod2 = 1e7 + 79) {
    this._base0 = base0
    this._base1 = base1
    this._mod0 = mod1
    this._mod1 = mod2
    this._power0 = [1]
    this._power1 = [1]
  }

  build(s: ArrayLike<number> | string): Hash[] {
    const n = s.length
    const hashTable: Hash[] = Array(n + 1)
    for (let i = 0; i < hashTable.length; i++) {
      hashTable[i] = [0, 0]
    }
    const base0 = this._base0
    const base1 = this._base1
    const mod0 = this._mod0
    const mod1 = this._mod1
    if (typeof s === 'string') {
      for (let i = 0; i < n; i++) {
        const v = s.charCodeAt(i)
        hashTable[i + 1] = [
          (hashTable[i][0] * base0 + v) % mod0,
          (hashTable[i][1] * base1 + v) % mod1
        ]
      }
    } else {
      for (let i = 0; i < n; i++) {
        const v = s[i]
        hashTable[i + 1] = [
          (hashTable[i][0] * base0 + v) % mod0,
          (hashTable[i][1] * base1 + v) % mod1
        ]
      }
    }
    return hashTable
  }

  query(sTable: Hash[], start: number, end: number): Hash {
    if (start >= end) return [0, 0]
    const diff = end - start
    this._expand(diff)
    const h1 = (sTable[end][0] - sTable[start][0] * this._power0[diff]) % this._mod0
    const h2 = (sTable[end][1] - sTable[start][1] * this._power1[diff]) % this._mod1
    return [h1 < 0 ? h1 + this._mod0 : h1, h2 < 0 ? h2 + this._mod1 : h2]
  }

  combine(h1: Hash, h2: Hash, h2len: number): Hash {
    this._expand(h2len)
    return [
      (h1[0] * this._power0[h2len] + h2[0]) % this._mod0,
      (h1[1] * this._power1[h2len] + h2[1]) % this._mod1
    ]
  }

  addChar(hash: Hash, c: string | number): Hash {
    const v = typeof c === 'string' ? c.charCodeAt(0) : c
    return [(hash[0] * this._base0 + v) % this._mod0, (hash[1] * this._base1 + v) % this._mod1]
  }

  lcp(
    sTable: Hash[],
    start1: number,
    end1: number,
    tTable: Hash[],
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
      const hash1 = this.query(sTable, start1, start1 + mid)
      const hash2 = this.query(tTable, start2, start2 + mid)
      if (hash1[0] === hash2[0] && hash1[1] === hash2[1]) {
        low = mid
      } else {
        high = mid
      }
    }
    return low
  }

  private _expand(size: number): void {
    if (this._power0.length < size + 1) {
      const power0 = this._power0
      const power1 = this._power1
      const base0 = this._base0
      const base1 = this._base1
      const mod0 = this._mod0
      const mod1 = this._mod1
      const preSize = this._power0.length
      const diff = size + 1 - preSize
      for (let i = 0; i < diff; i++) {
        power0.push(0)
        power1.push(0)
      }
      for (let i = preSize - 1; i < size; i++) {
        power0[i + 1] = (power0[i] * base0) % mod0
        power1[i + 1] = (power1[i] * base1) % mod1
      }
    }
  }
}

/**
 * !`注意这种方式很容易哈希冲突,需要使用双哈希.`
 * @example
 * ```ts
 * const H1 = new StringHasher1(131, 1e7 + 19)
 * const H2 = new StringHasher1(13331, 1e7 + 79)
 * const s = 'abcde'
 * const table1 = H1.build(s)
 * const table2 = H2.build(s)
 * const hash1 = H1.query(table1, 0, 3)
 * const hash2 = H2.query(table2, 0, 3)
 * ```
 * mod和base的选取使得运算过程中不超过2^53-1,这种方式比BigInt64要快.
 *
 *
 * 哈希值计算方法：
 * hash(s, p, m) = (val(s[0]) * p^k-1 + val(s[1]) * p^k-2 + ... + val(s[k-1]) * p^0) mod m.
 * 越靠左字符权重越大.
 */
class StringHasher1 {
  private readonly _base: number
  private readonly _mod: number
  private readonly _power: number[]

  /**
   * @param base
   * 131/13331/1333331/12821/12721/12421
   * @param mod
   * 1e7+19/1e7+79/1e7+103/1e7+121/1e7+139/1e7+141/1e7+169
   */
  constructor(base = 131, mod = 1e7 + 19) {
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

export { StringHasher2, StringHasher1, SafeStringHasher }
