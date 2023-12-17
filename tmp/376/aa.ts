function maxFrequencyScore(nums: number[], k: number): number {
  nums.sort((a, b) => a - b)

  const sl = new SortedListFastWithSum()
  const D = distSum(sl)
  let res = 0
  let left = 0
  for (let right = 0; right < nums.length; right++) {
    sl.add(nums[right])
    while (left <= right) {
      const median =
        sl.length % 2 === 0
          ? Math.floor((sl.at(sl.length / 2 - 1)! + sl.at(sl.length / 2)!) / 2)
          : sl.at((sl.length - 1) / 2)!

      const distSum = D(median)
      if (distSum <= k) break
      sl.discard(nums[left])
      left++
    }
    res = Math.max(res, right - left + 1)
  }

  return res

  function distSum(sl: SortedListFastWithSum): (k: number) => number {
    return (k: number): number => {
      const pos = sl.bisectRight(k)
      const leftSum = k * pos - sl.sumSlice(0, pos)
      const rightSum = sl.sumSlice(pos, sl.length) - k * (sl.length - pos)
      return leftSum + rightSum
    }
  }
}

interface Options<V> {
  values?: Iterable<V>
  compareFn?: (a: V, b: V) => number
  abelGroup?: {
    e: () => V
    op: (a: V, b: V) => V
    inv: (a: V) => V
  }
}

/**
 * 使用分块+树状数组维护的有序序列.
 * !如果数组较短(<2000),直接使用`bisectInsort`维护即可.
 */
class SortedListFast<V = number> {
  static setLoadFactor(load = 500): void {
    this._LOAD = load
  }

  /**
   * 负载因子，用于控制每个块的长度.
   * 长度为`1e5`的数组, 负载因子取`500`左右性能较好.
   */
  protected static _LOAD = 500
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

