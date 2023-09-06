// 一个list，长度固定为n，两种操作：
// 1. set i, v 设置某个下标为0或1
// 2. query i 查找从i下标开始的第一个1的位置

/**
 * 01线段树，支持 flip/indexOf/onesCount/kth，可用于模拟Bitset
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

  /**
   * 树上二分查询第k个0/1的位置
   * k >= 1
   */
  kth(searchDigit: 0 | 1, k: number): number {
    if (searchDigit === 0) {
      if (k > this._n - this._ones[1]) return -1
      return this._kthZero(1, k, 1, this._n)
    }
    if (k > this._ones[1]) return -1
    return this._kthOne(1, k, 1, this._n)
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
    this._lazyFlip[root] ^= 1
  }

  set(index: number): void {
    index++
    if (this.onesCount(index, index) === 1) return
    this.flip(index, index)
  }

  unset(index: number): void {
    index++
    if (this.onesCount(index, index) === 0) return
    this.flip(index, index)
  }
}

function lowbit(x: number) {
  return (x & -x) >>> 0
}

function bitPos(x: number) {
  let ret = 0
  for (let bits = 16; bits; bits >>= 1) {
    if (((x >> ret) >> bits) & ((1 << bits) - 1)) {
      ret |= bits
    }
  }
  return ret
}

class BitSet {
  private readonly _bitset: Uint32Array
  private readonly _fenwick: Uint32Array
  private readonly _bucketCount: number

  constructor(readonly n: number) {
    const size = (n + 31) >> 5
    this._bucketCount = size
    this._bitset = new Uint32Array(size)
    this._fenwick = new Uint32Array(size)
  }

  add(index: number): this {
    const id = index >> 5
    const mask = index & 31
    if (!this._bitset[id]) {
      this.fenwickAdd(id, 1)
    }
    this._bitset[id] |= 1 << mask
    return this
  }

  delete(index: number): void {
    const id = index >> 5
    const mask = index & 31
    this._bitset[id] &= -1 << mask
    if (!this._bitset[id]) {
      this.fenwickAdd(id, -1)
    }
  }

  has(index: number): boolean {
    const id = index >> 5
    const mask = index & 31
    return (this._bitset[id] & (1 << mask)) !== 0
  }

  indexOfOne(start = 0): number {
    const id = start >> 5
    const mask = start & 31
    const rest = this._bitset[id] & (-1 << mask)
    if (rest) {
      return (id << 5) | bitPos(lowbit(rest))
    }

    const pos = this.fenwickFind(id + 1)
    if (pos < 0) {
      return -1
    }

    return (pos << 5) | bitPos(lowbit(this._bitset[pos]))
  }

  private fenwickAdd(index: number, delta: number): void {
    while (index) {
      this._fenwick[index] += delta
      index &= ~lowbit(index)
    }
    this._fenwick[index] += delta
  }

  private fenwickFind(pos: number): number {
    // up phase, find first non-zero block
    const n = this._bucketCount
    let i = 1
    if (!this._fenwick[pos]) {
      for (; pos < n; i <<= 1) {
        if (!(pos & i)) {
          pos = (pos & -i) | i
          if (this._fenwick[pos]) {
            break
          }
        }
      }
    }
    if (pos >= n) {
      return -1
    }

    let blockCount = this._fenwick[pos]
    for (i = lowbit(pos) >> 1; i; i >>= 1) {
      // if left block is zero, turn to right block.
      const rcount = this._fenwick[pos | i]
      if (blockCount === rcount) {
        pos |= i
      } else {
        blockCount -= rcount
      }
    }
    return pos
  }
}

if (require.main === module) {
  const N = 2e5
  const UPDATE_TIMES = 1e5
  const QUERY_TIMES = 1e5
  const bitset = new BitSet(N)
  const segtree = new SegmentTree01(new Uint8Array(N))
  const updateIndexes = Array.from({ length: UPDATE_TIMES }, () => Math.floor(Math.random() * N))
  const queryIndexes = Array.from({ length: QUERY_TIMES }, () => Math.floor(Math.random() * N))

  // for (const i of updateIndexes) {
  //   bitset.set(i)
  //   segtree.set(i)
  // }

  // for (let i = 0; i < QUERY_TIMES; i++) {
  //   const index = queryIndexes[i]
  //   const bitsetResult = bitset.query(index)
  //   const segtreeResult = segtree.indexOf(1, index + 1)
  //   if (bitsetResult === -1 && segtreeResult === -1) continue
  //   if (bitsetResult !== segtreeResult - 1) {
  //     console.log(i)
  //     console.log('Wrong answer')
  //     console.log('index', index)
  //     console.log('bitsetResult', bitsetResult)
  //     console.log('segtreeResult', segtreeResult)
  //     break
  //   }
  // }

  console.time('bitset')
  for (let i = 0; i < UPDATE_TIMES; i++) {
    const x = updateIndexes[i]
    bitset.add(x)
  }
  for (let i = 0; i < QUERY_TIMES; i++) {
    const x = queryIndexes[i]
    bitset.indexOfOne(x)
  }
  console.timeEnd('bitset')

  console.time('segtree')
  for (let i = 0; i < UPDATE_TIMES; i++) {
    const x = updateIndexes[i]
    segtree.flip(x + 1, x + 1)
  }
  for (let i = 0; i < QUERY_TIMES; i++) {
    const x = queryIndexes[i]
    segtree.indexOf(1, x + 1)
  }
  console.timeEnd('segtree')
}

export {}
