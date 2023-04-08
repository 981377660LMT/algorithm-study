/* eslint-disable no-console */
/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */

// 如果需要支持set那样的 lowerBound/upperBound/erase 迭代器功能,
// !需要使用分块链表(所有元素之间有前驱后继,便于迭代器移动,且删除非迭代器所在元素后迭代器不会失效)

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
    const res: T[] = Array(count).fill(0)
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
      let [bid, endPos] = this._moveTo(end - 1)
      for (; ~bid && count > 0; bid--, ~bid && (endPos = this._blocks[bid].length)) {
        const block = this._blocks[bid]
        const startPos = Math.max(0, endPos - count)
        const curCount = endPos - startPos
        for (let j = endPos - 1; j >= startPos; j--) {
          yield block[j]
        }
        count -= curCount
        console.log(endPos)
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
    const newB: T[][] = Array(bCount).fill(0)
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
    const newB: T[][] = Array(bCount).fill(0)
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

  get length(): number {
    return this._size
  }
}

export { SortedList }

if (require.main === module) {
  // https://leetcode.cn/problems/count-the-number-of-fair-pairs/
  function countFairPairs(nums: number[], lower: number, upper: number): number {
    const sl = new SortedList<number>()
    let res = 0
    nums.forEach(x => {
      res += sl.bisectRight(upper - x) - sl.bisectLeft(lower - x)
      sl.add(x)
    })
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
      const sl = new SortedList<number>(workers.slice(-mid))
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

  // https://leetcode.cn/problems/find-score-of-an-array-after-marking-all-elements/
  function findScore(nums: number[]): number {
    const objects: [number, number][] = nums.map((v, i) => [v, i])
    const sl = new SortedList<[number, number]>((a, b) => a[0] - b[0] || a[1] - b[1], objects)

    let res = 0
    while (sl.length > 0) {
      const [v, i] = sl.pop(0)!
      res += v
      if (i - 1 >= 0) sl.discard(objects[i - 1])
      if (i + 1 < nums.length) sl.discard(objects[i + 1])
    }
    return res
  }

  const sl = new SortedList<number>([1, 4, 2, 14, 611, 3])
  console.log(sl.length)
  console.log(sl.toString())
  console.log(sl.length)
  console.log(sl.at(0))
  console.log(sl.at(-1))
  console.log(sl.toString(), sl.ceiling(-1), sl.has(4))
  console.log(...sl.islice(0, 3, true))
  console.log(...sl.irange(0, 300, true))
}
