/* eslint-disable arrow-body-style */
/* eslint-disable no-inner-declarations */
/* eslint-disable generator-star-spacing */
/* eslint-disable prefer-destructuring */

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

  toString(): string

  forEach(callbackfn: (value: V, index: number) => void | boolean, reverse?: boolean): void
  enumerate(start: number, end: number, f: (value: V) => void, erase?: boolean): void
  islice(start: number, end: number, reverse?: boolean): IterableIterator<V>
  irange(min: V, max: V, reverse?: boolean): IterableIterator<V>
  entries(): IterableIterator<[number, V]>
  [Symbol.iterator](): IterableIterator<V>

  iteratorAt(index: number): ISortedListIterator<V>
  lowerBound(value: V): ISortedListIterator<V>
  upperBound(value: V): ISortedListIterator<V>

  get length(): number
  get min(): V | undefined
  get max(): V | undefined
}

interface ISortedListIterator<V> {
  hasNext(): boolean
  next(): V | undefined
  hasPrev(): boolean
  prev(): V | undefined
  remove(): void
  get value(): V | undefined
}

/**
 * 使用分块+树状数组维护的有序序列.
 * !如果数组较短(<2000),直接使用`bisectInsort`维护即可.
 */
class SortedListFast<V = number> {
  static setLoadFactor(load = 1 << 9): void {
    this._LOAD = load
  }

  /**
   * 负载因子，用于控制每个块的长度.
   * 长度为`1e5`的数组, 负载因子取`500`左右性能较好.
   */
  private static _LOAD = 1 << 9
  private static _isPrimitive(
    o: unknown
  ): o is number | string | boolean | symbol | bigint | null | undefined {
    return o === null || (typeof o !== 'object' && typeof o !== 'function')
  }

  private readonly _compareFn: (a: V, b: V) => number
  private _len: number

  /**
   * 各个块.
   */
  private _blocks: V[][]

  /**
   * 各个块的长度.
   */
  private _blockLens: number[]

  /**
   * 各个块的最小值.
   */
  private _mins: V[]

  /**
   * 树状数组维护 {@link _blockLens} 的前缀和，便于根据索引定位到对应的元素.
   */
  private _tree: number[] = []
  private _shouldRebuild = true

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

    defaultData.sort(defaultCompareFn)
    const n = defaultData.length
    const lists = []
    for (let i = 0; i < n; i += SortedListFast._LOAD) {
      lists.push(defaultData.slice(i, i + SortedListFast._LOAD))
    }
    const listLens = Array(lists.length)
    const mins = Array(lists.length)
    for (let i = 0; i < lists.length; i++) {
      const cur = lists[i]
      listLens[i] = cur.length
      mins[i] = cur[0]
    }

