/**
 * W叉Trie树,这里W=32.
 */
class WAryTrie {
  private readonly _n: number
  private readonly _a1: Uint32Array
  private readonly _a2: Uint32Array
  private readonly _a3: Uint32Array
  private _a4 = 0
  private _size = 0

  /**
   * 建立一个元素范围为`[0,n)`的W叉Trie树.
   * @param n n<2**20.
   */
  constructor(n: number) {
    this._n = n
    this._a1 = new Uint32Array((n >>> 5) + 1)
    this._a2 = new Uint32Array((n >>> 10) + 1)
    this._a3 = new Uint32Array((n >>> 15) + 1)
  }

  insert(i: number): boolean {
    if (this.has(i)) return false
    this._a1[i >>> 5] |= 1 << (i & 31)
    this._a2[i >>> 10] |= 1 << ((i >>> 5) & 31)
    this._a3[i >>> 15] |= 1 << ((i >>> 10) & 31)
    this._a4 |= 1 << (i >>> 15)
    this._size++
    return true
  }

  has(i: number): boolean {
    return !!((this._a1[i >>> 5] >>> (i & 31)) & 1)
  }

  erase(i: number): boolean {
    const bit0 = 1 << (i & 31)
    if (!(this._a1[i >>> 5] & bit0)) return false
    this._size--
    this._a1[i >>> 5] -= bit0
    if (this._a1[i >>> 5]) return true
    const bit1 = 1 << ((i >>> 5) & 31)
    this._a2[i >>> 10] -= bit1
    if (this._a2[i >>> 10]) return true
    const bit2 = 1 << ((i >>> 10) & 31)
    this._a3[i >>> 15] -= bit2
    if (this._a3[i >>> 15]) return true
    this._a4 -= 1 << (i >>> 15)
    return true
  }

  /**
   * 返回大于等于i的最小元素.
   * 如果不存在,返回`n`.
   */
  next(i: number): number {
    if (i < 0) i = 0
    if (i >= this._n) return this._n
    if (this.has(i)) return i

    let a = this._a1[i >>> 5]
    if (WAryTrie._nextBit(a, i) > 1) {
      return i + 1 + WAryTrie._minBit(WAryTrie._nextBit(a, i + 1))
    }
    i >>>= 5
    a = this._a2[i >>> 5]
    if (WAryTrie._nextBit(a, i) > 1) {
      i += 1 + WAryTrie._minBit(WAryTrie._nextBit(a, i + 1))
      return (i << 5) + WAryTrie._minBit(this._a1[i])
    }
    i >>>= 5
    a = this._a3[i >>> 5]
    if (WAryTrie._nextBit(a, i) > 1) {
      i += 1 + WAryTrie._minBit(WAryTrie._nextBit(a, i + 1))
      i = (i << 5) + WAryTrie._minBit(this._a2[i])
      return (i << 5) + WAryTrie._minBit(this._a1[i])
    }
    i >>>= 5
    if (WAryTrie._nextBit(this._a4, i) > 1) {
      i += 1 + WAryTrie._minBit(WAryTrie._nextBit(this._a4, i + 1))
      i = (i << 5) + WAryTrie._minBit(this._a3[i])
      i = (i << 5) + WAryTrie._minBit(this._a2[i])
      return (i << 5) + WAryTrie._minBit(this._a1[i])
    }

    return this._n
  }

  /**
   * 返回小于等于i的最大元素.
   * 如果不存在,返回`-1`.
   */
  prev(i: number): number {
    if (i < 0) return -1
    if (i >= this._n) i = this._n - 1
    if (this.has(i)) return i

    const tmp1 = WAryTrie._prevBit(this._a1[i >>> 5], i)
    if (tmp1) {
      // 低5位设置为0
      return (i & 0xffffffe0) + WAryTrie._maxBit(tmp1)
    }
    i >>>= 5
    const tmp2 = WAryTrie._prevBit(this._a2[i >>> 5], i)
    if (tmp2) {
      i = (i & 0xffffffe0) + WAryTrie._maxBit(tmp2)
      return (i << 5) + WAryTrie._maxBit(this._a1[i])
    }
    i >>>= 5
    const tmp3 = WAryTrie._prevBit(this._a3[i >>> 5], i)
    if (tmp3) {
      i = (i & 0xffffffe0) + WAryTrie._maxBit(tmp3)
      i = (i << 5) + WAryTrie._maxBit(this._a2[i])
      return (i << 5) + WAryTrie._maxBit(this._a1[i])
    }
    i >>>= 5
    const tmp4 = WAryTrie._prevBit(this._a4, i)
    if (tmp4) {
      i = WAryTrie._maxBit(tmp4)
      i = (i << 5) + WAryTrie._maxBit(this._a3[i])
      i = (i << 5) + WAryTrie._maxBit(this._a2[i])
      return (i << 5) + WAryTrie._maxBit(this._a1[i])
    }

    return -1
  }

  /**
   * 遍历[start,end)区间内的元素.
   */
  enumerateRange(start: number, end: number, f: (v: number) => void): void {
    let x = start - 1
    while (true) {
      x = this.next(x + 1)
      if (x >= end) break
      f(x)
    }
  }

  toString(): string {
    const sb: string[] = []
    this.enumerateRange(0, this._n, v => sb.push(v.toString()))
    return `WAryTrie(${this.size}){${sb.join(', ')}}`
  }

  /**
   * 返回集合中的最小值.如果不存在, 返回-1.
   */
  get min(): number {
    if (!this._a4) return -1
    let x = WAryTrie._minBit(this._a4)
    x = (x << 5) + WAryTrie._minBit(this._a3[x])
    x = (x << 5) + WAryTrie._minBit(this._a2[x])
    return (x << 5) + WAryTrie._minBit(this._a1[x])
  }

  /**
   * 返回集合中的最大值.如果不存在, 返回n.
   */
  get max(): number {
    if (!this._a4) return this._n
    let x = WAryTrie._maxBit(this._a4)
    x = (x << 5) + WAryTrie._maxBit(this._a3[x])
    x = (x << 5) + WAryTrie._maxBit(this._a2[x])
    return (x << 5) + WAryTrie._maxBit(this._a1[x])
  }

  get size(): number {
    return this._size
  }

  private static _maxBit(n: number): number {
    return 31 - Math.clz32(n)
  }

  private static _minBit(n: number): number {
    if (!n) return 32
    return 31 - Math.clz32(n & -n)
  }

  private static _prevBit(x: number, y: number): number {
    return x & ((1 << (y & 31)) - 1)
  }

  private static _nextBit(x: number, y: number): number {
    return x >>> (y & 31)
  }
}

export { WAryTrie }

if (require.main === module) {
  const set = new WAryTrie(100)
  set.insert(1)
  set.insert(2)
  set.insert(30)
  console.log(set.toString(), set.next(4), set.min, set.max)
  set.erase(30)
  console.log(set.toString(), set.next(4), set.min, set.max)

  const n = 1e7
  const set2 = new WAryTrie(n)
  console.time('WAryTrie')
  for (let i = 0; i < n; i++) {
    set2.insert(i)
    set2.next(i)
    set2.prev(i)
    set2.has(i)
    set2.erase(i)
    set2.insert(i)
  }
  console.timeEnd('WAryTrie') // 250ms
}
