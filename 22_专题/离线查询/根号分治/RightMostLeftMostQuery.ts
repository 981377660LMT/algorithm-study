/* eslint-disable no-inner-declarations */

// 对每个下标，查询 最右侧/最左侧 lower/floor/ceiling/higher 的元素.

class RightMostLeftMostQuery {
  readonly _nums: number[]
  readonly _belong: Uint16Array
  private readonly _blockStart: Uint32Array
  private readonly _blockEnd: Uint32Array
  private readonly _blockCount: number

  readonly _blockMin: number[]
  readonly _blockMax: number[]
  readonly _blockLazy: number[]

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
  rightMostLower(index: number): number {
    const cur = this._nums[index] + this._blockLazy[this._belong[index]]
    return this._queryRightMost(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] < cur,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] < cur
    )
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] <= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostFloor(index: number): number {
    const cur = this._nums[index] + this._blockLazy[this._belong[index]]
    return this._queryRightMost(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] <= cur,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] <= cur
    )
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] >= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostCeiling(index: number): number {
    const cur = this._nums[index] + this._blockLazy[this._belong[index]]
    return this._queryRightMost(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] >= cur,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] >= cur
    )
  }

  /**
   * 查询`index`右侧最远的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  rightMostHigher(index: number): number {
    const cur = this._nums[index] + this._blockLazy[this._belong[index]]
    return this._queryRightMost(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] > cur,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] > cur
    )
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] < nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostLower(index: number): number {
    const cur = this._nums[index] + this._blockLazy[this._belong[index]]
    return this._queryLeftMost(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] < cur,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] < cur
    )
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] <= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostFloor(index: number): number {
    const cur = this._nums[index] + this._blockLazy[this._belong[index]]
    return this._queryLeftMost(
      index,
      bid => this._blockMin[bid] + this._blockLazy[bid] <= cur,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] <= cur
    )
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] >= nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostCeiling(index: number): number {
    const cur = this._nums[index] + this._blockLazy[this._belong[index]]
    return this._queryLeftMost(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] >= cur,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] >= cur
    )
  }

  /**
   * 查询`index`左侧最远的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftMostHigher(index: number): number {
    const cur = this._nums[index] + this._blockLazy[this._belong[index]]
    return this._queryLeftMost(
      index,
      bid => this._blockMax[bid] + this._blockLazy[bid] > cur,
      (eid, bid) => this._nums[eid] + this._blockLazy[bid] > cur
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

export { RightMostLeftMostQuery }

if (require.main === module) {
  // 962. 最大宽度坡
  // https://leetcode.cn/problems/maximum-width-ramp/
  function maxWidthRamp(nums: number[]): number {
    const Q = new RightMostLeftMostQuery(nums)
    let res = 0
    for (let i = 0; i < nums.length; i++) {
      const rightMostLower = Q.rightMostCeiling(i)
      if (rightMostLower !== -1) {
        res = Math.max(res, rightMostLower - i)
      }
    }
    return res
  }

  // 2863. 最长半递减数组
  // https://leetcode.cn/problems/maximum-length-of-semi-decreasing-subarrays/description/
  function maxSubarrayLength(nums: number[]): number {
    const Q = new RightMostLeftMostQuery(nums)
    let res = 0
    for (let i = 0; i < nums.length; i++) {
      const rightMostLower = Q.rightMostLower(i)
      if (rightMostLower !== -1) {
        res = Math.max(res, rightMostLower - i + 1)
      }
    }
    return res
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
    }

    const N = 1e5
    const MAX = 1e9
    const randomNums = Array.from({ length: N }, () => (Math.random() * MAX) | 0)
    const mocker = new Mocker(randomNums.slice())
    const real = new RightMostLeftMostQuery(randomNums)
    const debug = (
      mockerFunc: (index: number) => number,
      realFunc: (index: number) => number
    ): void => {
      const index = randint(0, randomNums.length - 1)
      const mockerRes = mockerFunc(index)
      const QRes = realFunc(index)
      if (mockerRes !== QRes) {
        console.log(realFunc.name, index, mockerRes, QRes)
        console.log(mocker._nums, real._nums, real._blockLazy)
        console.log('blockMin', real._blockMin)
        console.log('blockMax', real._blockMax)
        throw new Error(realFunc.name)
      }
    }

    const randint = (a: number, b: number) => Math.floor(Math.random() * (b - a + 1)) + a
    for (let i = 0; i < 1e4; i++) {
      const op = randint(0, 4)
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
      }
    }

    console.log('ok')
  }
}
