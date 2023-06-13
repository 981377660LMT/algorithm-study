/* eslint-disable no-cond-assign */
/* eslint-disable no-param-reassign */

// !单点修改+区间查询

class SegmentTreePointUpdateRangeQuery<E = number> {
  private readonly _n: number
  private readonly _size: number
  private readonly _seg: E[]
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E

  /**
   * 单点更新,区间查询的线段树.
   * @param nOrLeaves 大小或叶子节点的值.
   * @param e 幺元.
   * @param op 结合律.
   */
  constructor(nOrLeaves: number | ArrayLike<E>, e: () => E, op: (a: E, b: E) => E) {
    const n = typeof nOrLeaves === 'number' ? nOrLeaves : nOrLeaves.length
    let size = 1
    while (size < n) size <<= 1
    const seg = Array(size << 1)
    for (let i = 0; i < seg.length; i++) seg[i] = e()

    this._n = n
    this._size = size
    this._seg = seg
    this._e = e
    this._op = op

    if (typeof nOrLeaves !== 'number') this.build(nOrLeaves)
  }

  set(index: number, value: E): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    this._seg[index] = value
    while ((index >>= 1)) {
      this._seg[index] = this._op(this._seg[index << 1], this._seg[(index << 1) | 1])
    }
  }

  get(index: number): E {
    if (index < 0 || index >= this._n) return this._e()
    return this._seg[index + this._size]
  }

  /**
   * 将`index`处的值与作用素`value`结合.
   */
  update(index: number, value: E): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    this._seg[index] = this._op(this._seg[index], value)
    while ((index >>= 1)) {
      this._seg[index] = this._op(this._seg[index << 1], this._seg[(index << 1) | 1])
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

    let leftRes = this._e()
    let rightRes = this._e()
    for (start += this._size, end += this._size; start < end; start >>= 1, end >>= 1) {
      if (start & 1) leftRes = this._op(leftRes, this._seg[start++])
      if (end & 1) rightRes = this._op(this._seg[--end], rightRes)
    }
    return this._op(leftRes, rightRes)
  }

  queryAll(): E {
    return this._seg[1]
  }

  /**
   * 树上二分查询最大的`end`使得`[start,end)`内的值满足`predicate`.
   */
  maxRight(start: number, predicate: (value: E) => boolean): number {
    if (start === this._n) return this._n
    start += this._size
    let res = this._e()
    while (true) {
      while (!(start & 1)) start >>= 1
      if (!predicate(this._op(res, this._seg[start]))) {
        while (start < this._size) {
          start <<= 1
          if (predicate(this._op(res, this._seg[start]))) {
            res = this._op(res, this._seg[start])
            start++
          }
        }
        return start - this._size
      }
      res = this._op(res, this._seg[start])
      start++
      if ((start & -start) === start) break
    }
    return this._n
  }

  /**
   * 树上二分查询最小的`start`使得`[start,end)`内的值满足`predicate`
   */
  minLeft(end: number, predicate: (value: E) => boolean): number {
    if (!end) return 0
    end += this._size
    let res = this._e()
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      if (!predicate(this._op(this._seg[end], res))) {
        while (end < this._size) {
          end = (end << 1) | 1
          if (predicate(this._op(this._seg[end], res))) {
            res = this._op(this._seg[end], res)
            end--
          }
        }
        return end + 1 - this._size
      }
      res = this._op(this._seg[end], res)
      if ((end & -end) === end) break
    }
    return 0
  }

  build(arr: ArrayLike<E>): void {
    if (arr.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < arr.length; i++) {
      this._seg[i + this._size] = arr[i] // 叶子结点
    }
    for (let i = this._size - 1; ~i; i--) {
      this._seg[i] = this._op(this._seg[i << 1], this._seg[(i << 1) | 1])
    }
  }

  toString(): string {
    const sb: string[] = []
    sb.push('SegmentTreePointUpdateRangeQuery(')
    for (let i = 0; i < this._n; i++) {
      if (i) sb.push(', ')
      sb.push(JSON.stringify(this.get(i)))
    }
    sb.push(')')
    return sb.join('')
  }
}

export { SegmentTreePointUpdateRangeQuery }

if (require.main === module) {
  const seg = new SegmentTreePointUpdateRangeQuery(
    10,
    () => 0,
    (a, b) => a + b
  )
  console.log(seg.toString())
  seg.set(0, 1)
  seg.set(1, 2)
  console.log(seg.toString())
  seg.update(3, 4)
  console.log(seg.toString())
  console.log(seg.query(0, 4))
  seg.build([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
  console.log(seg.toString())
  console.log(seg.minLeft(10, x => x < 15))
  console.log(seg.maxRight(0, x => x <= 15))
  console.log(seg.queryAll())
}
