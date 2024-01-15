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
    return { pos: pos + 1, index: k }
  }
}

export { SortedListFast, ISortedList }

if (require.main === module) {
  const sl = new SortedListFast<number>()
  sl.add(2)
  sl.add(2)
  sl.add(1)
  sl.add(3)
  console.log(sl.toString())
  console.log(sl.floor(2), sl.lower(2), sl.ceiling(2), sl.higher(2))
  console.log(...sl.iRange(0, 2))

  // https://leetcode.cn/problems/sliding-subarray-beauty/
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

  // https://leetcode.cn/problems/count-nodes-that-are-great-enough/
  class TreeNode {
    val: number
    left: TreeNode | null
    right: TreeNode | null
    constructor(val?: number, left?: TreeNode | null, right?: TreeNode | null) {
      this.val = val === undefined ? 0 : val
      this.left = left === undefined ? null : left
      this.right = right === undefined ? null : right
    }
  }

  function countGreatEnoughNodes(root: TreeNode | null, k: number): number {
    let res = 0
    dfs(root)
    return res

    function dfs(cur: TreeNode | null): SortedListFast<number> {
      if (!cur) {
        return new SortedListFast()
      }
      let left = dfs(cur.left)
      let right = dfs(cur.right)
      if (left.length < right.length) {
        const tmp = left
        left = right
        right = tmp
      }
      left.merge(right)
      res += +(left.bisectLeft(cur.val) >= k)
      left.add(cur.val)
      return left
    }
  }

  // https://leetcode.cn/problems/minimum-absolute-difference-between-elements-with-constraint/description/
  function minAbsoluteDifference(nums: number[], x: number): number {
    const sl = new SortedListFast<number>()
    let res = 2e15
    for (let right = x; right < nums.length; right++) {
      sl.add(nums[right - x])
      const cur = nums[right]
      const floor = sl.floor(cur)
      if (floor != undefined) {
        res = Math.min(res, cur - floor)
      }
      const ceiling = sl.ceiling(cur)
      if (ceiling != undefined) {
        res = Math.min(res, ceiling - cur)
      }
    }
    return res
  }

  // 1649. 通过指令创建有序数组
  // https://leetcode.cn/problems/create-sorted-array-through-instructions/description/

  function createSortedArray(instructions: number[]): number {
    const MOD = 1e9 + 7
    let res = 0
    const sl = new SortedListFast<number>()
    instructions.forEach(num => {
      const smaller = sl.bisectLeft(num)
      const bigger = sl.length - sl.bisectRight(num)
      res += Math.min(smaller, bigger)
      res %= MOD
      sl.add(num)
    })
    return res % MOD
  }
}
