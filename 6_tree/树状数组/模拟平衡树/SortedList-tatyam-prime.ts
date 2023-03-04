// 内部使用分块的SortedList 无须离散化

/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */

/**
 * A fast SortedList with O(sqrt(n)) insertion and deletion.
 * @_see https://github.com/tatyam-prime/SortedSet/blob/main/SortedMultiset.py
 */
class SortedList<T = number> {
  /** Optimized for 1e5 elements in javascript. */
  private static readonly _BLOCK_RATIO = 500
  private static readonly _REBUILD_RATIO = 1700

  private readonly _compareFn: (a: T, b: T) => number
  private _size: number
  private _blocks: T[][]

  constructor()
  constructor(iterable: Iterable<T>)
  constructor(compareFn: (a: T, b: T) => number)
  constructor(iterable: Iterable<T>, compareFn: (a: T, b: T) => number)
  constructor(
    iterableOrCompareFn?: Iterable<T> | ((a: T, b: T) => number),
    compareFn?: (a: T, b: T) => number
  ) {
    let defaultCompareFn = (a: T, b: T) => (a as unknown as number) - (b as unknown as number)
    let defaultData: T[] = []
    if (iterableOrCompareFn !== void 0) {
      if (typeof iterableOrCompareFn === 'function') {
        defaultCompareFn = iterableOrCompareFn
      } else {
        defaultData = [...iterableOrCompareFn]
      }
    }
    if (compareFn !== void 0) {
      defaultCompareFn = compareFn
    }

    this._compareFn = defaultCompareFn
    defaultData.sort(defaultCompareFn)
    this._blocks = this._initBlocks(defaultData)
    this._size = defaultData.length
  }

  add(value: T): void {
    if (this._size === 0) {
      this._blocks = [[value]]
      this._size = 1
      return
    }
    const block = this._findBlock(value)
    const pos = this._bisectRight(block, value)
    block.splice(pos, 0, value)
    this._size += 1
    if (block.length > this._blocks.length * SortedList._REBUILD_RATIO) {
      this._rebuild()
    }
  }

  has(value: T): boolean {
    if (this._size === 0) return false
    const block = this._findBlock(value)
    const pos = this._bisectLeft(block, value)
    return pos < block.length && this._compareFn(block[pos], value) === 0
  }

  discard(value: T): boolean {
    if (this._size === 0) return false
    const block = this._findBlock(value)
    const pos = this._bisectLeft(block, value)
    if (pos === block.length || this._compareFn(block[pos], value) !== 0) {
      return false
    }
    block.splice(pos, 1)
    this._size -= 1
    if (block.length === 0) {
      this._rebuild()
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
        if (block.length === 0) {
          this._rebuild()
        }
        return res
      }
      index -= block.length
    }
    return void 0
  }

  count(value: T): number {
    return this.bisectRight(value) - this.bisectLeft(value)
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

  clear(): void {
    this._blocks = []
    this._size = 0
  }

  toString(): string {
    return `SortedList[${[...this].join(', ')}]`
  }

  // Find the block which should contain x. Block must not be empty.
  private _findBlock(x: T): T[] {
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (this._compareFn(x, block[block.length - 1]) <= 0) {
        return block
      }
    }
    return this._blocks[this._blocks.length - 1]
  }

  private _rebuild(): void {
    const count = ~~((this._size + SortedList._BLOCK_RATIO - 1) / SortedList._BLOCK_RATIO)
    const newBlocks: T[][] = Array.from({ length: count }, () => [])
    let ptr = 0
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      for (let j = 0; j < block.length; j++) {
        newBlocks[~~(ptr / SortedList._BLOCK_RATIO)].push(block[j])
        ptr++
      }
    }
    this._blocks = newBlocks
  }

  // eslint-disable-next-line class-methods-use-this
  private _initBlocks(sorted: T[]): T[][] {
    const count = ~~((sorted.length + SortedList._BLOCK_RATIO - 1) / SortedList._BLOCK_RATIO)
    const newBlocks: T[][] = Array.from({ length: count }, () => [])
    for (let i = 0; i < count; i++) {
      const start = ~~((i * sorted.length) / count)
      const end = ~~(((i + 1) * sorted.length) / count)
      for (let j = start; j < end; j++) {
        newBlocks[i].push(sorted[j])
      }
    }
    return newBlocks
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

  const sl = new SortedList<number>([1, 4, 2, 14, 611, 3])
  console.log(sl.length)
  console.log(sl.toString())
  console.log(sl.length)
  console.log(sl.at(0))
  console.log(sl.at(1))
  console.log(sl.at(3))
  console.log(sl.at(-1))
  console.log(sl.at(-2))
}
