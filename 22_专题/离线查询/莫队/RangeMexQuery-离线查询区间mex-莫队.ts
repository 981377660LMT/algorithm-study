import { MoAlgo } from './Moalgo'

/**
 * 离线查询区间mex.基于莫队算法实现.
 */
class RangeMexQueryMo {
  private readonly _mo: MoAlgo
  private readonly _nums: ArrayLike<number>
  private readonly _q: number

  constructor(nums: ArrayLike<number>, q: number) {
    const n = nums.length
    this._mo = new MoAlgo(n, q)
    this._nums = nums
    this._q = q
  }

  /**
   * [start, end).
   * 0 <= start <= end <= n.
   */
  addQuery(start: number, end: number): void {
    this._mo.addQuery(start, end)
  }

  /**
   * @param mexStart mex的起始值(从0开始还是从1开始).
   */
  run(mexStart: number): number[] {
    const res = Array(this._q)
    let mex = mexStart
    const counter = new Map<number, number>()
    const add = (index: number) => {
      const num = this._nums[index]
      counter.set(num, (counter.get(num) || 0) + 1)
      while (counter.get(mex)) mex++
    }
    const remove = (index: number) => {
      const num = this._nums[index]
      counter.set(num, (counter.get(num) || 0) - 1)
      if (!counter.get(num)) mex = Math.min(mex, num)
    }
    const query = (qi: number) => {
      res[qi] = mex
    }
    this._mo.run(add, remove, query)
    return res
  }
}

export { RangeMexQueryMo }

if (require.main === module) {
  const nums = [1, 2, 3, 5, 6, 6, 7, 8, 9]
  const M = new RangeMexQueryMo(nums, 2)
  M.addQuery(0, 3)
  M.addQuery(0, 4)
  console.log(M.run(1))

  const n = 1e5
  const arr = Array(n)
    .fill(0)
    .map((_, i) => i + 1)
  const M2 = new RangeMexQueryMo(arr, n)
  for (let i = 0; i < n; i++) {
    M2.addQuery(0, i)
  }
  console.time('mex')
  M2.run(1)
  console.timeEnd('mex')
}
