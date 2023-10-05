/* eslint-disable no-inner-declarations */

/**
 * 基于分块实现的`树状数组`.
 * `O(1)`单点加，`O(sqrt(n))`查询区间和.
 * 一般配合莫队算法使用.
 */
class BITRangeBlock {
  private readonly _n: number
  private readonly _belong: Uint16Array
  private readonly _blockStart: Uint32Array
  private readonly _blockEnd: Uint32Array
  private readonly _nums: number[]
  private readonly _blockSum: number[]

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
    const nums = Array(n).fill(0)
    const blockSum = Array(blockCount).fill(0)
    this._n = n
    this._belong = belong
    this._blockStart = blockStart
    this._blockEnd = blockEnd
    this._nums = nums
    this._blockSum = blockSum
    if (typeof lengthOrArrayLike !== 'number') {
      this.build(lengthOrArrayLike)
    }
  }

  add(index: number, delta: number): void {
    if (index < 0 || index >= this._n) {
      throw new RangeError(`index ${index} out of range`)
    }
    this._nums[index] += delta
    this._blockSum[this._belong[index]] += delta
  }

  queryRange(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return 0
    let res = 0
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    if (bid1 === bid2) {
      for (let i = start; i < end; i++) res += this._nums[i]
      return res
    }
    for (let i = start; i < this._blockEnd[bid1]; i++) res += this._nums[i]
    for (let bid = bid1 + 1; bid < bid2; bid++) res += this._blockSum[bid]
    for (let i = this._blockStart[bid2]; i < end; i++) res += this._nums[i]
    return res
  }

  build(arr: ArrayLike<number>): void {
    if (arr.length !== this._n) {
      throw new Error(`array length ${arr.length} mismatch n ${this._n}`)
    }
    this._nums.fill(0)
    this._blockSum.fill(0)
    for (let i = 0; i < arr.length; i++) this.add(i, arr[i])
  }

  toString(): string {
    const sb: string[] = []
    sb.push('BITRangeBlock{')
    sb.push(this._nums.join(', '))
    sb.push('}')
    return sb.join('')
  }
}

export { BITRangeBlock }

if (require.main === module) {
  const bitRangeBlock = new BITRangeBlock(10)
  bitRangeBlock.add(0, 1)
  bitRangeBlock.add(1, 2)
  console.log(bitRangeBlock.queryRange(0, 2))
  console.log(bitRangeBlock.toString())

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
      const bitLike = new BITRangeBlock(stations.length)
      stations.forEach((v, i) => bitLike.add(i, v))
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
