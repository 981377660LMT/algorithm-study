/* eslint-disable no-inner-declarations */
// BucketSort 桶排序
// 桶排序是计数排序的升级版。它利用函数的映射关系，将数据分到有限数量的桶里，每个桶再分别排序。
// 分桶=>排序=>合并

class BucketSort {
  private readonly _buckets: number[][]

  // hash, 将数值映射到桶内.
  private readonly _bucketSize: number
  private readonly _offset: number
  //

  /**
   * 所有数的范围必须在 `[min, max]` 内.
   */
  constructor(min: number, max: number, bucketCount: number) {
    if (min > max) throw new Error('min must be less than or equal to max')
    if (bucketCount <= 0) throw new Error('bucketCount must be greater than 0')
    const diff = max - min
    this._buckets = Array(bucketCount)
    for (let i = 0; i < bucketCount; i++) this._buckets[i] = []
    this._bucketSize = Math.floor(diff / bucketCount) + 1
    this._offset = min
  }

  /**
   * 返回一个新的排序后的数组.
   */
  sorted(arr: ArrayLike<number>, compareFn: (a: number, b: number) => number): number[] {
    for (let i = 0; i < this._buckets.length; i++) this._buckets[i] = []
    for (let i = 0; i < arr.length; i++) {
      const hash = Math.floor((arr[i] - this._offset) / this._bucketSize)
      this._buckets[hash].push(arr[i])
    }

    // !桶内排序策略
    for (let i = 0; i < this._buckets.length; i++) this._buckets[i].sort(compareFn)

    const res = Array(arr.length)
    for (let i = 0, ptr = 0; i < this._buckets.length; i++) {
      const bucket = this._buckets[i]
      for (let j = 0; j < bucket.length; j++) {
        res[ptr++] = bucket[j]
      }
    }
    return res
  }
}

if (require.main === module) {
  const arr = [3, 1, 4, 1, 5, 99]
  const B = new BucketSort(-10, 100, 10)
  const sorted = B.sorted(arr, (a, b) => a - b)
  console.log(sorted)

  // https://leetcode.cn/problems/sort-an-array/
  // 912. 排序数组

  function sortArray(nums: number[]): number[] {
    const B = new BucketSort(-5e4, 5e4, 20)
    return B.sorted(nums, (a, b) => a - b)
  }
}
