/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

const INF = 2e15

/**
 * 珂朵莉树，基于数据随机的颜色段均摊。
 * `SortedList`实现.
 * 初始时，区间为`[-INF,INF)`，值为`noneValue`.
 */
class ODTMap<S> {
  private _count = 0
  private _len = 0
  private readonly _leftLimit = -INF
  private readonly _rightLimit = INF
  private readonly _data: SortedDict<number, S> = new SortedDict()
  private readonly _noneValue: S

  /**
   * 指定哨兵值建立一个ODTMap.
   * @param noneValue 表示空值的哨兵值.
   */
  constructor(noneValue: S) {
    this._noneValue = noneValue
    this._data.set(this._leftLimit, noneValue)
    this._data.set(this._rightLimit, noneValue)
  }

  /**
   * 返回包含`x`的区间的信息.
   */
  get(x: number, erase = false): [start: number, end: number, value: S] {
    const pos2 = this._data.bisectRight(x)
    const pos1 = pos2 - 1
    const [l, vl] = this._data.peekItem(pos1)!
    const r = this._data.peekItem(pos2)![0]
    if (vl !== this._noneValue && erase) {
      this._len--
      this._count -= r - l
      this._data.set(l, this._noneValue)
      this._mergeAt(l)
      this._mergeAt(r)
    }
    return [l, r, vl]
  }

  /**
   * 将区间`[start,end)`的值置为`value`.
   */
  set(start: number, end: number, value: S): void {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerateRange(start, end, () => {}, true) // remove
    this._data.set(start, value)
    if (value !== this._noneValue) {
      this._len++
      this._count += end - start
    }
    this._mergeAt(start)
    this._mergeAt(end)
  }

  enumerateAll(f: (start: number, end: number, value: S) => void): void {
    this.enumerateRange(this._leftLimit, this._rightLimit, f, false)
  }

  /**
   * 遍历范围`[start, end)`内的所有区间.
   */
  enumerateRange(
    start: number,
    end: number,
    f: (start: number, end: number, value: S) => void,
    erase = false
  ): void {
    if (start >= end) return

    if (!erase) {
      let pos = this._data.bisectRight(start) - 1
      let [k1, v1] = this._data.peekItem(pos)!
      while (k1 < end) {
        pos++
        const [k2, v2] = this._data.peekItem(pos)!
        f(Math.max(k1, start), Math.min(k2, end), v1)
        k1 = k2
        v1 = v2
      }
      return
    }

    let pos = this._data.bisectRight(start) - 1
    let [k, v] = this._data.peekItem(pos)!
    if (k < start) {
      this._data.set(start, v)
      if (v !== this._noneValue) this._len++
    }

    pos = this._data.bisectLeft(end)
    // eslint-disable-next-line semi-style
    ;[k, v] = this._data.peekItem(pos)!
    if (k > end) {
      const v2 = this._data.peekItem(pos - 1)![1]
      this._data.set(end, v2)
      if (v2 !== this._noneValue) this._len++
    }

    pos = this._data.bisectLeft(start)
    let [k1, v1] = this._data.peekItem(pos)!
    while (k1 < end) {
      const [k2, v2] = this._data.peekItem(pos + 1)!
      f(k1, k2, v1)
      if (v1 !== this._noneValue) {
        this._len--
        this._count -= k2 - k1
      }
      this._data.popItem(pos)
      k1 = k2
      v1 = v2
    }

    this._data.set(start, this._noneValue)
  }

  toString(): string {
    const sb: string[] = [`ODTMap(${this.length}) {`]
    this.enumerateAll((start, end, value) => {
      const v = value === this._noneValue ? 'null' : value
      sb.push(`  [${start},${end}):${v}`)
    })
    sb.push('}')
    return sb.join('\n')
  }

  /**
   * 区间个数.
   */
  get length(): number {
    return this._len
  }

  /**
   * 区间内元素个数总和.
   */
  get count(): number {
    return this._count
  }

  private _mergeAt(p: number): void {
    if (p === this._leftLimit || p === this._rightLimit) return
    const pos1 = this._data.bisectLeft(p)
    const pos2 = pos1 - 1
    const v1 = this._data.peekItem(pos1)![1]
    const v2 = this._data.peekItem(pos2)![1]
    if (v1 === v2) {
      if (v1 !== this._noneValue) this._len--
      this._data.popItem(pos1)
    }
  }
}

