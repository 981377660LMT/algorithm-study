/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable generator-star-spacing */

import { SortedListFast } from './SortedListFast'

/**
 * 有序字典.
 * 模拟python的`sortedcontainers.SortedDict`.
 */
class SortedDictFast<K = number, V = unknown> {
  private readonly _sl: SortedListFast<K> = new SortedListFast()
  private readonly _dict: Map<K, V> = new Map()

  constructor()
  constructor(iterable: Iterable<readonly [K, V]>)
  constructor(compareFn: (a: K, b: K) => number)
  constructor(iterable: Iterable<readonly [K, V]>, compareFn: (a: K, b: K) => number)
  constructor(compareFn: (a: K, b: K) => number, iterable: Iterable<readonly [K, V]>)
  constructor(
    arg1?: Iterable<readonly [K, V]> | ((a: K, b: K) => number),
    arg2?: Iterable<readonly [K, V]> | ((a: K, b: K) => number)
  ) {
    let defaultCompareFn = (a: any, b: any) => a - b
    let defaultData: (readonly [K, V])[] = []

    if (arg1 !== void 0) {
      if (typeof arg1 === 'function') {
        defaultCompareFn = arg1
      } else {
        defaultData = [...arg1]
      }
    }
    if (arg2 !== void 0) {
      if (typeof arg2 === 'function') {
        defaultCompareFn = arg2
      } else {
        defaultData = [...arg2]
      }
    }

    this._sl = new SortedListFast(defaultCompareFn)
    defaultData.forEach(([k, v]) => {
      this.set(k, v)
    })
  }

  set(key: K, value: V): this {
    if (!this._dict.has(key)) this._sl.add(key)
    this._dict.set(key, value)
    return this
  }

  setDefault(key: K, defaultValue: V): V {
    if (this._dict.has(key)) return this._dict.get(key)!
    this._sl.add(key)
    this._dict.set(key, defaultValue)
    return defaultValue
  }

  has(key: K): boolean {
    return this._dict.has(key)
  }

  get(key: K): V | undefined {
    return this._dict.get(key)
  }

  delete(key: K, equals?: (a: K, b: K) => boolean): boolean {
    if (!this._dict.has(key)) return false
    this._sl.discard(key, equals)
    this._dict.delete(key)
    return true
  }

  pop(key: K, equals?: (a: K, b: K) => boolean): V | undefined {
    if (!this._dict.has(key)) return void 0
    this._sl.discard(key, equals)
    const value = this._dict.get(key)!
    this._dict.delete(key)
    return value
  }

  popItem(index = -1): [K, V] | undefined {
    if (!this._dict.size) return void 0
    const key = this._sl.pop(index)!
    const value = this._dict.get(key)!
    this._dict.delete(key)
    return [key, value]
  }

  peekItem(index = -1): [K, V] | undefined {
    if (!this._dict.size) return void 0
    const key = this._sl.at(index)!
    const value = this._dict.get(key)!
    return [key, value]
  }

  peekMinItem(): [K, V] | undefined {
    if (!this._dict.size) return void 0
    const key = this._sl.min!
    const value = this._dict.get(key)!
    return [key, value]
  }

  peekMaxItem(): [K, V] | undefined {
    if (!this._dict.size) return void 0
    const key = this._sl.max!
    const value = this._dict.get(key)!
    return [key, value]
  }

  forEach(callbackfn: (value: V, key: K) => void): void {
    this._sl.forEach(key => {
      callbackfn(this._dict.get(key)!, key)
    })
  }

  enumerate(start: number, end: number, f: (key: K, value: V) => void, erase = false): void {
    this._sl.enumerate(start, end, key => {
      f(key, this._dict.get(key)!)
      if (erase) this._dict.delete(key)
    })
  }

  /**
   * 返回`key`在有序字典中的位置.
   */
  bisectLeft(key: K): number {
    return this._sl.bisectLeft(key)
  }

  /**
   * 返回`key`在有序字典中的位置.
   */
  bisectRight(key: K): number {
    return this._sl.bisectRight(key)
  }

  floor(key: K): [K, V] | undefined {
    const floorKey = this._sl.floor(key)
    if (floorKey === void 0) return void 0
    return [floorKey, this._dict.get(floorKey)!]
  }

  ceiling(key: K): [K, V] | undefined {
    const ceilingKey = this._sl.ceiling(key)
    if (ceilingKey === void 0) return void 0
    return [ceilingKey, this._dict.get(ceilingKey)!]
  }

  lower(key: K): [K, V] | undefined {
    const lowerKey = this._sl.lower(key)
    if (lowerKey === void 0) return void 0
    return [lowerKey, this._dict.get(lowerKey)!]
  }

  higher(key: K): [K, V] | undefined {
    const higherKey = this._sl.higher(key)
    if (higherKey === void 0) return void 0
    return [higherKey, this._dict.get(higherKey)!]
  }

  clear(): void {
    this._sl.clear()
    this._dict.clear()
  }

  toString(): string {
    const sb: string[] = [`SortedDict(${this.size}) {`]
    this.forEach((v, k) => {
      sb.push(`  ${k} => ${v},`)
    })
    sb.push('}')
    return sb.join('\n')
  }

  /**
   * 返回一个迭代器，用于遍历键的范围在 `[start, end)` 内的所有键值对.
   */
  *islice(start: number, end: number, reverse = false): IterableIterator<[K, V]> {
    const keys = this._sl.iSlice(start, end, reverse)
    for (const key of keys) {
      yield [key, this._dict.get(key)!]
    }
  }

  /**
   * 返回一个迭代器，用于遍历键的范围在 `[min, max] 闭区间`内的所有键值对.
   */
  *irange(min: K, max: K, reverse = false): IterableIterator<[K, V]> {
    const keys = this._sl.iRange(min, max, reverse)
    for (const key of keys) {
      yield [key, this._dict.get(key)!]
    }
  }

  *keys(): IterableIterator<K> {
    yield* this._sl
  }

  *values(): IterableIterator<V> {
    for (const key of this._sl) {
      yield this._dict.get(key)!
    }
  }

  *entries(): IterableIterator<[K, V]> {
    for (const key of this._sl) {
      yield [key, this._dict.get(key)!]
    }
  }

  *[Symbol.iterator](): IterableIterator<[K, V]> {
    for (const key of this._sl) {
      yield [key, this._dict.get(key)!]
    }
  }

  get size(): number {
    return this._dict.size
  }
}

export { SortedDictFast }

if (require.main === module) {
  const mp = new SortedDictFast<number, number>()
  mp.set(1, 1)
  mp.set(2, 2)
  mp.set(3, 3)
  console.log(mp.toString())
}
