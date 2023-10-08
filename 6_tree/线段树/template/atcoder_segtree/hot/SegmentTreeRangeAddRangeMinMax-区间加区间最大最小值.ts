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
    // !1. set
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
  update(start: number, end: number, lazy: number): void {
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
  query(start: number, end: number): { min: number; max: number } {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return { min: INF, max: -INF }
    start += this._size
    end += this._size
    for (let i = this._height; i > 0; i--) {
      if ((start >> i) << i !== start) this._pushDown(start >> i)
      if ((end >> i) << i !== end) this._pushDown((end - 1) >> i)
    }
    let leftMin = INF
    let leftMax = -INF
    let rightMin = INF
    let rightMax = -INF
    for (; start < end; start >>= 1, end >>= 1) {
      if (start & 1) {
        leftMin = Math.min(leftMin, this._min[start])
        leftMax = Math.max(leftMax, this._max[start])
        start++
      }
      if (end & 1) {
        end--
        rightMin = Math.min(rightMin, this._min[end])
        rightMax = Math.max(rightMax, this._max[end])
      }
    }
    return { min: Math.min(leftMin, rightMin), max: Math.max(leftMax, rightMax) }
  }

  queryAll(): { min: number; max: number } {
    return { min: this._min[1], max: this._max[1] }
  }

  /**
   * 树上二分查询最大的`end`使得`[start,end)`内的值满足`predicate`.
   * @alias findFirst
   */
  maxRight(start: number, predicate: (min: number, max: number) => boolean): number {
    if (start < 0) start = 0
    if (start >= this._n) return this._n
    start += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(start >> i)
    let resMin = INF
    let resMax = -INF

    while (true) {
      while (!(start & 1)) start >>= 1
      const tmpMin1 = Math.min(resMin, this._min[start])
      const tmpMax1 = Math.max(resMax, this._max[start])
      if (!predicate(tmpMin1, tmpMax1)) {
        while (start < this._size) {
          this._pushDown(start)
          start <<= 1
          const tmpMin2 = Math.min(resMin, this._min[start])
          const tmpMax2 = Math.max(resMax, this._max[start])
          if (predicate(tmpMin2, tmpMax2)) {
            resMin = tmpMin2
            resMax = tmpMax2
            start++
          }
        }
        return start - this._size
      }
      resMin = Math.min(resMin, this._min[start])
      resMax = Math.max(resMax, this._max[start])
      start++
      if ((start & -start) === start) break
    }
    return this._n
  }

  /**
   * 树上二分查询最小的`start`使得`[start,end)`内的值满足`predicate`
   * @alias findLast
   */
  minLeft(end: number, predicate: (min: number, max: number) => boolean): number {
    if (end > this._n) end = this._n
    if (end <= 0) return 0
    end += this._size
    for (let i = this._height; i > 0; i--) this._pushDown((end - 1) >> i)
    let resMin = INF
    let resMax = -INF
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      const tmpMin1 = Math.min(resMin, this._min[end])
      const tmpMax1 = Math.max(resMax, this._max[end])
      if (!predicate(tmpMin1, tmpMax1)) {
        while (end < this._size) {
          this._pushDown(end)
          end = (end << 1) | 1
          const tmpMin2 = Math.min(resMin, this._min[end])
          const tmpMax2 = Math.max(resMax, this._max[end])
          if (predicate(tmpMin2, tmpMax2)) {
            resMin = tmpMin2
            resMax = tmpMax2
            end--
          }
        }
        return end + 1 - this._size
      }
      resMin = Math.min(resMin, this._min[end])
      resMax = Math.max(resMax, this._max[end])
      if ((end & -end) === end) break
    }
    return 0
  }

  build(leaves: ArrayLike<number>): void {
    if (leaves.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < this._n; i++) {
      this._min[this._size + i] = leaves[i]
      this._max[this._size + i] = leaves[i]
    }
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
    this._min[index] = Math.min(this._min[index << 1], this._min[(index << 1) | 1])
    this._max[index] = Math.max(this._max[index << 1], this._max[(index << 1) | 1])
  }

  private _pushDown(index: number): void {
    const lazy = this._lazy[index]
    if (!lazy) return
    this._propagate(index << 1, lazy)
    this._propagate((index << 1) | 1, lazy)
    this._lazy[index] = 0
  }

  private _propagate(index: number, lazy: number): void {
    this._min[index] += lazy
    this._max[index] += lazy
    if (index < this._size) this._lazy[index] += lazy
  }
}

export { SegmentTreeRangeAddRangeMinMax }
