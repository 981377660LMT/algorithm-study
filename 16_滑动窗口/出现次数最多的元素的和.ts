/* eslint-disable @typescript-eslint/no-non-null-assertion */

// class FreqManager

/**
 * 统计一个容器内 (最多元素出现的次数, 这些元素key的和).
 */
class MajorSum {
  private readonly _counter: Map<number, number> = new Map()
  private readonly _freqSum: Map<number, number> = new Map()
  private readonly _freqTypes: Map<number, number> = new Map()
  private _maxFreq = 0
  private _sum = 0

  add(x: number): void {
    this._counter.set(x, (this._counter.get(x) || 0) + 1)
    const xFreq = this._counter.get(x)!
    this._freqSum.set(xFreq, (this._freqSum.get(xFreq) || 0) + x)
    this._freqSum.set(xFreq - 1, (this._freqSum.get(xFreq - 1) || 0) - x)
    this._freqTypes.set(xFreq, (this._freqTypes.get(xFreq) || 0) + 1)
    this._freqTypes.set(xFreq - 1, (this._freqTypes.get(xFreq - 1) || 0) - 1)
    if (xFreq > this._maxFreq) {
      this._maxFreq = xFreq
      this._sum = x
    } else if (xFreq === this._maxFreq) {
      this._sum += x
    }
  }

  discard(x: number): void {
    if (this._counter.get(x) === 0) return
    this._counter.set(x, (this._counter.get(x) || 0) - 1)
    const xFreq = this._counter.get(x)!
    this._freqSum.set(xFreq, (this._freqSum.get(xFreq) || 0) + x)
    this._freqSum.set(xFreq + 1, (this._freqSum.get(xFreq + 1) || 0) - x)
    this._freqTypes.set(xFreq, (this._freqTypes.get(xFreq) || 0) + 1)
    this._freqTypes.set(xFreq + 1, (this._freqTypes.get(xFreq + 1) || 0) - 1)
    if (xFreq + 1 === this._maxFreq) {
      this._sum -= x
      if (this._freqTypes.get(this._maxFreq) === 0) {
        this._maxFreq -= 1
        this._sum = this._freqSum.get(this._maxFreq)!
      }
    }
    if (this._counter.get(x) === 0) {
      this._counter.delete(x)
    }
  }

  /**
   * 返回(最多元素出现的次数, 这些元素key的和).
   */
  query(): [maxFreq: number, maxFreqSum: number] {
    return [this._maxFreq, this._sum]
  }

  /**
   * 返回当前元素种类数.
   */
  size(): number {
    return this._counter.size
  }
}

if (require.main === module) {
  const majorSum = new MajorSum()
  majorSum.add(1)
  majorSum.add(1)
  majorSum.add(2)
  console.log(majorSum.query())
}

export { MajorSum }
