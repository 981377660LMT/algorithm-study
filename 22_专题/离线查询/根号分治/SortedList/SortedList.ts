/* eslint-disable no-console */
/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */

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
  private static readonly _BLOCK_RATIO = 50
  private static readonly _REBUILD_RATIO = 170

  private readonly _compareFn: (a: T, b: T) => number
  private _size: number
  private _blocks: T[][]

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
    if (this._size === 0) {
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
    if (this._size === 0) return false
    const blockIndex = this._findBlockIndex(value)
    if (blockIndex === void 0) return false
    const block = this._blocks[blockIndex]
    const pos = this._bisectLeft(block, value)
    return pos < block.length && value === block[pos]
  }

  /** 注意内部使用 `===` 来比较两个对象是否相等 */
  discard(value: T): boolean {
    if (this._size === 0) return false
    const blockIndex = this._findBlockIndex(value)
    if (blockIndex === -1) return false
    const block = this._blocks[blockIndex]
    const pos = this._bisectLeft(block, value)
    if (pos === block.length || block[pos] !== value) {
      return false
    }
    block.splice(pos, 1)
    this._size -= 1
    if (block.length === 0) {
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
        if (block.length === 0) {
          this._blocks.splice(i, 1) // !Splice When Empty, Do Not Rebuild
        }
        return res
      }
      index -= block.length
    }
    return void 0
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

  clear(): void {
    this._blocks = []
    this._size = 0
  }

  toString(): string {
    return `SortedList[${[...this].join(', ')}]`
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
    if (this._size === 0) {
      return
    }
    const bCount = Math.ceil(Math.sqrt(this._size / SortedList._BLOCK_RATIO))
    const bSize = ~~((this._size + bCount - 1) / bCount) // ceil
    const newB: T[][] = Array.from({ length: bCount }, () => [])
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
    const newB: T[][] = Array.from({ length: bCount }, () => [])
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
}
