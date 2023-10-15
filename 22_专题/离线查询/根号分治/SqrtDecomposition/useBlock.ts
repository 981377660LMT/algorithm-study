/* eslint-disable no-inner-declarations */

function useBlock(
  n: number,
  blockSize = Math.sqrt(n + 1) | 0
): {
  /** 下标所属的块. */
  belong: Uint16Array
  /** 每个块的起始下标(包含). */
  blockStart: Uint32Array
  /** 每个块的结束下标(不包含). */
  blockEnd: Uint32Array
  /** 块的数量. */
  blockCount: number
} {
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

  return {
    belong,
    blockStart,
    blockEnd,
    blockCount
  }
}

export { useBlock, useBlock as createBlock }

if (require.main === module) {
  // 2569. 更新数组后处理求和查询
  // https://leetcode.cn/problems/handling-sum-queries-after-update/
  //
  // 给你两个下标从 0 开始的数组 nums1 和 nums2 ，和一个二维数组 queries 表示一些操作。总共有 3 种类型的操作：
  // 操作类型 1 为 queries[i] = [1, l, r] 。
  // 你需要将 nums1 从下标 l 到下标 r 的所有 0 反转成 1 或将 1 反转成 0 。l 和 r 下标都从 0 开始。
  // 操作类型 2 为 queries[i] = [2, p, 0] 。对于 0 <= i < n 中的所有下标，令 nums2[i] = nums2[i] + nums1[i] * p 。
  // 操作类型 3 为 queries[i] = [3, 0, 0] 。求 nums2 中所有元素的和。
  // 请你返回一个数组，包含所有第三种操作类型的答案。
  // 区间反转+查询区间1的个数(RangeFlipRangeOnesCount)

  function handleQuery(nums1: number[], nums2: number[], queries: number[][]): number[] {
    const { belong, blockStart, blockEnd, blockCount } = useBlock(nums1.length)

    // !1.需要维护的属性.
    const blockOnes = new Uint16Array(blockCount)
    const blockFlip = new Uint8Array(blockCount)

    for (let i = 0; i < blockCount; i++) rebuild(i)

    // !2.处理询问.
    let nums2Sum = nums2.reduce((a, b) => a + b, 0)
    const res: number[] = []
    for (let qi = 0; qi < queries.length; qi++) {
      const kind = queries[qi][0]
      if (kind === 1) {
        // RangeFlip
        const { 1: left, 2: right } = queries[qi]
        update(left, right + 1)
      } else if (kind === 2) {
        // RangeOnesCount
        const start = 0
        const end = nums1.length
        const onesCount = query(start, end)
        nums2Sum += onesCount * queries[qi][1]
      } else if (kind === 3) {
        res.push(nums2Sum)
      }
    }
    return res

    // !3.初始化/更新零散块后重构整个块
    function rebuild(bid: number): void {
      blockOnes[bid] = 0
      for (let i = blockStart[bid]; i < blockEnd[bid]; i++) {
        blockOnes[bid] += nums1[i]
      }
    }

    function update(start: number, end: number): void {
      const bid1 = belong[start]
      const bid2 = belong[end - 1]
      if (bid1 === bid2) {
        for (let i = start; i < end; i++) {
          blockOnes[bid1] -= nums1[i] ^ blockFlip[bid1]
          nums1[i] ^= 1
          blockOnes[bid1] += nums1[i] ^ blockFlip[bid1]
        }
      } else {
        for (let i = start; i < blockEnd[bid1]; i++) {
          blockOnes[bid1] -= nums1[i] ^ blockFlip[bid1]
          nums1[i] ^= 1
          blockOnes[bid1] += nums1[i] ^ blockFlip[bid1]
        }
        for (let bid = bid1 + 1; bid < bid2; bid++) {
          blockFlip[bid] ^= 1
          blockOnes[bid] = blockEnd[bid] - blockStart[bid] - blockOnes[bid]
        }
        for (let i = blockStart[bid2]; i < end; i++) {
          blockOnes[bid2] -= nums1[i] ^ blockFlip[bid2]
          nums1[i] ^= 1
          blockOnes[bid2] += nums1[i] ^ blockFlip[bid2]
        }
      }
    }

    function query(start: number, end: number): number {
      const bid1 = belong[start]
      const bid2 = belong[end - 1]
      let res = 0
      if (bid1 === bid2) {
        for (let i = start; i < end; i++) {
          res += nums1[i] ^ blockFlip[bid1]
        }
      } else {
        for (let i = start; i < blockEnd[bid1]; i++) {
          res += nums1[i] ^ blockFlip[bid1]
        }
        for (let bid = bid1 + 1; bid < bid2; bid++) {
          res += blockOnes[bid]
        }
        for (let i = blockStart[bid2]; i < end; i++) {
          res += nums1[i] ^ blockFlip[bid2]
        }
      }
      return res
    }
  }

  const INF = 2e9 // !超过int32使用2e15

  // https://leetcode.cn/problems/range-sum-query-mutable/
  // RangeAssignRangeSum
  class NumArray {
    private readonly _nums: number[]
    private readonly _belong: Uint16Array
    private readonly _blockStart: Uint32Array
    private readonly _blockEnd: Uint32Array
    private readonly _blockCount: number
    private readonly _blockColor: number[]
    private readonly _blockSum: number[]

    constructor(nums: number[]) {
      const { belong, blockStart, blockEnd, blockCount } = useBlock(nums.length)
      this._nums = nums.slice()
      this._belong = belong
      this._blockStart = blockStart
      this._blockEnd = blockEnd
      this._blockCount = blockCount
      this._blockColor = Array(blockCount).fill(INF)
      this._blockSum = Array(blockCount).fill(0)

      for (let i = 0; i < blockCount; i++) this._rebuild(i)
    }

    /**
     * @param bid 初始化/更新零散块后重构整个块`bid`.
     */
    private _rebuild(bid: number): void {
      this._blockSum[bid] = 0
      for (let i = this._blockStart[bid]; i < this._blockEnd[bid]; i++) {
        this._blockSum[bid] += this._nums[i]
      }
    }

    /**
     * 区间赋值.
     */
    private _updateRange(start: number, end: number, target: number): void {
      const bid1 = this._belong[start]
      const bid2 = this._belong[end - 1]
      if (bid1 === bid2) {
        for (let i = start; i < end; i++) this._nums[i] = target
        this._rebuild(bid1)
      } else {
        for (let i = start; i < this._blockEnd[bid1]; i++) this._nums[i] = target
        this._rebuild(bid1)
        for (let bid = bid1 + 1; bid < bid2; bid++) this._blockColor[bid] = target
        for (let i = this._blockStart[bid2]; i < end; i++) this._nums[i] = target
        this._rebuild(bid2)
      }
    }

    /**
     * 区间求和.
     */
    private _queryRange(start: number, end: number): number {
      const bid1 = this._belong[start]
      const bid2 = this._belong[end - 1]
      let res = 0
      if (bid1 === bid2) {
        const color = this._blockColor[bid1]
        for (let i = start; i < end; i++) {
          res += color === INF ? this._nums[i] : color
        }
      } else {
        const color1 = this._blockColor[bid1]
        for (let i = start; i < this._blockEnd[bid1]; i++) {
          res += color1 === INF ? this._nums[i] : color1
        }
        for (let bid = bid1 + 1; bid < bid2; bid++) {
          const color = this._blockColor[bid]
          res +=
            color === INF
              ? this._blockSum[bid]
              : color * (this._blockEnd[bid] - this._blockStart[bid])
        }
        const color2 = this._blockColor[bid2]
        for (let i = this._blockStart[bid2]; i < end; i++) {
          res += color2 === INF ? this._nums[i] : color2
        }
      }

      return res
    }

    update(index: number, val: number): void {
      this._updateRange(index, index + 1, val)
    }

    sumRange(left: number, right: number): number {
      return this._queryRange(left, right + 1)
    }
  }
}
