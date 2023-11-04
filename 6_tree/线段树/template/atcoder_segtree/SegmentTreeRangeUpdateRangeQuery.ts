/* eslint-disable no-cond-assign */
/* eslint-disable no-param-reassign */
// !区间修改 + 区间查询

class SegmentTreeRangeUpdateRangeQuery<E = number, Id = number> {
  private readonly _n: number
  private readonly _size: number
  private readonly _height: number
  private readonly _data: E[]
  private readonly _lazy: Id[]
  private readonly _e: () => E
  private readonly _id: () => Id
  private readonly _op: (a: E, b: E) => E
  private readonly _mapping: (id: Id, data: E) => E
  private readonly _composition: (id1: Id, id2: Id) => Id
  private readonly _equalsToId: (o: Id) => boolean

  /**
   * 区间修改区间查询的懒标记线段树.维护幺半群.
   * @param nOrLeaves 大小或叶子节点的值.
   * @param operations 线段树的操作.
   */
  constructor(
    nOrLeaves: number | ArrayLike<E>,
    operations: {
      /**
       * 线段树维护的值的幺元.
       */
      e: () => E

      /**
       * 更新操作/懒标记的幺元.
       */
      id: () => Id

      /**
       * 合并左右区间的值.
       */
      op: (e1: E, e2: E) => E

      /**
       * 父结点的懒标记更新子结点的值.
       */
      mapping: (lazy: Id, data: E) => E

      /**
       * 父结点的懒标记更新子结点的懒标记(合并).
       */
      composition: (f: Id, g: Id) => Id

      /**
       * 判断两个懒标记是否相等.比较方式默认为`===`.
       */
      equalsId?: (id1: Id, id2: Id) => boolean
    } & ThisType<void>
  ) {
    const n = typeof nOrLeaves === 'number' ? nOrLeaves : nOrLeaves.length
    const { e, id, op, mapping, composition, equalsId } = operations
    if (!equalsId && !SegmentTreeRangeUpdateRangeQuery._isPrimitive(id())) {
      throw new Error('equalsId must be provided when id() returns an non-primitive value')
    }

    let size = 1
    let height = 0
    while (size < n) {
      size <<= 1
      height++
    }
    const data = Array(size << 1)
    for (let i = 0; i < data.length; i++) data[i] = e()
    const lazy = Array(size)
    for (let i = 0; i < lazy.length; i++) lazy[i] = id()
    this._n = n
    this._size = size
    this._height = height
    this._data = data
    this._lazy = lazy
    this._e = e
    this._id = id
    this._op = op
    this._mapping = mapping
    this._composition = composition

    const identity = id()
    this._equalsToId = equalsId ? (o: Id) => equalsId(o, identity) : (o: Id) => o === identity

    if (typeof nOrLeaves !== 'number') this._build(nOrLeaves)
  }

