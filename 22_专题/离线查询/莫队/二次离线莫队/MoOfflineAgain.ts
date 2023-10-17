/* eslint-disable no-inner-declarations */

// 二次离线莫队
//
// 一般莫队有O(n√n)次端点移动，如果要用数据结构维护信息的话，就有o(n√n)次修改和O(n√n)次查询。
// 而莫队二次离线能够优化为成O(n)次修改和O(n√n)次查询，
// !从而允许使用一些修改复杂度大而查询复杂度小的方式来维护信息。例如分块，如果能O(√n)修改和O(1)查询的话，总的复杂度就是O(n√n)。
//
// !https://github.com/Aeren1564/Algorithms/blob/master/Algorithm_Implementations_Cpp/Data_Structure/Range_Inversion_Query_Solver/range_inversion_query_solver_offline.sublime-snippet
// https://www.luogu.com.cn/blog/gxy001/mu-dui-er-ci-li-xian
// https://kewth.github.io/2019/10/16/%E8%8E%AB%E9%98%9F%E4%BA%8C%E6%AC%A1%E7%A6%BB%E7%BA%BF/
// https://www.luogu.com.cn/problem/P4887
// https://www.luogu.com.cn/problem/P5398
//
// 适用范围：
// !贡献是可交换群(阿贝尔群), 即 f(x,start,end)) = f(x,0,end) - f(x,0,start);

import { discretizeCompressed } from '../../../前缀与差分/差分数组/离散化/discretize'
import { BITRangeBlockFastQuery } from '../../根号分治/值域分块/BITRangeBlockFastQuery'

/** 可交换群(commutative group). */
interface AbelianGroup<V> {
  e: () => V
  op: (a: V, b: V) => V
  inv: (a: V) => V
}

class MoOfflineAgain<V = number> {
  private readonly _n: number
  private readonly _q: number
  private readonly _blockSize: number
  private readonly _queryBlocks: { qi: number; ql: number; qr: number }[][]
  private _queryOrder = 0

  private readonly _e: () => V
  private readonly _op: (a: V, b: V) => V
  private readonly _inv: (a: V) => V

  constructor(
    n: number,
    q: number,
    options?: {
      blockSize?: number
      abelianGroup?: AbelianGroup<V>
    }
  ) {
    let { abelianGroup, blockSize } = options || {}
    if (blockSize === void 0) {
      const sqrt = Math.sqrt((q * 2) / 3) | 0
      blockSize = Math.max(1, (n / Math.max(1, sqrt)) | 0)
    }
    const e = abelianGroup ? abelianGroup.e : () => 0 as V
    const op = abelianGroup ? abelianGroup.op : (a: any, b: any) => (a + b) as V
    const inv = abelianGroup ? abelianGroup.inv : (a: V) => -a as V
    const queryBlocks = Array(((n / blockSize) | 0) + 1)
    for (let i = 0; i < queryBlocks.length; i++) {
      queryBlocks[i] = []
    }
    this._n = n
    this._q = q
    this._blockSize = blockSize
    this._queryBlocks = queryBlocks
    this._e = e
    this._op = op
    this._inv = inv
  }

  /**
   * 添加一个查询，查询范围为`左闭右开区间` [start, end).
   * 0 <= start <= end <= n
   */
  addQuery(start: number, end: number): void {
    const bid = (start / this._blockSize) | 0
    this._queryBlocks[bid].push({ qi: this._queryOrder, ql: start, qr: end })
    this._queryOrder++
  }

