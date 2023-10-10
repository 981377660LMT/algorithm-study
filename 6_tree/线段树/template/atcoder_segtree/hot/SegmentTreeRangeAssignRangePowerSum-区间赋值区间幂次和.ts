/* eslint-disable no-inner-declarations */

type BuildPowerSum<
  Pow extends number,
  Res extends readonly number[] = []
> = Res['length'] extends Pow ? [...Res, number] : BuildPowerSum<Pow, [...Res, number]>

const INF = 2e9 // !超过int32使用2e15

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
    const lazy = new Float64Array(size).fill(INF)
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
    const n = this._pow + 1
    const leftRes = Array(n).fill(0)
    const rightRes = Array(n).fill(0)
    for (; start < end; start >>= 1, end >>= 1) {
      if (start & 1) {
        const offset = start * n
        for (let i = 0; i < n; i++) leftRes[i] += this._data[offset + i]
        start++
      }
      if (end & 1) {
        end--
        const offset = end * n
        for (let i = 0; i < n; i++) rightRes[i] += this._data[offset + i]
      }
    }
    for (let i = 0; i < n; i++) leftRes[i] += rightRes[i]
    return leftRes as BuildPowerSum<Pow>
  }

  queryAll(): BuildPowerSum<Pow> {
    const res = Array(this._pow + 1)
    const offset = this._pow + 1
    for (let i = 0; i < res.length; i++) res[i] = this._data[offset + i]
    return res as BuildPowerSum<Pow>
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
    const n = this._pow + 1
    let res = Array(n).fill(0)
    while (true) {
      while (!(start & 1)) start >>= 1
      const tmp1 = res.slice()
      const offset1 = start * n
      for (let i = 0; i < n; i++) tmp1[i] += this._data[offset1 + i]
      if (!predicate(tmp1 as BuildPowerSum<Pow>)) {
        while (start < this._size) {
          this._pushDown(start)
          start <<= 1
          const tmp2 = res.slice()
          const offset2 = start * n
          for (let i = 0; i < n; i++) tmp2[i] += this._data[offset2 + i]
          if (predicate(tmp2 as BuildPowerSum<Pow>)) {
            res = tmp2
            start++
          }
        }
        return start - this._size
      }
      const offset3 = start * n
      for (let i = 0; i < n; i++) res[i] += this._data[offset3 + i]
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
    const n = this._pow + 1
    let res = Array(n).fill(0)
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      const tmp1 = res.slice()
      const offset1 = end * n
      for (let i = 0; i < n; i++) tmp1[i] += this._data[offset1 + i]
      if (!predicate(tmp1 as BuildPowerSum<Pow>)) {
        while (end < this._size) {
          this._pushDown(end)
          end = (end << 1) | 1
          const tmp2 = res.slice()
          const offset2 = end * n
          for (let i = 0; i < n; i++) tmp2[i] += this._data[offset2 + i]
          if (predicate(tmp2 as BuildPowerSum<Pow>)) {
            res = tmp2
            end--
          }
        }
        return end + 1 - this._size
      }
      for (let i = 0; i < n; i++) res[i] += this._data[offset1 + i]
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
    for (let i = 0; i < this._n; i++) this._set(i + this._size, leaves[i])
    for (let i = this._size - 1; i > 0; i--) this._pushUp(i)
  }

  private _set(index: number, value: number): void {
    const n = this._pow + 1
    const offset = index * n
    let mul = 1
    for (let i = offset; i < offset + n; i++) {
      this._data[i] = mul
      mul *= value
    }
  }

  private _pushUp(index: number): void {
    const n = this._pow + 1
    const offset1 = index * n
    const offset2 = (index << 1) * n
    const offset3 = ((index << 1) | 1) * n
    for (let i = 0; i < n; i++) {
      this._data[offset1 + i] = this._data[offset2 + i] + this._data[offset3 + i]
    }
  }

  private _pushDown(index: number): void {
    const lazy = this._lazy[index]
    if (lazy === INF) return
    this._propagate(index << 1, lazy)
    this._propagate((index << 1) | 1, lazy)
    this._lazy[index] = INF
  }

  private _propagate(index: number, lazy: number): void {
    const n = this._pow + 1
    const offset = index * n
    const first = this._data[offset]
    let mul = 1
    for (let i = offset; i < offset + n; i++) {
      this._data[i] = first * mul
      mul *= lazy
    }
    if (index < this._size && lazy !== INF) this._lazy[index] = lazy
  }
}

export { SegmentTreeRangeAssignRangePowerSum }

