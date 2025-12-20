/* eslint-disable no-inner-declarations */
/* eslint-disable eqeqeq */
/* eslint-disable no-lone-blocks */
/* eslint-disable no-lonely-if */
/* eslint-disable no-nested-ternary */
/* eslint-disable generator-star-spacing */

/**
 * 参考 Python 的 range 类.
 * 实现了惰性求值、视图特性和 O(1) 的数学计算方法.
 */
class Range implements Iterable<number> {
  private static _empty = new Range(0)

  readonly start: number
  readonly stop: number
  readonly step: number
  readonly length: number

  constructor(stop: number)
  constructor(start: number, stop: number, step?: number)
  constructor(startOrStop: number, stop?: number, step = 1) {
    if (stop == undefined) {
      // range(stop)
      this.start = 0
      this.stop = startOrStop
      this.step = 1
    } else {
      // range(start, stop, step)
      this.start = startOrStop
      this.stop = stop
      this.step = step
    }

    if (this.step === 0) {
      throw new Error('range() arg 3 must not be zero')
    }

    if (this.step > 0 && this.start < this.stop) {
      this.length = Math.ceil((this.stop - this.start) / this.step)
    } else if (this.step < 0 && this.start > this.stop) {
      this.length = Math.ceil((this.start - this.stop) / -this.step)
    } else {
      this.length = 0
    }
  }

  at(index: number): number | undefined {
    if (index < 0) index += this.length
    if (index < 0 || index >= this.length) return undefined
    return this.start + index * this.step
  }

  slice(start?: number, stop?: number, step = 1): Range {
    if (step === 0) throw new Error('slice() step must not be zero')

    let s = start == undefined ? (step > 0 ? 0 : this.length - 1) : start
    let e = stop == undefined ? (step > 0 ? this.length : -1) : stop

    if (s < 0) s += this.length
    if (e < 0) e += this.length

    if (step > 0) {
      s = Math.max(0, Math.min(s, this.length))
      e = Math.max(0, Math.min(e, this.length))
      if (s >= e) return Range._empty
    } else {
      s = Math.max(-1, Math.min(s, this.length - 1))
      e = Math.max(-1, Math.min(e, this.length - 1))
      if (s <= e) return Range._empty
    }

    const newStart = this.start + s * this.step
    const newStep = this.step * step
    let sliceLen = 0
    if (step > 0 && s < e) {
      sliceLen = Math.ceil((e - s) / step)
    } else if (step < 0 && s > e) {
      sliceLen = Math.ceil((s - e) / -step)
    }
    const newStop = newStart + newStep * sliceLen
    return new Range(newStart, newStop, newStep)
  }

  includes(value: number): boolean {
    if (!this.length) return false
    if (this.step > 0) {
      if (value < this.start || value >= this.start + this.length * this.step) return false
    } else {
      if (value > this.start || value <= this.start + this.length * this.step) return false
    }
    // 这里假设是整数运算，如果是浮点数可能需要 epsilon
    return (value - this.start) % this.step === 0
  }

  indexOf(value: number): number {
    if (!this.includes(value)) return -1
    return (value - this.start) / this.step
  }

  count(value: number): number {
    return this.includes(value) ? 1 : 0
  }

  equals(other: Range): boolean {
    if (this.length !== other.length) return false
    if (!this.length) return true
    if (this.start !== other.start) return false
    if (this.length === 1) return true
    return this.step === other.step
  }

  reversed(): Range {
    if (!this.length) return this
    const last = this.start + (this.length - 1) * this.step
    const newStep = -this.step
    const newStop = last + this.length * newStep
    return new Range(last, newStop, newStep)
  }

  toString(): string {
    if (this.step === 1) return `range(${this.start}, ${this.stop})`
    return `range(${this.start}, ${this.stop}, ${this.step})`
  }

  toArray(): number[] {
    return Array.from(this)
  }

  *[Symbol.iterator](): Iterator<number> {
    const { start, step, length } = this
    for (let i = 0; i < length; i++) {
      yield start + i * step
    }
  }
}

/**
 * 工厂函数，模拟 Python 的 range() 调用.
 */
function range(stop: number): Range
function range(start: number, stop: number, step?: number): Range
function range(a: number, b?: number, c?: number): Range {
  if (b == undefined) return new Range(a)
  return new Range(a, b, c)
}

export { range }

{
  // 3782. 交替删除操作后最后剩下的整数
  // https://leetcode.cn/problems/last-remaining-integer-after-alternating-deletion-operations
  function lastInteger(n: number): number {
    let seq = range(1, n + 1)
    while (seq.length > 1) {
      seq = seq.slice(0, seq.length, 2).reversed()
    }
    return seq.start
  }
}
