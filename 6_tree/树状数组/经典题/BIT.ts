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

/**
 * Point add range sum, 1-indexed.
 */
class BITArray {
  /**
   * Build a tree from an array-like object using dp.
   * O(n) time.
   */
  private static _build(arrayLike: ArrayLike<number>): number[] {
    const tree = Array(arrayLike.length + 1).fill(0)
    for (let i = 1; i < tree.length; i++) {
      tree[i] += arrayLike[i - 1]
      const parent = i + (i & -i)
      if (parent < tree.length) tree[parent] += tree[i]
    }
    return tree
  }

  readonly length: number
  private readonly _tree: number[]

  constructor(lengthOrArrayLike: number | ArrayLike<number>) {
    if (typeof lengthOrArrayLike === 'number') {
      this.length = lengthOrArrayLike
      this._tree = Array(lengthOrArrayLike + 1).fill(0)
    } else {
      this.length = lengthOrArrayLike.length
      this._tree = BITArray._build(lengthOrArrayLike)
    }
  }

  /**
   * Add delta to the element at index.
   * @param index 1 <= index <= {@link length}
   */
  add(index: number, delta: number): void {
    if (index <= 0) throw new RangeError(`add index must be greater than 0, but got ${index}`)
    for (let i = index; i <= this.length; i += i & -i) {
      this._tree[i] += delta
    }
  }

  /**
   * Query the sum of [1, right].
   */
  query(right: number): number {
    if (right > this.length) right = this.length
    let res = 0
    for (let i = right; i > 0; i -= i & -i) {
      res += this._tree[i]
    }
    return res
  }

  /**
   * Query the sum of [left, right].
   */
  queryRange(left: number, right: number): number {
    return this.query(right) - this.query(left - 1)
  }

  toString(): string {
    const sb: string[] = []
    sb.push('BITArray: [')
    for (let i = 1; i < this._tree.length; i++) {
      sb.push(String(this.queryRange(i, i)))
      if (i < this._tree.length - 1) sb.push(', ')
    }
    sb.push(']')
    return sb.join('')
  }
}

// Implemented by Map. Slow.
class BIT1 {
  readonly _size: number
  private readonly _tree: Map<number, number> = new Map()

  constructor(size: number) {
    this._size = size
  }

  add(index: number, delta: number): void {
    if (index <= 0) throw RangeError(`add索引 ${index} 应为正整数`)
    for (let i = index; i <= this._size; i += i & -i) {
      this._tree.set(i, (this._tree.get(i) || 0) + delta)
    }
  }

  query(index: number): number {
    if (index > this._size) index = this._size
    let res = 0
    for (let i = index; i > 0; i -= i & -i) {
      res += this._tree.get(i) || 0
    }
    return res
  }

  queryRange(left: number, right: number): number {
    return this.query(right) - this.query(left - 1)
  }
}

/**
 * @description 区间修改 区间查询
 * @see {@link https://github.com/981377660LMT/algorithm-study/blob/master/6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.ts}
 */
class BIT2 {
  readonly _size: number
  private readonly _tree1: Map<number, number> = new Map()
  private readonly _tree2: Map<number, number> = new Map()

  constructor(size: number) {
    this._size = size
  }

  add(left: number, right: number, k: number): void {
    this._add(left, k)
    this._add(right + 1, -k)
  }

  query(left: number, right: number): number {
    return this._query(right) - this._query(left - 1)
  }

  private _add(index: number, delta: number): void {
    if (index <= 0) throw RangeError(`add索引 ${index} 应为正整数`)
    for (let i = index; i <= this._size; i += i & -i) {
      // 此处进行了差分操作，记录差分操作大小
      this._tree1.set(i, (this._tree1.get(i) || 0) + delta)
      // 前x-1个数没有进行差分操作，这里把总值记录下来
      this._tree2.set(i, (this._tree2.get(i) || 0) + (index - 1) * delta)
    }
  }

  private _query(index: number): number {
    if (index > this._size) index = this._size
    let res = 0
    for (let i = index; i > 0; i -= i & -i) {
      res += index * (this._tree1.get(i) || 0) - (this._tree2.get(i) || 0)
    }
    return res
  }
}

/**
 * @description !二维单点修改 区间查询
 * @see {@link https://github.com/981377660LMT/algorithm-study/blob/master/6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.ts}
 */
class BIT3 {
  private readonly _tree: Map<number, Map<number, number>> = new Map()
  private readonly _ROW: number
  private readonly _COL: number

  constructor(row: number, col: number) {
    this._ROW = row
    this._COL = col
  }

  /**
   * @description 单点修改 (row,col)的值为加上delta
   */
  update(row: number, col: number, delta: number): void {
    row++, col++
    for (let r = row; r <= this._ROW; r += r & -r) {
      for (let c = col; c <= this._COL; c += c & -c) {
        this._addDeep(this._tree, r, c, delta)
      }
    }
  }

  /**
   * @description 左上角 (0,0) 到 右下角 (row,col) 的矩形里所有数的和
   */
  query(row: number, col: number): number {
    row++, col++
    if (row > this._ROW) row = this._ROW
    if (col > this._COL) col = this._COL
    let res = 0
    for (let r = row; r > 0; r -= r & -r) {
      for (let c = col; c > 0; c -= c & -c) {
        res += this._getDeep(this._tree, r, c)
      }
    }
    return res
  }

  /**
   * @description 查询左上角 (row1,col1) 到右下角 (row2,col2) 的和
   */
  queryRange(row1: number, col1: number, row2: number, col2: number): number {
    return (
      this.query(row2, col2) -
      this.query(row1 - 1, col2) -
      this.query(row2, col1 - 1) +
      this.query(row1 - 1, col1 - 1)
    )
  }