/**
 * 有序字典.
 * 模拟python的`sortedcontainers.SortedDict`.
 */
class SortedDict<K = number, V = unknown> {
  private readonly _sl: SortedList<K> = new SortedList()
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

    this._sl = new SortedList(defaultCompareFn)
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

  /** 注意内部使用 `===` 来比较两个对象是否相等 */
  delete(key: K): boolean {
    if (!this._dict.has(key)) return false
    this._sl.discard(key)
    this._dict.delete(key)
    return true
  }

  /** 注意内部使用 `===` 来比较两个对象是否相等 */
  pop(key: K, defaultValue?: V): V | undefined {
    if (!this._dict.has(key)) return defaultValue
    this._sl.discard(key)
    const value = this._dict.get(key)!
    this._dict.delete(key)
    return value
  }

  popItem(index = -1): [K, V] | undefined {
    if (!this._dict.size) return undefined
    const key = this._sl.pop(index)!
    const value = this._dict.get(key)!
    this._dict.delete(key)
    return [key, value]
  }

  peekItem(index = -1): [K, V] | undefined {
    if (!this._dict.size) return undefined
    const key = this._sl.at(index)!
    const value = this._dict.get(key)!
    return [key, value]
  }

  peekMinItem(): [K, V] | undefined {
    if (!this._dict.size) return undefined
    const key = this._sl.min()!
    const value = this._dict.get(key)!
    return [key, value]
  }

  peekMaxItem(): [K, V] | undefined {
    if (!this._dict.size) return undefined
    const key = this._sl.max()!
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
    if (floorKey === void 0) return undefined
    return [floorKey, this._dict.get(floorKey)!]
  }

  ceiling(key: K): [K, V] | undefined {
    const ceilingKey = this._sl.ceiling(key)
    if (ceilingKey === void 0) return undefined
    return [ceilingKey, this._dict.get(ceilingKey)!]
  }

  lower(key: K): [K, V] | undefined {
    const lowerKey = this._sl.lower(key)
    if (lowerKey === void 0) return undefined
    return [lowerKey, this._dict.get(lowerKey)!]
  }

  higher(key: K): [K, V] | undefined {
    const higherKey = this._sl.higher(key)
    if (higherKey === void 0) return undefined
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
    const keys = this._sl.islice(start, end, reverse)
    for (const key of keys) {
      yield [key, this._dict.get(key)!]
    }
  }

