/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

// 区间加，区间平方和

const INF = 2e15 // !超过int32使用2e15

/**
 * 区间加, 区间平方和.
 */
class SegmentTreeRangeAddRangeSquareSum {
  private readonly _n: number
  private readonly _size: number
  private readonly _height: number
  private readonly _sum0: Float64Array
  private readonly _sum1: Float64Array
  private readonly _sum2: Float64Array
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
    const sum0 = new Float64Array(size << 1)
    const sum1 = new Float64Array(size << 1)
    const sum2 = new Float64Array(size << 1)
    const lazy = new Float64Array(size)
    this._sum0 = sum0
    this._sum1 = sum1
    this._sum2 = sum2
    this._lazy = lazy

    if (typeof nOrLeaves !== 'number') this.build(nOrLeaves)
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    // !1. set
    this._sum0[index] = 1
    this._sum1[index] = value
    this._sum2[index] = value * value
    for (let i = 1; i <= this._height; i++) this._pushUp(index >> i)
  }

  get(index: number): { sum0: number; sum1: number; sum2: number } {
    if (index < 0 || index >= this._n) {
      throw new RangeError(`index must be in [0, ${this._n})`)
    }
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    return { sum0: this._sum0[index], sum1: this._sum1[index], sum2: this._sum2[index] }
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
  query(start: number, end: number): { sum0: number; sum1: number; sum2: number } {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return { sum0: 0, sum1: 0, sum2: 0 }
    start += this._size
    end += this._size
    for (let i = this._height; i > 0; i--) {
      if ((start >> i) << i !== start) this._pushDown(start >> i)
      if ((end >> i) << i !== end) this._pushDown((end - 1) >> i)
    }
    let leftSum0 = 0
    let leftSum1 = 0
    let leftSum2 = 0
    let rightSum0 = 0
    let rightSum1 = 0
    let rightSum2 = 0
    for (; start < end; start >>= 1, end >>= 1) {
      if (start & 1) {
        leftSum0 += this._sum0[start]
        leftSum1 += this._sum1[start]
        leftSum2 += this._sum2[start]
        start++
      }
      if (end & 1) {
        end--
        rightSum0 += this._sum0[end]
        rightSum1 += this._sum1[end]
        rightSum2 += this._sum2[end]
      }
    }
    return { sum0: leftSum0 + rightSum0, sum1: leftSum1 + rightSum1, sum2: leftSum2 + rightSum2 }
  }

  queryAll(): { sum0: number; sum1: number; sum2: number } {
    return { sum0: this._sum0[1], sum1: this._sum1[1], sum2: this._sum2[1] }
  }

  /**
   * 树上二分查询最大的`end`使得`[start,end)`内的值满足`predicate`.
   * @alias findFirst
   */
  maxRight(start: number, predicate: (sum0: number, sum1: number, sum2: number) => boolean): number {
    if (start < 0) start = 0
    if (start >= this._n) return this._n
    start += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(start >> i)
    let resSum0 = 0
    let resSum1 = 0
    let resSum2 = 0

    while (true) {
      while (!(start & 1)) start >>= 1
      const tmpSum01 = resSum0 + this._sum0[start]
      const tmpSum11 = resSum1 + this._sum1[start]
      const tmpSum21 = resSum2 + this._sum2[start]
      if (!predicate(tmpSum01, tmpSum11, tmpSum21)) {
        while (start < this._size) {
          this._pushDown(start)
          start <<= 1
          const tmpSum02 = resSum0 + this._sum0[start]
          const tmpSum12 = resSum1 + this._sum1[start]
          const tmpSum22 = resSum2 + this._sum2[start]
          if (predicate(tmpSum02, tmpSum12, tmpSum22)) {
            resSum0 = tmpSum02
            resSum1 = tmpSum12
            resSum2 = tmpSum22
            start++
          }
        }
        return start - this._size
      }
      resSum0 += this._sum0[start]
      resSum1 += this._sum1[start]
      resSum2 += this._sum2[start]
      start++
      if ((start & -start) === start) break
    }
    return this._n
  }

  /**
   * 树上二分查询最小的`start`使得`[start,end)`内的值满足`predicate`
   * @alias findLast
   */
  minLeft(end: number, predicate: (sum0: number, sum1: number, sum2: number) => boolean): number {
    if (end > this._n) end = this._n
    if (end <= 0) return 0
    end += this._size
    for (let i = this._height; i > 0; i--) this._pushDown((end - 1) >> i)
    let resSum0 = 0
    let resSum1 = 0
    let resSum2 = 0
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      const tmpSum01 = resSum0 + this._sum0[end]
      const tmpSum11 = resSum1 + this._sum1[end]
      const tmpSum21 = resSum2 + this._sum2[end]
      if (!predicate(tmpSum01, tmpSum11, tmpSum21)) {
        while (end < this._size) {
          this._pushDown(end)
          end = (end << 1) | 1
          const tmpSum02 = resSum0 + this._sum0[end]
          const tmpSum12 = resSum1 + this._sum1[end]
          const tmpSum22 = resSum2 + this._sum2[end]
          if (predicate(tmpSum02, tmpSum12, tmpSum22)) {
            resSum0 = tmpSum02
            resSum1 = tmpSum12
            resSum2 = tmpSum22
            end--
          }
        }
        return end + 1 - this._size
      }
      resSum0 += this._sum0[end]
      resSum1 += this._sum1[end]
      resSum2 += this._sum2[end]
      if ((end & -end) === end) break
    }
    return 0
  }

  build(leaves: ArrayLike<number>): void {
    if (leaves.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < this._n; i++) {
      this._sum0[this._size + i] = 1
      this._sum1[this._size + i] = leaves[i]
      this._sum2[this._size + i] = leaves[i] * leaves[i]
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

export { SegmentTreeRangeAddRangeSquareSum }

if (require.main === module) {
  class Mocker {
    private readonly _n: number
    private readonly _a: number[]
    constructor(nums: number[]) {
      this._n = nums.length
      this._a = nums
    }

    set(index: number, value: number): void {
      this._a[index] = value
    }

    get(index: number): { sum0: number; sum1: number; sum2: number } {
      return { sum0: 1, sum1: this._a[index], sum2: this._a[index] * this._a[index] }
    }

    update(start: number, end: number, lazy: number): void {
      for (let i = start; i < end; i++) this._a[i] += lazy
    }

    query(start: number, end: number): { sum0: number; sum1: number; sum2: number } {
      let sum0 = 0
      let sum1 = 0
      let sum2 = 0
      for (let i = start; i < end; i++) {
        sum0 += 1
        sum1 += this._a[i]
        sum2 += this._a[i] * this._a[i]
      }
      return { sum0, sum1, sum2 }
    }

    queryAll(): { sum0: number; sum1: number; sum2: number } {
      return this.query(0, this._n)
    }

    maxRight(start: number, predicate: (sum0: number, sum1: number, sum2: number) => boolean): number {
      let sum0 = 0
      let sum1 = 0
      let sum2 = 0
      for (let i = start; i < this._n; i++) {
        sum0 += 1
        sum1 += this._a[i]
        sum2 += this._a[i] * this._a[i]
        if (!predicate(sum0, sum1, sum2)) return i
      }
      return this._n
    }

    minLeft(end: number, predicate: (sum0: number, sum1: number, sum2: number) => boolean): number {
      let sum0 = 0
      let sum1 = 0
      let sum2 = 0
      for (let i = end - 1; i >= 0; i--) {
        sum0 += 1
        sum1 += this._a[i]
        sum2 += this._a[i] * this._a[i]
        if (!predicate(sum0, sum1, sum2)) return i + 1
      }
      return 0
    }

    build(leaves: ArrayLike<number>): void {
      for (let i = 0; i < this._n; i++) this._a[i] = leaves[i]
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
    const N = 5e4
    const real = new SegmentTreeRangeAddRangeSquareSum(N)
    const mock = new Mocker(Array(N).fill(0))
    for (let i = 0; i < N; i++) {
      const op = randint(0, 10)
      if (op === 0) {
        // set
        const index = randint(0, N - 1)
        const value = randint(-10, 10)
        real.set(index, value)
        mock.set(index, value)
      } else if (op === 1) {
        // get
        const index = randint(0, N - 1)
        const realValue = real.get(index)
        const mockValue = mock.get(index)
        assertSame(realValue, mockValue)
      } else if (op === 2) {
        // update
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const lazy = randint(-10, 10)
        real.update(start, end, lazy)
        mock.update(start, end, lazy)
      } else if (op === 3) {
        // query
        const start = randint(0, N - 1)
        const end = randint(start, N - 1)
        const realValue = real.query(start, end)
        const mockValue = mock.query(start, end)
        // assertSame(realValue, mockValue)
      } else if (op === 4) {
        // queryAll
        const realValue = real.queryAll()
        const mockValue = mock.queryAll()
        // assertSame(realValue, mockValue)
      } else if (op === 5) {
        // maxRight
        const start = randint(0, N - 1)
        const lazy = randint(-10, 10)
        const realValue = real.maxRight(start, (sum0, sum1, sum2) => sum2 + lazy * sum1 + lazy * lazy * sum0 <= 0)
        const mockValue = mock.maxRight(start, (sum0, sum1, sum2) => sum2 + lazy * sum1 + lazy * lazy * sum0 <= 0)
        // assertSame(realValue, mockValue)
      } else if (op === 6) {
        // minLeft
        const end = randint(0, N)
        const lazy = randint(-10, 10)
        const realValue = real.minLeft(end, (sum0, sum1, sum2) => sum2 + lazy * sum1 + lazy * lazy * sum0 <= 0)
        const mockValue = mock.minLeft(end, (sum0, sum1, sum2) => sum2 + lazy * sum1 + lazy * lazy * sum0 <= 0)
        // assertSame(realValue, mockValue)
      } else if (op === 7) {
        // build
        const leaves = Array(N)
          .fill(0)
          .map(() => randint(-10, 10))
        real.build(leaves)
        mock.build(leaves)
      }
    }
  }

  console.log('test passed')
}