  set(index: number, value: E): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    this._data[index] = value
    for (let i = 1; i <= this._height; i++) this._pushUp(index >> i)
  }

  get(index: number): E {
    if (index < 0 || index >= this._n) return this._e()
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    return this._data[index]
  }

  /**
   * 区间`[start,end)`的值与`lazy`进行作用.
   * 0 <= start <= end <= n.
   */
  update(start: number, end: number, lazy: Id): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    start += this._size
    end += this._size
    for (let i = this._height; i > 0; i--) {
      if ((start >> i) << i !== start) this._pushDown(start >> i)
      if ((end >> i) << i !== end) this._pushDown((end - 1) >> i)
    }
    let start2 = start
    let end2 = end
    for (; start < end; start >>= 1, end >>= 1) {
      if (start & 1) this._propagate(start++, lazy)
      if (end & 1) this._propagate(--end, lazy)
    }
    start = start2
    end = end2
    for (let i = 1; i <= this._height; i++) {
      if ((start >> i) << i !== start) this._pushUp(start >> i)
      if ((end >> i) << i !== end) this._pushUp((end - 1) >> i)
    }
  }

  /**
   * 查询区间`[start,end)`的聚合值.
   * 0 <= start <= end <= n.
   */
  query(start: number, end: number): E {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return this._e()
    start += this._size
    end += this._size
    for (let i = this._height; i > 0; i--) {
      if ((start >> i) << i !== start) this._pushDown(start >> i)
      if ((end >> i) << i !== end) this._pushDown((end - 1) >> i)
    }
    let leftRes = this._e()
    let rightRes = this._e()
    for (; start < end; start >>= 1, end >>= 1) {
      if (start & 1) leftRes = this._op(leftRes, this._data[start++])
      if (end & 1) rightRes = this._op(this._data[--end], rightRes)
    }
    return this._op(leftRes, rightRes)
  }

  queryAll(): E {
    return this._data[1]
  }

  /**
   * 树上二分查询最大的`end`使得`[start,end)`内的值满足`predicate`.
   * @alias findFirst
   */
  maxRight(start: number, predicate: (value: E) => boolean): number {
    if (start < 0) start = 0
    if (start >= this._n) return this._n
    start += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(start >> i)
    let res = this._e()
    while (true) {
      while (!(start & 1)) start >>= 1
      if (!predicate(this._op(res, this._data[start]))) {
        while (start < this._size) {
          this._pushDown(start)
          start <<= 1
          const tmp = this._op(res, this._data[start])
          if (predicate(tmp)) {
            res = tmp
            start++
          }
        }
        return start - this._size
      }
      res = this._op(res, this._data[start])
      start++
      if ((start & -start) === start) break
    }
    return this._n
  }

  /**
   * 树上二分查询最小的`start`使得`[start,end)`内的值满足`predicate`
   * @alias findLast
   */
  minLeft(end: number, predicate: (value: E) => boolean): number {
    if (end > this._n) end = this._n
    if (end <= 0) return 0
    end += this._size
    for (let i = this._height; i > 0; i--) this._pushDown((end - 1) >> i)
    let res = this._e()
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      if (!predicate(this._op(this._data[end], res))) {
        while (end < this._size) {
          this._pushDown(end)
          end = (end << 1) | 1
          const tmp = this._op(this._data[end], res)
          if (predicate(tmp)) {
            res = tmp
            end--
          }
        }
        return end + 1 - this._size
      }
      res = this._op(this._data[end], res)
      if ((end & -end) === end) break
    }
    return 0
  }

  getAll(): E[] {
    for (let i = 1; i < this._size; i++) {
      this._pushDown(i)
    }
    return this._data.slice(this._size, this._size + this._n)
  }

  toString(): string {
    const sb: string[] = []
    sb.push('SegmentTreeRangeUpdateRangeQuery(')
    for (let i = 0; i < this._n; i++) {
      if (i) sb.push(', ')
      sb.push(JSON.stringify(this.get(i)))
    }
    sb.push(')')
    return sb.join('')
  }

  private _build(leaves: ArrayLike<E>): void {
    if (leaves.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < this._n; i++) this._data[this._size + i] = leaves[i]
    for (let i = this._size - 1; i > 0; i--) this._pushUp(i)
  }

  private _pushUp(index: number): void {
    this._data[index] = this._op(this._data[index << 1], this._data[(index << 1) | 1])
  }

  private _pushDown(index: number): void {
    const lazy = this._lazy[index]
    if (this._equalsToId(lazy)) return
    this._propagate(index << 1, lazy)
    this._propagate((index << 1) | 1, lazy)
    this._lazy[index] = this._id()
  }

  private _propagate(index: number, lazy: Id): void {
    this._data[index] = this._mapping(lazy, this._data[index])
    if (index < this._size) this._lazy[index] = this._composition(lazy, this._lazy[index])
  }

  private static _isPrimitive(
    o: unknown
  ): o is number | string | boolean | symbol | bigint | null | undefined {
    return o === null || (typeof o !== 'object' && typeof o !== 'function')
  }
}

export { SegmentTreeRangeUpdateRangeQuery }

if (require.main === module) {
  const seg = new SegmentTreeRangeUpdateRangeQuery(10, {
    e: () => 0,
    id: () => 0,
    op: (a, b) => Math.max(a, b),
    mapping: (lazy, data) => Math.max(lazy, data),
    composition: (f, g) => Math.max(f, g)
  })

  console.log(seg.getAll())
  seg.update(0, 10, 1)
  console.log(seg.getAll())
}
