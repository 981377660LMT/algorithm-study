/**
 * 动态区间频率查询.
 */
class RangeFreqQueryDynamic {
  private static _bisectLeft<T>(arr: ArrayLike<T>, value: T): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      if (arr[mid] < value) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private readonly _sqrt: SqrtDecomposition<number, number, [v: number, same: boolean]>

  constructor(nums: number[]) {
    const n = nums.length

    this._sqrt = new SqrtDecomposition<number, number, [v: number, same: boolean]>(
      n,
      (_, left, right) => {
        const curNums = nums.slice(left, right)
        let sortedNums: number[] = []
        let lazyAdd = 0

        const created = () => {
          updated()
        }
        const updated = () => {
          sortedNums = [...curNums].sort((a, b) => a - b)
        }

        // 区间加
        const updateAll = (lazy: number) => {
          lazyAdd += lazy
        }
        const updatePart = (left: number, right: number, lazy: number) => {
          for (let i = left; i < right; i++) {
            curNums[i] += lazy
          }
        }

        // 区间查询.
        const queryAll = (queryArg: [v: number, same: boolean]) => {
          const [v, same] = queryArg
          if (same) {
            return (
              RangeFreqQueryDynamic._bisectLeft(sortedNums, v - lazyAdd + 1) -
              RangeFreqQueryDynamic._bisectLeft(sortedNums, v - lazyAdd)
            )
          }
          return sortedNums.length - RangeFreqQueryDynamic._bisectLeft(sortedNums, v - lazyAdd)
        }

        const queryPart = (left: number, right: number, queryArg: [v: number, same: boolean]) => {
          const [v, same] = queryArg
          if (same) {
            let res = 0
            for (let i = left; i < right; i++) {
              if (curNums[i] + lazyAdd === v) {
                res++
              }
            }
            return res
          }

          let res = 0
          for (let i = left; i < right; i++) {
            if (curNums[i] + lazyAdd >= v) {
              res++
            }
          }
          return res
        }

        return {
          created,
          updated,
          updateAll,
          updatePart,
          queryAll,
          queryPart
        }
      }
    )
  }

  /**
   * 区间`[left, right)`每个元素加上`value`.
   */
  update(left: number, right: number, value: number): void {
    this._sqrt.update(left, right, value)
  }

  /**
   * 查询区间`[left, right)`中元素`target`出现的次数.
   */
  rangeFreq(left: number, right: number, target: number): number {
    let res = 0
    this._sqrt.query(
      left,
      right,
      blockRes => {
        res += blockRes
      },
      [target, true]
    )
    return res
  }

  /**
   * 查询区间`[left, right)`中大于等于`floor`的元素出现的次数.
   */
  rangeFreqWithFloor(left: number, right: number, floor: number): number {
    let res = 0
    this._sqrt.query(
      left,
      right,
      blockRes => {
        res += blockRes
      },
      [floor, false]
    )
    return res
  }
}

interface Block<E, Id, Q = unknown> {
  /**
   * 在创建块时调用，用于初始化块的值.
   */
  created(): void

  /**
   * 在部分更新块之后时调用，用于将`整个块`的值更新到最新.
   */
  updated(): void

  /**
   * 查询块内左闭右开区间 [start, end) 的值. 0 <= start <= end <= blockSize.
   *
   * !注意查询时需要把块内的懒标记计入影响.
   */
  queryPart(start: number, end: number, queryArg?: Q): E

  /**
   * 更新块内左闭右开区间 [start, end) 的值. 0 <= start <= end <= blockSize.
   */
  updatePart(start: number, end: number, lazy: Id): void

  /**
   * 查询块内所有值.
   */
  queryAll(queryArg?: Q): E

  /**
   * 更新块内所有值,打上懒标记.
   */
  updateAll(lazy: Id): void
}

/**
 * 当区间需要维护的数据难以用半群来描述时，可以考虑根号分块。
 */
class SqrtDecomposition<E, Id, Q = unknown> {
  private readonly _blockSize: number
  private readonly _blocks: Block<E, Id, Q>[]

  /**
   * @param n 区间长度.
   * @param createBlock 块的工厂函数.
   * @param blockSize 分块大小，一般取 `Math.sqrt(n) + 1`.
   */
  constructor(
    n: number,
    createBlock: (id: number, start: number, end: number) => Block<E, Id, Q>,
    blockSize = ~~Math.sqrt(n) + 1
  ) {
    this._blockSize = blockSize
    this._blocks = Array(1 + ~~(n / blockSize)).fill(null)
    for (let i = 0; i < this._blocks.length; i++) {
      this._blocks[i] = createBlock(i, i * blockSize, Math.min((i + 1) * blockSize, n))
      this._blocks[i].created()
    }
  }

  /**
   * 更新左闭右开区间 `[start, end)` 的值.
   *
   * 0 <= start <= end <= n.
   */
  update(start: number, end: number, lazy: Id): void {
    if (start >= end) {
      return
    }
    const id1 = ~~(start / this._blockSize)
    const id2 = ~~(end / this._blockSize)
    const pos1 = start % this._blockSize
    const pos2 = end % this._blockSize
    if (id1 === id2) {
      this._blocks[id1].updatePart(pos1, pos2, lazy)
      this._blocks[id1].updated()
    } else {
      this._blocks[id1].updatePart(pos1, this._blockSize, lazy)
      this._blocks[id1].updated()
      for (let i = id1 + 1; i < id2; i++) {
        this._blocks[i].updateAll(lazy)
      }
      this._blocks[id2].updatePart(0, pos2, lazy)
      this._blocks[id2].updated()
    }
  }

  /**
   * 查询左闭右开区间 `[start, end)` 的值.
   *
   * 0 <= start <= end <= n.
   * @param forEach 遍历每个块的结果.
   */
  query(start: number, end: number, forEach: (blockRes: E) => void, queryArg?: Q) {
    if (start >= end) {
      return
    }
    const id1 = ~~(start / this._blockSize)
    const id2 = ~~(end / this._blockSize)
    const pos1 = start % this._blockSize
    const pos2 = end % this._blockSize
    if (id1 === id2) {
      forEach(this._blocks[id1].queryPart(pos1, pos2, queryArg))
      return
    }
    forEach(this._blocks[id1].queryPart(pos1, this._blockSize, queryArg))
    for (let i = id1 + 1; i < id2; i++) {
      forEach(this._blocks[i].queryAll(queryArg))
    }
    forEach(this._blocks[id2].queryPart(0, pos2, queryArg))
  }
}

class RangeFreqQuery {
  private readonly _rf: RangeFreqQueryDynamic

  constructor(arr: number[]) {
    this._rf = new RangeFreqQueryDynamic(arr)
  }

  query(left: number, right: number, value: number): number {
    return this._rf.rangeFreq(left, right + 1, value)
  }
}

export {}
