/**
 * 查询区间`绝对众数`.
 */
class RangeMajorityQuery {
  private static _bisectRight(arr: ArrayLike<number>, x: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      const midElement = arr[mid]
      if (midElement <= x) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private static _bisectLeft(arr: ArrayLike<number>, x: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      const midElement = arr[mid]
      if (midElement < x) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private readonly _arr: ArrayLike<number>
  private readonly _repeat: number
  private readonly _pos: Map<number, number[]> = new Map()

  /**
   * O(n)预处理.
   * @param repeat 每次随机查询的重复次数.
   */
  constructor(arr: ArrayLike<number>, repeat = 30) {
    for (let i = 0; i < arr.length; i++) {
      const v = arr[i]
      if (!this._pos.has(v)) {
        this._pos.set(v, [])
      }
      this._pos.get(v)!.push(i)
    }
    this._arr = arr
    this._repeat = repeat
  }

  /**
   * O(log(n)*{@link _repeat}) 查询区间 [start, end) 中的`绝对众数`以及其出现次数.
   * 如果不存在,返回`undefined`.
   * @param threshold 阈值.绝对众数的定义为出现次数大于等于`threshold`的数.
   */
  query(
    start: number,
    end: number,
    threshold: number
  ): [majority: number, freq: number] | undefined {
    const len = end - start
    if (threshold <= len >>> 1) {
      throw new Error('threshold must be greater than half of the length of the array')
    }
    for (let _ = 0; _ < this._repeat; _++) {
      const randIndex = start + ((Math.random() * len) | 0)
      const num = this._arr[randIndex]
      const indices = this._pos.get(num) || []
      const right = RangeMajorityQuery._bisectRight(indices, end - 1)
      const left = RangeMajorityQuery._bisectLeft(indices, start)
      const freq = right - left
      if (freq >= threshold) {
        return [num, freq]
      }
    }
    return undefined
  }
}

export { RangeMajorityQuery }

if (require.main === module) {
  // https://leetcode.cn/problems/online-majority-element-in-subarray/
  class MajorityChecker {
    private readonly _rmq: RangeMajorityQuery
    constructor(arr: number[]) {
      this._rmq = new RangeMajorityQuery(arr)
    }

    query(left: number, right: number, threshold: number): number {
      const res = this._rmq.query(left, right + 1, threshold)
      return res ? res[0] : -1
    }
  }

  /**
   * Your MajorityChecker object will be instantiated and called as such:
   * var obj = new MajorityChecker(arr)
   * var param_1 = obj.query(left,right,threshold)
   */
}
