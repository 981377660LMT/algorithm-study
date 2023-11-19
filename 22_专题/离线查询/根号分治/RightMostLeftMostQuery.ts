/* eslint-disable no-inner-declarations */

// class MonoStack
// 对每个下标，查询 最右侧/最左侧/右侧第一个/左侧第一个 lower/floor/ceiling/higher 的元素.
// 动态单调栈(DynamicMonoStack).
// 分块实现.

class MonoStackDynamic {
  private readonly _nums: number[]
  private readonly _belong: Uint16Array
  private readonly _blockStart: Uint32Array
  private readonly _blockEnd: Uint32Array
  private readonly _blockCount: number

  private readonly _blockMin: number[]
  private readonly _blockMax: number[]
  private readonly _blockLazy: number[]

  constructor(arr: ArrayLike<number>) {
    const copy = Array(arr.length)
    for (let i = 0; i < arr.length; i++) copy[i] = arr[i]
    this._nums = copy

    const n = arr.length
    const blockSize = (Math.sqrt(n) + 1) | 0
    const blockCount = 1 + ((n / blockSize) | 0)
    const blockStart = new Uint32Array(blockCount)
    const blockEnd = new Uint32Array(blockCount)
    const belong = new Uint16Array(n)
    for (let i = 0; i < blockCount; i++) {
      blockStart[i] = i * blockSize
      blockEnd[i] = Math.min((i + 1) * blockSize, n)
    }
    for (let i = 0; i < n; i++) {
      belong[i] = (i / blockSize) | 0
    }

    this._belong = belong
    this._blockStart = blockStart
    this._blockEnd = blockEnd
    this._blockCount = blockCount
    this._blockMin = Array(blockCount).fill(Infinity)
    this._blockMax = Array(blockCount).fill(-Infinity)
    this._blockLazy = Array(blockCount).fill(0)
    for (let bid = 0; bid < blockCount; bid++) this._rebuildBlock(bid)
  }

