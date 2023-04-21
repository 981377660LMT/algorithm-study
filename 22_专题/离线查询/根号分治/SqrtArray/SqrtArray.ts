/* eslint-disable generator-star-spacing */
/* eslint-disable no-empty */

// https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/data_structure/sqrt_array.cc
//
// SQRT Array
//
// Description:
//   An array with o(n) deletion and insertion
//
// Algorithm:
//   Decompose array into O(sqrt(n)) subarrays.
//   Then the all operation is performed in O(sqrt(n)).
//
// Complexity:
//   O(sqrt(n)); however, due to the cheap constant factor,
//   it is comparable with binary search trees.
//   If only deletion is required, it is better choice.

/**
 * 分块数组.
 */
class SqrtArray<T = number> {
  private _n = 0
  private readonly _x: T[][] = []

  constructor(nOrNums: number | ArrayLike<T> = 0) {
    if (typeof nOrNums === 'number') {
      nOrNums = Array(nOrNums)
    }
    const n = nOrNums.length
    if (n === 0) {
      return
    }

    const bCount = ~~Math.sqrt(n)
    const bSize = ~~((n + bCount - 1) / bCount)
    const newB: T[][] = Array(bCount)
    for (let i = 0; i < bCount; i++) {
      newB[i] = []
      for (let j = i * bSize; j < Math.min((i + 1) * bSize, n); j++) {
        newB[i].push(nOrNums[j])
      }
    }
    this._n = n
    this._x = newB
  }

  /**
   * 0<= i < {@link length}
   */
  set(i: number, v: T): void {
    if (i === this._n - 1) {
      const bi = this._x.length - 1
      this._x[bi][this._x[bi].length - 1] = v
      return
    }
    let bi = 0
    for (; i >= this._x[bi].length; i -= this._x[bi++].length) {}
    this._x[bi][i] = v
  }

  /**
   * 0<= i < {@link length}
   */
  get(i: number): T | undefined {
    if (i === this._n - 1) {
      const bi = this._x.length - 1
      return this._x[bi][this._x[bi].length - 1]
    }
    let bi = 0
    for (; i >= this._x[bi].length; i -= this._x[bi++].length);
    return this._x[bi][i]
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
    let bi = 0
    let res: T | undefined
    if (i === this._n - 1) {
      bi = this._x.length - 1
      res = this._x[bi].pop()
    } else {
      for (; i >= this._x[bi].length; i -= this._x[bi++].length) {}
      res = this._x[bi][i]
      this._x[bi].splice(i, 1)
    }
    this._n--
    if (!this._x[bi].length) this._x.splice(bi, 1)
    return res
  }

  shift(): T | undefined {
    return this.pop(0)
  }

  unshift(v: T): void {
    this.insert(0, v)
  }

  /**
   * 删除区间 [start, end) 内的元素.
   * 0<= start <= end <= {@link length}
   */
  erase(start: number, end: number): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return

