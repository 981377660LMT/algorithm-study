/* eslint-disable no-inner-declarations */

class DiffArray {
  private readonly _diff: number[]
  private dirty = false

  constructor(n: number) {
    this._diff = Array(n + 1).fill(0)
  }

  /**
   * 区间 `[start,end)` 加上 `delta`.
   */
  addRange(start: number, end: number, delta: number): void {
    if (start < 0) start = 0
    if (end >= this._diff.length) end = this._diff.length - 1
    if (start >= end) return
    this.dirty = true
    this._diff[start] += delta
    this._diff[end] -= delta
  }

  build(): void {
    if (!this.dirty) return
    for (let i = 1; i < this._diff.length; i++) this._diff[i] += this._diff[i - 1]
    this.dirty = false
  }

  get(index: number): number {
    this.build()
    return this._diff[index]
  }

  getAll(): number[] {
    this.build()
    return this._diff.slice(0, this._diff.length - 1)
  }
}

class DiffMap {
  private readonly _diff: Map<number, number> = new Map()
  private _sortedKeys: number[] = []
  private _preSum: number[] = []
  private dirty = true

  addRange(start: number, end: number, delta: number): void {
    if (start >= end) return
    this.dirty = true
    this._diff.set(start, (this._diff.get(start) || 0) + delta)
    this._diff.set(end, (this._diff.get(end) || 0) - delta)
  }

  build(): void {
    if (!this.dirty) return
    this._sortedKeys = [...this._diff.keys()].sort((a, b) => a - b)
    const preSum: number[] = Array(this._sortedKeys.length + 1)
    preSum[0] = 0
    for (let i = 1; i < preSum.length; i++) {
      preSum[i] = preSum[i - 1] + (this._diff.get(this._sortedKeys[i - 1]) || 0)
    }
    this._preSum = preSum
    this.dirty = false
  }

  get(pos: number): number {
    this.build()
    return this._preSum[this._bisectRight(this._sortedKeys, pos)]
  }

  // eslint-disable-next-line class-methods-use-this
  private _bisectRight(arr: ArrayLike<number>, target: number): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] <= target) left = mid + 1
      else right = mid - 1
    }
    return left
  }
}

export { DiffArray, DiffMap }

if (require.main === module) {
  // 2251. 花期内花的数目
  // https://leetcode.cn/problems/number-of-flowers-in-full-bloom/description/
  function fullBloomFlowers(flowers: number[][], people: number[]): number[] {
    const diff = new DiffMap()
    flowers.forEach(([start, end]) => diff.addRange(start, end + 1, 1))
    return people.map(pos => diff.get(pos))
  }

  console.log(
    fullBloomFlowers(
      [
        [1, 6],
        [3, 7],
        [9, 12],
        [4, 13]
      ],
      [2, 3, 7, 7]
    )
  )
}
