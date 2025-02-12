/* eslint-disable generator-star-spacing */
/* eslint-disable no-empty */

// https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/data_structure/sqrt_array.cc
//
// 优化点：
// 1. 删除后，合并相邻的小块，以减少块的数量。

/**
 * 分块数组.
 */
class SqrtArray<T = number> {
  private _n = 0
  private readonly _blocks: T[][] = []
  private readonly _blockSize: number
  private readonly _threshold: number
  private _shouldRebuildTree = true
  private _tree: number[] = []

  constructor(n: number, f: (i: number) => T, blockSize = 1 << 7) {
    const blockCount = ((n + blockSize - 1) / blockSize) | 0
    const blocks: T[][] = Array(blockCount)
    for (let i = 0; i < blockCount; i++) {
      const start = i * blockSize
      const end = Math.min((i + 1) * blockSize, n)
      const cur = Array(end - start)
      for (let j = start; j < end; j++) cur[j - start] = f(j)
      blocks[i] = cur
    }

    this._n = n
    this._blocks = blocks
    this._blockSize = blockSize
    this._threshold = blockSize << 1
  }

  /**
   * 0<= i < {@link length}
   */
  set(i: number, v: T): void {
    const { bid, pos } = this._findKth(i)
    this._blocks[bid][pos] = v
  }

  /**
   * 0<= i < {@link length}
   */
  get(i: number): T | undefined {
    const { bid, pos } = this._findKth(i)
    return this._blocks[bid][pos]
  }

  /**
   * i可以是负数索引.
   */
  at(i: number): T | undefined {
    if (i < 0) i += this._n
    if (i < 0 || i >= this._n) return undefined
    return this.get(i)
  }

  push(v: T): void {
    this.insert(this._n, v)
  }

  /**
   * i可以是负数索引.
   */
  pop(i = this._n - 1): T | undefined {
    if (i < 0) i += this._n
    if (i < 0 || i >= this._n) return undefined
    const { bid, pos } = this._findKth(i)
    const res = this._blocks[bid].splice(pos, 1)[0]
    this._n--
    this._updateTree(bid, -1)
    // TODO: 过稀疏时，合并相邻块
    if (!this._blocks[bid].length) {
      this._blocks.splice(bid, 1)
      this._shouldRebuildTree = true
    }
    return res
  }

  shift(): T | undefined {
    return this.pop(0)
  }

  unshift(v: T): void {
    this.insert(0, v)
  }

  clear(): void {
    this._n = 0
    this._blocks.length = 0
    this._shouldRebuildTree = true
    this._tree.length = 0
  }

  /**
   * 删除区间 [start, end) 内的元素.
   * 0<= start <= end <= {@link length}
   */
  erase(start: number, end: number): void {
    this.enumerate(start, end, undefined, true)
  }

  /**
   * 在 i 位置`前`插入 v.
   * 0<= i <= {@link length}
   */
  insert(i: number, v: T): void {
    if (!this._n) {
      this._blocks.push([v])
      this._shouldRebuildTree = true
      this._n++
      return
    }

    if (i < 0) i += this._n
    if (i < 0) i = 0
    if (i > this._n) i = this._n

    const { bid, pos } = this._findKth(i)
    this._updateTree(bid, 1)
    this._blocks[bid].splice(pos, 0, v)

    // 定期重构
    if (this._blocks[bid].length > this._threshold) {
      const right = this._blocks[bid].splice(this._blockSize)
      this._blocks.splice(bid + 1, 0, right)
      this._shouldRebuildTree = true
    }

    this._n++
  }

  /**
   * 遍历区间 [start, end) 内的元素,并选择是否在遍历后删除.
   * 0<= start <= end <= {@link length}
   */
  enumerate(start: number, end: number, f: ((value: T) => void) | undefined, erase: boolean): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return

    let { bid, pos } = this._findKth(start)
    let count = end - start

    let eraseStart = -1
    let eraseCount = 0
    for (; bid < this._blocks.length && count > 0; bid++) {
      const block = this._blocks[bid]
      const endPos = Math.min(block.length, pos + count)
      if (f) {
        for (let j = pos; j < endPos; j++) {
          f(block[j])
        }
      }

      const curDeleteCount = endPos - pos
      if (erase) {
        if (curDeleteCount === block.length) {
          if (eraseStart === -1) eraseStart = bid
          eraseCount++
        } else {
          this._updateTree(bid, -curDeleteCount)
          block.splice(pos, curDeleteCount)
        }
        this._n -= curDeleteCount
      }

      count -= curDeleteCount
      pos = 0
    }

