/* eslint-disable @typescript-eslint/no-non-null-assertion */
// # !维护topK之和 (这里topK是最小值)
// # TopKSum

import { SortedListFastWithSum } from '../../22_专题/离线查询/根号分治/SortedList/SortedListWithSum'

import { Heap } from '../../8_heap/Heap'

/**
 * 两个堆维护topK之和.
 * 当维护的和更复杂时，使用 {@link SortedListFastWithSum} 代替.
 */
class TopKSum {
  private _sum = 0
  private _k: number
  private readonly _min: boolean
  private readonly _in = new Heap<number>({ data: [], less: (a, b) => a < b })
  private readonly _out = new Heap<number>({ data: [], less: (a, b) => a < b })
  private readonly _dIn = new Heap<number>({ data: [], less: (a, b) => a < b })
  private readonly _dOut = new Heap<number>({ data: [], less: (a, b) => a < b })
  private readonly _counter = new Map<number, number>()

  /**
   * 维护容器的topK之和.
   * @param k topK
   * @param min 是否是最小值.默认是最小k个数之和.
   */
  constructor(k: number, min = true) {
    this._k = k
    this._min = min
  }

  query(): number {
    return this._min ? this._sum : -this._sum
  }

  add(x: number): void {
    if (!this._min) x = -x
    this._counter.set(x, (this._counter.get(x) || 0) + 1)
    this._in.push(-x)
    this._sum += x
    this._modify()
  }

  discard(x: number): void {
    if (!this._min) x = -x
    if (!this._counter.has(x)) return
    this._counter.set(x, this._counter.get(x)! - 1)
    if (this._in.size && -this._in.top()! === x) {
      this._sum -= x
      this._in.pop()
    } else if (this._in.size && -this._in.top()! > x) {
      this._sum -= x
      this._dIn.push(-x)
    } else {
      this._dOut.push(x)
    }
    this._modify()
  }

  setK(k: number): void {
    this._k = k
    this._modify()
  }

  getK(): number {
    return this._k
  }

  has(x: number): boolean {
    if (!this._min) x = -x
    return this._counter.has(x)
  }

  get size(): number {
    return this._in.size + this._out.size - this._dIn.size - this._dOut.size
  }

  private _modify(): void {
    while (this._out.size && this._in.size - this._dIn.size < this._k) {
      const p = this._out.pop()!
      if (this._dOut.size && p === this._dOut.top()!) {
        this._dOut.pop()
      } else {
        this._sum += p
        this._in.push(-p)
      }
    }

    while (this._in.size - this._dIn.size > this._k) {
      const p = -this._in.pop()!
      if (this._dIn.size && p === -this._dIn.top()!) {
        this._dIn.pop()
      } else {
        this._sum -= p
        this._out.push(p)
      }
    }

    while (this._dIn.size && this._in.size && this._in.top()! === this._dIn.top()!) {
      this._in.pop()
      this._dIn.pop()
    }
  }
}

export { TopKSum }

if (require.main === module) {
  // https://leetcode.cn/problems/minimum-difference-in-sums-after-removal-of-elements/
  // eslint-disable-next-line no-inner-declarations
  function minimumDifference(nums: number[]): number {
    // 前面最小n个和后面大n个
    const n = nums.length / 3
    const minK = new TopKSum(n, true)
    const maxK = new TopKSum(n, false)
    for (let i = 0; i < n; i++) minK.add(nums[i])
    for (let i = n; i < 3 * n; i++) maxK.add(nums[i])
    let res = minK.query() - maxK.query()
    for (let i = n; i < 2 * n; i++) {
      minK.add(nums[i])
      maxK.discard(nums[i])
      console.log(minK.query(), maxK.query())
      res = Math.min(res, minK.query() - maxK.query())
    }
    return res
  }
}
