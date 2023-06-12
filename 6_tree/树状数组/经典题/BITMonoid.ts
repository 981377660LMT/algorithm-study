/* eslint-disable no-inner-declarations */
// https://leetcode.cn/circle/discuss/9n7Hnx/

// !set/update/get/queryPrefix/queryRange

/**
 * 维护幺半群的树状数组.
 * 支持单点更新,单点修改,前缀查询,区间查询.
 * !内部由Map实现,无需离散化.
 */
class BITMonoidMap<E> {
  private readonly _n: number
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E
  private readonly _data: Map<number, E> = new Map()
  private readonly _sum: Map<number, E> = new Map()

  constructor(n: number, e: () => E, op: (a: E, b: E) => E) {
    this._n = n + 5
    this._e = e
    this._op = op
  }

  get(index: number): E {
    return this.queryRange(index, index + 1)
  }

  /**
   * 单点更新,时间复杂度O(log n).
   * 0<=index<n.
   */
  update(index: number, value: E): void {
    index++
    this._data.set(index, this._op(this._data.get(index) ?? this._e(), value))
    for (; index <= this._n; index += index & -index) {
      this._sum.set(index, this._op(this._sum.get(index) ?? this._e(), value))
    }
  }

  /**
   * 查询前缀`[0,right)`聚合值,时间复杂度O(log n).
   * 0<=right<=n.
   */
  queryPrefix(right: number): E {
    if (right > this._n) right = this._n
    let res = this._e()
    for (; right > 0; right -= right & -right) {
      res = this._op(res, this._sum.get(right) ?? this._e())
    }
    return res
  }

  /**
   * 查询区间`[left,right)`聚合值,时间复杂度O(log^2 n).
   * 0<=left<=right<=n.
   */
  queryRange(left: number, right: number): E {
    if (right > this._n) right = this._n
    left++
    let res = this._e()
    while (right >= left) {
      if ((right & (right - 1)) >= left - 1) {
        res = this._op(res, this._sum.get(right) ?? this._e())
        right &= right - 1
      } else {
        res = this._op(res, this._data.get(right) ?? this._e())
        right--
      }
    }
    return res
  }
}

/**
 * 维护幺半群的树状数组.
 * 支持单点更新,单点修改,前缀查询,区间查询.
 */
class BITMonoidArray<E> {
  private readonly _n: number
  private readonly _data: E[]
  private readonly _sum: E[]
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E

  constructor(nOrArr: number | ArrayLike<E>, e: () => E, op: (a: E, b: E) => E) {
    const n = typeof nOrArr === 'number' ? nOrArr : nOrArr.length
    this._n = n
    this._data = Array(n + 1)
    this._sum = Array(n + 1)
    this._e = e
    this._op = op
    if (typeof nOrArr !== 'number') this.build(nOrArr)
  }

  /**
   * 单点修改,时间复杂度O(log^2 n).
   * 0<=index<n.
   */
  set(index: number, value: E): void {
    index++
    this._data[index] = value
    for (; index <= this._n; index += index & -index) {
      this._sum[index] = value
      for (let i = 1; i < (index & -index); i <<= 1) {
        this._sum[index] = this._op(this._sum[index], this._sum[index - i])
      }
    }
  }

  get(index: number): E {
    return this.queryRange(index, index + 1)
  }

  /**
   * 单点更新,时间复杂度O(log n).
   * 0<=index<n.
   */
  update(index: number, value: E): void {
    index++
    this._data[index] = this._op(this._data[index], value)
    for (; index <= this._n; index += index & -index) {
      this._sum[index] = this._op(this._sum[index], value)
    }
  }

  /**
   * 查询前缀`[0,right)`聚合值,时间复杂度O(log n).
   * 0<=right<=n.
   */
  queryPrefix(right: number): E {
    if (right > this._n) right = this._n
    let res = this._e()
    for (; right > 0; right &= right - 1) res = this._op(res, this._sum[right])
    return res
  }

  /**
   * 查询区间`[left,right)`聚合值,时间复杂度O(log^2 n).
   * 0<=left<=right<=n.
   */
  queryRange(left: number, right: number): E {
    if (right > this._n) right = this._n
    left++
    let res = this._e()
    while (right >= left) {
      if ((right & (right - 1)) >= left - 1) {
        res = this._op(res, this._sum[right])
        right &= right - 1
      } else {
        res = this._op(res, this._data[right])
        right--
      }
    }
    return res
  }

  /**
   * O(n)建树.
   */
  build(arr: ArrayLike<E>): void {
    if (arr.length !== this._n) throw new RangeError(`arr length must be equal to ${this._n}`)
    for (let i = 1; i <= this._n; i++) {
      this._data[i] = arr[i - 1]
      this._sum[i] = arr[i - 1]
      for (let j = 1; j < (i & -i); j <<= 1) {
        this._sum[i] = this._op(this._sum[i], this._sum[i - j])
      }
    }
  }

  toString(): string {
    const res: E[] = []
    for (let i = 0; i < this._n; i++) {
      res.push(this.queryRange(i, i + 1))
    }
    return `BITMonoid{${res.join(',')}}`
  }
}

export { BITMonoidArray, BITMonoidMap }

if (require.main === module) {
  const bit = new BITMonoidArray([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], () => 0, Math.max)

  console.log(bit.toString())
  bit.set(0, 3)
  bit.set(1, 4)
  bit.set(2, 5)
  console.log(bit.toString())

  console.log(bit.toString())
  console.log(bit.queryRange(0, 3))
  console.log(bit.toString())

  bit.update(1, 0)
  console.log(bit.toString())

  console.time('queryRange')
  const n = 2e5
  const tree = new BITMonoidArray(n, () => 0, Math.max)
  for (let i = 0; i <= n; i++) {
    tree.queryRange(i, n)
  }
  console.timeEnd('queryRange') // queryRange: 69.958ms

  console.time('queryPrefix')
  for (let i = 0; i <= n; i++) {
    tree.queryPrefix(i)
  }
  console.timeEnd('queryPrefix') // queryPrefix: 14.018ms

  console.time('set')
  for (let i = 0; i <= n; i++) {
    tree.set(i, i)
  }
  console.timeEnd('set') // set: 124.716ms

  console.time('update')
  for (let i = 0; i <= n; i++) {
    tree.update(i, i)
  }
  console.timeEnd('update') // update: 12.925ms

  console.time('get')
  for (let i = 0; i <= n; i++) {
    tree.get(i)
  }
  console.timeEnd('get') // get: 7.697ms

  const INF = 2e15
  // https://leetcode.cn/problems/sliding-window-maximum/submissions/
  function maxSlidingWindow(nums: number[], k: number): number[] {
    const tree = new BITMonoidArray(nums.length, () => -INF, Math.max)
    const res: number[] = []
    for (let i = 0; i < nums.length; i++) {
      tree.set(i, nums[i])
      if (i >= k - 1) {
        res.push(tree.queryRange(i - k + 1, i + 1))
      }
    }
    return res
  }
}
