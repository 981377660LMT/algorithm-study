/* eslint-disable max-len */
/* eslint-disable arrow-body-style */
/* eslint-disable no-inner-declarations */

/**
 * 区间赋值区间加，区间和.
 */
class SegmentTreeRangeAssignRangeAddRangeSum {
  private readonly _n: number
  private readonly _m: number
  private readonly _height: number
  private readonly _sum: Float64Array
  private readonly _size: Uint32Array
  private readonly _lazyMul: Int8Array
  private readonly _lazyAdd: Float64Array

  constructor(nOrArr: number | ArrayLike<number>) {
    const n = typeof nOrArr === 'number' ? nOrArr : nOrArr.length
    let m = 1
    let height = 0
    while (m < n) {
      m <<= 1
      height++
    }
    this._n = n
    this._m = m
    this._height = height

    // !0.init data and lazy
    const sum = new Float64Array(m << 1)
    const size = new Uint32Array(m << 1)
    const lazyMul = new Int8Array(m).fill(1)
    const lazyAdd = new Float64Array(m)
    this._sum = sum
    this._size = size
    this._lazyMul = lazyMul
    this._lazyAdd = lazyAdd

    if (typeof nOrArr === 'number') nOrArr = new Uint8Array(nOrArr)
    this._build(nOrArr)
  }

  /**
   * 将区间`[start,end)`的值全部加上`delta`.
   * 0 <= start <= end <= n.
   */
  addRange(start: number, end: number, delta: number): void {
    this.update(start, end, 1, delta)
  }

