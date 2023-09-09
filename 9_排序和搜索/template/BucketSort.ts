/* eslint-disable no-shadow */
/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */

// BucketSort 桶排序
// 桶排序是计数排序的升级版。它利用函数的映射关系，将数据分到有限数量的桶里，每个桶再分别排序。
// 分桶=>排序=>合并
// https://zhuanlan.zhihu.com/p/649120343

import { sortRange } from './sortRange'
import { timSort } from './timSort'

/**
 * 如果值域范围不大(<=2e7), 采用 {@link CountingSort} 更快.
 * 如果值域范围很大，且元素个数较多(>5e6)，可以用桶排序分桶避免大数组的排序.
 * @param bucketCount 桶的数量.
 * @param sorter 桶内排序策略.默认为 {@link Array.prototype.sort}.
 * @returns 返回一个`新的`排序后的数组.
 */
function bucketSort(
  arr: ArrayLike<number>,
  compareFn: (a: number, b: number) => number,
  options?: {
    bucketCount?: number
    sorter?: (arr: number[], compareFn: (a: number, b: number) => number) => number[] | void
  }
): number[] {
  const n = arr.length
  if (n <= 1) return Array.from(arr)

  let { bucketCount, sorter } = options ?? {}

  let min = arr[0]
  let max = arr[0]
  for (let i = 1; i < n; i++) {
    min = Math.min(min, arr[i])
    max = Math.max(max, arr[i])
  }
  const diff = max - min

  if (bucketCount === void 0) {
    const cand1 = Math.floor(diff / n)
    const cand2 = Math.ceil(n / 100)
    bucketCount = Math.min(1e5, Math.max(cand1, cand2))
  }

  const buckets = Array(bucketCount)
  for (let i = 0; i < bucketCount; i++) buckets[i] = []
  const gap = Math.floor(diff / bucketCount) + 1 // 桶的大小
  for (let i = 0; i < n; i++) {
    const id = ((arr[i] - min) / gap) | 0
    buckets[id].push(arr[i])
  }

  // !桶内排序策略
  if (sorter === void 0) {
    for (let i = 0; i < bucketCount; i++) {
      buckets[i].sort(compareFn)
    }
  } else {
    for (let i = 0; i < bucketCount; i++) {
      const res = sorter(buckets[i], compareFn)
      if (res !== void 0) buckets[i] = res
    }
  }

  const res = Array(n)
  for (let i = 0, ptr = 0; i < bucketCount; i++) {
    const bucket = buckets[i]
    for (let j = 0; j < bucket.length; j++) {
      res[ptr++] = bucket[j]
    }
  }
  return res
}

export { bucketSort }

if (require.main === module) {
  // https://leetcode.cn/problems/sort-an-array/
  // 912. 排序数组

  function sortArray(nums: number[]): number[] {
    return bucketSort(nums, (a, b) => a - b)
  }

  const bigArr1e7 = Array(1e7)
  for (let i = 0; i < bigArr1e7.length; i++) bigArr1e7[i] = (Math.random() * 1e9) | 0
  const copy1 = bigArr1e7.slice()
  const copy2 = bigArr1e7.slice()
  const copy3 = bigArr1e7.slice()
  const copy4 = bigArr1e7.slice()

  console.time('BucketSort')
  bucketSort(copy1, (a, b) => a - b)
  console.timeEnd('BucketSort')

  console.time('Array.sort')
  copy2.sort((a, b) => a - b)
  console.timeEnd('Array.sort')

  // timeSort
  console.time('TimSort')
  timSort(copy3)
  console.timeEnd('TimSort')

  console.time('pdqSort')
  bucketSort(copy4, (a, b) => a - b, { sorter: sortRange })
  console.timeEnd('pdqSort')
}
