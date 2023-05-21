/* eslint-disable no-param-reassign */

//
// Persistent Array (sqrt decomposition)
//
// Description:
//   An array with O(sqrt(n)) operations (inc. copy)
//
// Algorithm:
//   Store base aray and operation sequence.
//   If the length of operation sequence exceeds sqrt(n),
//   update base array and clear the operation sequence.
//
// Complexity:
//   Copy O(sqrt(n))
//   Get O(sqrt(n))
//   Set O(sqrt(n)) time and space per operation
//
// Comment:
//   This implementation is much faster than the implementation
//   based on a complete binary tree representation,
//   which runs in O(log n) time / extra space per operation.

// 非常通用的实现:
// 1. 记录每个版本的修改:[index,value,version]
// 2. 查询时,从最新版本开始,找到最近的修改,返回修改后的值.
// 3. 更新时，如果修改的数量超过了2*sqrt(n),就更新基础数组并清空操作序列，保证修改的记录不会太长.
// !整体的时间复杂度是O(sqrt(n)).

/**
 * 基于分块的完全可持久化数组.
 */
class PersistentArraySqrt<T> {
  private _arr: ArrayLike<T>
  private _opIndex: number[] = []
  private _opValue: T[] = []
  private _opVersion = 0

  constructor(nOrNums: number | ArrayLike<T>) {
    if (typeof nOrNums === 'number') {
      nOrNums = Array(nOrNums)
    }
    this._arr = nOrNums
  }

  get(i: number): T | undefined {
    // 查询上一个修改操作
    for (let j = this._opVersion - 1; ~j; j--) {
      if (this._opIndex[j] === i) {
        return this._opValue[j]
      }
    }
    return this._arr[i]
  }

  set(i: number, v: T): PersistentArraySqrt<T> {
    this._opIndex.push(i)
    this._opValue.push(v)
    const n = this._arr.length
    if (this._opIndex.length * this._opIndex.length <= 4 * n) {
      return this._update()
    }
    // !如果操作序列的长度超过了2*sqrt(n)，就更新基础数组并清空操作序列.
    const newArr = new Array(n)
    for (let j = 0; j < n; j++) {
      newArr[j] = this._arr[j]
    }
    this._opIndex.forEach((v, i) => {
      newArr[v] = this._opValue[i]
    })
    return new PersistentArraySqrt(newArr)
  }

  private _update(): PersistentArraySqrt<T> {
    const copy = new PersistentArraySqrt(this._arr)
    copy._opIndex = this._opIndex
    copy._opValue = this._opValue
    copy._opVersion = this._opVersion + 1
    return copy
  }
}

class PersistentArraySqrtInt32 {
  private _arr: Int32Array
  private _opIndex: number[] = []
  private _opValue: number[] = []
  private _opVersion = 0

  constructor(nOrNums: number | Int32Array) {
    if (typeof nOrNums === 'number') {
      nOrNums = new Int32Array(nOrNums)
    }
    this._arr = nOrNums
  }

  get(i: number): number | undefined {
    for (let j = this._opVersion - 1; ~j; j--) {
      if (this._opIndex[j] === i) {
        return this._opValue[j]
      }
    }
    return this._arr[i]
  }

  set(i: number, v: number): PersistentArraySqrtInt32 {
    this._opIndex.push(i)
    this._opValue.push(v)
    const n = this._arr.length
    if (this._opIndex.length * this._opIndex.length <= 4 * n) {
      return this._update()
    }

    const newArr = new Int32Array(n)
    for (let j = 0; j < n; j++) {
      newArr[j] = this._arr[j]
    }
    this._opIndex.forEach((v, i) => {
      newArr[v] = this._opValue[i]
    })
    return new PersistentArraySqrtInt32(newArr)
  }

  private _update(): PersistentArraySqrtInt32 {
    const copy = new PersistentArraySqrtInt32(this._arr)
    copy._opIndex = this._opIndex
    copy._opValue = this._opValue
    copy._opVersion = this._opVersion + 1
    return copy
  }
}

export { PersistentArraySqrt, PersistentArraySqrtInt32 }

if (require.main === module) {
  // https://leetcode.cn/problems/snapshot-array/
  class SnapshotArray {
    private readonly _gits: PersistentArraySqrtInt32[] = []
    private _root: PersistentArraySqrtInt32
    constructor(length: number) {
      this._root = new PersistentArraySqrtInt32(length)
    }

    set(index: number, val: number): void {
      this._root = this._root.set(index, val)
    }

    snap(): number {
      this._gits.push(this._root)
      return this._gits.length - 1
    }

    get(index: number, snapId: number): number {
      return this._gits[snapId].get(index)!
    }
  }
}
