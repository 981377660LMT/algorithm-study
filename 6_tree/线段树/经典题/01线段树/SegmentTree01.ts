/* eslint-disable no-param-reassign */

/**
 * flip/indexOf/onesCount
 * @see https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
 */
class SegmentTree01 {
  private readonly _n: number
  private readonly _ones: Uint32Array
  private readonly _lazyFlip: Uint8Array

  /**
   * little-endian
   */
  constructor(bits: ArrayLike<number>) {
    this._n = bits.length
    const cap = 1 << (32 - Math.clz32(this._n - 1) + 1)
    this._ones = new Uint32Array(cap)
    this._lazyFlip = new Uint8Array(cap)
    this._build(1, 1, this._n, bits)
  }

  /**
   * 1 <= left <= right <= n
   */
  flip(left: number, right: number): void {
    this._flip(1, left, right, 1, this._n)
  }

  /**
   * 1 <= position <= n
   */
  indexOf(searchDigit: 0 | 1, position = 1): number {
    if (position > this._n) return -1
    if (searchDigit === 0) return this._indexofZero(1, position, 1, this._n)
    return this._indexofOne(1, position, 1, this._n)
  }

  /**
   * 1 <= left <= right <= n
   */
  onesCount(left: number, right: number): number {
    return this._onesCount(1, left, right, 1, this._n)
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

  private _onesCount(root: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this._ones[root]
    this._pushDown(root, l, r)
    const mid = (l + r) >>> 1
    let res = 0
    if (L <= mid) res += this._onesCount(root << 1, L, R, l, mid)
    if (mid < R) res += this._onesCount((root << 1) | 1, L, R, mid + 1, r)
    return res
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
      const mid = Math.floor((left + right) / 2)
      this._propagateFlip(root << 1, left, mid)
      this._propagateFlip((root << 1) | 1, mid + 1, right)
      this._lazyFlip[root] = 0
    }
  }

  private _propagateFlip(root: number, left: number, right: number): void {
    this._ones[root] = right - left + 1 - this._ones[root]
    this._lazyFlip[root] ^= 1
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
}

if (require.main === module) {
  const tree01 = new SegmentTree01([0, 1, 0, 1, 0, 1, 0, 1, 0, 1])
  console.log(tree01.indexOf(0, 2))
  console.log(tree01.indexOf(1, 1))
  console.log(tree01.toString())
  tree01.flip(2, 5)
  console.log(tree01.toString())
  console.log(tree01.onesCount(1, 10))

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
      idx++
      if (this.tree01.onesCount(idx, idx) === 1) return
      this.tree01.flip(idx, idx)
    }

    unfix(idx: number): void {
      idx++
      if (this.tree01.onesCount(idx, idx) === 0) return
      this.tree01.flip(idx, idx)
    }

    flip(): void {
      this.tree01.flip(1, this.size)
    }

    all(): boolean {
      return this.tree01.onesCount(1, this.size) === this.size
    }

    one(): boolean {
      return this.tree01.onesCount(1, this.size) > 0
    }

    count(): number {
      return this.tree01.onesCount(1, this.size)
    }

    toString(): string {
      return this.tree01.toString()
    }
  }

  /**
   * Your Bitset object will be instantiated and called as such:
   * var obj = new Bitset(size)
   * obj.fix(idx)
   * obj.unfix(idx)
   * obj.flip()
   * var param_4 = obj.all()
   * var param_5 = obj.one()
   * var param_6 = obj.count()
   * var param_7 = obj.toString()
   */
}

export { SegmentTree01 }
