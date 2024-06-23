/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
// more details:https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta

/**
 * 01线段树，支持 flip/indexOf/lastIndexOf/onesCount/kth，可用于模拟Bitset.
 */
class SegmentTree01 {
  private readonly _n: number
  private readonly _ones: Uint32Array
  private readonly _lazyFlip: Uint8Array

  /**
   * little-endian
   * @param bitsOrLength 0/1数组或者是长度.注意必须要是正整数.
   * @example
   * ```ts
   * const seg01 = new SegmentTree01([1, 0, 1, 1, 0, 1])
   * seg01.toString() // 101101
   * ```
   */
  constructor(bitsOrLength: ArrayLike<number> | number) {
    if (typeof bitsOrLength === 'number') bitsOrLength = Array(bitsOrLength).fill(0)
    if (bitsOrLength.length === 0) throw new Error('empty bits')
    this._n = bitsOrLength.length
    const log = 32 - Math.clz32(this._n - 1)
    const size = 1 << log
    this._ones = new Uint32Array(size << 1)
    this._lazyFlip = new Uint8Array(size) // 叶子结点不需要更新lazy (composition)
    this._build(1, 1, this._n, bitsOrLength)
  }

  /**
   * 0 <= start <= end <= n
   * 翻转[start,end)区间的bit.
   */
  flip(start: number, end: number): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return