if (require.main === module) {
  const powerSum = new SegmentTreeRangeAssignRangePowerSum(3, 2)
  console.log(powerSum.toString())
  powerSum.update(0, 3, 2)
  console.log(powerSum.queryAll())
  console.log(powerSum.toString())

  class Mocker {
    readonly _n: number
    private readonly _a: number[]
    constructor(nums: number[]) {
      this._n = nums.length
      this._a = nums.slice()
    }

    set(index: number, value: number): void {
      this._a[index] = value
    }

    get(index: number): [number, number] {
      return [this._a[index], this._a[index] * this._a[index]]
    }

    update(start: number, end: number, lazy: number): void {
      this._a.fill(lazy, start, end)
    }

    query(start: number, end: number): [number, number] {
      let sum1 = 0
      let sum2 = 0
      for (let i = start; i < end; i++) {
        sum1 += this._a[i]
        sum2 += this._a[i] * this._a[i]
      }
      return [sum1, sum2]
    }

    queryAll(): [number, number] {
      return this.query(0, this._n)
    }

    maxRight(start: number, predicate: (a: [number, number]) => boolean): number {
      let sum1 = 0
      let sum2 = 0
      for (let i = start; i < this._n; i++) {
        sum1 += this._a[i]
        sum2 += this._a[i] * this._a[i]
        if (!predicate([sum1, sum2])) return i
      }
      return this._n
    }

    minLeft(end: number, predicate: (a: [number, number]) => boolean): number {
      let sum1 = 0
      let sum2 = 0
      for (let i = end - 1; i >= 0; i--) {
        sum1 += this._a[i]
        sum2 += this._a[i] * this._a[i]
        if (!predicate([sum1, sum2])) return i + 1
      }
      return 0
    }

    build(leaves: ArrayLike<number>): void {
      for (let i = 0; i < this._a.length; i++) this._a[i] = leaves[i]
    }

    toString(): string {
      return `Mocker(${this._a})`
    }
  }

  // TODO:FIXME
  checkWithBruteForce()
  function assertSame(obj1: unknown, obj2: unknown) {
    if (JSON.stringify(obj1) !== JSON.stringify(obj2)) {
      throw new Error(`expect ${JSON.stringify(obj2)}, got ${JSON.stringify(obj1)}`)
    }
  }

  function checkWithBruteForce(): void {
    const randint = (min: number, max: number) => Math.floor(Math.random() * (max - min + 1)) + min
    const N = 2e4
    const real = new SegmentTreeRangeAssignRangePowerSum(N, 2)
    const mock = new Mocker(Array(N).fill(0))
    for (let i = 0; i < N; i++) {
      const op = randint(0, 6)
      if (op === 0) {
        // set
        const index = randint(0, N - 1)
        const value = randint(0, 10)
        real.set(index, value)
        mock.set(index, value)
        // console.log('set', index, value)
      } else if (op === 1) {
        // get
        const index = randint(0, N - 1)
        const realValue = real.get(index)
        const mockValue = mock.get(index)
        // console.log(realValue, mockValue, index)
        // console.log('get', index, realValue, mockValue)
        assertSame(realValue.slice(1), mockValue)
      } else if (op === 2) {
        // update
        const start = randint(0, N - 1)
        const end = randint(start, N)
        const lazy = randint(0, 2)
        real.update(start, end, lazy)
        mock.update(start, end, lazy)
        // console.log('update', start, end, lazy)
      } else if (op === 3) {
        // query
        const start = randint(0, N - 1)
        const end = randint(start, N)
        const realValue = real.query(start, end)
        const mockValue = mock.query(start, end)
        // console.log('query', start, end, realValue, mockValue)
        assertSame(realValue.slice(1), mockValue)
      } else if (op === 4) {
        // queryAll
        const realValue = real.queryAll()
        const mockValue = mock.queryAll()
        assertSame(realValue.slice(1), mockValue)
      } else if (op === 5) {
        // maxRight
        const start = randint(0, N - 1)
        const target = randint(0, N)
        const realValue = real.maxRight(start, ([_, sum1, sum2]) => sum2 + sum1 <= target)
        const mockValue = mock.maxRight(start, ([sum1, sum2]) => sum2 + sum1 <= target)
        assertSame(realValue, mockValue)
      } else if (op === 6) {
        // minLeft
        const end = randint(0, N)
        const target = randint(0, N)
        const realValue = real.minLeft(end, ([_, sum1, sum2]) => sum2 + sum1 <= target)
        const mockValue = mock.minLeft(end, ([sum1, sum2]) => sum2 + sum1 <= target)
        assertSame(realValue, mockValue)
      }
    }
  }

  console.log('All tests passed')
}
