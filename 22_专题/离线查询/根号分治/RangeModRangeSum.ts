/* eslint-disable no-inner-declarations */

class RangeModRangeSum {
  private readonly _nums: number[]
  private readonly _belong: Uint32Array
  private readonly _blockStart: Uint32Array
  private readonly _blockEnd: Uint32Array
  private readonly _blockSum: number[]
  private readonly _blockMax: number[]

  /**
   * @param arr 非负整数数组.
   */
  constructor(arr: ArrayLike<number>, blockSize = 60) {
    const n = arr.length
    const copy = Array(n)
    for (let i = 0; i < n; ++i) copy[i] = arr[i]
    const belong = new Uint32Array(n)
    for (let i = 0; i < n; ++i) belong[i] = (i / blockSize) | 0
    const blockCount = (1 + n / blockSize) | 0
    const blockStart = new Uint32Array(blockCount)
    const blockEnd = new Uint32Array(blockCount)
    for (let i = 0; i < blockCount; ++i) {
      blockStart[i] = i * blockSize
      blockEnd[i] = Math.min((i + 1) * blockSize, n)
    }
    const blockSum = Array(blockCount)
    const blockMax = Array(blockCount)
    this._nums = copy
    this._belong = belong
    this._blockStart = blockStart
    this._blockEnd = blockEnd
    this._blockSum = blockSum
    this._blockMax = blockMax
    for (let bid = 0; bid < blockCount; bid++) this._rebuild(bid)
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._nums.length) return
    const pre = this._nums[index]
    if (pre === value) return
    this._nums[index] = value
    this._rebuild(this._belong[index])
  }

  /**
   * 区间`[start, end)`模`mod`.
   */
  update(start: number, end: number, mod: number): void {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    if (bid1 === bid2) {
      for (let i = start; i < end; ++i) this._nums[i] %= mod
      this._rebuild(bid1)
    } else {
      for (let i = start; i < this._blockEnd[bid1]; ++i) this._nums[i] %= mod
      this._rebuild(bid1)
      for (let bid = bid1 + 1; bid < bid2; ++bid) {
        if (this._blockMax[bid] < mod) continue
        for (let i = this._blockStart[bid]; i < this._blockEnd[bid]; ++i) this._nums[i] %= mod
        this._rebuild(bid)
      }
      for (let i = this._blockStart[bid2]; i < end; ++i) this._nums[i] %= mod
      this._rebuild(bid2)
    }
  }

  /**
   * 查询区间`[start, end)`的和.
   */
  query(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this._nums.length) end = this._nums.length
    if (start >= end) return 0
    const bid1 = this._belong[start]
    const bid2 = this._belong[end - 1]
    let res = 0
    if (bid1 === bid2) {
      for (let i = start; i < end; ++i) res += this._nums[i]
    } else {
      for (let i = start; i < this._blockEnd[bid1]; ++i) res += this._nums[i]
      for (let bid = bid1 + 1; bid < bid2; ++bid) res += this._blockSum[bid]
      for (let i = this._blockStart[bid2]; i < end; ++i) res += this._nums[i]
    }
    return res
  }

  private _rebuild(bid: number): void {
    this._blockSum[bid] = 0
    this._blockMax[bid] = 0
    for (let i = this._blockStart[bid]; i < this._blockEnd[bid]; ++i) {
      this._blockSum[bid] += this._nums[i]
      this._blockMax[bid] = Math.max(this._blockMax[bid], this._nums[i])
    }
  }
}

export { RangeModRangeSum }

if (require.main === module) {
  // check()
  testTime()
  function testTime(): void {
    const n = 1e5
    const arr = Array(n)
    for (let i = 0; i < n; ++i) arr[i] = (Math.random() * 1e9) | 0
    console.time('RangeModRangeSum')
    const rangeModRangeSum = new RangeModRangeSum(arr)
    for (let i = 0; i < 1e5; ++i) {
      rangeModRangeSum.set(i, i)
      rangeModRangeSum.update(0, i, i)
      rangeModRangeSum.query(0, n)
    }
    console.timeEnd('RangeModRangeSum')
  }

  function check(): void {
    class Mocker {
      private readonly _arr: number[]
      constructor(arr: number[]) {
        this._arr = arr.slice()
      }

      set(index: number, value: number): void {
        this._arr[index] = value
      }

      update(start: number, end: number, mod: number): void {
        for (let i = start; i < end; ++i) this._arr[i] %= mod
      }

      query(start: number, end: number): number {
        let res = 0
        for (let i = start; i < end; ++i) res += this._arr[i]
        return res
      }
    }

    const n = 1e5
    const arr = Array(n)
    for (let i = 0; i < n; ++i) arr[i] = (Math.random() * 1e9) | 0

    const mocker = new Mocker(arr)
    const rangeModRangeSum = new RangeModRangeSum(arr)
    for (let i = 0; i < 1e5; ++i) {
      const op = (Math.random() * 3) | 0
      const start = (Math.random() * n) | 0
      const end = (Math.random() * n) | 0
      const mod = (Math.random() * 1e9) | 0
      if (op === 0) {
        mocker.set(start, mod)
        rangeModRangeSum.set(start, mod)
      } else if (op === 1) {
        mocker.update(start, end, mod)
        rangeModRangeSum.update(start, end, mod)
      } else {
        const ans = mocker.query(start, end)
        const res = rangeModRangeSum.query(start, end)
        if (ans !== res) {
          console.error('oops')
          console.log(arr)
          console.log(start, end, mod)
          console.log(ans, res)
          throw new Error()
        }
      }
    }
    console.log('pass')
  }
}
