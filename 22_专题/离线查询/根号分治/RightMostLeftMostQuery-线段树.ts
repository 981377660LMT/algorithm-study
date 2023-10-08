/* eslint-disable no-inner-declarations */

// 对每个下标，查询 最右侧/最左侧/右侧第一个/左侧第一个 lower/floor/ceiling/higher 的元素.
// 动态单调栈(DynamicMonoStack).
// 线段树实现(非常快).

import { SegmentTreeRangeAddRangeMinMax } from '../../../6_tree/线段树/template/atcoder_segtree/hot/SegmentTreeRangeAddRangeMinMax-区间加区间最大最小值'

class RightMostLeftMostQuerySegmentTree {
  private readonly _n: number
  private readonly _rangeAddRangeMinMax: SegmentTreeRangeAddRangeMinMax

  constructor(arr: ArrayLike<number>) {
    this._n = arr.length
    this._rangeAddRangeMinMax = new SegmentTreeRangeAddRangeMinMax(arr)
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._n) return
    this._rangeAddRangeMinMax.set(index, value)
  }

  addRange(start: number, end: number, delta: number): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    this._rangeAddRangeMinMax.update(start, end, delta)
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] < nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostLower(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.minLeft(this._n, min => min >= cur) - 1
    return cand > index ? cand : -1
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] <= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostFloor(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.minLeft(this._n, min => min > cur) - 1
    return cand > index ? cand : -1
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] >= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostCeiling(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.minLeft(this._n, (_, max) => max < cur) - 1
    return cand > index ? cand : -1
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostHigher(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.minLeft(this._n, (_, max) => max <= cur) - 1
    return cand > index ? cand : -1
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] < nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostLower(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.maxRight(0, min => min >= cur)
    return cand < index ? cand : -1
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] <= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostFloor(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.maxRight(0, min => min > cur)
    return cand < index ? cand : -1
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] >= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostCeiling(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.maxRight(0, (_, max) => max < cur)
    return cand < index ? cand : -1
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostHigher(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.maxRight(0, (_, max) => max <= cur)
    return cand < index ? cand : -1
  }

  /**
   * 查询`index`右侧最近的下标`j`，使得 `nums[j] < nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightNearestLower(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.maxRight(index + 1, min => min >= cur)
    return cand === this._n ? -1 : cand
  }

  /**
   * 查询`index`右侧最近的下标`j`，使得 `nums[j] <= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightNearestFloor(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.maxRight(index + 1, min => min > cur)
    return cand === this._n ? -1 : cand
  }

  /**
   * 查询`index`右侧最近的下标`j`，使得 `nums[j] >= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightNearestCeiling(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.maxRight(index + 1, (_, max) => max < cur)
    return cand === this._n ? -1 : cand
  }

  /**
   * 查询`index`右侧最近的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightNearestHigher(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.maxRight(index + 1, (_, max) => max <= cur)
    return cand === this._n ? -1 : cand
  }

  /**
   * 查询`index`左侧最近的下标`j`，使得 `nums[j] < nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftNearestLower(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.minLeft(index, min => min >= cur) - 1
    return cand
  }

  /**
   * 查询`index`左侧最近的下标`j`，使得 `nums[j] <= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftNearestFloor(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.minLeft(index, min => min > cur) - 1
    return cand
  }

  /**
   * 查询`index`左侧最近的下标`j`，使得 `nums[j] >= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftNearestCeiling(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.minLeft(index, (_, max) => max < cur) - 1
    return cand
  }

  /**
   * 查询`index`左侧最近的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftNearestHigher(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.minLeft(index, (_, max) => max <= cur) - 1
    return cand
  }
}

if (require.main === module) {
  // 962. 最大宽度坡
  // https://leetcode.cn/problems/maximum-width-ramp/
  function maxWidthRamp(nums: number[]): number {
    const Q = new RightMostLeftMostQuerySegmentTree(nums)
    let res = 0
    for (let i = 0; i < nums.length; i++) {
      const cand = Q.rightMostCeiling(i)
      if (cand !== -1) {
        res = Math.max(res, cand - i)
      }
    }
    return res
  }

  // 2863. 最长半递减数组
  // https://leetcode.cn/problems/maximum-length-of-semi-decreasing-subarrays/description/
  function maxSubarrayLength(nums: number[]): number {
    const Q = new RightMostLeftMostQuerySegmentTree(nums)
    let res = 0
    for (let i = 0; i < nums.length; i++) {
      const cand = Q.rightMostLower(i)
      if (cand !== -1) {
        res = Math.max(res, cand - i + 1)
      }
    }
    return res
  }

  // 901. 股票价格跨度
  // 当日股票价格的 跨度 被定义为股票价格小于或等于今天价格的最大连续日数
  class StockSpanner {
    private readonly _Q: RightMostLeftMostQuerySegmentTree
    private _ptr = 0

    constructor() {
      this._Q = new RightMostLeftMostQuerySegmentTree(Array(1e5 + 10).fill(0))
    }

    next(price: number): number {
      const pos = this._ptr++
      this._Q.set(pos, price)
      const higherPos = this._Q.leftNearestHigher(pos)
      return higherPos === -1 ? pos + 1 : pos - higherPos
    }
  }

  // checkWithBruteForce()
  function checkWithBruteForce(): void {
    class Mocker {
      readonly _nums: number[]
      constructor(nums: number[]) {
        this._nums = nums
      }

      set(index: number, value: number): void {
        this._nums[index] = value
      }

      addRange(start: number, end: number, delta: number): void {
        for (let i = start; i < end; i++) this._nums[i] += delta
      }

      rightMostLower(index: number): number {
        for (let i = this._nums.length - 1; i > index; i--) {
          if (this._nums[i] < this._nums[index]) return i
        }
        return -1
      }

      rightMostFloor(index: number): number {
        for (let i = this._nums.length - 1; i > index; i--) {
          if (this._nums[i] <= this._nums[index]) return i
        }
        return -1
      }

      rightMostCeiling(index: number): number {
        for (let i = this._nums.length - 1; i > index; i--) {
          if (this._nums[i] >= this._nums[index]) return i
        }
        return -1
      }

      rightMostHigher(index: number): number {
        for (let i = this._nums.length - 1; i > index; i--) {
          if (this._nums[i] > this._nums[index]) return i
        }
        return -1
      }

      leftMostLower(index: number): number {
        for (let i = 0; i < index; i++) {
          if (this._nums[i] < this._nums[index]) return i
        }
        return -1
      }

      leftMostFloor(index: number): number {
        for (let i = 0; i < index; i++) {
          if (this._nums[i] <= this._nums[index]) return i
        }
        return -1
      }

      leftMostCeiling(index: number): number {
        for (let i = 0; i < index; i++) {
          if (this._nums[i] >= this._nums[index]) return i
        }
        return -1
      }

      leftMostHigher(index: number): number {
        for (let i = 0; i < index; i++) {
          if (this._nums[i] > this._nums[index]) return i
        }
        return -1
      }

      rightNearestLower(index: number): number {
        for (let i = index + 1; i < this._nums.length; i++) {
          if (this._nums[i] < this._nums[index]) return i
        }
        return -1
      }

      rightNearestFloor(index: number): number {
        for (let i = index + 1; i < this._nums.length; i++) {
          if (this._nums[i] <= this._nums[index]) return i
        }
        return -1
      }

      rightNearestCeiling(index: number): number {
        for (let i = index + 1; i < this._nums.length; i++) {
          if (this._nums[i] >= this._nums[index]) return i
        }
        return -1
      }

      rightNearestHigher(index: number): number {
        for (let i = index + 1; i < this._nums.length; i++) {
          if (this._nums[i] > this._nums[index]) return i
        }
        return -1
      }

      leftNearestLower(index: number): number {
        for (let i = index - 1; i >= 0; i--) {
          if (this._nums[i] < this._nums[index]) return i
        }
        return -1
      }

      leftNearestFloor(index: number): number {
        for (let i = index - 1; i >= 0; i--) {
          if (this._nums[i] <= this._nums[index]) return i
        }
        return -1
      }

      leftNearestCeiling(index: number): number {
        for (let i = index - 1; i >= 0; i--) {
          if (this._nums[i] >= this._nums[index]) return i
        }
        return -1
      }

      leftNearestHigher(index: number): number {
        for (let i = index - 1; i >= 0; i--) {
          if (this._nums[i] > this._nums[index]) return i
        }
        return -1
      }
    }

    const N = 1e5
    const MAX = 1e9
    const randomNums = Array.from({ length: N }, () => (Math.random() * MAX) | 0)
    const mocker = new Mocker(randomNums.slice())
    const real = new RightMostLeftMostQuerySegmentTree(randomNums)
    const debug = (
      mockerFunc: (index: number) => number,
      realFunc: (index: number) => number
    ): void => {
      const index = randint(0, randomNums.length - 1)
      const mockerRes = mockerFunc(index)
      const QRes = realFunc(index)
      if (mockerRes !== QRes) {
        console.log(realFunc.name, index, mockerRes, QRes)
        // console.log(mocker._nums, real._nums, real._blockLazy)
        // console.log('blockMin', real._blockMin)
        // console.log('blockMax', real._blockMax)
        throw new Error(realFunc.name)
      }
    }

    const randint = (a: number, b: number) => Math.floor(Math.random() * (b - a + 1)) + a
    for (let i = 0; i < 1e5; i++) {
      const op = randint(0, 17)
      // set
      if (op === 0) {
        const index = randint(0, randomNums.length - 1)
        const value = (Math.random() * MAX) | 0
        mocker.set(index, value)
        real.set(index, value)
      } else if (op === 1) {
        // addRange
        const start = randint(0, randomNums.length - 1)
        const end = randint(start, randomNums.length - 1)
        const delta = (Math.random() * MAX) | 0
        mocker.addRange(start, end, delta)
        real.addRange(start, end, delta)
      } else if (op === 2) {
        debug(mocker.rightMostLower.bind(mocker), real.rightMostLower.bind(real))
      } else if (op === 3) {
        debug(mocker.rightMostFloor.bind(mocker), real.rightMostFloor.bind(real))
      } else if (op === 4) {
        debug(mocker.rightMostCeiling.bind(mocker), real.rightMostCeiling.bind(real))
      } else if (op === 5) {
        debug(mocker.rightMostHigher.bind(mocker), real.rightMostHigher.bind(real))
      } else if (op === 6) {
        debug(mocker.leftMostLower.bind(mocker), real.leftMostLower.bind(real))
      } else if (op === 7) {
        debug(mocker.leftMostFloor.bind(mocker), real.leftMostFloor.bind(real))
      } else if (op === 8) {
        debug(mocker.leftMostCeiling.bind(mocker), real.leftMostCeiling.bind(real))
      } else if (op === 9) {
        debug(mocker.leftMostHigher.bind(mocker), real.leftMostHigher.bind(real))
      } else if (op === 10) {
        debug(mocker.rightNearestLower.bind(mocker), real.rightNearestLower.bind(real))
      } else if (op === 11) {
        debug(mocker.rightNearestFloor.bind(mocker), real.rightNearestFloor.bind(real))
      } else if (op === 12) {
        debug(mocker.rightNearestCeiling.bind(mocker), real.rightNearestCeiling.bind(real))
      } else if (op === 13) {
        debug(mocker.rightNearestHigher.bind(mocker), real.rightNearestHigher.bind(real))
      } else if (op === 14) {
        debug(mocker.leftNearestLower.bind(mocker), real.leftNearestLower.bind(real))
      } else if (op === 15) {
        debug(mocker.leftNearestFloor.bind(mocker), real.leftNearestFloor.bind(real))
      } else if (op === 16) {
        debug(mocker.leftNearestCeiling.bind(mocker), real.leftNearestCeiling.bind(real))
      } else if (op === 17) {
        debug(mocker.leftNearestHigher.bind(mocker), real.leftNearestHigher.bind(real))
      }
    }

    console.log('ok')
  }

  testTime()
  function testTime(): void {
    const bigArr = Array(1e5)
      .fill(0)
      .map(() => (Math.random() * 1e9) | 0)
    const Q = new RightMostLeftMostQuerySegmentTree(bigArr)
    console.time('bigArr')
    for (let i = 0; i < 1e5; i++) {
      Q.set(i, i)
      Q.addRange(i, i + 1, i)
      Q.rightMostLower(i)
      Q.rightMostFloor(i)
      Q.rightMostCeiling(i)
      Q.rightMostHigher(i)
      Q.leftMostLower(i)
      Q.leftMostFloor(i)
      Q.leftMostCeiling(i)
      Q.leftMostHigher(i)
      Q.rightNearestLower(i)
      Q.rightNearestFloor(i)
      Q.rightNearestCeiling(i)
      Q.rightNearestHigher(i)
      Q.leftNearestLower(i)
      Q.leftNearestFloor(i)
      Q.leftNearestCeiling(i)
      Q.leftNearestHigher(i)
    }
    console.timeEnd('bigArr') // bigArr: 330.583ms (比分块的2s快)
  }
}
