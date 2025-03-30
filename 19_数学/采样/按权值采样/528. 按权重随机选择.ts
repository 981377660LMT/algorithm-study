/* eslint-disable prefer-destructuring */
// 彩票调度

import { bisectLeft } from '../../../9_排序和搜索/二分/bisect'

function randint(start: number, end: number) {
  if (start > end) throw new Error('invalid interval')
  const amplitude = end - start
  return start + Math.floor((amplitude + 1) * Math.random())
}

class Solution {
  private readonly _presum: number[]
  private readonly _sum: number

  /**
   * @param w 给定一个正整数数组 w,其中 w[i] 代表下标 i 的权重.
   */
  constructor(w: number[]) {
    this._presum = Array(w.length)
    this._presum[0] = w[0]
    for (let i = 1; i < this._presum.length; i++) {
      this._presum[i] = this._presum[i - 1] + w[i]
    }
    this._sum = this._presum[this._presum.length - 1]
  }

  /**
   * @return 返回下标 i 的概率为 w[i] / sum(w).
   */
  pickIndex(): number {
    const rand = randint(1, this._sum)
    return bisectLeft(this._presum, rand)
  }
}

const solution = new Solution([1, 3])
console.log(solution.pickIndex())
console.log(solution.pickIndex())
console.log(solution.pickIndex())
console.log(solution.pickIndex())
// solution.pickIndex(); // 返回 1，返回下标 1，返回该下标概率为 3/4 。
// solution.pickIndex(); // 返回 1
// solution.pickIndex(); // 返回 1
// solution.pickIndex(); // 返回 1
// solution.pickIndex(); // 返回 0，返回下标 0，返回该下标概率为 1/4 。

export {}