  get(index: number): number | undefined {
    if (index < 0 || index >= this._nums.length) return undefined
    return this._nums[index] + this._blockLazy[this._belong[index]]
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._nums.length) return
    const bid = this._belong[index]
    const lazy = this._blockLazy[bid]
    const pre = this._nums[index] + lazy
    if (pre === value) return
    this._nums[index] = value - lazy
    this._rebuildBlock(bid)
  }

  addRange(start: number, end: number, delta: number): void {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    if (bid1 === bid2) {
      for (let i = start; i < end; i++) this._nums[i] += delta
      this._rebuildBlock(bid1)
    } else {
      for (let i = start; i < this._blockEnd[bid1]; i++) this._nums[i] += delta
      this._rebuildBlock(bid1)
      for (let bid = bid1 + 1; bid < bid2; bid++) this._blockLazy[bid] += delta
      for (let i = this._blockStart[bid2]; i < end; i++) this._nums[i] += delta
      this._rebuildBlock(bid2)
    }
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] < nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostLower(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryRightMost(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] < target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] < target!
    )
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] <= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostFloor(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryRightMost(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] <= target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] <= target!
    )
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] >= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostCeiling(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryRightMost(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] >= target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] >= target!
    )
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostHigher(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryRightMost(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] > target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] > target!
    )
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] < nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostLower(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryLeftMost(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] < target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] < target!
    )
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] <= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostFloor(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryLeftMost(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] <= target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] <= target!
    )
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] >= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostCeiling(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryLeftMost(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] >= target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] >= target!
    )
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostHigher(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryLeftMost(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] > target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] > target!
    )
  }

  /**
   * 查询`index`右侧最近的下标`j`，使得 `nums[j] < nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightNearestLower(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryRightNearest(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] < target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] < target!
    )
  }

  /**
   * 查询`index`右侧最近的下标`j`，使得 `nums[j] <= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightNearestFloor(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryRightNearest(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] <= target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] <= target!
    )
  }

  /**
   * 查询`index`右侧最近的下标`j`，使得 `nums[j] >= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightNearestCeiling(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryRightNearest(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] >= target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] >= target!
    )
  }

  /**
   * 查询`index`右侧最近的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightNearestHigher(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryRightNearest(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] > target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] > target!
    )
  }

  /**
   * 查询`index`左侧最近的下标`j`，使得 `nums[j] < nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftNearestLower(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryLeftNearest(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] < target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] < target!
    )
  }

  /**
   * 查询`index`左侧最近的下标`j`，使得 `nums[j] <= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftNearestFloor(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryLeftNearest(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] <= target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] <= target!
    )
  }

  /**
   * 查询`index`左侧最近的下标`j`，使得 `nums[j] >= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftNearestCeiling(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryLeftNearest(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] >= target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] >= target!
    )
  }

  /**
   * 查询`index`左侧最近的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftNearestHigher(index: number, target?: number): number {
    if (target == undefined) {
      target = this.get(index)
    }
    return this._queryLeftNearest(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] > target!,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] > target!
    )
  }

  private _queryRightMost(
    pos: number,
    predicateBlock: (bid: number) => boolean,
    predicateElement: (eid: number, bid: number) => boolean
  ): number {
    const bid = this._belong[pos]
    for (let i = this._blockCount - 1; i > bid; i--) {
      if (!predicateBlock(i)) continue
      for (let j = this._blockEnd[i] - 1; j >= this._blockStart[i]; j--) {
        if (predicateElement(j, i)) return j
      }
    }
    for (let i = this._blockEnd[bid] - 1; i > pos; i--) {
      if (predicateElement(i, bid)) return i
    }
    return -1
  }

  private _queryLeftMost(
    pos: number,
    predicateBlock: (bid: number) => boolean,
    predicateElement: (eid: number, bid: number) => boolean
  ): number {
    const bid = this._belong[pos]
    for (let i = 0; i < bid; i++) {
      if (!predicateBlock(i)) continue
      for (let j = this._blockStart[i]; j < this._blockEnd[i]; j++) {
        if (predicateElement(j, i)) return j
      }
    }
    for (let i = this._blockStart[bid]; i < pos; i++) {
      if (predicateElement(i, bid)) return i
    }
    return -1
  }

  private _queryRightNearest(
    pos: number,
    predicateBlock: (bid: number) => boolean,
    predicateElement: (eid: number, bid: number) => boolean
  ): number {
    const bid = this._belong[pos]
    for (let i = pos + 1; i < this._blockEnd[bid]; i++) {
      if (predicateElement(i, bid)) return i
    }
    for (let i = bid + 1; i < this._blockCount; i++) {
      if (!predicateBlock(i)) continue
      for (let j = this._blockStart[i]; j < this._blockEnd[i]; j++) {
        if (predicateElement(j, i)) return j
      }
    }
    return -1
  }

  private _queryLeftNearest(
    pos: number,
    predicateBlock: (bid: number) => boolean,
    predicateElement: (eid: number, bid: number) => boolean
  ): number {
    const bid = this._belong[pos]
    for (let i = pos - 1; i >= this._blockStart[bid]; i--) {
      if (predicateElement(i, bid)) return i
    }
    for (let i = bid - 1; i >= 0; i--) {
      if (!predicateBlock(i)) continue
      for (let j = this._blockEnd[i] - 1; j >= this._blockStart[i]; j--) {
        if (predicateElement(j, i)) return j
      }
    }
    return -1
  }

  private _rebuildBlock(bid: number): void {
    this._blockMin[bid] = Infinity
    this._blockMax[bid] = -Infinity
    const lazy = this._blockLazy[bid]
    this._blockLazy[bid] = 0
    for (let i = this._blockStart[bid]; i < this._blockEnd[bid]; i++) {
      this._nums[i] += lazy
      this._blockMin[bid] = Math.min(this._blockMin[bid], this._nums[i])
      this._blockMax[bid] = Math.max(this._blockMax[bid], this._nums[i])
    }
  }
}

export {
  MonoStackDynamic,
  MonoStackDynamic as MomoStack,
  MonoStackDynamic as RightMostLeftMostQuery
}

if (require.main === module) {
  // 962. 最大宽度坡
  // https://leetcode.cn/problems/maximum-width-ramp/
  function maxWidthRamp(nums: number[]): number {
    const Q = new MonoStackDynamic(nums)
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
    const Q = new MonoStackDynamic(nums)
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
    private readonly _Q: MonoStackDynamic
    private _ptr = 0

    constructor() {
      this._Q = new MonoStackDynamic(Array(1e5 + 10).fill(0))
    }

    next(price: number): number {
      const pos = this._ptr++
      this._Q.set(pos, price)
      const higherPos = this._Q.leftNearestHigher(pos)
      return higherPos === -1 ? pos + 1 : pos - higherPos
    }
  }

  checkWithBruteForce()
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
    const real = new MonoStackDynamic(randomNums)
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
    const Q = new MonoStackDynamic(bigArr)
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
    console.timeEnd('bigArr') // bigArr: 2.088s
  }
}
