/**
 * 利用位运算寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
 */
class FastSet {
  private readonly _n: number
  private readonly _lg: number
  private readonly _seg: Uint32Array[]
  private _size = 0

  /**
   * @param n [0, n).
   * @param initValue 初始化时是否全部置为1.默认初始时所有位置都未被访问过.
   */
  constructor(n: number, initValue = false) {
    this._n = n
    const seg: Uint32Array[] = []
    while (true) {
      if (initValue) {
        seg.push(new Uint32Array((n + 31) >>> 5).fill(-1))
      } else {
        seg.push(new Uint32Array((n + 31) >>> 5))
      }
      n = (n + 31) >>> 5
      if (n <= 1) {
        break
      }
    }
    this._lg = seg.length
    this._seg = seg
    this._size = initValue ? this._n : 0
  }

  insert(i: number): boolean {
    if (this.has(i)) return false
    for (let h = 0; h < this._lg; h++) {
      this._seg[h][i >>> 5] |= 1 << (i & 31)
      i >>>= 5
    }
    this._size++
    return true
  }

  has(i: number): boolean {
    return !!(this._seg[0][i >>> 5] & (1 << (i & 31)))
  }

  erase(i: number): boolean {
    if (!this.has(i)) return false
    for (let h = 0; h < this._lg; h++) {
      const cache = this._seg[h]
      cache[i >>> 5] &= ~(1 << (i & 31))
      if (cache[i >>> 5]) break
      i >>>= 5
    }
    this._size--
    return true
  }

  /**
   * 返回x右侧第一个未被访问过的位置(包含x).
   * 如果不存在,返回`n`.
   */
  next(i: number): number {
    if (i < 0) i = 0
    if (i >= this._n) return this._n

    for (let h = 0; h < this._lg; h++) {
      const cacheH = this._seg[h]
      if (i >>> 5 === cacheH.length) break
      let d = cacheH[i >>> 5] >>> (i & 31)
      if (d === 0) {
        i = (i >>> 5) + 1
        continue
      }
      // !trailingZeros32: 31 - Math.clz32(x & -x)
      i += 31 - Math.clz32(d & -d)
      for (let g = h - 1; ~g; g--) {
        i <<= 5
        const tmp = this._seg[g][i >>> 5]
        i += 31 - Math.clz32(tmp & -tmp)
      }
      return i
    }

    return this._n
  }

  /**
   * 返回x左侧第一个未被访问过的位置(包含x).
   * 如果不存在,返回`-1`.
   */
  prev(i: number): number {
    if (i < 0) return -1
    if (i >= this._n) i = this._n - 1

    for (let h = 0; h < this._lg; h++) {
      if (i === -1) break
      let d = this._seg[h][i >>> 5] << (31 - (i & 31))
      if (d === 0) {
        i = (i >>> 5) - 1
        continue
      }

      i -= Math.clz32(d)
      for (let g = h - 1; ~g; g--) {
        i <<= 5
        i += 31 - Math.clz32(this._seg[g][i >>> 5])
      }
      return i
    }

    return -1
  }

  /**
   * 遍历[start,end)区间内的元素.
   */
  enumerateRange(start: number, end: number, f: (v: number) => void): void {
    for (let x = this.next(start); x < end; x = this.next(x + 1)) {
      f(x)
    }
  }

  toString(): string {
    const sb: string[] = []
    this.enumerateRange(0, this._n, v => sb.push(v.toString()))
    return `FastSet(${this.size}){${sb.join(', ')}}`
  }

  get min(): number {
    return this.next(-1)
  }

  get max(): number {
    return this.prev(this._n)
  }

  get size(): number {
    return this._size
  }
}

/**
 * 珂朵莉树，基于数据随机的颜色段均摊。
 * `FastSet`实现.
 */
class ODT<S> {
  private _len = 0
  private _count = 0
  private readonly _leftLimit: number
  private readonly _rightLimit: number
  private readonly _noneValue: S
  private readonly _data: S[]
  private readonly _fs: FastSet

