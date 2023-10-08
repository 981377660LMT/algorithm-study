/* eslint-disable no-param-reassign */
/* eslint-disable no-inner-declarations */
/* eslint-disable class-methods-use-this */

/**
 * @summary
 * 高效计算数列的前缀和，区间和
 * 树状数组或二叉索引树（Binary Indexed Tree, Fenwick Tree）
 * 性质
 * 1. tree[x]保存以x为根的子树中叶节点值的和
 * 2. tree[x]的父节点为tree[x+lowbit(x)]
 * 3. tree[x]节点覆盖的长度等于lowbit(x)
 * 4. 树的高度为logn+1
 */

import assert from 'assert'

// BITArray
// BIT1 (Map)
// BIT2
// BIT2Map
// BIT2D
// BIT3
// BIT4

/**
 * Point add range sum, 0-indexed.
 */
class BITArray {
  /**
   * Build a tree from an array-like object using dp.
   * O(n) time.
   */
  private static _buildTree(arr: ArrayLike<number>): Float64Array {
    const tree = new Float64Array(arr.length + 1)
    for (let i = 1; i < tree.length; i++) {
      tree[i] += arr[i - 1]
      const parent = i + (i & -i)
      if (parent < tree.length) tree[parent] += tree[i]
    }
    return tree
  }

  readonly length: number
  private readonly _tree: Float64Array

  /**
   * 指定长度或者从类数组建立树状数组.
   *
   * @warning
   * !如果需要使用`值域树状数组`，需要在构造函数中传入`长度n(值域1-n)`而不是类数组.
   */
  constructor(lengthOrArrayLike: number | ArrayLike<number>) {
    if (typeof lengthOrArrayLike === 'number') {
      this.length = lengthOrArrayLike
      this._tree = new Float64Array(lengthOrArrayLike + 1)
    } else {
      this.length = lengthOrArrayLike.length
      this._tree = BITArray._buildTree(lengthOrArrayLike)
    }
  }

  /**
   * Add delta to the element at index.
   * @param index 0 <= index < {@link length}.
   */
  add(index: number, delta: number): void {
    index++
    for (let i = index; i <= this.length; i += i & -i) {
      this._tree[i] += delta
    }
  }

  /**
   * Query the sum of [0, end).
   */
  query(end: number): number {
    if (end > this.length) end = this.length
    let res = 0
    for (let i = end; i > 0; i &= i - 1) {
      res += this._tree[i]
    }
    return res
  }

  /**
   * Query the sum of [start, end).
   */
  queryRange(start: number, end: number): number {
    return this.query(end) - this.query(start)
  }

  toString(): string {
    const sb: string[] = []
    sb.push('BITArray: [')
    for (let i = 1; i < this._tree.length; i++) {
      sb.push(String(this.queryRange(i, i + 1)))
      if (i < this._tree.length - 1) sb.push(', ')
    }
    sb.push(']')
    return sb.join('')
  }
}

/**
 * Point add range sum, 0-indexed.
 * Implemented by Map. Slow.
 */
class BIT1 {
  readonly size: number
  private readonly _tree: Map<number, number> = new Map()

  constructor(size: number) {
    this.size = size + 5
  }

  add(index: number, delta: number): void {
    index++
    for (let i = index; i <= this.size; i += i & -i) {
      this._tree.set(i, (this._tree.get(i) || 0) + delta)
    }
  }

  /**
   * [0,index).
   */
  query(index: number): number {
    if (index > this.size) index = this.size
    let res = 0
    for (let i = index; i > 0; i &= i - 1) {
      res += this._tree.get(i) || 0
    }
    return res
  }

  /**
   * [left,right).
   */
  queryRange(left: number, right: number): number {
    return this.query(right) - this.query(left)
  }

  toString(): string {
    const sb: string[] = []
    sb.push('BIT1: [')
    for (let i = 0; i < this.size; i++) {
      sb.push(String(this.queryRange(i, i + 1)))
      if (i < this.size - 1) sb.push(',')
    }
    sb.push(']')
    return sb.join('')
  }
}

const BITMap = BIT1

/**
 * 区间修改 区间查询, 0-indexed.
 */
class BIT2 {
  readonly size: number
  private readonly _tree1: number[]
  private readonly _tree2: number[]

  constructor(size: number) {
    this.size = size
    this._tree1 = Array(size + 1).fill(0)
    this._tree2 = Array(size + 1).fill(0)
  }

  addRange(left: number, right: number, k: number): void {
    right--
    this._add(left, k)
    this._add(right + 1, -k)
  }