  /**
   * 返回一个迭代器，用于遍历键的范围在 `[min, max] 闭区间`内的所有键值对.
   */
  *irange(min: K, max: K, reverse = false): IterableIterator<[K, V]> {
    const keys = this._sl.irange(min, max, reverse)
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

/**
 * A fast SortedList with O(sqrt(n)) insertion and deletion.
 *
 * @_see {@link https://github.com/981377660LMT/algorithm-study/blob/master/22_%E4%B8%93%E9%A2%98/%E7%A6%BB%E7%BA%BF%E6%9F%A5%E8%AF%A2/%E6%A0%B9%E5%8F%B7%E5%88%86%E6%B2%BB/SortedList/SortedList.ts}
 * @_see {@link https://github.com/tatyam-prime/SortedSet/blob/main/SortedMultiset.py}
 * @_see {@link https://qiita.com/tatyam/items/492c70ac4c955c055602}
 * @_see {@link https://speakerdeck.com/tatyam_prime/python-dezui-qiang-falseping-heng-er-fen-tan-suo-mu-wozuo-ru}
 */
class SortedList<T = number> {
  /** Optimized for 1e5 elements in javascript. Do not change it. */
  protected static readonly _BLOCK_RATIO = 50
  protected static readonly _REBUILD_RATIO = 170

  protected readonly _compareFn: (a: T, b: T) => number
  protected _size: number
  protected _blocks: T[][]

  constructor()
  constructor(iterable: Iterable<T>)
  constructor(compareFn: (a: T, b: T) => number)
  constructor(iterable: Iterable<T>, compareFn: (a: T, b: T) => number)
  constructor(compareFn: (a: T, b: T) => number, iterable: Iterable<T>)
  constructor(
    arg1?: Iterable<T> | ((a: T, b: T) => number),
    arg2?: Iterable<T> | ((a: T, b: T) => number)
  ) {
    let defaultCompareFn = (a: T, b: T) => (a as unknown as number) - (b as unknown as number)
    let defaultData: T[] = []
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

    this._compareFn = defaultCompareFn
    defaultData.sort(defaultCompareFn)
    this._blocks = this._initBlocks(defaultData)
    this._size = defaultData.length
  }

  add(value: T): void {
    if (!this._size) {
      this._blocks = [[value]]
      this._size = 1
      return
    }

    const blockIndex = this._findBlockIndex(value)
    if (blockIndex === -1) {
      this._blocks[this._blocks.length - 1].push(value)
      this._size++
      if (
        this._blocks[this._blocks.length - 1].length >
        this._blocks.length * SortedList._REBUILD_RATIO
      ) {
        this._rebuild()
      }
      return
    }

    const block = this._blocks[blockIndex]
    const pos = this._bisectRight(block, value)
    block.splice(pos, 0, value)
    this._size += 1
    if (block.length > this._blocks.length * SortedList._REBUILD_RATIO) {
      this._rebuild()
    }
  }

  /** 注意内部使用 `===` 来比较两个对象是否相等 */
  has(value: T): boolean {
    if (!this._size) return false
    const blockIndex = this._findBlockIndex(value)
    if (blockIndex === void 0) return false
    const block = this._blocks[blockIndex]
    const pos = this._bisectLeft(block, value)
    return pos < block.length && value === block[pos]
  }

  /** 注意内部使用 `===` 来比较两个对象是否相等 */
  discard(value: T): boolean {
    if (!this._size) return false
    const blockIndex = this._findBlockIndex(value)
    if (blockIndex === -1) return false
    const block = this._blocks[blockIndex]
    const pos = this._bisectLeft(block, value)
    if (pos === block.length || block[pos] !== value) {
      return false
    }
    block.splice(pos, 1)
    this._size -= 1
    if (!block.length) {
      this._blocks.splice(blockIndex, 1) // !Splice When Empty, Do Not Rebuild
    }
    return true
  }

  pop(index = -1): T | undefined {
    if (index < 0) index += this._size
    if (index < 0 || index >= this._size) {
      return void 0
    }
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (index < block.length) {
        const res = block[index]
        block.splice(index, 1)
        this._size -= 1
        if (!block.length) {
          this._blocks.splice(i, 1) // !Splice When Empty, Do Not Rebuild
        }
        return res
      }
      index -= block.length
    }
    return void 0
  }

  /**
   * 删除区间 [start, end) 内的元素.
   */
  erase(start: number, end: number): void {
    if (start < 0) start = 0
    if (end > this._size) end = this._size
    if (start >= end) return

    let [bid, startPos] = this._moveTo(start)
    let deleteCount = end - start
    for (; bid < this._blocks.length && deleteCount > 0; bid++) {
      const block = this._blocks[bid]
      const endPos = Math.min(block.length, startPos + deleteCount)
      const curDeleteCount = endPos - startPos
      if (curDeleteCount === block.length) {
        this._blocks.splice(bid, 1)
        bid--
      } else {
        block.splice(startPos, curDeleteCount)
      }
      deleteCount -= curDeleteCount
      this._size -= curDeleteCount
      startPos = 0
    }
  }

  at(index: number): T | undefined {
    if (index < 0) index += this._size
    if (index < 0 || index >= this._size) {
      return void 0
    }
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (index < block.length) {
        return block[index]
      }
      index -= block.length
    }
    return void 0
  }

  lower(value: T): T | undefined {
    for (let i = this._blocks.length - 1; ~i; i--) {
      const block = this._blocks[i]
      if (this._compareFn(block[0], value) < 0) {
        const pos = this._bisectLeft(block, value)
        return block[pos - 1]
      }
    }
    return void 0
  }

