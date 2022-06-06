import assert from 'assert'

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
class BIT1 {
  readonly size: number
  private readonly tree: Map<number, number> = new Map()

  private static lowbit(x: number): number {
    return x & -x
  }

  constructor(size: number) {
    this.size = size
  }

  add(index: number, delta: number): void {
    if (index <= 0) throw Error('index 应为正整数')
    for (let i = index; i <= this.size; i += BIT1.lowbit(i)) {
      this.tree.set(i, (this.tree.get(i) ?? 0) + delta)
    }
  }

  query(index: number): number {
    if (index > this.size) index = this.size
    let res = 0
    for (let i = index; i > 0; i -= BIT1.lowbit(i)) {
      res += this.tree.get(i) ?? 0
    }
    return res
  }

  queryRange(left: number, right: number): number {
    return this.query(right) - this.query(left - 1)
  }
}

/**
 * @description 区间修改 区间查询
 */
class BIT2 {
  readonly size: number
  private readonly tree1: Map<number, number> = new Map()
  private readonly tree2: Map<number, number> = new Map()

  private static lowbit(x: number): number {
    return x & -x
  }

  constructor(size: number) {
    this.size = size
  }

  add(left: number, right: number, k: number): void {
    this._add(left, k)
    this._add(right + 1, -k)
  }

  query(left: number, right: number): number {
    return this._query(right) - this._query(left - 1)
  }

  private _add(index: number, delta: number): void {
    if (index <= 0) throw Error('查询索引应为正整数')
    for (let i = index; i <= this.size; i += BIT2.lowbit(i)) {
      this.tree1.set(i, (this.tree1.get(i) ?? 0) + delta) // 此处进行了差分操作，记录差分操作大小
      this.tree2.set(i, (this.tree2.get(i) ?? 0) + (index - 1) * delta) // 前x-1个数没有进行差分操作，这里把总值记录下来
    }
  }

  private _query(index: number): number {
    if (index > this.size) index = this.size
    let res = 0
    for (let i = index; i > 0; i -= BIT2.lowbit(i)) {
      res += index * (this.tree1.get(i) ?? 0) - (this.tree2.get(i) ?? 0)
    }

    return res
  }
}

/**
 * @description 二维单点修改 区间查询
 */
class BIT3 {
  private readonly tree: Map<number, Map<number, number>> = new Map()
  private readonly ROW: number
  private readonly COL: number

  constructor(row: number, col: number) {
    this.ROW = row
    this.COL = col
  }

  private static lowbit(x: number): number {
    return x & -x
  }

  /**
   * @description 单点修改 (row,col)的值为加上delta
   */
  update(row: number, col: number, delta: number): void {
    row++, col++
    for (let r = row; r <= this.ROW; r += BIT3.lowbit(r)) {
      for (let c = col; c <= this.COL; c += BIT3.lowbit(c)) {
        this.addDeep(this.tree, r, c, delta)
      }
    }
  }

  /**
   * @description 左上角 (0,0) 到 右下角 (row,col) 的矩形里所有数的和
   */
  query(row: number, col: number): number {
    row++, col++
    if (row > this.ROW) row = this.ROW
    if (col > this.COL) col = this.COL
    let res = 0
    for (let r = row; r > 0; r -= BIT3.lowbit(r)) {
      for (let c = col; c > 0; c -= BIT3.lowbit(c)) {
        res += this.getDeep(this.tree, r, c)
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

/**
 * @description 二维区间修改 区间查询
 */
class BIT4 {
  private readonly ROW: number
  private readonly COL: number
  private readonly tree1: Map<number, Map<number, number>> = new Map()
  private readonly tree2: Map<number, Map<number, number>> = new Map()
  private readonly tree3: Map<number, Map<number, number>> = new Map()
  private readonly tree4: Map<number, Map<number, number>> = new Map()

  constructor(row: number, col: number) {
    this.ROW = row
    this.COL = col
  }

  private static lowbit(x: number): number {
    return x & -x
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
    const [preRow, preCol] = [row, col]

    for (let r = row; r <= this.ROW; r += BIT4.lowbit(r)) {
      for (let c = col; c <= this.COL; c += BIT4.lowbit(c)) {
        this.addDeep(this.tree1, r, c, delta)
        this.addDeep(this.tree2, r, c, (preRow - 1) * delta)
        this.addDeep(this.tree3, r, c, (preCol - 1) * delta)
        this.addDeep(this.tree4, r, c, (preRow - 1) * (preCol - 1) * delta)
      }
    }
  }

  /**
   * @description 左上角 (0,0) 到 右下角 (row,col) 的矩形里所有数的和
   */
  private query(row: number, col: number): number {
    row++, col++
    if (row > this.ROW) row = this.ROW
    if (col > this.COL) col = this.COL

    const [preRow, preCol] = [row, col]

    let res = 0
    for (let r = row; r > 0; r -= BIT4.lowbit(r)) {
      for (let c = col; c > 0; c -= BIT4.lowbit(c)) {
        res +=
          preRow * preCol * this.getDeep(this.tree1, r, c) -
          preCol * this.getDeep(this.tree2, r, c) -
          preRow * this.getDeep(this.tree3, r, c) +
          this.getDeep(this.tree4, r, c)
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
}

export { BIT1, BIT2, BIT3, BIT4 }