    this._build(defaultData, defaultCompareFn)
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
      _blocks.splice(pos + 1, 0, block.slice(load))
      _mins.splice(pos + 1, 0, block[load])
      block.splice(load)
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
    this._build(data, this._compareFn)
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
    this._build(data, this._compareFn)
    return this
  }

  has(value: V, equals?: (a: V, b: V) => boolean): boolean {
    if (!equals && !SortedListFast._isPrimitive(value)) {
      throw new Error('equals must be provided when value is a non-primitive value')
    }
    const block = this._blocks
    if (!block.length) return false
    equals = equals || ((a: V, b: V) => a === b)
    const { pos, index } = this._locLeft(value)
    return index < block[pos].length && equals(block[pos][index], value)
  }

  discard(value: V, equals?: (a: V, b: V) => boolean): boolean {
    if (!equals && !SortedListFast._isPrimitive(value)) {
      throw new Error('equals must be provided when value is a non-primitive value')
    }
    const block = this._blocks
    if (!block.length) return false
    equals = equals || ((a: V, b: V) => a === b)
    const { pos, index } = this._locRight(value)
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

    let { pos, index: startIndex } = this._findKth(start)
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
          for (let i = startIndex; i < endIndex; i++) this._updateTree(pos, -1)
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

  protected _build(data: V[], compareFn: (a: V, b: V) => number): void {
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

  protected _locBlock(value: V): number {
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
    return { pos: pos + 1, index: k }
  }
}

/**
 * 支持区间求和的有序列表.
 * {@link sumSlice} 和 {@link sumRange} 的时间复杂度为 `O(sqrt(n))`.
 */
class SortedListFastWithSum<V = number> extends SortedListFast<V> {
  private readonly _e: () => V
  private readonly _op: (a: V, b: V) => V
  private readonly _inv: (a: V) => V
  private _sums: V[] = []

  constructor(options?: Options<V>) {
    super()

    const {
      values = [],
      compareFn = (a: any, b: any) => a - b,
      abelGroup = {
        e: () => 0 as any,
        op: (a: any, b: any) => a + b,
        inv: (a: any) => -a as any
      }
    } = options ?? {}

    this._e = abelGroup.e
    this._op = abelGroup.op
    this._inv = abelGroup.inv
    this._build([...values], compareFn)
  }

  /**
   * 返回区间 `[start, end)` 的和.
   */
  sumSlice(start = 0, end = this.length): V {
    if (start < 0) start += this._len
    if (start < 0) start = 0
    if (end < 0) end += this._len
    if (end > this._len) end = this._len
    if (start >= end) return this._e()

    let res = this._e()
    let { pos, index } = this._findKth(start)
    let count = end - start
    for (; count && pos < this._blocks.length; pos++) {
      const block = this._blocks[pos]
      const endPos = Math.min(block.length, index + count)
      const curCount = endPos - index
      if (curCount === block.length) {
        res = this._op(res, this._sums[pos])
      } else {
        for (let j = index; j < endPos; j++) res = this._op(res, block[j])
      }
      count -= curCount
      index = 0
    }

    return res
  }

  /**
   * 返回范围 `[min, max]` 的和.
   */
  sumRange(min: V, max: V): V {
    if (this._compareFn(min, max) > 0) return this._e()

    let res = this._e()
    let { pos, index: start } = this._locLeft(min)
    for (let i = pos; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (this._compareFn(max, block[0]) < 0) break
      if (start === 0 && this._compareFn(block[block.length - 1], max) <= 0) {
        res = this._op(res, this._sums[i])
      } else {
        for (let j = start; j < block.length; j++) {
          const cur = block[j]
          if (this._compareFn(cur, max) > 0) break
          res = this._op(res, cur)
        }
      }
      start = 0
    }
    return res
  }

  override add(value: V): this {
    const { _blocks, _mins, _sums } = this
    this._len++
    if (!_blocks.length) {
      _blocks.push([value])
      _mins.push(value)
      _sums.push(value)
      this._shouldRebuildTree = true
      return this
    }

    const load = SortedListFast._LOAD
    const { pos, index } = this._locRight(value)
    this._updateTree(pos, 1)
    const block = _blocks[pos]
    block.splice(index, 0, value)
    _mins[pos] = block[0]
    _sums[pos] = this._op(_sums[pos], value)

    // !-> [x]*load + [x]*(block.length - load)
    if (load + load < block.length) {
      const oldSum = _sums[pos]
      _blocks.splice(pos + 1, 0, block.slice(load))
      _mins.splice(pos + 1, 0, block[load])
      block.splice(load)
      this._shouldRebuildTree = true

      this._rebuildSum(pos)
      this._sums.splice(pos + 1, 0, this._op(oldSum, this._inv(this._sums[pos])))
    }

    return this
  }

  override enumerate(start: number, end: number, f?: (value: V) => void, erase?: boolean): void {
    if (start < 0) start = 0
    if (end > this._len) end = this._len
    if (start >= end) return

    let { pos, index: startIndex } = this._findKth(start)
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
          this._sums.splice(pos, 1)
          this._shouldRebuildTree = true
          pos--
        } else {
          // !delete [index, end)
          for (let i = startIndex; i < endIndex; i++) {
            this._updateTree(pos, -1)
            this._sums[pos] = this._op(this._sums[pos], this._inv(block[i]))
          }
          block.splice(startIndex, deleted)
          this._mins[pos] = block[0]
        }
        this._len -= deleted
      }

      count -= deleted
      startIndex = 0
    }
  }

  override clear(): void {
    super.clear()
    this._sums = []
  }

  protected override _build(data: V[], compareFn: (a: V, b: V) => number): void {
    data.sort(compareFn)
    const n = data.length
    const blocks = []
    const sums = []
    for (let i = 0; i < n; i += SortedListFast._LOAD) {
      const newBlock = data.slice(i, i + SortedListFast._LOAD)
      blocks.push(newBlock)
      let cur = this._e()
      for (let j = 0; j < newBlock.length; j++) cur = this._op(cur, newBlock[j])
      sums.push(cur)
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
    this._sums = sums
  }

  protected override _delete(pos: number, index: number): void {
    const { _blocks, _mins, _sums } = this

    // !delete element
    this._len--
    this._updateTree(pos, -1)
    const block = _blocks[pos]
    const deleted = block[index]
    block.splice(index, 1)
    if (block.length) {
      _mins[pos] = block[0]
      _sums[pos] = this._op(_sums[pos], this._inv(deleted))
      return
    }

    // !delete block
    _blocks.splice(pos, 1)
    _mins.splice(pos, 1)
    _sums.splice(pos, 1)
    this._shouldRebuildTree = true
  }

  private _rebuildSum(pos: number): void {
    let cur = this._e()
    const block = this._blocks[pos]
    for (let i = 0; i < block.length; i++) cur = this._op(cur, block[i])
    this._sums[pos] = cur
  }
}