  /**
   * 指定区间长度和哨兵值建立一个ODT.初始时,所有位置的值为 {@link noneValue}.
   * @param n 区间范围为`[0, n)`.
   * @param noneValue 表示空值的哨兵值.
   */
  constructor(n: number, noneValue: S) {
    const data = Array(n)
    for (let i = 0; i < n; i++) data[i] = noneValue
    const fs = new FastSet(n)
    fs.insert(0)

    this._leftLimit = 0
    this._rightLimit = n
    this._noneValue = noneValue
    this._data = data
    this._fs = fs
  }

  /**
   * 返回包含`x`的区间的信息.
   * 0 <= x < n.
   */
  get(x: number, erase = false): [start: number, end: number, value: S] | undefined {
    if (x < this._leftLimit || x >= this._rightLimit) return undefined
    const start = this._fs.prev(x)
    const end = this._fs.next(x + 1)
    const value = this._data[start]
    if (erase && value !== this._noneValue) {
      this._len--
      this._count -= end - start
      this._data[start] = this._noneValue
      this._mergeAt(start)
      this._mergeAt(end)
    }
    return [start, end, value]
  }

  set(start: number, end: number, value: S): void {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerateRange(start, end, () => {}, true) // remove
    this._fs.insert(start)
    this._data[start] = value
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
    if (start < this._leftLimit) start = this._leftLimit
    if (end > this._rightLimit) end = this._rightLimit
    if (start >= end) return

    const none = this._noneValue
    if (!erase) {
      let left = this._fs.prev(start)
      while (left < end) {
        const right = this._fs.next(left + 1)
        f(Math.max(left, start), Math.min(right, end), this._data[left])
        left = right
      }
      return
    }

    let p = this._fs.prev(start)
    if (p < start) {
      this._fs.insert(start)
      const v = this._data[p]
      this._data[start] = v
      if (v !== none) {
        this._len++
      }
    }

    p = this._fs.next(end)
    if (end < p) {
      const v = this._data[this._fs.prev(end)]
      this._data[end] = v
      this._fs.insert(end)
      if (v !== none) {
        this._len++
      }
    }

    p = start
    while (p < end) {
      const q = this._fs.next(p + 1)
      const x = this._data[p]
      f(p, q, x)
      if (x !== none) {
        this._len--
        this._count -= q - p
      }
      this._fs.erase(p)
      p = q
    }

    this._fs.insert(start)
    this._data[start] = none
  }

  toString(): string {
    const sb: string[] = [`ODT(${this.length}) {`]
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
   * 区间内元素个数之和.
   */
  get count(): number {
    return this._count
  }

  private _mergeAt(p: number): void {
    if (p <= 0 || this._rightLimit <= p) return
    const q = this._fs.prev(p - 1)
    const dataP = this._data[p]
    const dataQ = this._data[q]
    if (dataP === dataQ) {
      if (dataP !== this._noneValue) this._len--
      this._fs.erase(p)
    }
  }
}

// API.
interface ISortedList<V> {
  add(value: V): void
  has(value: V, equals?: (a: V, b: V) => boolean): boolean
  discard(value: V, equals?: (a: V, b: V) => boolean): boolean
  pop(index?: number): V | undefined
  at(index: number): V | undefined
  erase(start?: number, end?: number): void

  lower(value: V): V | undefined
  higher(value: V): V | undefined
  floor(value: V): V | undefined
  ceiling(value: V): V | undefined

  bisectLeft(value: V): number
  bisectRight(value: V): number
  count(value: V): number

  slice(start?: number, end?: number): V[]
  clear(): void
  update(...values: V[]): void
  merge(other: ISortedList<V>): void

  toString(): string

  forEach(callbackfn: (value: V, index: number) => void | boolean, reverse?: boolean): void
  enumerate(start: number, end: number, f?: (value: V) => void, erase?: boolean): void
  iSlice(start: number, end: number, reverse?: boolean): IterableIterator<V>
  iRange(min: V, max: V, reverse?: boolean): IterableIterator<V>
  entries(): IterableIterator<[number, V]>
  [Symbol.iterator](): IterableIterator<V>

