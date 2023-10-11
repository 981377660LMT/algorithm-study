// https://leetcode.cn/problems/fancy-sequence/

/* eslint-disable no-param-reassign */

// !区间修改+单点查询 => DualSegmentTree
// Update(start, end, value) => [start, end)的值都与value结合.
// Get(index) => index的值

class SegmentTreeRangeUpdatePointGet<Id = number> {
  private readonly _n: number
  private readonly _size: number
  private readonly _height: number
  private readonly _lazy: Id[]
  private readonly _id: () => Id
  private readonly _composition: (f: Id, g: Id) => Id
  private readonly _equalsToId: (o: Id) => boolean
  private readonly _commutative: boolean

  /**
   * 区间修改,单点查询的线段树.
   * @param n 线段树的大小.
   * @param id 单位元.
   * @param composition 父结点`f`与子结点`g`的结合函数.
   * @param equals 判断两个值是否相等的函数.比较方式默认为`===`.
   * @param commutative 群的结合是否可交换顺序.默认为`false`.为'true'时可以加速区间修改.
   *
   * @alias DualSegmentTree
   */
  constructor(
    n: number,
    id: () => Id,
    composition: (f: Id, g: Id) => Id,
    equals: (a: Id, b: Id) => boolean = (a, b) => a === b,
    commutative = false
  ) {
    if (!equals(id(), id())) {
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
    this._commutative = commutative
  }

  get(index: number): Id {
    if (index < 0 || index >= this._n) return this._id()
    index += this._size
    for (let i = this._height; i > 0; i--) this._propagate(index >> i)
    return this._lazy[index]
  }

  /**
   * 将区间`[left, right)`的值与`lazy`作用.
   */
  update(start: number, end: number, lazy: Id): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    start += this._size
    end += this._size
    if (!this._commutative) {
      for (let i = this._height; i > 0; i--) {
        if ((start >> i) << i !== start) this._propagate(start >> i)
        if ((end >> i) << i !== end) this._propagate((end - 1) >> i)
      }
    }
    while (start < end) {
      if (start & 1) {
        this._lazy[start] = this._composition(lazy, this._lazy[start])
        start++
      }
      if (end & 1) {
        end--
        this._lazy[end] = this._composition(lazy, this._lazy[end])
      }
      start >>= 1
      end >>= 1
    }
  }

  toString(): string {
    const sb: string[] = []
    sb.push('SegmentTreeRangeUpdatePointGet(')
    for (let i = 0; i < this._size; i++) {
      this._propagate(i)
    }
    for (let i = this._size; i < this._size + this._n; i++) {
      if (i !== this._size) sb.push(',')
      sb.push(String(this._lazy[i]))
    }
    sb.push(')')
    return sb.join('')
  }

  private _propagate(k: number): void {
    if (this._equalsToId(this._lazy[k])) return
    this._lazy[k << 1] = this._composition(this._lazy[k], this._lazy[k << 1])
    this._lazy[(k << 1) | 1] = this._composition(this._lazy[k], this._lazy[(k << 1) | 1])
    this._lazy[k] = this._id()
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

  // eslint-disable-next-line no-inner-declarations
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

  // 1622. 奇妙序列
  // 区间修改,单点查询

  const BIGMOD = BigInt(1e9 + 7)

  class Fancy {
    private readonly _seg: SegmentTreeRangeUpdatePointGet<[mul: bigint, add: bigint]> = new SegmentTreeRangeUpdatePointGet(
      1e5 + 10,
      () => [1n, 0n],
      (f, g) => [(f[0] * g[0]) % BIGMOD, (f[0] * g[1] + f[1]) % BIGMOD],
      (a, b) => a[0] === b[0] && a[1] === b[1]
    )
    private _length = 0

    append(val: number): void {
      this._seg.update(this._length, this._length + 1, [1n, BigInt(val)])
      this._length++
    }

    addAll(inc: number): void {
      this._seg.update(0, this._length, [1n, BigInt(inc)])
    }

    multAll(m: number): void {
      this._seg.update(0, this._length, [BigInt(m), 0n])
    }

    getIndex(idx: number): number {
      if (idx >= this._length) return -1
      return Number(this._seg.get(idx)[1])
    }
  }
}
