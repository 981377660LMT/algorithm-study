// https://natsugiri.hatenablog.com/entry/2016/10/10/035445

import assert from 'assert'

type Comparator<T> = (a: T, b: T) => number

/**
 * 维护最大值和最小值的堆.
 */
class IntervalHeap<E = number> {
  private readonly _data: E[]
  private readonly _comparator: Comparator<E>

  constructor()
  constructor(array: E[])
  constructor(comparator: Comparator<E>)
  constructor(array: E[], comparator: Comparator<E>)
  constructor(comparator: Comparator<E>, array: E[])
  constructor(arrayOrComparator1?: E[] | Comparator<E>, arrayOrComparator2?: E[] | Comparator<E>) {
    let defaultArray: E[] = []
    let defaultComparator = (a: E, b: E) => (a as unknown as number) - (b as unknown as number)

    if (arrayOrComparator1) {
      if (Array.isArray(arrayOrComparator1)) {
        defaultArray = arrayOrComparator1
      } else {
        defaultComparator = arrayOrComparator1
      }
    }

    if (arrayOrComparator2) {
      if (Array.isArray(arrayOrComparator2)) {
        defaultArray = arrayOrComparator2
      } else {
        defaultComparator = arrayOrComparator2
      }
    }

    this._comparator = defaultComparator
    this._data = defaultArray
    if (this._data.length) this._heapify()
  }

  push(value: E): void {
    const k = this._data.length
    this._data.push(value)
    this._pushUp(k)
  }

  popMax(): E | undefined {
    const res = this.max
    if (this._data.length < 2) {
      this._data.pop()
      return res
    }
    this._data[0] = this._data.pop()!
    const k = this._pushDown(0)
    this._pushUp(k)
    return res
  }

  popMin(): E | undefined {
    const res = this.min
    if (this._data.length < 3) {
      this._data.pop()
      return res
    }
    this._data[1] = this._data.pop()!
    const k = this._pushDown(1)
    this._pushUp(k)
    return res
  }

  get size(): number {
    return this._data.length
  }

  get max(): E | undefined {
    return this._data[0]
  }

  get min(): E | undefined {
    return this._data.length < 2 ? this._data[0] : this._data[1]
  }

  private _heapify(): void {
    for (let i = this._data.length - 1; ~i; i--) {
      if (i & 1 && this._comparator(this._data[i - 1], this._data[i]) < 0) {
        const tmp = this._data[i - 1]
        this._data[i - 1] = this._data[i]
        this._data[i] = tmp
      }
      const k = this._pushDown(i)
      this._pushUp(k, i)
    }
  }

  private _pushUp(k: number, root = 1): number {
    const a = k & ~1
    const b = k | 1
    const data = this._data
    if (b < data.length && this._comparator(data[a], data[b]) < 0) {
      const tmp = data[a]
      data[a] = data[b]
      data[b] = tmp
      k ^= 1
    }
    let p = 0
    for (; root < k; k = p) {
      p = ((k >> 1) - 1) & ~1 // parent(k)
      if (this._comparator(data[p], data[k]) >= 0) break
      const tmp = data[p]
      data[p] = data[k]
      data[k] = tmp
    }
    for (; root < k; k = p) {
      p = (((k >> 1) - 1) & ~1) | 1
      if (this._comparator(data[k], data[p]) >= 0) break
      const tmp = data[p]
      data[p] = data[k]
      data[k] = tmp
    }
    return k
  }

  private _pushDown(k: number): number {
    const n = this._data.length
    const data = this._data
    if (k & 1) {
      while (((k << 1) | 1) < n) {
        let c = (k << 1) + 3
        if (n <= c || this._comparator(data[c - 2], data[c]) < 0) {
          c -= 2
        }
        if (c < n && this._comparator(data[c], data[k]) < 0) {
          const tmp = data[k]
          data[k] = data[c]
          data[c] = tmp
          k = c
        } else {
          break
        }
      }
    } else {
      while ((k << 1) + 2 < n) {
        let c = (k << 1) + 4
        if (n <= c || this._comparator(data[c], data[c - 2]) < 0) {
          c -= 2
        }
        if (c < n && this._comparator(data[k], data[c]) < 0) {
          const tmp = data[k]
          data[k] = data[c]
          data[c] = tmp
          k = c
        } else {
          break
        }
      }
    }
    return k
  }
}

export { IntervalHeap }

if (require.main === module) {
  const pq = new IntervalHeap<number>((a, b) => a - b, [-3, 0, 1, 3])
  pq.push(3)
  assert.strictEqual(pq.popMax(), 3)
  assert.strictEqual(pq.popMax(), 3)
  pq.push(-2)
  pq.push(1)
  assert.strictEqual(pq.popMin(), -3)
  assert.strictEqual(pq.popMin(), -2)
  assert.strictEqual(pq.popMax(), 1)
  assert.strictEqual(pq.popMin(), 0)
  assert.strictEqual(pq.popMax(), 1)
  console.log('OK')
}