  /**
   * @param add 将`A[index]`加入窗口中.
   * @param queryLeft 窗口最左侧的`A[index]`对答案的贡献.
   * @param queryRight 窗口最右侧的`A[index]`对答案的贡献.
   * @complexity `add` 操作次数为`O(n)`，`query` 操作次数为`O(nsqrt(n))`.
   */
  run(add: (i: number) => void, queryLeft: (i: number) => V, queryRight: (i: number) => V): V[] {
    const { _n: n, _q: q, _queryBlocks: blocks, _e: e, _op: op, _inv: inv } = this
    const eventGroups: { start: number; end: number; qi: number; type: number }[][] = Array(n + 1)
    for (let i = 0; i < eventGroups.length; i++) eventGroups[i] = []
    let left = 0
    let right = 0

    for (let i = 0; i < blocks.length; i++) {
      const block = blocks[i]
      if (i & 1) {
        block.sort((a, b) => a.qr - b.qr)
      } else {
        block.sort((a, b) => b.qr - a.qr)
      }

      for (let j = 0; j < block.length; j++) {
        const { qi, ql, qr } = block[j]
        if (ql < left) {
          eventGroups[right].push({ qi, start: ql, end: left, type: 2 })
          left = ql
        }
        if (right < qr) {
          eventGroups[left].push({ qi, start: right, end: qr, type: 1 })
          right = qr
        }
        if (left < ql) {
          eventGroups[right].push({ qi, start: left, end: ql, type: 0 })
          left = ql
        }
        if (qr < right) {
          eventGroups[left].push({ qi, start: qr, end: right, type: 3 })
          right = qr
        }
      }
    }

    const rightSum: V[] = Array(n + 1)
    const leftSum: V[] = Array(n + 1)
    rightSum[0] = e()
    leftSum[0] = e()
    const res = Array(q)
    for (let i = 0; i < res.length; i++) res[i] = e()

    for (let i = 0; i <= n; i++) {
      const events = eventGroups[i]
      for (let j = 0; j < events.length; j++) {
        const { qi, start, end, type } = events[j]
        let sum = e()
        if (type & 1) {
          for (let k = start; k < end; k++) {
            sum = op(sum, queryRight(k))
          }
        } else {
          for (let k = start; k < end; k++) {
            sum = op(sum, queryLeft(k))
          }
        }
        res[qi] = type & 2 ? op(res[qi], inv(sum)) : op(res[qi], sum)
      }

      if (i < n) {
        rightSum[i + 1] = op(rightSum[i], queryRight(i))
        add(i)
        leftSum[i + 1] = op(leftSum[i], queryLeft(i))
      }
    }

    let curSum = e()
    for (let i = 0; i < blocks.length; i++) {
      const block = blocks[i]
      for (let j = 0; j < block.length; j++) {
        const { qi, ql, qr } = block[j]
        curSum = op(curSum, res[qi])
        res[qi] = op(op(leftSum[ql], rightSum[qr]), inv(curSum))
      }
    }

    return res
  }
}

export { MoOfflineAgain }

if (require.main === module) {
  // 静态区间逆序对-离线.
  // 时间复杂度O(nsqrt(n)),空间复杂度O(n).
  // https://judge.yosupo.jp/problem/static_range_inversions_query
  function staticRangeInversionsQuery(nums: number[], ranges: [start: number, end: number][]): number[] {
    const n = nums.length
    const q = ranges.length
    const [rank, newNums] = discretizeCompressed(nums)
    const bit = new BITRangeBlockFastQuery(rank.size)
    const mo = new MoOfflineAgain(n, q)
    ranges.forEach(([start, end]) => {
      mo.addQuery(start, end)
    })
    const res = mo.run(
      (index: number) => {
        bit.add(newNums[index], 1)
      },
      (index: number) => {
        const res = bit.queryRange(0, newNums[index])
        return res
      },
      (index: number) => {
        const res = bit.queryRange(newNums[index] + 1, rank.size)
        return res
      }
    )
    return res
  }

  console.log(
    staticRangeInversionsQuery(
      [4, 1, 4, 0],
      [
        [1, 3],
        [0, 4]
      ]
    )
  )

  // testTime()
  function testTime(): void {
    const n = 1e5
    const q = 1e5
    const nums2 = Array(n)
      .fill(0)
      .map((_, i) => i)
    const queries2 = Array(q)
    for (let i = 0; i < q; ++i) {
      const start = Math.floor(Math.random() * n)
      const end = Math.floor(Math.random() * n)
      queries2[i] = [Math.min(start, end), Math.max(start, end)]
    }
    const M2 = new MoOfflineAgain(n, q)
    queries2.forEach(([start, end]) => M2.addQuery(start, end))
    const bitLike2 = new BITRangeBlockFastQuery(n + 10)
    console.time('perf')
    M2.run(
      (index: number) => {
        bitLike2.add(nums2[index], 1)
      },
      (index: number) => {
        const res = bitLike2.queryRange(0, nums2[index])
        return res
      },
      (index: number) => {
        const res = bitLike2.queryRange(nums2[index] + 1, n)
        return res
      }
    )
    console.timeEnd('perf') // perf: 498.812ms
  }
}
