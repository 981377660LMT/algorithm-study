// 设计一个数据结构,有效地找到给定子数组的 出现 threshold 次数或次数以上的元素 。
// 1 <= arr.length <= 2e4
// threshold <= right - left + 1
// 2 * threshold > right - left + 1
// 调用 query 的次数最多为 1e4

// !按照`区间长度`根号分治
// !针对不同的询问区间长度，使用两种不同的方法。
// 记 SQRT = sqrt(2n)
// - 区间长度小于 SQRT ，使用暴力计算
// - 区间长度大于 SQRT ，则绝对众数出现次数 大于 SQRT/2
//   可能的候选人个数不超过 2n/SQRT ，使用前缀和统计频率大于SQRT/2的数的出现次数
// 考虑到常数，选择4*SQRT作为阈值.

class MajorityChecker {
  private readonly _arr: number[]
  private readonly _threshold: number
  private readonly _preSumRecord: Map<number, Uint32Array>

  constructor(arr: number[]) {
    const n = arr.length
    this._arr = arr.slice()
    this._threshold = 4 * (Math.sqrt(2 * n) | 0)
    this._preSumRecord = new Map()

    const counter = new Map<number, number>()
    arr.forEach(num => {
      counter.set(num, (counter.get(num) || 0) + 1)
    })

    const half = this._threshold >>> 1
    counter.forEach((count, num) => {
      if (count > half) {
        const preSum = new Uint32Array(n + 1)
        for (let i = 0; i < n; ++i) {
          preSum[i + 1] = preSum[i] + +(arr[i] === num)
        }
        this._preSumRecord.set(num, preSum)
      }
    })
  }

  query(left: number, right: number, threshold: number): number {
    const len = right - left + 1
    if (len <= this._threshold) {
      const counter = new Map<number, number>()
      for (let i = left; i <= right; ++i) {
        counter.set(this._arr[i], (counter.get(this._arr[i]) || 0) + 1)
      }
      for (const [num, count] of counter) {
        if (count >= threshold) {
          return num
        }
      }
      return -1
    }

    for (const [num, preSum] of this._preSumRecord) {
      const count = preSum[right + 1] - preSum[left]
      if (count >= threshold) {
        return num
      }
    }
    return -1
  }
}
