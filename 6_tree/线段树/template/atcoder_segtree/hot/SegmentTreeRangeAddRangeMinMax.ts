const INF = 2e15 // !超过int32使用2e15

/**
 * 区间加, 区间最大最小值.
 */
class SegmentTreeRangeAddRangeMinMax {
  private readonly _n: number
  private readonly _size: number
  private readonly _height: number
  private readonly _min: Float64Array
  private readonly _max: Float64Array
  private readonly _lazy: Float64Array

  constructor(nOrLeaves: number | ArrayLike<number>) {
    const n = typeof nOrLeaves === 'number' ? nOrLeaves : nOrLeaves.length
    let size = 1
    let height = 0
    while (size < n) {
      size <<= 1
      height++
    }
    this._n = n
    this._size = size
    this._height = height

    // !0.init data and lazy
    const min = new Float64Array(size << 1).fill(INF)
    const max = new Float64Array(size << 1).fill(-INF)
    const lazy = new Float64Array(size)
    this._min = min
    this._max = max
    this._lazy = lazy

    if (typeof nOrLeaves !== 'number') this.build(nOrLeaves)
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    // !1.
    this._min[index] = value
    this._max[index] = value
    for (let i = 1; i <= this._height; i++) this._pushUp(index >> i)
  }

  get(index: number): number {
    if (index < 0 || index >= this._n) {
      throw new RangeError(`index must be in [0, ${this._n})`)
    }
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    return this._max[index]
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
          if (predicate(this._op(res, this._data[start]))) {
            res = this._op(res, this._data[start])
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
          if (predicate(this._op(this._data[end], res))) {
            res = this._op(this._data[end], res)
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

  build(leaves: ArrayLike<E>): void {
    if (leaves.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < this._n; i++) this._data[this._size + i] = leaves[i]
    for (let i = this._size - 1; i > 0; i--) this._pushUp(i)
  }

  toString(): string {
    const sb: string[] = []
    sb.push('SegmentTreeRangeUpdateRangeQuery(')
    for (let i = 0; i < this._n; i++) {
      if (i) sb.push(', ')
      sb.push(String(this.get(i)))
    }
    sb.push(')')
    return sb.join('')
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

export { SegmentTreeRangeAddRangeMinMax }
