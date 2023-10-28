const MOD = 998244353

class MoAlgo {
  private readonly _chunkSize: number
  private readonly _buckets: { qi: number; left: number; right: number }[][]
  private _queryOrder = 0

  constructor(n: number, q: number) {
    const sqrt = Math.sqrt((q * 2) / 3) | 0
    const chunkSize = Math.max(1, (n / Math.max(1, sqrt)) | 0)
    const buckets = Array(((n / chunkSize) | 0) + 1)
    for (let i = 0; i < buckets.length; i++) {
      buckets[i] = []
    }
    this._chunkSize = chunkSize
    this._buckets = buckets
  }

  /**
   * 添加一个查询，查询范围为`左闭右开区间` [start, end).
   * 0 <= start <= end <= n
   */
  addQuery(start: number, end: number): void {
    const index = (start / this._chunkSize) | 0
    this._buckets[index].push({ qi: this._queryOrder, left: start, right: end })
    this._queryOrder++
  }

  run(
    addLeft: (index: number) => void,
    addRight: (index: number) => void,
    delLeft: (index: number) => void,
    delRight: (index: number) => void,
    query: (qid: number) => void
  ): void {
    let left = 0
    let right = 0

    this._buckets.forEach((bucket, i) => {
      if (i & 1) {
        bucket.sort((a, b) => a.right - b.right)
      } else {
        bucket.sort((a, b) => b.right - a.right)
      }

      bucket.forEach(({ qi, left: ql, right: qr }) => {
        // !窗口扩张
        while (left > ql) {
          addLeft(--left)
        }
        while (right < qr) {
          addRight(right++)
        }

        // !窗口收缩
        while (left < ql) {
          delLeft(left++)
        }
        while (right > qr) {
          delRight(--right)
        }

        query(qi)
      })
    })
  }
}

function modAdd(num1: number, num2: number, mod = 998244353): number {
  let cand = (num1 + num2) % mod
  if (cand < 0) cand += mod
  return cand
}

function modSub(num1: number, num2: number, mod = 998244353): number {
  return modAdd(num1, -num2, mod)
}

/**
 * 效率不如bigint，但是对数组存储更友好.
 */
function modMul(num1: number, num2: number, mod = 998244353): number {
  return (((Math.floor(num1 / 65536) * num2) % mod) * 65536 + (num1 & 65535) * num2) % mod
}

function qpow(base: number, exp: number, mod: number): number {
  let res = 1
  while (exp) {
    if (exp & 1) res = modMul(res, base, mod)
    base = modMul(base, base, mod)
    exp = Math.floor(exp / 2)
  }
  return res
}

class Enumeration {
  private readonly _fac: number[] = [1]
  private readonly _ifac: number[] = [1]
  private readonly _inv: number[] = [1]
  private readonly _mod: number

  constructor(size: number, mod: number) {
    this._mod = mod
    this._expand(size)
  }
  fac(k: number): number {
    this._expand(k)
    return this._fac[k]
  }

  ifac(k: number): number {
    this._expand(k)
    return this._ifac[k]
  }

  inv(k: number): number {
    this._expand(k)
    return this._inv[k]
  }

  comb(n: number, k: number): number {
    if (n < 0 || k < 0 || n < k) {
      return 0
    }
    const mod = this._mod
    return modMul(modMul(this.fac(n), this.ifac(k), mod), this.ifac(n - k), mod)
  }

  private _expand(size: number): void {
    size = Math.min(size, this._mod - 1)
    if (this._fac.length < size + 1) {
      const mod = this._mod
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

const E = new Enumeration(2e5 + 10, MOD)

function multiPointBinomimalSum(qs: [number, number][]): number[] {
  let N = 2
  qs.forEach(([first]) => {
    N = Math.max(N, first)
  })

  const Q = qs.length
  const mo = new MoAlgo(N + 1, Q)
  qs.forEach(([first, second]) => {
    mo.addQuery(second, first)
  })

  const res: number[] = Array(Q).fill(0)
  let cur = 1
  let n = 0
  let m = 0
  const al = () => {
    cur -= E.comb(n, m--)
    cur %= MOD
  }
  const ar = () => {
    cur += cur - E.comb(n++, m)
    cur %= MOD
  }
  const el = () => {
    cur += E.comb(n, ++m)
    cur %= MOD
  }
  const er = () => {
    cur = modMul(modAdd(cur, E.comb(--n, m), MOD), E.inv(2), MOD)
  }
  const q = (i: number) => {
    res[i] = cur
  }
  mo.run(al, ar, el, er, q)
  return res
}

function beautifulString(s: string): number {
  let curOne = 0
  let curZero = 0
  const todo: [number, number][] = []
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

  const naive = (n: number, m: number): number => {
    let res = 0
    for (let i = 0; i <= m; i++) {
      res += E.comb(n, i)
      res %= MOD
    }
    return res
  }

  const multiRes = multiPointBinomimalSum(todo)
  for (let i = 0; i < multiRes.length; i++) {
    res -= multiRes[i]
    // res -= naive(todo[i][0], todo[i][1])
    res %= MOD
  }
  if (res < 0) res += MOD
  return res
}

export {}
// "11"
// "01001"

console.log(beautifulString('11'))
console.log(beautifulString('01001'))
// "0100111000100001111101111011010011110101001101110111001101110000110001110101110000000000010101111110"
console.log(beautifulString('0100111000100001111101111011010011110101001101110111001101110000110001110101110000000000010101111110'))
// 954957892