  /**
   * 将区间`[start,end)`的值全部设置为`value`.
   * 0 <= start <= end <= n.
   */
  assignRange(start: number, end: number, value: number): void {
    this.update(start, end, 0, value)
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._n) return
    index += this._m
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    this._set(index, value)
    for (let i = 1; i <= this._height; i++) this._pushUp(index >> i)
  }

  get(index: number): number {
    if (index < 0 || index >= this._n) {
      throw new RangeError(`index must be in [0, ${this._n})`)
    }
    index += this._m
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    return this._sum[index]
  }

  /**
   * 区间`[start,end)`的值与`lazy`进行作用.
   * 0 <= start <= end <= n.
   */
  update(start: number, end: number, lazyMul: number, lazyAdd: number): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    start += this._m
    end += this._m
    for (let i = this._height; i > 0; i--) {
      if ((start >> i) << i !== start) this._pushDown(start >> i)
      if ((end >> i) << i !== end) this._pushDown((end - 1) >> i)
    }
    let start2 = start
    let end2 = end
    for (; start < end; start >>= 1, end >>= 1) {
      if (start & 1) this._propagate(start++, lazyMul, lazyAdd)
      if (end & 1) this._propagate(--end, lazyMul, lazyAdd)
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
  query(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return 0
    start += this._m
    end += this._m
    for (let i = this._height; i > 0; i--) {
      if ((start >> i) << i !== start) this._pushDown(start >> i)
      if ((end >> i) << i !== end) this._pushDown((end - 1) >> i)
    }
    let leftSum = 0
    let rightSum = 0
    for (; start < end; start >>= 1, end >>= 1) {
      if (start & 1) leftSum += this._sum[start++]
      if (end & 1) rightSum += this._sum[--end]
    }
    return leftSum + rightSum
  }

  queryAll(): number {
    return this._sum[1]
  }

  /**
   * 树上二分查询最大的`end`使得`[start,end)`内的值满足`predicate`.
   * @alias findFirst
   */
  maxRight(start: number, predicate: (sum: number) => boolean): number {
    if (start < 0) start = 0
    if (start >= this._n) return this._n
    start += this._m
    for (let i = this._height; i > 0; i--) this._pushDown(start >> i)
    let resSum = 0
    while (true) {
      while (!(start & 1)) start >>= 1
      const tmpSum1 = resSum + this._sum[start]
      if (!predicate(tmpSum1)) {
        while (start < this._m) {
          this._pushDown(start)
          start <<= 1
          const tmpSum2 = resSum + this._sum[start]
          if (predicate(tmpSum2)) {
            resSum = tmpSum2
            start++
          }
        }
        return start - this._m
      }
      resSum += this._sum[start]
      start++
      if ((start & -start) === start) break
    }
    return this._n
  }

  /**
   * 树上二分查询最小的`start`使得`[start,end)`内的值满足`predicate`
   * @alias findLast
   */
  minLeft(end: number, predicate: (sum: number) => boolean): number {
    if (end > this._n) end = this._n
    if (end <= 0) return 0
    end += this._m
    for (let i = this._height; i > 0; i--) this._pushDown((end - 1) >> i)
    let resSum = 0
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      const tmpSum1 = this._sum[end] + resSum
      if (!predicate(tmpSum1)) {
        while (end < this._m) {
          this._pushDown(end)
          end = (end << 1) | 1
          const tmpSum2 = this._sum[end] + resSum
          if (predicate(tmpSum2)) {
            resSum = tmpSum2
            end--
          }
        }
        return end + 1 - this._m
      }
      resSum += this._sum[end]
      if ((end & -end) === end) break
    }
    return 0
  }

  toString(): string {
    const sb: string[] = []
    sb.push('SegmentTreeRangeAssignRangeAddRangeSum(')
    for (let i = 0; i < this._n; i++) {
      if (i) sb.push(', ')
      sb.push(JSON.stringify(this.get(i)))
    }
    sb.push(')')
    return sb.join('')
  }

  private _build(leaves: ArrayLike<number>): void {
    if (leaves.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < this._n; i++) this._set(i + this._m, leaves[i])
    for (let i = this._m - 1; i > 0; i--) this._pushUp(i)
  }

  private _set(index: number, value: number): void {
    this._sum[index] = value
    this._size[index] = 1
  }

  private _pushUp(index: number): void {
    this._sum[index] = this._sum[index << 1] + this._sum[(index << 1) | 1]
    this._size[index] = this._size[index << 1] + this._size[(index << 1) | 1]
  }

  private _pushDown(index: number): void {
    const mul = this._lazyMul[index]
    const add = this._lazyAdd[index]
    if (mul === 1 && !add) return
    this._propagate(index << 1, mul, add)
    this._propagate((index << 1) | 1, mul, add)
    this._lazyMul[index] = 1
    this._lazyAdd[index] = 0
  }

  private _propagate(index: number, mul: number, add: number): void {
    this._sum[index] = mul * this._sum[index] + add * this._size[index]
    if (index < this._m) {
      this._lazyMul[index] &= mul
      this._lazyAdd[index] = mul * this._lazyAdd[index] + add
    }
  }
}

export { SegmentTreeRangeAssignRangeAddRangeSum }

