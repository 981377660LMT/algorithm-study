/* eslint-disable no-inner-declarations */
/* eslint-disable no-else-return */

import { SortedListFastWithSum } from '../../../22_专题/离线查询/根号分治/SortedList/SortedListWithSum'

/**
 * SortedList 动态维护中位数信息.
 * `Proxy 的内部持有一个对 SortedList 的引用`.
 */
class MedianFinderSortedList {
  readonly sortedList: SortedListFastWithSum<number>

  /**
   * @param sortedListRef 传入一个已经初始化好的 SortedList 引用.
   * 如果不传入, 则内部会初始化一个.
   */
  constructor(sortedListRef?: SortedListFastWithSum<number>) {
    this.sortedList = sortedListRef || new SortedListFastWithSum<number>()
  }

  /**
   * 返回向下取整的中位数.
   */
  median(): number | undefined {
    const len = this.sortedList.length
    if (!len) return undefined
    if (len & 1) {
      return this.sortedList.at(len >>> 1)!
    } else {
      const mid = len >>> 1
      return Math.floor((this.sortedList.at(mid - 1)! + this.sortedList.at(mid)!) / 2)
    }
  }

  /**
   * 返回所有数到`to`的距离和.
   */
  distSum(to: number): number {
    const sl = this.sortedList
    const pos = sl.bisectRight(to)
    const allSum = sl.sumAll()
    let sum1: number
    let sum2: number
    if (pos < sl.length >>> 1) {
      sum1 = sl.sumSlice(0, pos)
      sum2 = allSum - sum1
    } else {
      sum2 = sl.sumSlice(pos, sl.length)
      sum1 = allSum - sum2
    }
    const leftSum = to * pos - sum1
    const rightSum = sum2 - to * (sl.length - pos)
    return leftSum + rightSum
  }

  /**
   * 返回切片`[start,end)`中所有数到`to`的距离和.
   */
  distSumRange(to: number, start: number, end: number): number {
    const sl = this.sortedList
    if (start < 0) start = 0
    if (end > sl.length) end = sl.length
    if (start >= end) return 0
    const pos = sl.bisectLeft(to)
    if (pos <= start) {
      return sl.sumSlice(start, end) - to * (end - start)
    }
    if (pos >= end) {
      return to * (end - start) - sl.sumSlice(start, end)
    }
    const leftSum = to * (pos - start) - sl.sumSlice(start, pos)
    const rightSum = sl.sumSlice(pos, end) - to * (end - pos)
    return leftSum + rightSum
  }

  /**
   * 返回所有数到中位数的距离和.
   */
  distSumToMedian(): number {
    if (!this.sortedList.length) return 0
    return this.distSum(this.median()!)
  }

  /**
   * 返回切片`[start,end)`中所有数到中位数的距离和.
   */
  distSumToMedianRange(start: number, end: number): number {
    if (!this.sortedList.length) return 0
    return this.distSumRange(this.median()!, start, end)
  }
}

export { MedianFinderSortedList }

if (require.main === module) {
  // 100123. 执行操作使频率分数最大
  // https://leetcode.cn/problems/apply-operations-to-maximize-frequency-score/description/
  function maxFrequencyScore(nums: number[], k: number): number {
    nums.sort((a, b) => a - b)

    const sl = new SortedListFastWithSum<number>()
    const proxy = new MedianFinderSortedList(sl)

    let res = 0
    let left = 0
    for (let right = 0; right < nums.length; right++) {
      sl.add(nums[right])
      while (left <= right && proxy.distSumToMedian()! > k) {
        sl.discard(nums[left])
        left++
      }
      res = Math.max(res, right - left + 1)
    }

    return res
  }
}
