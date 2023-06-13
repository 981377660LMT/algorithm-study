/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */

// !区间修改+单点查询 => DualSegmentTree
// Update(start, end, value) => [start, end)的值都与value结合.
// Get(index) => index的值

/**
 * @warning Not Verified.
 */
class SegmentTreeRangeUpdatePointGet<Id = number> {
  private readonly _n: number
  private readonly _size: number
  private readonly _height: number
  private readonly _lazy: Id[]
  private readonly _id: () => Id
  private readonly _composition: (f: Id, g: Id) => Id
  private readonly _equalsToId: (o: Id) => boolean

  /**
   * 区间修改,单点查询的线段树.
   * @param n 线段树的大小.
   * @param id 单位元.
   * @param composition 父结点`f`与子结点`g`的合成函数.
   * @param equals 判断两个值是否相等的函数.默认为`===`.
   *
   * @alias DualSegmentTree
   */
  constructor(
    n: number,
    id: () => Id,
    composition: (f: Id, g: Id) => Id,
    equals: (a: Id, b: Id) => boolean = (a, b) => a === b
  ) {
    if (!equals && !SegmentTreeRangeUpdatePointGet._isPrimitive(id())) {
      throw new Error('equals must be provided when id() returns an non-primitive value')
    }

    let size = 1
    let height = 0
    while (size < n) {
      size <<= 1
      height++
    }
    const lazy = Array(size << 1)
    for (let i = 0; i < lazy.length; i++) lazy[i] = id()
    this._n = n
    this._size = size
    this._height = height
    this._lazy = lazy
    this._id = id
    this._composition = composition

    const identity = id()
    this._equalsToId = equals ? (o: Id) => equals(o, identity) : (o: Id) => o === identity
  }

  get(index: number): Id {
    if (index < 0 || index >= this._n) return this._id()
    this._thrust((index += this._size))
    return this._lazy[index]
  }

  /**
   * 将区间`[left, right)`的值与`lazy`作用.
   */
  update(start: number, end: number, lazy: Id): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    this._thrust((start += this._size))
    this._thrust((end += this._size - 1))
    for (let l = start, r = end + 1; l < r; l >>= 1, r >>= 1) {
      if (l & 1) {
        this._lazy[l] = this._composition(lazy, this._lazy[l])
        l++
      }
      if (r & 1) {
        r--
        this._lazy[r] = this._composition(lazy, this._lazy[r])
      }
    }
  }

  toString(): string {
    const sb: string[] = []
    sb.push('SegmentTreeRangeUpdatePointGet(')
    for (let i = 0; i < this._n; i++) {
      if (i) sb.push(', ')
      sb.push(JSON.stringify(this.get(i)))
    }
    sb.push(')')
    return sb.join('')
  }

  private _thrust(k: number): void {
    for (let i = this._height; i > 0; i--) this._propagate(k >> i)
  }

  private _propagate(k: number): void {
    if (this._equalsToId(this._lazy[k])) return
    this._lazy[k << 1] = this._composition(this._lazy[k], this._lazy[k << 1])
    this._lazy[(k << 1) | 1] = this._composition(this._lazy[k], this._lazy[(k << 1) | 1])
    this._lazy[k] = this._id()
  }

  private static _isPrimitive(
    o: unknown
  ): o is number | string | boolean | symbol | bigint | null | undefined {
    return o === null || (typeof o !== 'object' && typeof o !== 'function')
  }
}

export { SegmentTreeRangeUpdatePointGet }

if (require.main === module) {
  // 区间染色(区间赋值),单点查询.
  const seg = new SegmentTreeRangeUpdatePointGet<number>(
    10,
    () => -1,
    (parent, child) => parent
  )
  seg.update(0, 8, 1)
  console.log(seg.get(0))
  console.log(seg.toString())

  benchMark()

  function benchMark(): void {
    const n = 2e5
    const seg = new SegmentTreeRangeUpdatePointGet<number>(
      n,
      () => 0,
      (parent, child) => parent + child
    )
    console.time('update')
    for (let i = 0; i < n; i++) {
      seg.update(0, i, 1)
      seg.get(i)
    }
    console.timeEnd('update')
  }
}
