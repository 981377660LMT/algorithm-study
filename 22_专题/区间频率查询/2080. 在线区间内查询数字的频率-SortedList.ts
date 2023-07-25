import { SortedList } from '../离线查询/根号分治/SortedList/_SortedList'

/**
 * 静态区间频率查询.
 * 基于二分实现,单次查询时间复杂度`O(logn)`.
 */
class RangeFreqQueryStatic {
  private static _bisectLeft(arr: number[], target: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      if (arr[mid] < target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private static _bisectRight(arr: number[], target: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      if (arr[mid] <= target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private readonly _mp: Map<number, number[]> = new Map()

  constructor(arr: number[]) {
    arr.forEach((v, i) => {
      if (!this._mp.has(v)) this._mp.set(v, [])
      this._mp.get(v)!.push(i)
    })
  }

  /**
   * 查询闭区间`[left,right]`中`target`出现的次数.
   * 0 <= left <= right < arr.length
   */
  query(left: number, right: number, target: number): number {
    return (
      RangeFreqQueryStatic._bisectRight(this._mp.get(target) || [], right) -
      RangeFreqQueryStatic._bisectLeft(this._mp.get(target) || [], left)
    )
  }
}

/**
 * 支持单点更新的区间频率查询.
 * 基于`SortedList`实现,单次查询时间复杂度`O(sqrt(n))`.
 */
class RangeFreqQueryPointUpdate {
  private readonly _mp: Map<number, SortedList> = new Map()
  private readonly _arr: number[]

  constructor(arr: number[]) {
    this._arr = arr
    arr.forEach((v, i) => {
      if (!this._mp.has(v)) this._mp.set(v, new SortedList())
      this._mp.get(v)!.add(i)
    })
  }

  /**
   * 查询闭区间`[left,right]`中`target`出现的次数.
   * 0 <= left <= right < arr.length
   */
  query(left: number, right: number, target: number): number {
    const sl = this._mp.get(target)
    if (!sl) return 0
    return sl.bisectRight(right) - sl.bisectLeft(left)
  }

  /**
   * 将下标`index`处的值更新为`value`.
   */
  update(index: number, value: number): void {
    const pre = this._arr[index]
    if (pre === value) return
    this._mp.get(pre)!.discard(index)
    if (!this._mp.has(value)) this._mp.set(value, new SortedList())
    this._mp.get(value)!.add(index)
  }
}

export { RangeFreqQueryStatic, RangeFreqQueryPointUpdate }
