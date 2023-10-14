/* eslint-disable no-inner-declarations */

import { Heap } from '../../8_heap/Heap'

/**
 * 从数组不相邻选择 `k(0<=k<=(n+1/2))` 个数,最大化和/最小化和.
 */
class NonAdjacentSelection {
  private readonly _n: number
  private _nums: number[]
  private readonly _minimize: boolean
  private readonly _history: number[] = [] //  history[i<<1], history[(i<<1)|1] -> {start, end}
  private _solved = false

  constructor(nums: number[], minimize = false) {
    this._n = nums.length
    this._nums = nums
    this._minimize = minimize
  }

  solve(): number[] {
    if (this._minimize) {
      const copy = Array(this._n)
      for (let i = 0; i < this._n; i++) copy[i] = -this._nums[i]
      this._nums = copy
    }

    const { _n: n, _nums: nums, _history: history } = this
    const rest = new Uint8Array(n + 2)
    rest.fill(1, 1, n + 1)
    const left = new Int32Array(n + 2)
    const right = new Int32Array(n + 2)
    for (let i = 0; i < n + 2; i++) {
      left[i] = i - 1
      right[i] = i + 1
    }

    // i 和 i+1 表示 start 和 end 的下标
    const range = new Int32Array(2 * (n + 2))
    for (let i = 1; i < n + 1; i++) {
      range[i << 1] = i - 1
      range[(i << 1) | 1] = i
    }
    const val = Array<number>(n + 2).fill(0)
    for (let i = 1; i < n + 1; i++) val[i] = nums[i - 1]

    const pairs: { add: number; index: number }[] = Array(n)
    for (let i = 0; i < n; i++) pairs[i] = { add: val[i + 1], index: i + 1 }
    const pq = new Heap<{ add: number; index: number }>((a, b) => b.add - a.add, pairs)

    const res: number[] = [0]
    while (pq.size) {
      const { add, index } = pq.pop()!
      if (!rest[index]) continue
      res.push(res[res.length - 1] + add)
      const L = left[index]
      const R = right[index]
      history.push(range[index << 1], range[(index << 1) | 1])
      if (L >= 1) {
        right[left[L]] = index
        left[index] = left[L]
      }
      if (R <= n) {
        left[right[R]] = index
        right[index] = right[R]
      }

      if (rest[L] && rest[R]) {
        val[index] = val[L] + val[R] - val[index]
        pq.push({ add: val[index], index })
        range[index << 1] = range[L << 1]
        range[(index << 1) | 1] = range[(R << 1) | 1]
      } else {
        rest[index] = 0
      }
      rest[L] = 0
      rest[R] = 0
    }

    if (this._minimize) {
      for (let i = 0; i < res.length; i++) res[i] = -res[i]
    }

    this._solved = true
    return res
  }

  /**
   * @param k 选择 k 个数使得和最大/最小, 返回选择的数的下标.
   * 0 <= k <= (n+1)/2.
   */
  restore(k: number): number[] {
    if (k < 0 || k > (this._n + 1) >>> 1) throw new Error('k is out of range')
    if (!this._solved) this.solve()
    const diff = Array<number>(this._n + 1).fill(0)
    for (let i = 0; i < k; i++) {
      const start = this._history[i << 1]
      const end = this._history[(i << 1) | 1]
      diff[start]++
      diff[end]--
    }
    for (let i = 1; i < diff.length; i++) diff[i] += diff[i - 1]
    const res: number[] = []
    for (let i = 0; i < this._n; i++) {
      if (diff[i] & 1) res.push(i)
    }
    return res
  }
}

export { NonAdjacentSelection }

if (require.main === module) {
  const nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
  const nas = new NonAdjacentSelection(nums)
  const res = nas.solve()
  console.log(res, nas.restore(4))

  testTime()

  function testTime() {
    const nums = Array(2e5)
      .fill(0)
      .map(() => Math.floor(Math.random() * 2e5))

    console.time('NonAdjacentSelection')
    const nas = new NonAdjacentSelection(nums)
    nas.solve()
    console.timeEnd('NonAdjacentSelection')
  }
}
