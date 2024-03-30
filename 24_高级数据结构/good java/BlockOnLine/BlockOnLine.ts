import { SortedListFast } from '../../../22_专题/离线查询/根号分治/SortedList/SortedListFast'

class BlockOnLine {
  private readonly _left: Int32Array
  private readonly _right: Int32Array
  private readonly _n: number

  constructor(n: number) {
    this._left = new Int32Array(n)
    this._right = new Int32Array(n)
    this._n = n
    this.init()
  }

  init(): void {
    this._left.fill(this._n)
    this._right.fill(-1)
  }

  /**
   * 添加`index`位置的元素.
   * 如果左边或者右边存在元素, 则删除原有块，合并为一个新块.
   */
  add(
    index: number,
    options: {
      onAddBlock?: (start: number, end: number) => void
      onRemoveBlock?: (start: number, end: number) => void
    } = {}
  ): boolean {
    if (!(0 <= index && index < this._n)) return false
    if (!(this._left[index] > this._right[index])) return false

    const { onAddBlock, onRemoveBlock } = options
    let from = index
    let to = index
    if (index > 0 && this._left[index - 1] <= this._right[index - 1]) {
      from = this._left[index - 1]
      onRemoveBlock && onRemoveBlock(from, index)
    }
    if (index + 1 < this._n && this._left[index + 1] <= this._right[index + 1]) {
      to = this._right[index + 1]
      onRemoveBlock && onRemoveBlock(index + 1, to + 1)
    }
    this._left[from] = from
    this._right[from] = to
    this._left[to] = from
    this._right[to] = to
    onAddBlock && onAddBlock(from, to + 1)
    return true
  }
}

export { BlockOnLine }

if (require.main === module) {
  // 2382. 删除操作后的最大子段和
  // https://leetcode.cn/problems/maximum-segment-sum-after-removals/
  // nums[i] >= 0
  function maximumSegmentSum(nums: number[], removeQueries: number[]): number[] {
    const n = nums.length
    const preSum = Array(n + 1)
    preSum[0] = 0
    for (let i = 0; i < n; i++) preSum[i + 1] = preSum[i] + nums[i]

    const sl = new SortedListFast()
    const res = Array<number>(removeQueries.length).fill(0)
    const B = new BlockOnLine(n)

    for (let i = removeQueries.length - 1; i >= 0; i--) {
      const index = removeQueries[i]
      res[i] = sl.length ? sl.max! : 0
      B.add(index, {
        onAddBlock: (start, end) => {
          const sum = preSum[end] - preSum[start]
          sl.add(sum)
        },
        onRemoveBlock: (start, end) => {
          const sum = preSum[end] - preSum[start]
          sl.discard(sum)
        }
      })
    }

    return res
  }

  function demo(): void {
    const B = new BlockOnLine(4)
    const hooks = {
      onAddBlock: (start: number, end: number) => {
        console.log(`add block: [${start}, ${end})`)
      },
      onRemoveBlock: (start: number, end: number) => {
        console.log(`remove block: [${start}, ${end})`)
      }
    }
    B.add(2, hooks)
    B.add(0, hooks)
    B.add(3, hooks)
    B.add(1, hooks)
  }
}