  readonly length: number
  readonly min: V | undefined
  readonly max: V | undefined
}

/**
 * 使用分块+树状数组维护的有序序列.
 * !如果数组较短(<2000),直接使用`bisectInsort`维护即可.
 */
class SortedListFast<V = number> implements ISortedList<V> {
  static setLoadFactor(load = 200): void {
    this._LOAD = load
  }

  /**
   * 负载因子，用于控制每个块的长度.
   * 长度为`1e5`的数组, 负载因子取`200`左右性能较好.
   */
  protected static _LOAD = 200
  private static _isPrimitive(
    o: unknown
  ): o is number | string | boolean | symbol | bigint | null | undefined {
    return o === null || (typeof o !== 'object' && typeof o !== 'function')
  }

  protected _compareFn!: (a: V, b: V) => number
  protected _len!: number

  /**
   * 各个块.
   */
  protected _blocks!: V[][]

  /**
   * 各个块的最小值.
   */
  protected _mins!: V[]

  /**
   * 树状数组维护各个块长度的前缀和，便于根据索引定位到对应的元素.
   */
  protected _tree!: number[]
  protected _shouldRebuildTree!: boolean

  constructor()
  constructor(iterable: Iterable<V>)
  constructor(compareFn: (a: V, b: V) => number)
  constructor(iterable: Iterable<V>, compareFn: (a: V, b: V) => number)
  constructor(compareFn: (a: V, b: V) => number, iterable: Iterable<V>)
  constructor(
    arg1?: Iterable<V> | ((a: V, b: V) => number),
    arg2?: Iterable<V> | ((a: V, b: V) => number)
  ) {
    let defaultCompareFn = (a: V, b: V) => (a as unknown as number) - (b as unknown as number)
    let defaultData: V[] = []
    if (arg1 !== undefined) {
      if (typeof arg1 === 'function') {
        defaultCompareFn = arg1
      } else {
        defaultData = [...arg1]
      }
    }
    if (arg2 !== undefined) {
      if (typeof arg2 === 'function') {
        defaultCompareFn = arg2
      } else {
        defaultData = [...arg2]
      }
    }

    this._init(defaultData, defaultCompareFn)
  }

  add(value: V): this {
    const { _blocks, _mins } = this
    this._len++
    if (!_blocks.length) {
      _blocks.push([value])
      _mins.push(value)
      this._shouldRebuildTree = true
      return this
    }

    const load = SortedListFast._LOAD
    const { pos, index } = this._locRight(value)
    this._updateTree(pos, 1)
    const block = _blocks[pos]
    block.splice(index, 0, value)
    _mins[pos] = block[0]
    // !-> [x]*load + [x]*(block.length - load)
    if (load + load < block.length) {
      const rightBlock = block.splice(load)
      _blocks.splice(pos + 1, 0, rightBlock)
      _mins.splice(pos + 1, 0, rightBlock[0])
      this._shouldRebuildTree = true
    }
    return this
  }

  update(...values: V[]): this {
    const n = values.length
    if (n < this._len << 2) {
      for (let i = 0; i < n; i++) this.add(values[i])
      return this
    }

    const data = Array(this._len + n)
    let ptr = 0
    this.forEach(v => {
      data[ptr++] = v
    })

    for (let i = 0; i < n; i++) data[ptr++] = values[i]
    this._reBuild(data, this._compareFn)
    return this
  }

  merge(other: SortedListFast<V>): this {
    const n = other.length
    if (n < this._len << 2) {
      other.forEach(v => {
        this.add(v)
      })
      return this
    }

    const data = Array(this._len + n)
    let ptr = 0
    this.forEach(v => {
      data[ptr++] = v
    })

    other.forEach(v => {
      data[ptr++] = v
    })
    this._reBuild(data, this._compareFn)
    return this
  }