  higher(value: T): T | undefined {
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (this._compareFn(block[block.length - 1], value) > 0) {
        const pos = this._bisectRight(block, value)
        return block[pos]
      }
    }
    return void 0
  }

  floor(value: T): T | undefined {
    for (let i = this._blocks.length - 1; ~i; i--) {
      const block = this._blocks[i]
      if (this._compareFn(block[0], value) <= 0) {
        const pos = this._bisectRight(block, value)
        return block[pos - 1]
      }
    }
    return void 0
  }

  ceiling(value: T): T | undefined {
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (this._compareFn(block[block.length - 1], value) >= 0) {
        const pos = this._bisectLeft(block, value)
        return block[pos]
      }
    }
    return void 0
  }

  /**
   * Count the number of elements < value or
   * returns the index of the first element >= value.
   */
  bisectLeft(value: T): number {
    let res = 0
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (this._compareFn(value, block[block.length - 1]) <= 0) {
        return res + this._bisectLeft(block, value)
      }
      res += block.length
    }
    return res
  }

  /**
   * Count the number of elements <= value or
   * returns the index of the first element > value.
   */
  bisectRight(value: T): number {
    let res = 0
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (this._compareFn(value, block[block.length - 1]) < 0) {
        return res + this._bisectRight(block, value)
      }
      res += block.length
    }
    return res
  }

  slice(start: number, end: number): T[] {
    if (start < 0) start = 0
    if (end > this._size) end = this._size
    if (start >= end) return []
    let count = end - start
    const res: T[] = Array(count)
    let [bid, startPos] = this._moveTo(start)
    let ptr = 0
    for (; bid < this._blocks.length && count > 0; bid++) {
      const block = this._blocks[bid]
      const endPos = Math.min(block.length, startPos + count)
      const curCount = endPos - startPos
      for (let j = startPos; j < endPos; j++) {
        res[ptr++] = block[j]
      }
      count -= curCount
      startPos = 0
    }
    return res
  }

  clear(): void {
    this._blocks = []
    this._size = 0
  }

  min(): T | undefined {
    if (!this._size) return void 0
    return this._blocks[0][0]
  }

  max(): T | undefined {
    if (!this._size) return void 0
    const last = this._blocks[this._blocks.length - 1]
    return last[last.length - 1]
  }

  toString(): string {
    return `SortedList[${[...this].join(', ')}]`
  }

  /**
   * 返回一个迭代器，用于遍历区间 [start, end) 内的元素.
   */
  *islice(start: number, end: number, reverse = false): IterableIterator<T> {
    if (start < 0) start = 0
    if (end > this._size) end = this._size
    if (start >= end) return
    let count = end - start

    if (reverse) {
      let [bid, endPos] = this._moveTo(end)
      for (; ~bid && count > 0; bid--, ~bid && (endPos = this._blocks[bid].length)) {
        const block = this._blocks[bid]
        const startPos = Math.max(0, endPos - count)
        const curCount = endPos - startPos
        for (let j = endPos - 1; j >= startPos; j--) {
          yield block[j]
        }
        count -= curCount
      }
    } else {
      let [bid, startPos] = this._moveTo(start)
      for (; bid < this._blocks.length && count > 0; bid++) {
        const block = this._blocks[bid]
        const endPos = Math.min(block.length, startPos + count)
        const curCount = endPos - startPos
        for (let j = startPos; j < endPos; j++) {
          yield block[j]
        }
        count -= curCount
        startPos = 0
      }
    }
  }

  /**
   * 返回一个迭代器，用于遍历范围在 `[min, max] 闭区间`内的元素.
   */
  *irange(min: T, max: T, reverse = false): IterableIterator<T> {
    if (this._compareFn(min, max) > 0) {
      return
    }

    if (reverse) {
      let bi = this._findBlockIndex(max)
      if (bi === -1) {
        bi = this._blocks.length - 1
      }
      for (let i = bi; ~i; i--) {
        const block = this._blocks[i]
        for (let j = block.length - 1; ~j; j--) {
          const x = block[j]
          if (this._compareFn(x, min) < 0) {
            return
          }
          if (this._compareFn(x, max) <= 0) {
            yield x
          }
        }
      }
    } else {
      const bi = this._findBlockIndex(min)
      for (let i = bi; i < this._blocks.length; i++) {
        const block = this._blocks[i]
        for (let j = 0; j < block.length; j++) {
          const x = block[j]
          if (this._compareFn(x, max) > 0) {
            return
          }
          if (this._compareFn(x, min) >= 0) {
            yield x
          }
        }
      }
    }
  }

  enumerate(start: number, end: number, f: (value: T) => void, erase = false): void {
    let [bid, startPos] = this._moveTo(start)
    let count = end - start

    for (; bid < this._blocks.length && count > 0; bid++) {
      const block = this._blocks[bid]
      const endPos = Math.min(block.length, startPos + count)
      for (let j = startPos; j < endPos; j++) {
        f(block[j])
      }

      const curDeleteCount = endPos - startPos
      if (erase) {
        if (curDeleteCount === block.length) {
          this._blocks.splice(bid, 1)
          bid--
        } else {
          block.splice(startPos, curDeleteCount)
        }
        this._size -= curDeleteCount
      }

      count -= curDeleteCount
      startPos = 0
    }
  }

  forEach(callbackfn: (value: T, index: number) => void): void {
    let pos = 0
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      for (let j = 0; j < block.length; j++) {
        callbackfn(block[j], pos++)
      }
    }
  }

  *entries(): IterableIterator<[number, T]> {
    let pos = 0
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      for (let j = 0; j < block.length; j++) {
        yield [pos++, block[j]]
      }
    }
  }

  *[Symbol.iterator](): Iterator<T> {
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      for (let j = 0; j < block.length; j++) {
        yield block[j]
      }
    }
  }

  get length(): number {
    return this._size
  }

  // Find the block which should contain x. Block must not be empty.
  private _findBlockIndex(x: T): number {
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (this._compareFn(x, block[block.length - 1]) <= 0) {
        return i
      }
    }
    return -1
  }

  private _rebuild(): void {
    if (!this._size) {
      return
    }
    const bCount = Math.ceil(Math.sqrt(this._size / SortedList._BLOCK_RATIO))
    const bSize = ~~((this._size + bCount - 1) / bCount) // ceil
    const newB: T[][] = Array(bCount)
    for (let i = 0; i < bCount; i++) {
      newB[i] = []
    }
    let ptr = 0
    for (let i = 0; i < this._blocks.length; i++) {
      const b = this._blocks[i]
      for (let j = 0; j < b.length; j++) {
        newB[~~(ptr / bSize)].push(b[j])
        ptr++
      }
    }
    this._blocks = newB
  }

  // eslint-disable-next-line class-methods-use-this
  private _initBlocks(sorted: T[]): T[][] {
    const bCount = Math.ceil(Math.sqrt(sorted.length / SortedList._BLOCK_RATIO))
    const bSize = ~~((sorted.length + bCount - 1) / bCount) // ceil
    const newB: T[][] = Array(bCount)
    for (let i = 0; i < bCount; i++) {
      newB[i] = []
    }
    for (let i = 0; i < bCount; i++) {
      for (let j = i * bSize; j < Math.min((i + 1) * bSize, sorted.length); j++) {
        newB[i].push(sorted[j])
      }
    }
    return newB
  }

  private _bisectLeft(nums: T[], value: T): number {
    let left = 0
    let right = nums.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      if (this._compareFn(value, nums[mid]) <= 0) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left
  }

  private _bisectRight(nums: T[], value: T): number {
    let left = 0
    let right = nums.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      if (this._compareFn(value, nums[mid]) < 0) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left
  }

  private _moveTo(index: number): [blockId: number, startPos: number] {
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (index < block.length) {
        return [i, index]
      }
      index -= block.length
    }
    return [this._blocks.length, 0]
  }
}

