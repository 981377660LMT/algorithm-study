/* eslint-disable no-inner-declarations */

class DigitSwap {
  private static readonly _base = (() => {
    const res = Array<number>(16)
    res[0] = 1
    for (let i = 1; i < res.length; i++) res[i] = res[i - 1] * 10
    return res
  })()

  private readonly _num: number
  private readonly _digits: number[]

  constructor(num: number) {
    if (num < 0) throw new Error('num must be non-negative')
    if (num > Number.MAX_SAFE_INTEGER) throw new Error('num must be no more than 2^53-1')
    this._num = num
    this._digits = []
    let x = num
    while (x > 0) {
      this._digits.push(x % 10)
      x = Math.floor(x / 10)
    }
  }

  /**
   * 交换数字的两个位置.
   *
   * @example
   * ```ts
   * const D = new DigitSwap(1234)
   * D.swap(0, 1) // 1234 -> 1243
   * ```
   */
  swap(i: number, j: number): number {
    if (i === j || this._digits[i] === this._digits[j]) return this._num
    return (
      this._num + (this._digits[j] - this._digits[i]) * (DigitSwap._base[i] - DigitSwap._base[j])
    )
  }

  at(i: number): number {
    return this._digits[i]
  }

  get length(): number {
    return this._digits.length
  }
}

export { DigitSwap, DigitSwap as SwapDigit }

if (require.main === module) {
  // 3267. 统计近似相等数对 II (ts有8秒时限)
  // https://leetcode.cn/problems/count-almost-equal-pairs-ii/
  // 给你一个正整数数组 nums 。
  // 如果我们执行以下操作 `至多两次` 可以让两个整数 x 和 y 相等，那么我们称这个数对是 近似相等 的：
  // 选择 x 或者 y  之一，将这个数字中的两个数位交换。
  // 请你返回 nums 中，下标 i 和 j 满足 i < j 且 nums[i] 和 nums[j] 近似相等 的数对数目。
  // 注意 ，执行操作后得到的整数可以有前导 0 。
  //
  // !先对原数组排序，根据数组中的数的数位个数，按小到大排序
  function countPairs(nums: number[]): number {
    const swap2 = (x: number): Set<number> => {
      const res = new Set<number>([x])
      let queue = new Set<number>([x])
      for (let t = 0; t < 2; t++) {
        const nextQueue = new Set<number>()
        queue.forEach(num => {
          const D = new DigitSwap(num)
          for (let i = 0; i < D.length; i++) {
            for (let j = i + 1; j < D.length; j++) {
              const next = D.swap(i, j)
              nextQueue.add(next)
              res.add(next)
            }
          }
        })
        queue = nextQueue
      }
      return res
    }

    nums.sort((a, b) => a - b) // 排序后，一定是换后面那个数
    let res = 0
    const preCounter = new Map<number, number>()
    for (let i = 0; i < nums.length; i++) {
      const swapRes = swap2(nums[i])
      // eslint-disable-next-line no-loop-func
      swapRes.forEach(v => {
        res += preCounter.get(v) || 0
      })
      preCounter.set(nums[i], (preCounter.get(nums[i]) || 0) + 1)
    }
    return res
  }
}