  has(value: V, equals?: (a: V, b: V) => boolean): boolean {
    if (!equals && !SortedListFast._isPrimitive(value)) {
      throw new Error('equals must be provided when value is a non-primitive value')
    }
    const block = this._blocks
    if (!block.length) return false
    const { pos, index } = this._locLeft(value)
    equals = equals || ((a: V, b: V) => a === b)
    return index < block[pos].length && equals(block[pos][index], value)
  }

  discard(value: V, equals?: (a: V, b: V) => boolean): boolean {
    if (!equals && !SortedListFast._isPrimitive(value)) {
      throw new Error('equals must be provided when value is a non-primitive value')
    }
    const block = this._blocks
    if (!block.length) return false
    const { pos, index } = this._locRight(value)
    equals = equals || ((a: V, b: V) => a === b)

    if (index && equals(block[pos][index - 1], value)) {
      this._delete(pos, index - 1)
      return true
    }
    return false
  }

  pop(index = -1): V | undefined {
    if (index < 0) index += this._len
    if (index < 0 || index >= this._len) return undefined
    const { pos, index: index_ } = this._findKth(index)
    const value = this._blocks[pos][index_]
    this._delete(pos, index_)
    return value
  }

  /**
   * @complexity O(log(blockCount))
   */
  at(index: number): V | undefined {
    if (index < 0) index += this._len
    if (index < 0 || index >= this._len) return undefined
    const { pos, index: index_ } = this._findKth(index)
    return this._blocks[pos][index_]
  }

  /**
   * 删除区间 `[start, end)` 内的元素.
   */
  erase(start = 0, end = this.length): void {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerate(start, end, undefined, true)
  }

  lower(value: V): V | undefined {
    const pos = this.bisectLeft(value)
    return pos ? this.at(pos - 1) : undefined
  }

  higher(value: V): V | undefined {
    const pos = this.bisectRight(value)
    return pos < this._len ? this.at(pos) : undefined
  }

  floor(value: V): V | undefined {
    const pos = this.bisectRight(value)
    return pos ? this.at(pos - 1) : undefined
  }

  ceiling(value: V): V | undefined {
    const pos = this.bisectLeft(value)
    return pos < this._len ? this.at(pos) : undefined
  }

  /**
   * 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
   */
  bisectLeft(value: V): number {
    const { pos, index } = this._locLeft(value)
    return this._queryTree(pos) + index
  }

  /**
   * 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
   */
  bisectRight(value: V): number {
    const { pos, index } = this._locRight(value)
    return this._queryTree(pos) + index
  }

  count(value: V): number {
    return this.bisectRight(value) - this.bisectLeft(value)
  }

  slice(start = 0, end = this.length): V[] {
    if (start < 0) start += this._len
    if (start < 0) start = 0
    if (end < 0) end += this._len
    if (end > this._len) end = this._len
    if (start >= end) return []
    const res = Array(end - start)
    let count = 0
    this.enumerate(
      start,
      end,
      value => {
        res[count++] = value
      },
      false
    )
    return res
  }

  clear(): void {
    this._len = 0
    this._blocks = []
    this._mins = []
    this._tree = []
    this._shouldRebuildTree = true
  }

  toString(): string {
    const res = Array(this._len)
    this.forEach((value, index) => {
      res[index] = value
    })
    return `SortedList{${JSON.stringify(res)}}`
  }

  /**
   * @param callbackfn 回调函数, 返回 `true` 时终止遍历.
   * @param reverse 是否逆序遍历.
   */
  forEach(callbackfn: (value: V, index: number) => void | boolean, reverse = false): void {
    if (!reverse) {
      let count = 0
      for (let i = 0; i < this._blocks.length; i++) {
        const block = this._blocks[i]
        for (let j = 0; j < block.length; j++) {
          if (callbackfn(block[j], count++)) return
        }
      }
      return
    }

    let count = 0
    for (let i = this._blocks.length - 1; ~i; i--) {
      const block = this._blocks[i]
      for (let j = block.length - 1; ~j; j--) {
        if (callbackfn(block[j], count++)) return
      }
    }
  }