    this._compareFn = defaultCompareFn
    this._len = n
    this._blocks = lists
    this._blockLens = listLens
    this._mins = mins
  }

  add(value: V): this {
    const { _blocks, _mins, _blockLens } = this
    this._len++
    if (!_blocks.length) {
      _blocks.push([value])
      _mins.push(value)
      _blockLens.push(1)
      this._shouldRebuild = true
      return this
    }

    const load = SortedListFast._LOAD
    const pair = this._locRight(value)
    const pos = pair[0]
    const index = pair[1]
    this._update(pos, 1)
    const list = _blocks[pos]
    list.splice(index, 0, value)
    _blockLens[pos]++
    _mins[pos] = list[0]
    if (load + load < list.length) {
      _blocks.splice(pos + 1, 0, list.slice(load))
      _blockLens.splice(pos + 1, 0, list.length - load)
      _mins.splice(pos + 1, 0, list[load])
      _blockLens[pos] = load
      list.splice(load)
      this._shouldRebuild = true
    }

    return this
  }

  has(value: V, equals?: (a: V, b: V) => boolean): boolean {
    if (!equals && !SortedListFast._isPrimitive(value)) {
      throw new Error('equals must be provided when value is a non-primitive value')
    }
    const list = this._blocks
    if (!list.length) return false
    equals = equals || ((a: V, b: V) => a === b)
    const pair = this._locLeft(value)
    const pos = pair[0]
    const index = pair[1]
    return index < list[pos].length && equals(list[pos][index], value)
  }

  discard(value: V, equals?: (a: V, b: V) => boolean): boolean {
    if (!equals && !SortedListFast._isPrimitive(value)) {
      throw new Error('equals must be provided when value is a non-primitive value')
    }
    const list = this._blocks
    if (!list.length) return false
    equals = equals || ((a: V, b: V) => a === b)
    const pair = this._locRight(value)
    const pos = pair[0]
    const index = pair[1]
    if (index && equals(list[pos][index - 1], value)) {
      this._delete(pos, index - 1)
      return true
    }
    return false
  }

  pop(index = -1): V | undefined {
    if (index < 0) index += this._len
    if (index < 0 || index >= this._len) return void 0
    const pair = this._findKth(index)
    const pos = pair[0]
    const index_ = pair[1]
    const value = this._blocks[pos][index_]
    this._delete(pos, index_)
    return value
  }

  at(index: number): V | undefined {
    if (index < 0) index += this._len
    if (index < 0 || index >= this._len) return void 0
    const pair = this._findKth(index)
    const pos = pair[0]
    const index_ = pair[1]
    return this._blocks[pos][index_]
  }

  /**
   * 删除区间 `[start, end)` 内的元素.
   */
  erase(start = 0, end = this.length): void {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerate(start, end, () => {}, true)
  }

  lower(value: V): V | undefined {
    const pos = this.bisectLeft(value)
    return pos ? this.at(pos - 1) : void 0
  }

  higher(value: V): V | undefined {
    const pos = this.bisectRight(value)
    return pos < this._len ? this.at(pos) : void 0
  }

  floor(value: V): V | undefined {
    const pos = this.bisectRight(value)
    return pos ? this.at(pos - 1) : void 0
  }

  ceiling(value: V): V | undefined {
    const pos = this.bisectLeft(value)
    return pos < this._len ? this.at(pos) : void 0
  }

  /**
   * 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
   */
  bisectLeft(value: V): number {
    const pair = this._locLeft(value)
    const pos = pair[0]
    const index = pair[1]
    return this._query(pos) + index
  }

  /**
   * 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
   */
  bisectRight(value: V): number {
    const pair = this._locRight(value)
    const pos = pair[0]
    const index = pair[1]
    return this._query(pos) + index
  }

  count(value: V): number {
    return this.bisectRight(value) - this.bisectLeft(value)
  }

  slice(start = 0, end = this.length): V[] {
    if (start < 0) start = 0
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
    this._blockLens = []
    this._mins = []
    this._tree = []
    this._shouldRebuild = true
  }

  toString(): string {
    const res = Array(this._len)
    this.forEach((value, index) => {
      res[index] = value
    })
    return `SortedList{${res.join(',')}}`
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

  enumerate(start: number, end: number, f: (value: V) => void, erase = false): void {
    if (start < 0) start = 0
    if (end > this._len) end = this._len
    if (start >= end) return

    const pair = this._findKth(start)
    let pos = pair[0]
    let startIndex = pair[1]
    let count = end - start
    for (; count && pos < this._blocks.length; pos++) {
      const block = this._blocks[pos]
      const endIndex = Math.min(block.length, startIndex + count)
      for (let j = startIndex; j < endIndex; j++) f(block[j])
      const deleted = endIndex - startIndex

      if (erase) {
        if (deleted === block.length) {
          // !delete block
          this._blocks.splice(pos, 1)
          this._blockLens.splice(pos, 1)
          this._mins.splice(pos, 1)
          this._shouldRebuild = true
          pos--
        } else {
          // !delete [index, end)
          for (let i = startIndex; i < endIndex; i++) this._update(pos, -1)
          block.splice(startIndex, deleted)
          this._blockLens[pos] -= deleted
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
  *islice(start: number, end: number, reverse = false): IterableIterator<V> {
    if (start < 0) start = 0
    if (end > this._len) end = this._len
    if (start >= end) return
    let count = end - start

    if (reverse) {
      let [pos, index] = this._findKth(end)
      for (; count && ~pos; pos--, ~pos && (index = this._blocks[pos].length)) {
        const block = this._blocks[pos]
        const startPos = Math.max(0, index - count)
        const curCount = index - startPos
        for (let j = index - 1; j >= startPos; j--) yield block[j]
        count -= curCount
      }
    } else {
      let [pos, index] = this._findKth(start)
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
  *irange(min: V, max: V, reverse = false): IterableIterator<V> {
    if (this._compareFn(min, max) > 0) return

    if (reverse) {
      let pos = this._locBlock(max)
      if (pos === -1) pos = this._blocks.length - 1
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

  iteratorAt(index: number): ISortedListIterator<V> {
    if (index < 0) index += this._len
    if (index < 0 || index >= this._len) throw new RangeError('Index out of range')
    const pair = this._findKth(index)
    return this._iteratorAt(pair[0], pair[1])
  }

  lowerBound(value: V): ISortedListIterator<V> {
    const [pos, index] = this._locLeft(value)
    return this._iteratorAt(pos, index)
  }

  upperBound(value: V): ISortedListIterator<V> {
    const [pos, index] = this._locRight(value)
    return this._iteratorAt(pos, index)
  }

  get length(): number {
    return this._len
  }

  get min(): V | undefined {
    if (!this._len) return void 0
    return this._mins[0]
  }

  get max(): V | undefined {
    if (!this._len) return void 0
    const { _blocks, _blockLens } = this
    const last = _blocks.length - 1
    return _blocks[last][_blockLens[last] - 1]
  }

  private _delete(pos: number, index: number): void {
    const { _blocks, _mins, _blockLens } = this

    // !delete element
    this._len--
    this._update(pos, -1)
    _blocks[pos].splice(index, 1)
    _blockLens[pos]--
    if (_blockLens[pos]) {
      _mins[pos] = _blocks[pos][0]
      return
    }

    // !delete block
    _blocks.splice(pos, 1)
    _blockLens.splice(pos, 1)
    _mins.splice(pos, 1)
    this._shouldRebuild = true
  }

  private _locLeft(value: V): [pos: number, index: number] {
    if (!this._len) return [0, 0]
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

    return [pos, right]
  }

  private _locRight(value: V): [pos: number, index: number] {
    if (!this._len) return [0, 0]
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

    return [pos, right]
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

  private _build(): void {
    this._tree = this._blockLens.slice()
    const tree = this._tree
    for (let i = 0; i < tree.length; i++) {
      const j = i | (i + 1)
      if (j < tree.length) {
        tree[j] += tree[i]
      }
    }
    this._shouldRebuild = false
  }

  private _update(index: number, delta: number): void {
    if (!this._shouldRebuild) {
      const tree = this._tree
      while (index < tree.length) {
        tree[index] += delta
        index |= index + 1
      }
    }
  }

  private _query(end: number): number {
    if (this._shouldRebuild) this._build()
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
  private _findKth(k: number): [pos: number, index: number] {
    const listLens = this._blockLens
    if (k < listLens[0]) return [0, k]
    if (k >= this._len - listLens[listLens.length - 1]) {
      return [listLens.length - 1, k + listLens[listLens.length - 1] - this._len]
    }
    if (this._shouldRebuild) this._build()
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
    return [pos + 1, k]
  }

  private _iteratorAt(pos: number, index: number): ISortedListIterator<V> {
    const hasNext = (): boolean => {
      return pos < this._blocks.length - 1 || index < this._blocks[pos].length - 1
    }

    const next = (): V | undefined => {
      if (!hasNext()) return void 0
      index++
      if (index === this._blocks[pos].length) {
        pos++
        index = 0
      }
      return this._blocks[pos][index]
    }

    const hasPrev = (): boolean => {
      return pos > 0 || index > 0
    }

    const prev = (): V | undefined => {
      if (!hasPrev()) return void 0
      index--
      if (index === -1) {
        pos--
        index = this._blocks[pos].length - 1
      }
      return this._blocks[pos][index]
    }

    const remove = (): void => {
      this._delete(pos, index)
    }

    // eslint-disable-next-line @typescript-eslint/no-this-alias
    const sl = this
    return {
      hasNext,
      next,
      hasPrev,
      prev,
      remove,
      get value(): V | undefined {
        if (pos < 0 || pos >= sl.length) return void 0
        const block = sl._blocks[pos]
        if (index < 0 || index >= block.length) return void 0
        return block[index]
      }
    }
  }
}

export { SortedListFast }

if (require.main === module) {
  const sl = new SortedListFast<number>()
  sl.add(2)
  sl.add(2)
  sl.add(1)
  sl.add(3)
  console.log(sl.toString())
  console.log(sl.floor(2), sl.lower(2), sl.ceiling(2), sl.higher(2))
  console.log(...sl.irange(0, 2))

  const iter = sl.iteratorAt(0)
  console.log(iter.value)
  while (iter.hasNext()) {
    console.log(iter.next())
  }
  console.log(iter.value)
  while (iter.hasPrev()) {
    console.log(iter.prev())
  }
  console.log(iter.value, iter.remove())
  sl.add(1).add(8).add(5).add(3)
  console.log(sl.toString())
  const lower = sl.lowerBound(3)
  console.log(lower.value)
  const upper = sl.upperBound(3)
  console.log(upper.value)
  lower.prev()
  console.log(lower.value)

  // https://leetcode.cn/problems/sliding-subarray-beauty/ 2200ms
  function getSubarrayBeauty(nums: number[], k: number, x: number): number[] {
    const res: number[] = []
    const sl = new SortedListFast()
    const n = nums.length
    for (let right = 0; right < n; right++) {
      sl.add(nums[right])
      if (right >= k) {
        sl.discard(nums[right - k])
      }
      if (right >= k - 1) {
        const xth = sl.at(x - 1)!
        res.push(xth < 0 ? xth : 0)
      }
    }
    return res
  }

  // https://leetcode.cn/problems/sum-of-imbalance-numbers-of-all-subarrays/
  function sumImbalanceNumbers(nums: number[]): number {
    let res = 0
    const n = nums.length
    for (let left = 0; left < n; left++) {
      const sl = new SortedListFast<number>()
      for (let right = left; right < n; right++) {
        sl.add(nums[right])
        const cur = sl.slice()
        for (let i = 1; i < cur.length; i++) {
          res += +(cur[i] - cur[i - 1] > 1)
        }
      }
    }
    return res
  }

  // https://leetcode.cn/problems/maximum-number-of-tasks-you-can-assign/
  function maxTaskAssign(
    tasks: number[],
    workers: number[],
    pills: number,
    strength: number
  ): number {
    tasks.sort((a, b) => a - b)
    workers.sort((a, b) => a - b)
    let left = 0
    let right = Math.min(tasks.length, workers.length)
    while (left <= right) {
      const mid = (left + right) >> 1
      if (check(mid)) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }

    return right

    function check(mid: number): boolean {
      let remain = pills
      const sl = new SortedListFast<number>(workers.slice(-mid))
      // const wls = useSortedList(workers.slice(-mid))
      for (let i = mid - 1; i >= 0; i--) {
        const t = tasks[i]
        if (sl.at(sl.length - 1)! >= t) {
          sl.pop()
        } else {
          if (remain === 0) {
            return false
          }
          const cand = sl.bisectLeft(t - strength)
          if (cand === sl.length) {
            return false
          }
          remain -= 1
          sl.pop(cand)
        }
      }

      return true
    }
  }

  // https://leetcode.cn/problems/count-the-number-of-fair-pairs/
  function countFairPairs(nums: number[], lower: number, upper: number): number {
    const sl = new SortedListFast<number>()
    let res = 0
    nums.forEach(x => {
      res += sl.bisectRight(upper - x) - sl.bisectLeft(lower - x)
      sl.add(x)
    })
    return res
  }

  // https://leetcode.cn/problems/find-score-of-an-array-after-marking-all-elements/
  function findScore(nums: number[]): number {
    const sl = new SortedListFast<[number, number]>((a, b) => a[0] - b[0] || a[1] - b[1])
    const equals = (a: [number, number], b: [number, number]) => a[0] === b[0] && a[1] === b[1]
    nums.forEach((x, i) => sl.add([x, i]))

    let res = 0
    while (sl.length > 0) {
      const [v, i] = sl.pop(0)!
      res += v
      if (i - 1 >= 0) sl.discard([nums[i - 1], i - 1], equals)
      if (i + 1 < nums.length) sl.discard([nums[i + 1], i + 1], equals)
    }
    return res
  }
}
