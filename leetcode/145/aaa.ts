function minLength(s: string, numOps: number): number {
  let n = s.length
  if (n <= 1) return n

  function buildOdt(str: string) {
    let odt = new ODT<string>(n, 'x')
    let start = 0
    for (let i = 1; i <= n; i++) {
      if (i === n || str[i] !== str[i - 1]) {
        odt.set(start, i, str[start])
        start = i
      }
    }
    return odt
  }

  function canAchieve(m: number): boolean {
    let odt = buildOdt(s)
    let queue: [number, number, string][] = []
    let flips = 0

    function pushLargeInterval(start: number, end: number, val: string) {
      let length = end - start
      if (val !== 'x' && length > m) {
        queue.push([start, end, val])
      }
    }

    function scanAll() {
      queue = []
      odt.enumerateAll((st, ed, v) => {
        pushLargeInterval(st, ed, v)
      })
    }

    scanAll()
    while (queue.length) {
      if (flips === numOps) return false
      flips++
      let [st, ed, val] = queue.shift()!
      let mid = (st + ed) >>> 1
      let flipVal = val === '0' ? '1' : '0'
      odt.set(mid, mid + 1, flipVal)
      scanAll()
    }
    return true
  }

  let left = 1
  let right = n
  while (left < right) {
    let mid = (left + right) >>> 1
    if (canAchieve(mid)) {
      right = mid
    } else {
      left = mid + 1
    }
  }
  return left
}

class ODT<S> {
  private _len = 0
  private _count = 0
  private readonly _leftLimit: number
  private readonly _rightLimit: number
  private readonly _noneValue: S
  private readonly _data: S[]
  private readonly _fs: FastSet

  constructor(n: number, noneValue: S) {
    let data = Array<S>(n)
    for (let i = 0; i < n; i++) data[i] = noneValue
    let fs = new FastSet(n)
    fs.insert(0)
    this._leftLimit = 0
    this._rightLimit = n
    this._noneValue = noneValue
    this._data = data
    this._fs = fs
  }

  get(x: number, erase = false): [number, number, S] | undefined {
    if (x < this._leftLimit || x >= this._rightLimit) return undefined
    let start = this._fs.prev(x)
    let end = this._fs.next(x + 1)
    let value = this._data[start]
    if (erase && value !== this._noneValue) {
      this._len--
      this._count -= end - start
      this._data[start] = this._noneValue
      this._mergeAt(start)
      this._mergeAt(end)
    }
    return [start, end, value]
  }

  set(start: number, end: number, value: S): void {
    this.enumerateRange(start, end, () => {}, true)
    this._fs.insert(start)
    this._data[start] = value
    if (value !== this._noneValue) {
      this._len++
      this._count += end - start
    }
    this._mergeAt(start)
    this._mergeAt(end)
  }

  enumerateAll(f: (start: number, end: number, value: S) => void): void {
    this.enumerateRange(this._leftLimit, this._rightLimit, f, false)
  }

  enumerateRange(
    start: number,
    end: number,
    f: (start: number, end: number, value: S) => void,
    erase = false
  ): void {
    if (start < this._leftLimit) start = this._leftLimit
    if (end > this._rightLimit) end = this._rightLimit
    if (start >= end) return
    let none = this._noneValue
    if (!erase) {
      let left = this._fs.prev(start)
      while (left < end) {
        let right = this._fs.next(left + 1)
        f(Math.max(left, start), Math.min(right, end), this._data[left])
        left = right
      }
      return
    }
    let p = this._fs.prev(start)
    if (p < start) {
      this._fs.insert(start)
      let v = this._data[p]
      this._data[start] = v
      if (v !== none) this._len++
    }
    p = this._fs.next(end)
    if (end < p) {
      let v = this._data[this._fs.prev(end)]
      this._data[end] = v
      this._fs.insert(end)
      if (v !== none) this._len++
    }
    p = start
    while (p < end) {
      let q = this._fs.next(p + 1)
      let x = this._data[p]
      f(p, q, x)
      if (x !== none) {
        this._len--
        this._count -= q - p
      }
      this._fs.erase(p)
      p = q
    }
    this._fs.insert(start)
    this._data[start] = none
  }

  get length(): number {
    return this._len
  }

  get count(): number {
    return this._count
  }

  private _mergeAt(p: number): void {
    if (p <= 0 || this._rightLimit <= p) return
    let q = this._fs.prev(p - 1)
    let dataP = this._data[p]
    let dataQ = this._data[q]
    if (dataP === dataQ) {
      if (dataP !== this._noneValue) this._len--
      this._fs.erase(p)
    }
  }
}

class FastSet {
  private readonly _n: number
  private readonly _lg: number
  private readonly _seg: Uint32Array[][]

  constructor(n: number) {
    this._n = n
    let seg: Uint32Array[][] = []
    let size = n
    while (true) {
      seg.push([new Uint32Array((size + 31) >>> 5)])
      size = (size + 31) >>> 5
      if (size <= 1) break
    }
    this._lg = seg.length
    this._seg = seg
    this.insert(0)
  }

  insert(i: number): void {
    for (let h = 0; h < this._lg; h++) {
      if (!this._seg[h][0][i >>> 5]) this._seg[h][0][i >>> 5] = 0
      this._seg[h][0][i >>> 5] |= 1 << (i & 31)
      i >>>= 5
    }
  }

  has(i: number): boolean {
    return !!(this._seg[0][0][i >>> 5] & (1 << (i & 31)))
  }

  erase(i: number): void {
    for (let h = 0; h < this._lg; h++) {
      this._seg[h][0][i >>> 5] &= ~(1 << (i & 31))
      if (this._seg[h][0][i >>> 5]) break
      i >>>= 5
    }
  }

  next(i: number): number {
    if (i < 0) i = 0
    if (i >= this._n) return this._n
    for (let h = 0; h < this._lg; h++) {
      let d = this._seg[h][0][i >>> 5] >>> (i & 31)
      if (!d) {
        i = (i >>> 5) + 1
        continue
      }
      i += 31 - Math.clz32(d & -d)
      for (let g = h - 1; ~g; g--) {
        i <<= 5
        let tmp = this._seg[g][0][i >>> 5]
        i += 31 - Math.clz32(tmp & -tmp)
      }
      return i
    }
    return this._n
  }

  prev(i: number): number {
    if (i < 0) return -1
    if (i >= this._n) i = this._n - 1
    for (let h = 0; h < this._lg; h++) {
      let d = this._seg[h][0][i >>> 5] << (31 - (i & 31))
      if (!d) {
        i = (i >>> 5) - 1
        continue
      }
      i -= Math.clz32(d)
      for (let g = h - 1; ~g; g--) {
        i <<= 5
        let tmp = this._seg[g][0][i >>> 5]
        i += 31 - Math.clz32(tmp & -tmp)
      }
      return i
    }
    return -1
  }

  enumerateRange(start: number, end: number, f: (v: number) => void): void {
    let x = start - 1
    while (true) {
      x = this.next(x + 1)
      if (x >= end) break
      f(x)
    }
  }

  get min(): number | null {
    let p = this.next(-1)
    return p >= this._n ? null : p
  }

  get max(): number | null {
    let p = this.prev(this._n)
    return p < 0 ? null : p
  }
}