  enumerate(start: number, end: number, f?: (value: V) => void, erase = false): void {
    if (start < 0) start = 0
    if (end > this._len) end = this._len
    if (start >= end) return

    const pair = this._findKth(start)
    let pos = pair.pos
    let startIndex = pair.index
    let count = end - start
    for (; count && pos < this._blocks.length; pos++) {
      const block = this._blocks[pos]
      const endIndex = Math.min(block.length, startIndex + count)
      if (f) {
        for (let j = startIndex; j < endIndex; j++) f(block[j])
      }
      const deleted = endIndex - startIndex

      if (erase) {
        if (deleted === block.length) {
          // !delete block
          this._blocks.splice(pos, 1)
          this._mins.splice(pos, 1)
          this._shouldRebuildTree = true
          pos--
        } else {
          // !delete [index, end)
          this._updateTree(pos, -deleted)
          block.splice(startIndex, deleted)
          this._mins[pos] = block[0]
        }
        this._len -= deleted
      }

      count -= deleted
      startIndex = 0
    }
  }

  /**
   * 返回一个迭代器，用于遍历区间 `[start, end)` 内的元素.
   */
  *iSlice(start = 0, end = this.length, reverse = false): IterableIterator<V> {
    if (start < 0) start += this._len
    if (start < 0) start = 0
    if (end < 0) end += this._len
    if (end > this._len) end = this._len
    if (start >= end) return
    let count = end - start

    if (reverse) {
      let { pos, index } = this._findKth(end)
      for (; count && ~pos; pos--, ~pos && (index = this._blocks[pos].length)) {
        const block = this._blocks[pos]
        const startPos = Math.max(0, index - count)
        const curCount = index - startPos
        for (let j = index - 1; j >= startPos; j--) yield block[j]
        count -= curCount
      }
    } else {
      let { pos, index } = this._findKth(start)
      for (; count && pos < this._blocks.length; pos++) {
        const block = this._blocks[pos]
        const endPos = Math.min(block.length, index + count)
        const curCount = endPos - index
        for (let j = index; j < endPos; j++) yield block[j]
        count -= curCount
        index = 0
      }
    }
  }

  /**
   * 返回一个迭代器，用于遍历范围在 `[min, max] 闭区间`内的元素.
   */
  *iRange(min: V, max: V, reverse = false): IterableIterator<V> {
    if (this._compareFn(min, max) > 0) return

    if (reverse) {
      let pos = this._locBlock(max)
      for (let i = pos; ~i; i--) {
        const block = this._blocks[i]
        for (let j = block.length - 1; ~j; j--) {
          const x = block[j]
          if (this._compareFn(x, min) < 0) return
          if (this._compareFn(x, max) <= 0) yield x
        }
      }
    } else {
      const pos = this._locBlock(min)
      for (let i = pos; i < this._blocks.length; i++) {
        const block = this._blocks[i]
        for (let j = 0; j < block.length; j++) {
          const x = block[j]
          if (this._compareFn(x, max) > 0) return
          if (this._compareFn(x, min) >= 0) yield x
        }
      }
    }
  }

  *entries(): IterableIterator<[number, V]> {
    let count = 0
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      for (let j = 0; j < block.length; j++) {
        yield [count++, block[j]]
      }
    }
  }

