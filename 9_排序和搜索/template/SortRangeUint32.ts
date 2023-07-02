/**
 * 所有元素都在[0, 2^32)范围内的数组排序.
 */
class SortRangeUint32 {
  private readonly _origin: Uint32Array
  private readonly _worker: Uint32Array

  constructor(nums: Uint32Array) {
    this._origin = nums.slice()
    this._worker = nums.slice()
  }

  /**
   * 返回一个新的排序后的数组.
   */
  sorted(start = 0, end = this._worker.length, reverse = false): number[] {
    const res = Array(end - start).fill(0)
    this._worker.subarray(start, end).sort()
    if (reverse) {
      for (let i = 0; i < end - start; i++) res[i] = this._worker[end - i - 1]
    } else {
      for (let i = 0; i < end - start; i++) res[i] = this._worker[start + i]
    }
    this._worker.set(this._origin.subarray(start, end), start)
    return res
  }
}

export {}

if (require.main === module) {
  const arr = new Uint32Array([1, 4, 2, 5, 3, 6, 7])
  const sorter = new SortRangeUint32(arr)
  console.log(sorter.sorted(0, 4, true))
  console.log(sorter.sorted())
}