    let [bid, startPos] = this._moveTo(start)
    let deleteCount = end - start
    for (; bid < this._x.length && deleteCount > 0; bid++) {
      const block = this._x[bid]
      const endPos = Math.min(block.length, startPos + deleteCount)
      const curDeleteCount = endPos - startPos
      if (curDeleteCount === block.length) {
        this._x.splice(bid, 1)
        bid--
      } else {
        block.splice(startPos, curDeleteCount)
      }
      deleteCount -= curDeleteCount
      this._n -= curDeleteCount
      startPos = 0
    }
  }

  /**
   * 在 i 位置`前`插入 v.
   * 0<= i <= {@link length}
   */
  insert(i: number, v: T): void {
    if (!this._n) {
      this._x.push([v])
      this._n++
      return
    }

    let bi = 0
    if (i >= this._n) {
      bi = this._x.length - 1
      this._x[bi].push(v)
    } else {
      for (; bi < this._x.length && i >= this._x[bi].length; i -= this._x[bi++].length) {}
      this._x[bi].splice(i, 0, v)
    }
    this._n++

    const sqrtn2 = ~~Math.sqrt(this._n) * 3
    // !rebuild when block size > 6 * sqrt(n), about 2000 when n = 1e5
    if (this._x[bi].length > 2 * sqrtn2) {
      const y = this._x[bi].splice(sqrtn2)
      this._x.splice(bi + 1, 0, y)
    }
  }

  /**
   * 遍历区间 [start, end) 内的元素,并选择是否在遍历后删除.
   * 0<= start <= end <= {@link length}
   */
  enumerate(start: number, end: number, f: (value: T) => void, erase = false): void {
    let [bid, startPos] = this._moveTo(start)
    let count = end - start

    for (; bid < this._x.length && count > 0; bid++) {
      const block = this._x[bid]
      const endPos = Math.min(block.length, startPos + count)
      for (let j = startPos; j < endPos; j++) {
        f(block[j])
      }

      const curDeleteCount = endPos - startPos
      if (erase) {
        if (curDeleteCount === block.length) {
          this._x.splice(bid, 1)
          bid--
        } else {
          block.splice(startPos, curDeleteCount)
        }
        this._n -= curDeleteCount
      }

      count -= curDeleteCount
      startPos = 0
    }
  }

  slice(start: number, end: number): T[] {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return []
    let count = end - start
    const res: T[] = Array(count)
    let [bid, startPos] = this._moveTo(start)
    let ptr = 0
    for (; bid < this._x.length && count > 0; bid++) {
      const block = this._x[bid]
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

  fill(v: T): this {
    this._x.forEach(b => b.fill(v))
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
      let [bid, endPos] = this._moveTo(end - 1)
      for (; ~bid && count > 0; bid--, ~bid && (endPos = this._x[bid].length)) {
        const block = this._x[bid]
        const startPos = Math.max(0, endPos - count)
        const curCount = endPos - startPos
        for (let j = endPos - 1; j >= startPos; j--) {
          yield block[j]
        }
        count -= curCount
      }
    } else {
      let [bid, startPos] = this._moveTo(start)
      for (; bid < this._x.length && count > 0; bid++) {
        const block = this._x[bid]
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

  forEach(callback: (value: T, index: number) => void): void {
    let ptr = 0
    for (let bi = 0; bi < this._x.length; ++bi) {
      for (let j = 0; j < this._x[bi].length; ++j) {
        callback(this._x[bi][j], ptr++)
      }
    }
  }

  *entries(): IterableIterator<[number, T]> {
    let ptr = 0
    for (let i = 0; i < this._x.length; i++) {
      const block = this._x[i]
      for (let j = 0; j < block.length; j++) {
        yield [ptr++, block[j]]
      }
    }
  }

  *[Symbol.iterator](): Iterator<T> {
    for (let i = 0; i < this._x.length; i++) {
      const block = this._x[i]
      for (let j = 0; j < block.length; j++) {
        yield block[j]
      }
    }
  }

  toString(): string {
    return `SqrtArray{${this._x}}`
  }

  private _moveTo(index: number): [blockId: number, startPos: number] {
    for (let i = 0; i < this._x.length; i++) {
      const block = this._x[i]
      if (index < block.length) {
        return [i, index]
      }
      index -= block.length
    }
    return [this._x.length, 0]
  }

  get length(): number {
    return this._n
  }
}

export { SqrtArray }

if (require.main === module) {
  const arr = new SqrtArray()
  const rands = Array(4e5)
    .fill(0)
    .map((_, i) => ~~(Math.random() * i))
  console.time('insert')
  for (let i = 0; i < 4e5; i++) {
    arr.insert(rands[i], i)
    arr.get(rands[i])
    arr.set(rands[i], i)
  }

  for (let i = 0; i < 4e5; i++) {
    arr.get(rands[i])
  }
  console.timeEnd('insert')
  console.log(rands.slice(0, 10))
  arr.erase(0, 100)
  console.log(arr.length)

  // https://leetcode.cn/problems/design-circular-deque/
  class MyCircularDeque {
    private readonly _sqrt: SqrtArray<number>
    private readonly _k: number
    constructor(k: number) {
      this._sqrt = new SqrtArray()
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
