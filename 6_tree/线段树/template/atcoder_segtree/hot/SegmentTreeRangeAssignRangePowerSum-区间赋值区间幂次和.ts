type BuildPowerSum<Pow extends number, Res extends readonly number[] = []> = Res['length'] extends Pow
  ? [...Res, number]
  : BuildPowerSum<Pow, [...Res, number]>

/**
 * 区间赋值，区间幂次和.
 */
class SegmentTreeRangeAssignRangePowerSum<Pow extends number = 2> {
  private readonly _n: number
  private readonly _size: number
  private readonly _height: number
  private readonly _data: Float64Array // size<<1 个 [0]*(Pow+1) 数组连接而成
  private readonly _lazy: Float64Array
  private readonly _pow: Pow

  /**
   * 区间赋值，区间幂次和的线段树.
   * @param nOrArr 数组长度或数组.
   * @param pow 维护区间 0, 1, ..., pow 次幂和.
   */
  constructor(nOrArr: number | ArrayLike<number>, pow: Pow = 2 as Pow) {
    const n = typeof nOrArr === 'number' ? nOrArr : nOrArr.length
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
    const data = new Float64Array((size << 1) * (pow + 1))
    const lazy = new Float64Array(size)
    this._data = data
    this._lazy = lazy
    this._pow = pow

    if (typeof nOrArr === 'number') nOrArr = new Uint8Array(nOrArr)
    this._build(nOrArr)
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    this._set(index, value)
    for (let i = 1; i <= this._height; i++) this._pushUp(index >> i)
  }

  get(index: number): BuildPowerSum<Pow> {
    if (index < 0 || index >= this._n) {
      throw new RangeError(`index must be in [0, ${this._n})`)
    }
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    const res = Array(this._pow + 1)
    const start = index * (this._pow + 1)
    for (let i = 0; i < res.length; i++) res[i] = this._data[start + i]
    return res as BuildPowerSum<Pow>
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
  query(start: number, end: number): BuildPowerSum<Pow> {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return Array(this._pow + 1).fill(0) as BuildPowerSum<Pow>
    start += this._size
    end += this._size
    for (let i = this._height; i > 0; i--) {
      if ((start >> i) << i !== start) this._pushDown(start >> i)
      if ((end >> i) << i !== end) this._pushDown((end - 1) >> i)
    }
    const pow = this._pow
    const leftRes = Array(pow + 1)
    for (; start < end; start >>= 1, end >>= 1) {
      if (start & 1) {
        leftSum1 += this._sum1[start]
        leftSum2 += this._sum2[start]
        start++
      }
      if (end & 1) {
        end--
        rightSum1 += this._sum1[end]
        rightSum2 += this._sum2[end]
      }
    }
    return { sum1: leftSum1 + rightSum1, sum2: leftSum2 + rightSum2 }
  }

  queryAll(): BuildPowerSum<Pow> {
    return { sum1: this._sum1[1], sum2: this._sum2[1] }
  }

  /**
   * 树上二分查询最大的`end`使得`[start,end)`内的值满足`predicate`.
   * @alias findFirst
   */
  maxRight(start: number, predicate: (powerSum: BuildPowerSum<Pow>) => boolean): number {
    if (start < 0) start = 0
    if (start >= this._n) return this._n
    start += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(start >> i)
    let resSum1 = 0
    let resSum2 = 0
    while (true) {
      while (!(start & 1)) start >>= 1
      const tmpSum11 = resSum1 + this._sum1[start]
      const tmpSum21 = resSum2 + this._sum2[start]
      if (!predicate(tmpSum11, tmpSum21)) {
        while (start < this._size) {
          this._pushDown(start)
          start <<= 1
          const tmpSum12 = resSum1 + this._sum1[start]
          const tmpSum22 = resSum2 + this._sum2[start]
          if (predicate(tmpSum12, tmpSum22)) {
            resSum1 = tmpSum12
            resSum2 = tmpSum22
            start++
          }
        }
        return start - this._size
      }
      resSum1 += this._sum1[start]
      resSum2 += this._sum2[start]
      start++
      if ((start & -start) === start) break
    }
    return this._n
  }

  /**
   * 树上二分查询最小的`start`使得`[start,end)`内的值满足`predicate`.
   * @alias findLast
   */
  minLeft(end: number, predicate: (powerSum: BuildPowerSum<Pow>) => boolean): number {
    if (end > this._n) end = this._n
    if (end <= 0) return 0
    end += this._size
    for (let i = this._height; i > 0; i--) this._pushDown((end - 1) >> i)
    let resSum1 = 0
    let resSum2 = 0
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      const tmpSum11 = resSum1 + this._sum1[end]
      const tmpSum21 = resSum2 + this._sum2[end]
      if (!predicate(tmpSum11, tmpSum21)) {
        while (end < this._size) {
          this._pushDown(end)
          end = (end << 1) | 1
          const tmpSum12 = resSum1 + this._sum1[end]
          const tmpSum22 = resSum2 + this._sum2[end]
          if (predicate(tmpSum12, tmpSum22)) {
            resSum1 = tmpSum12
            resSum2 = tmpSum22
            end--
          }
        }
        return end + 1 - this._size
      }
      resSum1 += this._sum1[end]
      resSum2 += this._sum2[end]
      if ((end & -end) === end) break
    }
    return 0
  }

  toString(): string {
    const sb: string[] = []
    sb.push('SegmentTreeRangeAssignRangePowerSum(')
    for (let i = 0; i < this._n; i++) {
      if (i) sb.push(', ')
      sb.push(JSON.stringify(this.get(i)))
    }
    sb.push(')')
    return sb.join('')
  }

  private _build(leaves: ArrayLike<number>): void {
    if (leaves.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < this._n; i++) {
      this._sum0[this._size + i] = 1
      this._sum1[this._size + i] = leaves[i]
      this._sum2[this._size + i] = leaves[i] * leaves[i]
    }
    for (let i = this._size - 1; i > 0; i--) this._pushUp(i)
  }

  private _set(index: number, value: number): void {
    let cur = 1
    const pow = this._pow
    const start = index * (pow + 1)
    for (let i = start; i < start + pow + 1; i++) {
      this._data[i] = cur
      cur *= value
    }
  }

  private _pushUp(index: number): void {
    this._sum0[index] = this._sum0[index << 1] + this._sum0[(index << 1) | 1]
    this._sum1[index] = this._sum1[index << 1] + this._sum1[(index << 1) | 1]
    this._sum2[index] = this._sum2[index << 1] + this._sum2[(index << 1) | 1]
  }

  private _pushDown(index: number): void {
    const lazy = this._lazy[index]
    if (!lazy) return
    this._propagate(index << 1, lazy)
    this._propagate((index << 1) | 1, lazy)
    this._lazy[index] = 0
  }

  private _propagate(index: number, lazy: number): void {
    this._sum2[index] += 2 * lazy * this._sum1[index] + lazy * lazy * this._sum0[index]
    this._sum1[index] += lazy * this._sum0[index]
    if (index < this._size) this._lazy[index] += lazy
  }
}

export { SegmentTreeRangeAssignRangePowerSum }