  queryRange(left: number, right: number): number {
    right--
    return this._query(right) - this._query(left - 1)
  }

  toString(): string {
    const sb: string[] = []
    sb.push('BIT2: [')
    for (let i = 0; i < this.size; i++) {
      sb.push(String(this.queryRange(i, i + 1)))
      if (i < this.size - 1) sb.push(',')
    }
    sb.push(']')
    return sb.join('')
  }

  private _add(index: number, delta: number): void {
    index++
    for (let i = index; i <= this.size; i += i & -i) {
      this._tree1[i] += delta
      this._tree2[i] += (index - 1) * delta
    }
  }

  private _query(index: number): number {
    index++
    if (index > this.size) index = this.size
    let res = 0
    for (let i = index; i > 0; i &= i - 1) {
      res += index * this._tree1[i] - this._tree2[i]
    }
    return res
  }
}

/**
 * 区间修改 区间查询, 0-indexed.
 */
class BIT2Map {
  readonly size: number
  private readonly _tree1: Map<number, number> = new Map()
  private readonly _tree2: Map<number, number> = new Map()

  constructor(size: number) {
    this.size = size + 5
  }

  /**
   * [left,right)
   */
  addRange(left: number, right: number, delta: number): void {
    right--
    this._add(left, delta)
    this._add(right + 1, -delta)
  }

  /**
   * [left,right)
   */
  queryRange(left: number, right: number): number {
    right--
    return this._query(right) - this._query(left - 1)
  }

  toString(): string {
    const sb: string[] = []
    sb.push('BIT2Map: [')
    for (let i = 0; i < this.size; i++) {
      sb.push(String(this.queryRange(i, i + 1)))
      if (i < this.size - 1) sb.push(',')
    }
    sb.push(']')
    return sb.join('')
  }

  private _add(index: number, delta: number): void {
    index++
    for (let i = index; i <= this.size; i += i & -i) {
      this._tree1.set(i, (this._tree1.get(i) || 0) + delta)
      this._tree2.set(i, (this._tree2.get(i) || 0) + (index - 1) * delta)
    }
  }

  private _query(index: number): number {
    index++
    if (index > this.size) index = this.size
    let res = 0
    for (let i = index; i > 0; i &= i - 1) {
      res += index * (this._tree1.get(i) || 0) - (this._tree2.get(i) || 0)
    }
    return res
  }
}

/**
 * 二维单点修改 区间查询
 */
class BIT3 {
  private readonly _ROW: number
  private readonly _COL: number
  private readonly _tree: number[][]

  constructor(row: number, col: number) {
    this._ROW = row
    this._COL = col
    this._tree = Array(row + 1)
    for (let i = 0; i <= row; i++) {
      this._tree[i] = Array(col + 1).fill(0)
    }
  }

  /**
   * 单点修改 (row,col)的值为加上delta
   * 0 <= row < {@link _ROW}
   * 0 <= col < {@link _COL}
   */
  add(row: number, col: number, delta: number): void {
    row++, col++
    for (let r = row; r <= this._ROW; r += r & -r) {
      for (let c = col; c <= this._COL; c += c & -c) {
        this._tree[r][c] += delta
      }
    }
  }

  /**
   * 查询左上角 (row1,col1) 到右下角 (row2,col2) 的和
   * 0 <= row1 <= row2 < {@link _ROW}
   * 0 <= col1 <= col2 < {@link _COL}
   */
  queryRange(row1: number, col1: number, row2: number, col2: number): number {
    return (
      this._query(row2, col2) -
      this._query(row1 - 1, col2) -
      this._query(row2, col1 - 1) +
      this._query(row1 - 1, col1 - 1)
    )
  }

  private _query(row: number, col: number): number {
    row++, col++
    if (row > this._ROW) row = this._ROW
    if (col > this._COL) col = this._COL
    let res = 0
    for (let r = row; r > 0; r -= r & -r) {
      for (let c = col; c > 0; c -= c & -c) {
        res += this._tree[r][c]
      }
    }
    return res
  }
}

/**
 * 二维区间修改 区间查询
 */
class BIT4 {
  private readonly _ROW: number
  private readonly _COL: number
  private readonly _tree1: number[][]
  private readonly _tree2: number[][]
  private readonly _tree3: number[][]
  private readonly _tree4: number[][]

  constructor(row: number, col: number) {
    this._ROW = row
    this._COL = col
    this._tree1 = Array(row + 1)
    this._tree2 = Array(row + 1)
    this._tree3 = Array(row + 1)
    this._tree4 = Array(row + 1)
    for (let i = 0; i <= row; i++) {
      this._tree1[i] = Array(col + 1).fill(0)
      this._tree2[i] = Array(col + 1).fill(0)
      this._tree3[i] = Array(col + 1).fill(0)
      this._tree4[i] = Array(col + 1).fill(0)
    }
  }

