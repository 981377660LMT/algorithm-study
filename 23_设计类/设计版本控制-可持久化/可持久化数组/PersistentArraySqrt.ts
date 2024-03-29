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

export {}
