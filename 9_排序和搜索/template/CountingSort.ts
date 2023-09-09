/* eslint-disable no-inner-declarations */

// CountingSort 计数排序
// 如果要对其他类型的数值进行排序，要先进行预处理。

class CountingSort {
  private readonly _counter: Uint32Array

  /**
   * {@link Uint32Array.prototype.fill} 是遍历数组清除的5倍速度左右.
   * 当数组长度超过这个值时, 用 {@link Uint32Array.prototype.fill} 清除.
   * 否则用遍历清除.
   */
  private readonly _clearThreshold: number

  /**
   * 所有数的范围必须在 `[min, min + gap]` 之间.
   * gap 不能超过 2e7.
   */
  constructor(gap: number) {
    if (gap > 2e7) throw new Error('gap must be less than or equal to 2e7')
    this._counter = new Uint32Array(gap + 1)
    this._clearThreshold = ((gap + 5) / 5) | 0
  }

  /**
   * 返回一个新的排序后的数组.
   * @complexity `O(max+n)`
   */
  sorted(arr: ArrayLike<number>, reverse = false): number[] {
    const n = arr.length
    if (n <= 1) return Array.from(arr)
    const res = Array(n)
    const counter = this._counter
    let min = arr[0]
    for (let i = 1; i < n; i++) min = Math.min(min, arr[i])
    for (let i = 0; i < n; i++) counter[arr[i] - min]++

    if (reverse) {
      for (let i = counter.length - 1, ptr = 0; ~i; i--) {
        for (let j = 0; j < counter[i]; j++) res[ptr++] = i + min
      }
    } else {
      for (let i = 0, ptr = 0; i < counter.length; i++) {
        for (let j = 0; j < counter[i]; j++) res[ptr++] = i + min
      }
    }

    if (n >= this._clearThreshold) counter.fill(0)
    else for (let i = 0; i < n; i++) counter[arr[i] - min]--

    return res
  }
}

export { CountingSort }

if (require.main === module) {
  const B = new CountingSort(100)
  const arr = [1, 2, 3, 4, 5, 6, 7, 8, 9]
  console.log(B.sorted(arr, true))

  const asass = new Uint32Array(1e7)
  console.time('clear')
  for (let i = 0; i < 1e3; i++) asass.fill(0)
  console.timeEnd('clear')

  const curNums = Array(2e6).fill(0)
  console.time('curNums')
  for (let i = 0; i < 1e3; i++) {
    for (let j = 0; j < curNums.length; j++) curNums[j]--
  }
  console.timeEnd('curNums')

  // https://leetcode.cn/problems/sort-an-array/
  // 912. 排序数组
  const C = new CountingSort(1e5)
  function sortArray(nums: number[]): number[] {
    return C.sorted(nums)
  }

  // https://leetcode.cn/problems/sum-of-imbalance-numbers-of-all-subarrays/
  // 6894. 所有子数组中不平衡数字之和
  function sumImbalanceNumbers(N: number[]): number {
    const nums = new Uint16Array(N)
    const C = new CountingSort(Math.max(...nums) - Math.min(...nums))
    let res = 0
    for (let i = 0; i < nums.length; i++) {
      for (let j = i; j < nums.length; j++) {
        const sorted = C.sorted(nums.subarray(i, j + 1))
        for (let k = 0; k < sorted.length - 1; k++) res += +(sorted[k + 1] - sorted[k] > 1)
      }
    }
    return res
  }
}
