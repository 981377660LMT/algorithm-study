/* eslint-disable @typescript-eslint/no-non-null-assertion */

// class FreqManager

/**
 * 统计一个容器内出现次数最多的元素的`出现次数`.
 * @deprecated Map非常慢，使用其他有序数据结构代替.
 */
class MajorFreq {
  private readonly _counter: Map<number, number> = new Map()
  private readonly _freqTypes: Map<number, number> = new Map()
  private _maxFreq = 0
  private _size = 0

  add(value: number): this {
    const preFreq = this._counter.get(value) || 0
    this._counter.set(value, preFreq + 1)
    const curFreq = preFreq + 1
    this._freqTypes.set(curFreq, (this._freqTypes.get(curFreq) || 0) + 1)
    this._freqTypes.set(preFreq, (this._freqTypes.get(preFreq) || 0) - 1)
    if (curFreq > this._maxFreq) this._maxFreq = curFreq
    this._size++
    return this
  }

  discard(value: number): boolean {
    const preFreq = this._counter.get(value)
    if (!preFreq) return false
    const curFreq = preFreq - 1
    this._counter.set(value, curFreq)
    this._freqTypes.set(curFreq, (this._freqTypes.get(curFreq) || 0) + 1)
    this._freqTypes.set(preFreq, (this._freqTypes.get(preFreq) || 0) - 1)
    if (preFreq === this._maxFreq && !this._freqTypes.get(this._maxFreq)) this._maxFreq--
    if (!curFreq) this._counter.delete(value)
    this._size--
    return true
  }

  maxFreq(): number {
    return this._maxFreq
  }

  getFreq(value: number): number {
    return this._counter.get(value) || 0
  }

  get size(): number {
    return this._size
  }
}

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
    const preFreq = this._counter.get(x) || 0
    this._counter.set(x, preFreq + 1)
    const curFreq = preFreq + 1
    this._freqSum.set(curFreq, (this._freqSum.get(curFreq) || 0) + x)
    this._freqSum.set(preFreq, (this._freqSum.get(preFreq) || 0) - x)
    this._freqTypes.set(curFreq, (this._freqTypes.get(curFreq) || 0) + 1)
    this._freqTypes.set(preFreq, (this._freqTypes.get(preFreq) || 0) - 1)
    if (curFreq > this._maxFreq) {
      this._maxFreq = curFreq
      this._sum = x
    } else if (curFreq === this._maxFreq) {
      this._sum += x
    }
  }

  discard(x: number): boolean {
    const preFreq = this._counter.get(x)
    if (!preFreq) return false
    const curFreq = preFreq - 1
    this._counter.set(x, curFreq)
    this._freqSum.set(curFreq, (this._freqSum.get(curFreq) || 0) + x)
    this._freqSum.set(preFreq, (this._freqSum.get(preFreq) || 0) - x)
    this._freqTypes.set(curFreq, (this._freqTypes.get(curFreq) || 0) + 1)
    this._freqTypes.set(preFreq, (this._freqTypes.get(preFreq) || 0) - 1)
    if (preFreq === this._maxFreq) {
      this._sum -= x
      if (!this._freqTypes.get(this._maxFreq)) {
        this._maxFreq -= 1
        this._sum = this._freqSum.get(this._maxFreq)!
      }
    }
    if (!curFreq) this._counter.delete(x)
    return true
  }

  /**
   * 返回(最多元素出现的次数, 这些元素key的和).
   */
  query(): [maxFreq: number, maxFreqSum: number] {
    return [this._maxFreq, this._sum]
  }

  getFreq(x: number): number {
    return this._counter.get(x) || 0
  }
}

export { MajorSum, MajorFreq }

if (require.main === module) {
  type Public<T extends object> = { [K in keyof T]: T[K] }

  const mocker = {
    counter: new Map<number, number>(),
    add(value: number) {
      this.counter.set(value, (this.counter.get(value) || 0) + 1)
      return this
    },
    discard(value: number): boolean {
      if (!this.counter.get(value)) return false
      this.counter.set(value, (this.counter.get(value) || 0) - 1)
      if (!this.counter.get(value)) {
        this.counter.delete(value)
      }
      return true
    },
    maxFreq(): number {
      return Math.max(0, ...this.counter.values())
    },
    getFreq(value: number): number {
      return this.counter.get(value) || 0
    },
    get size() {
      let size = 0
      for (const freq of this.counter.values()) {
        size += freq
      }
      return size
    }
  }

  const max = 10000
  const majorFreq = new MajorFreq()
  for (let i = 0; i < max; i++) {
    const add = Math.random() < 0.5
    if (add) {
      const x = Math.floor(Math.random() * 100)
      majorFreq.add(x)
      mocker.add(x)
    } else {
      const x = Math.floor(Math.random() * 100)
      majorFreq.discard(x)
      mocker.discard(x)
    }

    if (majorFreq.maxFreq() !== mocker.maxFreq()) {
      console.log(majorFreq.maxFreq(), 111)
      console.log(mocker.maxFreq(), 222)
      throw new Error()
    }
  }

  console.log('ok')
}