if (require.main === module) {
  // demo()
  checkWithBruteForce()
  timeit()
  function demo(): void {
    const tree = new SegmentTreeRangeAssignRangeAddRangeSum(10)
    console.log(tree.toString())
    tree.addRange(0, 4, 2)
    console.log(tree.toString())
    tree.assignRange(2, 6, 3)
    console.log(tree.toString())
    tree.set(1, 4)
    console.log(tree.toString())
  }

  function checkWithBruteForce(): void {
    class Mocker {
      readonly _n: number
      private readonly _a: number[]
      constructor(nums: number[]) {
        this._n = nums.length
        this._a = nums.slice()
      }

      assignRange(start: number, end: number, value: number): void {
        for (let i = start; i < end; i++) this._a[i] = value
      }

      addRange(start: number, end: number, delta: number): void {
        for (let i = start; i < end; i++) this._a[i] += delta
      }

      set(index: number, value: number): void {
        this._a[index] = value
      }

      get(index: number): number {
        return this._a[index]
      }

      update(start: number, end: number, lazyMul: number, lazyAdd: number): void {
        for (let i = start; i < end; i++) {
          this._a[i] = this._a[i] * lazyMul + lazyAdd
        }
      }

      query(start: number, end: number): number {
        let sum = 0
        for (let i = start; i < end; i++) {
          sum += this._a[i]
        }
        return sum
      }

      queryAll(): number {
        return this.query(0, this._n)
      }

      maxRight(start: number, predicate: (sum: number) => boolean): number {
        let sum = 0
        for (let i = start; i < this._n; i++) {
          sum += this._a[i]
          if (!predicate(sum)) return i
        }
        return this._n
      }

      minLeft(end: number, predicate: (sum: number) => boolean): number {
        let sum = 0
        for (let i = end - 1; i >= 0; i--) {
          sum += this._a[i]
          if (!predicate(sum)) return i + 1
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

    function assertSame(obj1: unknown, obj2: unknown) {
      if (JSON.stringify(obj1) !== JSON.stringify(obj2)) {
        throw new Error(`expect ${JSON.stringify(obj2)}, got ${JSON.stringify(obj1)}`)
      }
    }

    const randint = (min: number, max: number) => Math.floor(Math.random() * (max - min + 1)) + min
    const N = 5e4
    const real = new SegmentTreeRangeAssignRangeAddRangeSum(N)
    const mock = new Mocker(Array(N).fill(0))
    for (let i = 0; i < N; i++) {
      const op = randint(0, 8)
      if (op === 0) {
        // sets
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
        assertSame(realValue, mockValue)
      } else if (op === 2) {
        // update
        // const start = randint(0, N - 1)
        // const end = randint(start, N)
        // const lazyMul = randint(0, 2)
        // const lazyAdd = randint(0, 10)
        // real.update(start, end, lazyMul, lazyAdd)
        // mock.update(start, end, lazyMul, lazyAdd)
        // // console.log('update', start, end, lazy)
      } else if (op === 3) {
        // query
        const start = randint(0, N - 1)
        const end = randint(start, N)
        const realValue = real.query(start, end)
        const mockValue = mock.query(start, end)
        // console.log('query', start, end, realValue, mockValue)
        assertSame(realValue, mockValue)
      } else if (op === 4) {
        // queryAll
        const realValue = real.queryAll()
        const mockValue = mock.queryAll()
        assertSame(realValue, mockValue)
      } else if (op === 5) {
        // maxRight
        const start = randint(0, N - 1)
        const target = randint(0, N)
        const realValue = real.maxRight(start, sum => sum <= target)
        const mockValue = mock.maxRight(start, sum => sum <= target)
        assertSame(realValue, mockValue)
      } else if (op === 6) {
        // minLeft
        const end = randint(0, N)
        const target = randint(0, N)
        const realValue = real.minLeft(end, sum => sum <= target)
        const mockValue = mock.minLeft(end, sum => sum <= target)
        assertSame(realValue, mockValue)
      } else if (op === 7) {
        // addRange
        const start = randint(0, N - 1)
        const end = randint(start, N)
        const delta = randint(0, 10)
        real.addRange(start, end, delta)
        mock.addRange(start, end, delta)
      } else if (op === 8) {
        // assignRange
        const start = randint(0, N - 1)
        const end = randint(start, N)
        const value = randint(0, 10)
        real.assignRange(start, end, value)
        mock.assignRange(start, end, value)
      }
    }
    console.log('test passed')
  }

  function timeit(): void {
    const n = 2e5
    const arr = Array(n)
    for (let i = 0; i < n; i++) arr[i] = Math.floor(Math.random() * 10)
    const seg = new SegmentTreeRangeAssignRangeAddRangeSum(arr)
    console.time('SegmentTreeRangeAssignRangeAddRangeSum')
    for (let i = 0; i < n; i++) {
      seg.query(i, n)
      seg.update(i, n, 1, 1)
      seg.set(i, 1)
      seg.maxRight(i, sum => sum <= i)
      seg.minLeft(i, sum => sum <= i)
    }
    console.timeEnd('SegmentTreeRangeAssignRangeAddRangeSum') // SegmentTreeRangeAssignRangeAddRangeSum: 215.476ms
  }
}