    if (erase && eraseStart !== -1) {
      this._blocks.splice(eraseStart, eraseCount)
      this._shouldRebuildTree = true
    }
  }

  slice(start: number, end: number): T[] {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return []
    let count = end - start
    const res: T[] = Array(count)
    let { bid, pos } = this._findKth(start)
    let ptr = 0
    for (; bid < this._blocks.length && count > 0; bid++) {
      const block = this._blocks[bid]
      const endPos = Math.min(block.length, pos + count)
      const curCount = endPos - pos
      for (let j = pos; j < endPos; j++) res[ptr++] = block[j]
      count -= curCount
      pos = 0
    }
    return res
  }

  fill(v: T): this {
    this._blocks.forEach(b => b.fill(v))
    return this
  }

  /**
   * 返回一个迭代器，用于遍历区间 [start, end) 内的元素.
   */
  *islice(start: number, end: number, reverse = false): IterableIterator<T> {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    let count = end - start

    if (reverse) {
      let { bid, pos } = this._findKth(end - 1)
      for (; ~bid && count > 0; bid--, ~bid && (pos = this._blocks[bid].length)) {
        const block = this._blocks[bid]
        const startPos = Math.max(0, pos - count)
        const curCount = pos - startPos
        for (let j = pos - 1; j >= startPos; j--) {
          yield block[j]
        }
        count -= curCount
      }
    } else {
      let { bid, pos } = this._findKth(start)
      for (; bid < this._blocks.length && count > 0; bid++) {
        const block = this._blocks[bid]
        const endPos = Math.min(block.length, pos + count)
        const curCount = endPos - pos
        for (let j = pos; j < endPos; j++) {
          yield block[j]
        }
        count -= curCount
        pos = 0
      }
    }
  }

  forEach(callback: (value: T, index: number) => void): void {
    let ptr = 0
    for (let bi = 0; bi < this._blocks.length; ++bi) {
      for (let j = 0; j < this._blocks[bi].length; ++j) {
        callback(this._blocks[bi][j], ptr++)
      }
    }
  }

  toString(): string {
    return `SqrtArray{${this._blocks}}`
  }

  private _findKth(index: number): { bid: number; pos: number } {
    if (index < this._blocks[0].length) return { bid: 0, pos: index }
    const last = this._blocks.length - 1
    const lastLen = this._blocks[last].length
    if (index >= this._n) return { bid: last, pos: lastLen }
    if (index >= this._n - lastLen) return { bid: last, pos: index + lastLen - this._n }
    if (this._shouldRebuildTree) this._rebuildTree()
    let pos = -1
    const tree = this._tree
    const m = tree.length
    const bitLen = 32 - Math.clz32(m)
    for (let d = bitLen - 1; ~d; d--) {
      const next = pos + (1 << d)
      if (next < m && index >= tree[next]) {
        pos = next
        index -= tree[pos]
      }
    }
    return { bid: pos + 1, pos: index }
  }

  private _updateTree(bid: number, delta: number): void {
    if (this._shouldRebuildTree) return
    for (let i = bid; i < this._tree.length; i |= i + 1) {
      this._tree[i] += delta
    }
  }

  private _rebuildTree(): void {
    const n = this._blocks.length
    const tree = Array(n)
    for (let i = 0; i < n; i++) tree[i] = this._blocks[i].length
    for (let i = 0; i < n; i++) {
      const j = i | (i + 1)
      if (j < n) tree[j] += tree[i]
    }

    this._tree = tree
    this._shouldRebuildTree = false
  }

  get length(): number {
    return this._n
  }
}

export { SqrtArray }

if (require.main === module) {
  const arr = new SqrtArray<number>(0, i => i, 128)
  const n = 1e6

  const time1 = Date.now()
  for (let i = 0; i < n; i++) {
    arr.insert(n - i, i)
    arr.get(i)
    arr.set(i, i)
    arr.unshift(i)
    arr.shift()
  }

  for (let i = 0; i < n; i++) {
    arr.get(i)
  }
  for (let i = 0; i < n; i++) {
    arr.pop(0)
  }

  const cur = Date.now()
  console.log(cur - time1)

  // https://leetcode.cn/problems/design-circular-deque/
  class MyCircularDeque {
    private readonly _sqrt: SqrtArray<number>
    private readonly _k: number
    constructor(k: number) {
      this._sqrt = new SqrtArray(0, () => 0)
      this._k = k
    }

    insertFront(value: number): boolean {
      if (this.isFull()) return false
      this._sqrt.unshift(value)
      return true
    }

    insertLast(value: number): boolean {
      if (this.isFull()) return false
      this._sqrt.push(value)
      return true
    }

    deleteFront(): boolean {
      if (this.isEmpty()) return false
      this._sqrt.shift()
      return true
    }

    deleteLast(): boolean {
      if (this.isEmpty()) return false
      this._sqrt.pop()
      return true
    }

    getFront(): number {
      return this.isEmpty() ? -1 : this._sqrt.get(0)!
    }

    getRear(): number {
      return this.isEmpty() ? -1 : this._sqrt.get(this._sqrt.length - 1)!
    }

    isEmpty(): boolean {
      return this._sqrt.length === 0
    }

    isFull(): boolean {
      return this._sqrt.length === this._k
    }
  }
}