  *[Symbol.iterator](): IterableIterator<V> {
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      for (let j = 0; j < block.length; j++) {
        yield block[j]
      }
    }
  }

  get length(): number {
    return this._len
  }

  get min(): V | undefined {
    if (!this._len) return undefined
    return this._mins[0]
  }

  get max(): V | undefined {
    if (!this._len) return undefined
    const { _blocks } = this
    const lastBlock = _blocks[_blocks.length - 1]
    return lastBlock[lastBlock.length - 1]
  }

  protected _init(data: V[], compareFn: (a: V, b: V) => number): void {
    this._reBuild(data, compareFn)
  }

  protected _reBuild(data: V[], compareFn: (a: V, b: V) => number): void {
    data.sort(compareFn)
    const n = data.length
    const blocks = []
    for (let i = 0; i < n; i += SortedListFast._LOAD) {
      blocks.push(data.slice(i, i + SortedListFast._LOAD))
    }
    const mins = Array(blocks.length)
    for (let i = 0; i < blocks.length; i++) {
      const cur = blocks[i]
      mins[i] = cur[0]
    }

    this._compareFn = compareFn
    this._len = n
    this._blocks = blocks
    this._mins = mins
    this._tree = []
    this._shouldRebuildTree = true
  }

  protected _delete(pos: number, index: number): void {
    const { _blocks, _mins } = this

    // !delete element
    this._len--
    this._updateTree(pos, -1)
    const block = _blocks[pos]
    block.splice(index, 1)
    if (block.length) {
      _mins[pos] = block[0]
      return
    }

    // !delete block
    _blocks.splice(pos, 1)
    _mins.splice(pos, 1)
    this._shouldRebuildTree = true
  }

  protected _locLeft(value: V): { pos: number; index: number } {
    if (!this._len) return { pos: 0, index: 0 }
    const { _blocks, _mins } = this

    // find pos
    let left = -1
    let right = _blocks.length - 1
    while (left + 1 < right) {
      const mid = (left + right) >>> 1
      if (this._compareFn(value, _mins[mid]) <= 0) {
        right = mid
      } else {
        left = mid
      }
    }
    if (right) {
      const block = _blocks[right - 1]
      if (this._compareFn(value, block[block.length - 1]) <= 0) {
        right--
      }
    }
    const pos = right

    // find index
    const cur = _blocks[pos]
    left = -1
    right = cur.length
    while (left + 1 < right) {
      const mid = (left + right) >>> 1
      if (this._compareFn(value, cur[mid]) <= 0) {
        right = mid
      } else {
        left = mid
      }
    }

    return { pos, index: right }
  }

  protected _locRight(value: V): { pos: number; index: number } {
    if (!this._len) return { pos: 0, index: 0 }
    const { _blocks, _mins } = this

    // find pos
    let left = 0
    let right = _blocks.length
    while (left + 1 < right) {
      const mid = (left + right) >>> 1
      if (this._compareFn(value, _mins[mid]) < 0) {
        right = mid
      } else {
        left = mid
      }
    }
    const pos = left

    // find index
    const cur = _blocks[pos]
    left = -1
    right = cur.length
    while (left + 1 < right) {
      const mid = (left + right) >>> 1
      if (this._compareFn(value, cur[mid]) < 0) {
        right = mid
      } else {
        left = mid
      }
    }

    return { pos, index: right }
  }

  private _locBlock(value: V): number {
    let left = -1
    let right = this._blocks.length - 1
    while (left + 1 < right) {
      const mid = (left + right) >>> 1
      if (this._compareFn(value, this._mins[mid]) <= 0) {
        right = mid
      } else {
        left = mid
      }
    }
    if (right) {
      const block = this._blocks[right - 1]
      if (this._compareFn(value, block[block.length - 1]) <= 0) {
        right--
      }
    }
    return right
  }

  private _buildTree(): void {
    this._tree = Array(this._blocks.length)
    for (let i = 0; i < this._blocks.length; i++) {
      this._tree[i] = this._blocks[i].length
    }
    const tree = this._tree
    for (let i = 0; i < tree.length; i++) {
      const j = i | (i + 1)
      if (j < tree.length) {
        tree[j] += tree[i]
      }
    }
    this._shouldRebuildTree = false
  }

  protected _updateTree(index: number, delta: number): void {
    if (!this._shouldRebuildTree) {
      const tree = this._tree
      while (index < tree.length) {
        tree[index] += delta
        index |= index + 1
      }
    }
  }

  private _queryTree(end: number): number {
    if (this._shouldRebuildTree) this._buildTree()
    const tree = this._tree
    let x = 0
    while (end) {
      x += tree[end - 1]
      end &= end - 1
    }
    return x
  }

  /**
   * 树状数组树上二分, 找到索引为 `k` 的元素的`(所在块的索引pos, 块内的索引index)`.
   * 内部对头部块和尾部块做了特殊处理.
   */
  protected _findKth(k: number): { pos: number; index: number } {
    if (k < this._blocks[0].length) return { pos: 0, index: k }
    const last = this._blocks.length - 1
    const lastLen = this._blocks[last].length
    if (k >= this._len - lastLen) return { pos: last, index: k + lastLen - this._len }
    if (this._shouldRebuildTree) this._buildTree()
    const tree = this._tree
    let pos = -1
    const bitLength = 32 - Math.clz32(tree.length)
    for (let d = bitLength - 1; ~d; d--) {
      const next = pos + (1 << d)
      if (next < tree.length && k >= tree[next]) {
        pos = next
        k -= tree[pos]
      }
    }
    pos++
    return { pos, index: k }
  }
}

