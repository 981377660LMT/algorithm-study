// https://www.acwing.com/solution/content/1016/ 你能回答这些问题吗?
// !查询数组的最大子段和
// ![l,r]区间的最大子段和就是左区间的最大子段和,右区间最大子段和,以及左右两区间结合在一起中间的最大子段和.

// !单点修改 不需要pushDown/lazy标签

const INF = 2e15

/**
 * 查询区间的最大子段和
 */
class SegmentTree {
  private readonly _size: number
  private readonly _max: number[] // 区间最大子段和
  private readonly _preMax: number[] // 前缀最大子段和
  private readonly _sufMax: number[] // 后缀最大子段和
  private readonly _hasRemoved: Uint8Array // 是否被删除

  /**
   * @param sizeOrNums 数组长度或数组
   */
  constructor(sizeOrNums: number | number[]) {
    this._size = typeof sizeOrNums === 'number' ? sizeOrNums : sizeOrNums.length
    this._max = Array(this._size << 2).fill(-INF)
    this._preMax = Array(this._size << 2).fill(-INF)
    this._sufMax = Array(this._size << 2).fill(-INF)
    this._hasRemoved = new Uint8Array(this._size << 2)
    if (Array.isArray(sizeOrNums)) {
      this._build(1, 1, this._size, sizeOrNums)
    }
  }

  split(index: number): void {
    // this._checkRange(index, index)
    this._update(1, index, index, 1, this._size)
  }

  queryAll(): number {
    return this._max[1]
  }

  private _build(rt: number, l: number, r: number, nums: number[]): void {
    if (l === r) {
      this._max[rt] = nums[l - 1]
      this._preMax[rt] = nums[l - 1]
      this._sufMax[rt] = nums[l - 1]
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._build(rt << 1, l, mid, nums)
    this._build((rt << 1) | 1, mid + 1, r, nums)
    this._pushUp(rt)
  }

  private _update(rt: number, L: number, R: number, l: number, r: number): void {
    if (L <= l && r <= R) {
      this._max[rt] = 0
      this._preMax[rt] = 0
      this._sufMax[rt] = 0
      this._hasRemoved[rt] = 1
      return
    }

    const mid = Math.floor((l + r) / 2)
    if (L <= mid) this._update(rt << 1, L, R, l, mid)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r)
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this._hasRemoved[rt] = this._hasRemoved[rt << 1] || this._hasRemoved[(rt << 1) | 1]

    // 是否可以左右合并
    this._preMax[rt] = this._hasRemoved[rt << 1]
      ? this._preMax[rt << 1]
      : this._preMax[rt << 1] + this._preMax[(rt << 1) | 1]

    // 是否可以左右合并
    this._sufMax[rt] = this._hasRemoved[(rt << 1) | 1]
      ? this._sufMax[(rt << 1) | 1]
      : this._sufMax[(rt << 1) | 1] + this._sufMax[rt << 1]

    this._max[rt] = Math.max(
      this._max[rt << 1],
      this._max[(rt << 1) | 1],
      this._sufMax[rt << 1] + this._preMax[(rt << 1) | 1]
    )
  }

  private _checkRange(l: number, r: number): void {
    if (l < 1 || r > this._size) {
      throw new RangeError(`[${l}, ${r}] out of range: [1, ${this._size}]`)
    }
  }
}

function maximumSegmentSum(nums: number[], removeQueries: number[]): number[] {
  const n = nums.length
  const tree = new SegmentTree(nums)
  const res = Array<number>(n).fill(0)

  for (let i = 0; i < n; i++) {
    const qi = removeQueries[i] + 1
    tree.split(qi)
    res[i] = tree.queryAll()
  }

  return res
}

if (require.main === module) {
  console.log(maximumSegmentSum([1, 2, 5, 6, 1], [0, 3, 2, 4, 1]))
}

export {}