  private _addDeep(
    map: Map<number, Map<number, number>>,
    key1: number,
    key2: number,
    delta: number
  ): void {
    if (!map.has(key1)) map.set(key1, new Map())
    const innerMap = map.get(key1)!
    innerMap.set(key2, (innerMap.get(key2) || 0) + delta)
  }

  private _getDeep(map: Map<number, Map<number, number>>, key1: number, key2: number): number {
    if (!map.has(key1)) return 0
    const innerMap = map.get(key1)!
    return innerMap.get(key2) || 0
  }
}

/**
 * @description 二维区间修改 区间查询
 * @see {@link https://github.com/981377660LMT/algorithm-study/blob/master/6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.ts}
 */
class BIT4 {
  private readonly _ROW: number
  private readonly _COL: number
  private readonly _tree1: Map<number, Map<number, number>> = new Map()
  private readonly _tree2: Map<number, Map<number, number>> = new Map()
  private readonly _tree3: Map<number, Map<number, number>> = new Map()
  private readonly _tree4: Map<number, Map<number, number>> = new Map()

  constructor(row: number, col: number) {
    this._ROW = row
    this._COL = col
  }

  /**
   * @description 区间修改 (row1,col1) 到 (row2,col2) 里的每一个点的值加上delta
   */
  updateRange(row1: number, col1: number, row2: number, col2: number, delta: number): void {
    this.update(row1, col1, delta)
    this.update(row2 + 1, col1, -delta)
    this.update(row1, col2 + 1, -delta)
    this.update(row2 + 1, col2 + 1, delta)
  }

  /**
   * @description 查询左上角 (row1,col1) 到右下角 (row2,col2) 的和
   */
  queryRange(row1: number, col1: number, row2: number, col2: number): number {
    return (
      this.query(row2, col2) -
      this.query(row1 - 1, col2) -
      this.query(row2, col1 - 1) +
      this.query(row1 - 1, col1 - 1)
    )
  }

  /**
   * @description 单点修改 (row,col)的值为加上delta
   */
  private update(row: number, col: number, delta: number): void {
    row++, col++
    const preRow = row
    const preCol = col
    for (let r = row; r <= this._ROW; r += r & -r) {
      for (let c = col; c <= this._COL; c += c & -c) {
        this.addDeep(this._tree1, r, c, delta)
        this.addDeep(this._tree2, r, c, (preRow - 1) * delta)
        this.addDeep(this._tree3, r, c, (preCol - 1) * delta)
        this.addDeep(this._tree4, r, c, (preRow - 1) * (preCol - 1) * delta)
      }
    }
  }

  /**
   * @description 左上角 (0,0) 到 右下角 (row,col) 的矩形里所有数的和
   */
  private query(row: number, col: number): number {
    row++, col++
    if (row > this._ROW) row = this._ROW
    if (col > this._COL) col = this._COL

    const preRow = row
    const preCol = col

    let res = 0
    for (let r = row; r > 0; r -= r & -r) {
      for (let c = col; c > 0; c -= c & -c) {
        res +=
          preRow * preCol * this.getDeep(this._tree1, r, c) -
          preCol * this.getDeep(this._tree2, r, c) -
          preRow * this.getDeep(this._tree3, r, c) +
          this.getDeep(this._tree4, r, c)
      }
    }

    return res
  }

  private addDeep(
    map: Map<number, Map<number, number>>,
    key1: number,
    key2: number,
    delta: number
  ): void {
    if (!map.has(key1)) map.set(key1, new Map())
    const innerMap = map.get(key1)!
    innerMap.set(key2, (innerMap.get(key2) ?? 0) + delta)
  }

  private getDeep(map: Map<number, Map<number, number>>, key1: number, key2: number): number {
    if (!map.has(key1)) return 0
    const innerMap = map.get(key1)!
    return innerMap.get(key2) ?? 0
  }
}

if (require.main === module) {
  const bit1 = new BIT1(5)
  assert.strictEqual(bit1.query(1), 0)
  bit1.add(1, 3)
  assert.strictEqual(bit1.query(1), 3)

  const bit2 = new BIT2(10)
  bit2.add(2, 4, 1) // 区间更新
  bit2.add(2, 2, 1) // 单点更新
  assert.strictEqual(bit2.query(2, 4), 4)
  assert.strictEqual(bit2.query(2, 2), 2)
  function maximumWhiteTiles(tiles: number[][], carpetLen: number): number {
    const bit = new BIT2(Math.max(...tiles.flat()) + 10)
    for (const [left, right] of tiles) bit.add(left, right, 1)
    let res = 0
    for (const [left] of tiles) res = Math.max(res, bit.query(left, left + carpetLen - 1))
    return res
  }

  const bit3 = new BIT3(3, 3)
  bit3.update(1, 1, 1)
  assert.strictEqual(bit3.queryRange(0, 0, 4, 4), 1)
  bit3.update(1, 1, 2)
  assert.strictEqual(bit3.queryRange(0, 0, 4, 4), 3)
  bit3.update(0, 0, 3)
  assert.strictEqual(bit3.queryRange(0, 0, 4, 4), 6)

  const bit4 = new BIT4(3, 3)
  bit4.updateRange(1, 1, 1, 1, 1)
  assert.strictEqual(bit4.queryRange(0, 0, 4, 4), 1)
  bit4.updateRange(1, 1, 1, 1, 2)
  assert.strictEqual(bit4.queryRange(0, 0, 4, 4), 3)
  bit4.updateRange(0, 0, 0, 0, 3)
  assert.strictEqual(bit4.queryRange(0, 0, 4, 4), 6)

  const bitArray = new BITArray([1, 2, 3])
  console.log(bitArray.toString())
}

export { BIT1, BIT2, BIT3, BIT4, BITArray }