    start++
    this._flip(1, start, end, 1, this._n)
  }

  /**
   * 0 <= position <= n-1.
   * @param searchDigit 0/1
   * @param position 查找的起始位置, 0 <= position < n.
   */
  indexOf(searchDigit: 0 | 1, position = 0): number {
    position++
    if (position > this._n) return -1
    if (searchDigit === 0) {
      const cand = this._indexofZero(1, position, 1, this._n)
      return cand === -1 ? cand : cand - 1
    }
    const cand = this._indexofOne(1, position, 1, this._n)
    return cand === -1 ? cand : cand - 1
  }

  /**
   * 0 <= position <= n-1.
   * @param searchDigit 0/1
   * @param position 查找的起始位置, 0 <= position < n.
   */
  lastIndexOf(searchDigit: 0 | 1, position = this._n - 1): number {
    position++
    if (position < 1) return -1
    if (searchDigit === 0) {
      const cand = this._lastIndexOfZero(1, position, 1, this._n)
      return cand === -1 ? cand : cand - 1
    }
    const cand = this._lastIndexOfOne(1, position, 1, this._n)
    return cand === -1 ? cand : cand - 1
  }

  /**
   * 0 <= left <= right <= n
   * 返回[left,right)区间内1的个数.
   */
  onesCount(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return 0
    start++
    return this._onesCount(1, start, end, 1, this._n)
  }

  /**
   * 树上二分查询第k个0/1的位置.如果不存在第k个0/1，返回-1.
   * !k >= 1
   * @returns -1<=pos<n.
   */
  kth(searchDigit: 0 | 1, k: number): number {
    if (searchDigit === 0) {
      if (k > this._n - this._ones[1]) return -1
      return this._kthZero(1, k, 1, this._n) - 1
    }
    if (k > this._ones[1]) return -1
    return this._kthOne(1, k, 1, this._n) - 1
  }

  get(index: number): 0 | 1 {
    return this.onesCount(index, index + 1) as 0 | 1
  }

  set(index: number, value: 0 | 1): void {
    this._set(1, index + 1, value, 1, this._n)
  }

  toString(): string {
    const sb: string[] = []
    this._toString(1, 1, this._n, sb)
    return sb.join('')
  }

  private _flip(root: number, L: number, R: number, l: number, r: number): void {
    if (L <= l && r <= R) {
      this._propagateFlip(root, l, r)
      return
    }
    this._pushDown(root, l, r)
    const mid = (l + r) >>> 1
    if (L <= mid) this._flip(root << 1, L, R, l, mid)
    if (mid < R) this._flip((root << 1) | 1, L, R, mid + 1, r)
    this._pushUp(root)
  }

  private _indexofOne(root: number, position: number, left: number, right: number): number {
    if (left === right) {
      if (this._ones[root] > 0) return left
      return -1
    }
    this._pushDown(root, left, right)
    const mid = (left + right) >>> 1
    if (position <= mid && this._ones[root << 1] > 0) {
      const leftPos = this._indexofOne(root << 1, position, left, mid)
      if (leftPos > 0) return leftPos
    }
    return this._indexofOne((root << 1) | 1, position, mid + 1, right)
  }

  private _indexofZero(root: number, position: number, left: number, right: number): number {
    if (left === right) {
      if (this._ones[root] === 0) return left
      return -1
    }
    this._pushDown(root, left, right)
    const mid = (left + right) >>> 1
    if (position <= mid && this._ones[root << 1] < mid - left + 1) {
      const leftPos = this._indexofZero(root << 1, position, left, mid)
      if (leftPos > 0) return leftPos
    }
    return this._indexofZero((root << 1) | 1, position, mid + 1, right)
  }

  private _lastIndexOfOne(root: number, position: number, left: number, right: number): number {
    if (left === right) {
      if (this._ones[root] > 0) return left
      return -1
    }
    this._pushDown(root, left, right)
    const mid = (left + right) >>> 1
    if (position > mid && this._ones[(root << 1) | 1] > 0) {
      const rightPos = this._lastIndexOfOne((root << 1) | 1, position, mid + 1, right)
      if (rightPos > 0) return rightPos
    }
    return this._lastIndexOfOne(root << 1, position, left, mid)
  }

  private _lastIndexOfZero(root: number, position: number, left: number, right: number): number {
    if (left === right) {
      if (this._ones[root] === 0) return left
      return -1
    }
    this._pushDown(root, left, right)
    const mid = (left + right) >>> 1
    if (position > mid && this._ones[(root << 1) | 1] < right - mid) {
      const rightPos = this._lastIndexOfZero((root << 1) | 1, position, mid + 1, right)
      if (rightPos > 0) return rightPos
    }
    return this._lastIndexOfZero(root << 1, position, left, mid)
  }

  private _onesCount(root: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this._ones[root]
    this._pushDown(root, l, r)
    const mid = (l + r) >>> 1
    let res = 0
    if (L <= mid) res += this._onesCount(root << 1, L, R, l, mid)
    if (mid < R) res += this._onesCount((root << 1) | 1, L, R, mid + 1, r)
    return res
  }

  private _kthOne(root: number, k: number, left: number, right: number): number {
    if (left === right) return left
    this._pushDown(root, left, right)
    const mid = (left + right) >>> 1
    if (this._ones[root << 1] >= k) return this._kthOne(root << 1, k, left, mid)
    return this._kthOne((root << 1) | 1, k - this._ones[root << 1], mid + 1, right)
  }

  private _kthZero(root: number, k: number, left: number, right: number): number {
    if (left === right) return left
    this._pushDown(root, left, right)
    const mid = (left + right) >>> 1
    const leftZero = mid - left + 1 - this._ones[root << 1]
    if (leftZero >= k) return this._kthZero(root << 1, k, left, mid)
    return this._kthZero((root << 1) | 1, k - leftZero, mid + 1, right)
  }

  private _set(root: number, position: number, value: 0 | 1, left: number, right: number): void {
    if (left === right) {
      this._ones[root] = value
      return
    }
    this._pushDown(root, left, right)
    const mid = (left + right) >>> 1
    if (position <= mid) this._set(root << 1, position, value, left, mid)
    else this._set((root << 1) | 1, position, value, mid + 1, right)
    this._pushUp(root)
  }

  private _toString(root: number, left: number, right: number, sb: string[]): void {
    if (left === right) {
      sb.push(this._ones[root] === 1 ? '1' : '0')
      return
    }
    this._pushDown(root, left, right)
    const mid = (left + right) >>> 1
    this._toString(root << 1, left, mid, sb)
    this._toString((root << 1) | 1, mid + 1, right, sb)
  }

  private _build(root: number, left: number, right: number, leaves: ArrayLike<number>): void {
    if (left === right) {
      this._ones[root] = leaves[left - 1]
      return
    }
    const mid = (left + right) >>> 1
    this._build(root << 1, left, mid, leaves)
    this._build((root << 1) | 1, mid + 1, right, leaves)
    this._pushUp(root)
  }

  private _pushUp(root: number): void {
    this._ones[root] = this._ones[root << 1] + this._ones[(root << 1) | 1]
  }

  private _pushDown(root: number, left: number, right: number): void {
    if (this._lazyFlip[root] !== 0) {
      const mid = (left + right) >>> 1
      this._propagateFlip(root << 1, left, mid)
      this._propagateFlip((root << 1) | 1, mid + 1, right)
      this._lazyFlip[root] = 0
    }
  }

  private _propagateFlip(root: number, left: number, right: number): void {
    this._ones[root] = right - left + 1 - this._ones[root]
    if (root < this._lazyFlip.length) {
      this._lazyFlip[root] ^= 1
    }
  }
}

if (require.main === module) {
  // 01线段树模拟位集
  // https://leetcode.cn/problems/design-bitset/

  class Bitset {
    private readonly size: number
    private readonly tree01: SegmentTree01

    constructor(size: number) {
      this.size = size
      this.tree01 = new SegmentTree01(new Uint8Array(size))
    }

    fix(idx: number): void {
      this.tree01.set(idx, 1)
    }

    unfix(idx: number): void {
      this.tree01.set(idx, 0)
    }

    flip(): void {
      this.tree01.flip(0, this.size)
    }

    all(): boolean {
      return this.tree01.onesCount(0, this.size) === this.size
    }

    one(): boolean {
      return this.tree01.onesCount(0, this.size) > 0
    }

    count(): number {
      return this.tree01.onesCount(0, this.size)
    }

    toString(): string {
      return this.tree01.toString()
    }
  }
}

export { SegmentTree01 }