  /**
   * 区间修改 (row1,col1) 到 (row2,col2) 里的每一个点的值加上delta
   * 0<=row1<=row2<=ROW-1, 0<=col1<=col2<=COL-1
   */
  addRange(row1: number, col1: number, row2: number, col2: number, delta: number): void {
    this._add(row1, col1, delta)
    this._add(row2 + 1, col1, -delta)
    this._add(row1, col2 + 1, -delta)
    this._add(row2 + 1, col2 + 1, delta)
  }

  /**
   * 查询左上角 (row1,col1) 到右下角 (row2,col2) 的和
   * 0<=row1<=row2<=ROW-1, 0<=col1<=col2<=COL-1
   */
  queryRange(row1: number, col1: number, row2: number, col2: number): number {
    return (
      this._query(row2, col2) -
      this._query(row1 - 1, col2) -
      this._query(row2, col1 - 1) +
      this._query(row1 - 1, col1 - 1)
    )
  }

  private _add(row: number, col: number, delta: number): void {
    row++, col++
    for (let r = row; r <= this._ROW; r += r & -r) {
      for (let c = col; c <= this._COL; c += c & -c) {
        this._tree1[r][c] += delta
        this._tree2[r][c] += (row - 1) * delta
        this._tree3[r][c] += (col - 1) * delta
        this._tree4[r][c] += (row - 1) * (col - 1) * delta
      }
    }
  }

  private _query(row: number, col: number): number {
    row++, col++
    if (row > this._ROW) row = this._ROW
    if (col > this._COL) col = this._COL
    let res = 0
    for (let r = row; r > 0; r -= r & -r) {
      for (let c = col; c > 0; c -= c & -c) {
        res +=
          row * col * this._tree1[r][c] -
          col * this._tree2[r][c] -
          row * this._tree3[r][c] +
          this._tree4[r][c]
      }
    }
    return res
  }
}

if (require.main === module) {
  const bit1 = new BIT1(5)
  assert.strictEqual(bit1.query(1), 0)
  bit1.add(0, 3)
  assert.strictEqual(bit1.query(1), 3)

  const bit2 = new BIT2(10)
  bit2.addRange(2, 5, 1) // 区间更新
  bit2.addRange(2, 5, 1) // 单点更新
  assert.strictEqual(bit2.queryRange(2, 4), 4)
  assert.strictEqual(bit2.queryRange(2, 3), 2)

  const bit3 = new BIT3(3, 3)
  bit3.add(1, 1, 1)
  assert.strictEqual(bit3.queryRange(0, 0, 4, 4), 1)
  bit3.add(1, 1, 2)
  assert.strictEqual(bit3.queryRange(0, 0, 4, 4), 3)
  bit3.add(0, 0, 3)
  assert.strictEqual(bit3.queryRange(0, 0, 4, 4), 6)

  const bit4 = new BIT4(3, 3)
  bit4.addRange(1, 1, 1, 1, 1)
  assert.strictEqual(bit4.queryRange(0, 0, 4, 4), 1)
  bit4.addRange(1, 1, 1, 1, 2)
  assert.strictEqual(bit4.queryRange(0, 0, 4, 4), 3)
  bit4.addRange(0, 0, 0, 0, 3)
  assert.strictEqual(bit4.queryRange(0, 0, 4, 4), 6)

  const bitArray = new BITArray([1, 2, 3])
  console.log(bitArray.toString())

  const bit2Map = new BIT2Map(10)
  bit2Map.addRange(2, 5, 1) // 区间更新
  bit2Map.addRange(2, 5, 1) // 单点更新
  assert.strictEqual(bit2Map.queryRange(2, 4), 4)
  assert.strictEqual(bit2Map.queryRange(2, 3), 2)
}

export { BIT1, BIT2, BIT2Map, BIT3, BIT4, BITArray, BITMap }

if (require.main === module) {
  // https://leetcode.cn/problems/maximum-white-tiles-covered-by-a-carpet/
  function maximumWhiteTiles(tiles: number[][], carpetLen: number): number {
    const bit = new BIT2(1e9 + 10)
    let res = 0
    tiles.forEach(([left, right]) => {
      bit.addRange(left, right + 1, 1)
    })
    tiles.forEach(([left]) => {
      res = Math.max(res, bit.queryRange(left, left + carpetLen))
    })
    return res
  }
}