export { ODTMap }

if (require.main === module) {
  const odtMap = new ODTMap<number>(-1)
  odtMap.set(1, 3, 1)
  console.log(odtMap.count, 999, odtMap.length)
  odtMap.set(3, 5, 2)
  console.log(odtMap.toString())
  console.log(odtMap.get(2))
  console.log(odtMap.length, odtMap.count)

  // https://leetcode.cn/problems/count-integers-in-intervals/
  class CountIntervals {
    private readonly _odtMap = new ODTMap<number>(-1)

    add(left: number, right: number): void {
      this._odtMap.set(left, right + 1, 1)
    }

    count(): number {
      return this._odtMap.count
    }
  }

  // https://leetcode.cn/problems/data-stream-as-disjoint-intervals/
  // 352. 将数据流变为多个不相交区间
  class SummaryRanges {
    private readonly _odt = new ODTMap(-1)

    addNum(value: number): void {
      this._odt.set(value, value + 1, 0)
    }

    getIntervals(): number[][] {
      const res: number[][] = []
      this._odt.enumerateAll((start, end, value) => {
        if (value === 0) res.push([start, end - 1])
      })
      return res
    }
  }

  // https://leetcode.cn/problems/range-module/
  class RangeModule {
    private readonly _odt = new ODTMap(-1)

    addRange(left: number, right: number): void {
      this._odt.set(left, right, 0)
    }

    queryRange(left: number, right: number): boolean {
      const [start, end, value] = this._odt.get(left)
      return start <= left && right <= end && value === 0
    }

    removeRange(left: number, right: number): void {
      this._odt.set(left, right, -1)
    }
  }
}