interface Interval {
  start: number
  end: number
}

const equals = (a: Interval, b: Interval) => a.start === b.start && a.end === b.end
const compare = (a: Interval, b: Interval) => {
  const len1 = a.end - a.start
  const len2 = b.end - b.start
  return len1 - len2 || a.start - b.start
}

/**
 * 维护相同元素的最长连续长度.
 */
class LongestRepeating<T> {
  private static readonly _NONE: any = Symbol('NONE')
  private readonly _intervals = new SortedListFast<Interval>(compare)
  private readonly _n: number
  private readonly _odt: ODT<T>

  constructor(arr: ArrayLike<T>) {
    const n = arr.length
    this._n = n
    this._odt = new ODT(n, LongestRepeating._NONE)

    let pre = 0
    let pos = 0
    while (pos < n) {
      const leader = arr[pos]
      pos++
      while (pos < n && arr[pos] === leader) pos++
      this._intervals.add({ start: pre, end: pos })
      this._odt.set(pre, pos, leader)
      pre = pos
    }
  }

  queryAll(): Interval {
    return this._intervals.at(-1)!
  }

  update(start: number, end: number, value: T): void {
    const leftStart = this._odt.get(start)![0]
    const rightEnd = this._odt.get(end - 1)![1]
    const leftSeg = this._odt.get(leftStart - 1)
    const rightSeg = this._odt.get(rightEnd)
    const first = leftSeg ? leftSeg[0] : 0
    const last = rightSeg ? rightSeg[1] : this._n
    this._odt.enumerateRange(first, last, (s, e) => {
      this._intervals.discard({ start: s, end: e }, equals)
    })
    this._odt.set(start, end, value)
    this._odt.enumerateRange(first, last, (s, e) => {
      this._intervals.add({ start: s, end: e })
    })
  }

  set(index: number, value: T): void {
    const [start, end] = this._odt.get(index)!
    const leftSeg = this._odt.get(start - 1)
    const rightSeg = this._odt.get(end)
    const first = leftSeg ? leftSeg[0] : 0
    const last = rightSeg ? rightSeg[1] : this._n
    this._odt.enumerateRange(first, last, (s, e) => {
      this._intervals.discard({ start: s, end: e }, equals)
    })
    this._odt.set(index, index + 1, value)
    this._odt.enumerateRange(first, last, (s, e) => {
      this._intervals.add({ start: s, end: e })
    })
  }

  toString(): string {
    return this._intervals.toString()
  }
}

function minLength(s: string, numOps: number): number {
  const bits = s.split('').map(Number)
  const L = new LongestRepeating(bits)
  for (let i = 0; i < numOps; i++) {
    const { start, end } = L.queryAll()
    if (end - start <= 1) break

    let mid = (start + end) >>> 1
    if (mid + 1 < bits.length && bits[mid] !== bits[mid + 1]) {
      mid--
    }
    if (mid - 1 >= 0 && bits[mid] !== bits[mid - 1]) {
      mid++
    }

    const preV = bits[mid]
    const newV = preV ^ 1
    L.update(mid, mid + 1, newV)
    bits[mid] = newV
    console.log(mid, preV, newV, L.toString())
  }

  const { start, end } = L.queryAll()
  return end - start
}

export {}

// "00000"
// 2

console.log(minLength('00000', 2))

// 2 1 1 2 1 1 1
