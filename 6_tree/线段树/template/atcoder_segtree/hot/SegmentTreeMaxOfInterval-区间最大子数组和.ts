/* eslint-disable max-len */
/* eslint-disable no-cond-assign */
/* eslint-disable no-inner-declarations */

// SegmentTreeMaxOfInterval-区间最大子数组和

import { createPointSetRangeMaxSumMinSum } from '../SegmentTreeUtils'

export { createPointSetRangeMaxSumMinSum as SegmentTreeMaxOfInterval }

if (require.main === module) {
  const nums = [1, 2, 3, 4, 5]
  const { tree, fromElement } = createPointSetRangeMaxSumMinSum(nums)
  console.log(tree.query(0, 5))

  timeit()
  function timeit(): void {
    const n = 2e5
    const arr = Array(n)
    for (let i = 0; i < n; i++) arr[i] = Math.floor(Math.random() * 10)
    const { tree, fromElement } = createPointSetRangeMaxSumMinSum(arr)
    console.time('PointSetRangeMaxSumMinSum')
    for (let i = 0; i < n; i++) {
      tree.query(i, n)
      tree.update(i, fromElement(1))
      tree.set(i, fromElement(1))
      tree.maxRight(i, interval => interval.minSum >= i)
      tree.minLeft(i, interval => interval.minSum >= i)
    }
    console.timeEnd('PointSetRangeMaxSumMinSum') // PointSetRangeMaxSumMinSum: 600.276ms
  }

  function maxSubArray(nums: number[]): number {
    const { tree } = createPointSetRangeMaxSumMinSum(nums)
    return tree.queryAll().maxSum
  }
}
