/* eslint-disable no-inner-declarations */

/**
 * 基于分块实现的`树状数组`.
 * `O(sqrt(n))`单点加，`O(1)`查询区间和.
 * 一般配合莫队算法使用.
 */
class BITRangeBlockFastQuery {
  private readonly _n: number
  private readonly _belong: Uint16Array
  private readonly _blockStart: Uint32Array
  private readonly _blockEnd: Uint32Array
  private readonly _blockCount: number
  private readonly _partPreSum: number[]
  private readonly _blockPreSum: number[]

  constructor(lengthOrArrayLike: number | ArrayLike<number>) {
    const n = typeof lengthOrArrayLike === 'number' ? lengthOrArrayLike : lengthOrArrayLike.length
    const blockSize = (Math.sqrt(n) + 1) | 0
    const blockCount = 1 + ((n / blockSize) | 0)
    const belong = new Uint16Array(n)
    for (let i = 0; i < n; i++) {
      belong[i] = (i / blockSize) | 0
    }
    const blockStart = new Uint32Array(blockCount)
    const blockEnd = new Uint32Array(blockCount)
    for (let i = 0; i < blockCount; i++) {
      blockStart[i] = i * blockSize
      blockEnd[i] = Math.min((i + 1) * blockSize, n)
    }
    const partPreSum = Array(n).fill(0)
    const blockPreSum = Array(blockCount).fill(0)
    this._n = n
    this._belong = belong
    this._blockStart = blockStart
    this._blockEnd = blockEnd
    this._blockCount = blockCount
    this._partPreSum = partPreSum
    this._blockPreSum = blockPreSum
    if (typeof lengthOrArrayLike !== 'number') {
      this.build(lengthOrArrayLike)
    }
  }

  add(index: number, delta: number): void {
    if (index < 0 || index >= this._n) {
      throw new RangeError(`index ${index} out of range`)
    }
    const bid = this._belong[index]
    for (let i = index; i < this._blockEnd[bid]; i++) {
      this._partPreSum[i] += delta
    }
    for (let id = bid + 1; id < this._blockCount; id++) {
      this._blockPreSum[id] += delta
    }
  }

  queryRange(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return 0
    return this._query(end) - this._query(start)
  }

  build(arr: ArrayLike<number>): void {
    if (arr.length !== this._n) {
      throw new Error(`array length ${arr.length} mismatch n ${this._n}`)
    }
    let curBlockSum = 0
    for (let bid = 0; bid < this._blockCount; bid++) {
      let curPartSum = 0
      for (let i = this._blockStart[bid]; i < this._blockEnd[bid]; i++) {
        curPartSum += arr[i]
        this._partPreSum[i] = curPartSum
      }
      this._blockPreSum[bid] = curBlockSum
      curBlockSum += curPartSum
    }
  }

  toString(): string {
    const sb: string[] = []
    sb.push('BITRangeBlockFastQuery{')
    for (let i = 0; i < this._partPreSum.length; i++) {
      sb.push(String(this.queryRange(i, i + 1)))
      if (i !== this._partPreSum.length - 1) sb.push(',')
    }
    sb.push('}')
    return sb.join('')
  }

  private _query(end: number): number {
    if (end <= 0) return 0
    return this._partPreSum[end - 1] + this._blockPreSum[this._belong[end - 1]]
  }
}

export { BITRangeBlockFastQuery }

if (require.main === module) {
  const bit = new BITRangeBlockFastQuery([1, 2, 3])
  bit.add(0, 1)
  bit.add(1, 2)
  console.log(bit.queryRange(0, 2), bit.toString())
  bit.build([2, 3, 4])
  console.log(bit.queryRange(0, 2), bit.toString())

  // https://leetcode.cn/problems/maximize-the-minimum-powered-city/
  function maxPower(stations: number[], r: number, k: number): number {
    const n = stations.length
    let left = 1
    let right = 2e15
    while (left <= right) {
      const mid = Math.floor((left + right) / 2)
      if (check(mid)) left = mid + 1
      else right = mid - 1
    }
    return right

    function check(mid: number): boolean {
      const bitLike = new BITRangeBlockFastQuery(stations)
      let curK = k
      for (let i = 0; i < n; i++) {
        const cur = bitLike.queryRange(Math.max(0, i - r), Math.min(i + r + 1, n))
        if (cur < mid) {
          const diff = mid - cur
          bitLike.add(Math.min(i + r, n - 1), diff)
          curK -= diff
          if (curK < 0) return false
        }
      }
      return true
    }
  }
}
