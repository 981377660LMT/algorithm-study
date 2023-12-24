/* eslint-disable no-inner-declarations */
/* eslint-disable prefer-destructuring */

/**
 * `O(nlogn)` 静态查询每个子树内的信息, 空间复杂度优于启发式合并.
 * @see {@link https://blog.csdn.net/qq_43472263/article/details/104150940}
 */
class DsuOnTree {
  private readonly _tree: number[][]
  private readonly _root: number
  private readonly _subSize: Uint32Array
  private readonly _euler: Uint32Array
  private readonly _down: Uint32Array
  private readonly _up: Uint32Array
  private _order = 0

  constructor(n: number, tree: number[][], root = 0) {
    this._tree = tree
    this._root = root
    this._subSize = new Uint32Array(n)
    this._euler = new Uint32Array(n)
    this._down = new Uint32Array(n)
    this._up = new Uint32Array(n)

    this._dfs1(root, -1)
    this._dfs2(root, -1)
  }

  /**
   * @param add 添加root处的贡献.
   * @param remove 移除root处的贡献.
   * @param query 查询root的子树的贡献并更新答案.
   * @param reset 退出轻儿子时的回调函数.
   */
  run(
    add: (root: number) => void,
    remove: (root: number) => void,
    query: (root: number) => void,
    reset?: () => void
  ) {
    const dsu = (cur: number, pre: number, keep: boolean): void => {
      const nexts = this._tree[cur]
      for (let i = 1; i < nexts.length; i++) {
        const next = nexts[i]
        if (next !== pre) {
          dsu(next, cur, false)
        }
      }

      if (this._subSize[cur] !== 1) {
        dsu(nexts[0], cur, true)
      }

      if (this._subSize[cur] !== 1) {
        for (let i = this._up[nexts[0]]; i < this._up[cur]; i++) {
          add(this._euler[i])
        }
      }

      add(cur)
      query(cur)
      if (!keep) {
        for (let i = this._down[cur]; i < this._up[cur]; i++) {
          remove(this._euler[i])
        }
        reset && reset()
      }
    }

    dsu(this._root, -1, false)
  }

  private _dfs1(cur: number, pre: number): number {
    this._subSize[cur] = 1
    const nexts = this._tree[cur]
    if (nexts.length >= 2 && nexts[0] === pre) {
      nexts[0] ^= nexts[1]
      nexts[1] ^= nexts[0]
      nexts[0] ^= nexts[1]
    }
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (next === pre) continue
      this._subSize[cur] += this._dfs1(next, cur)
      if (this._subSize[next] > this._subSize[nexts[0]]) {
        nexts[0] ^= nexts[i]
        nexts[i] ^= nexts[0]
        nexts[0] ^= nexts[i]
      }
    }
    return this._subSize[cur]
  }

  private _dfs2(cur: number, pre: number): void {
    this._euler[this._order] = cur
    this._down[cur] = this._order++
    const nexts = this._tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (next === pre) continue
      this._dfs2(next, cur)
    }
    this._up[cur] = this._order
  }
}

// 给你一棵 n 个节点的 无向 树，节点编号为 0 到 n - 1 ，树的根节点在节点 0 处。同时给你一个长度为 n - 1 的二维整数数组 edges ，其中 edges[i] = [ai, bi] 表示树中节点 ai 和 bi 之间有一条边。

// 给你一个长度为 n 下标从 0 开始的整数数组 cost ，其中 cost[i] 是第 i 个节点的 开销 。

// 你需要在树中每个节点都放置金币，在节点 i 处的金币数目计算方法如下：

// 如果节点 i 对应的子树中的节点数目小于 3 ，那么放 1 个金币。
// 否则，计算节点 i 对应的子树内 3 个不同节点的开销乘积的 最大值 ，并在节点 i 处放置对应数目的金币。如果最大乘积是 负数 ，那么放置 0 个金币。
// 请你返回一个长度为 n 的数组 coin ，coin[i]是节点 i 处的金币数目。
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

interface ISortedListIterator<V> {
  hasNext(): boolean
  next(): V | undefined
  hasPrev(): boolean
  prev(): V | undefined
  /** 删除后会使所有迭代器失效. */
  remove(): void
  readonly value: V | undefined
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

function placedCoins(edges: number[][], cost: number[]): number[] {
  const n = cost.length
  const tree: number[][] = Array(n)
  for (let i = 0; i < n; i++) tree[i] = []
  for (const [u, v] of edges) {
    tree[u].push(v)
    tree[v].push(u)
  }

  const dsu = new DsuOnTree(n, tree)
  const res: number[] = Array(n).fill(0)
  const sl = new SortedListFast<number>()

  dsu.run(add, remove, query)
  return res

  function add(root: number): void {
    sl.add(cost[root])
  }
  function remove(root: number): void {
    sl.discard(cost[root])
  }
  function query(root: number): void {
    if (sl.length < 3) {
      res[root] = 1
    } else {
      const cand1 = sl.at(sl.length - 1)! * sl.at(sl.length - 2)! * sl.at(sl.length - 3)!
      const cand2 = sl.at(0)! * sl.at(1)! * sl.at(sl.length - 1)!
      res[root] = Math.max(cand1, cand2, 0)
    }
  }
}
