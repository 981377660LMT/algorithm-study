// https://github.com/tdzl2003/leetcode_live/blob/e7434ae4a5e849608f171531edea0f06635e018b/leetcode/others/1476.md#L1
// 采用动态四叉树维护，当一个区间被完整设置成某个值时，可以砍掉这个区间的所有子树
// !每次操作的复杂度为O(logrow+logcol) 为加法而不是乘法
//
// !1.节点的分配和回收使用对象池：这种allocNode和freeNode的写法很有用，减少了分配对象的压力.
// !2.奇技淫巧：利用属性值的不同范围表示不同含义，复用了属性，减少了空间占用.

type SNode = {
  leftTop: number
  /**
   * !当`leftBottom`为0时,`leftTop`表示当前结点的值.
   */
  leftBottom: number
  rightTop: number
  rightBottom: number
}

class SubrectangleQueries<S> {
  private readonly _row: number
  private readonly _col: number
  private readonly _origin: S[][]
  private readonly _nodes: SNode[] = [{ leftTop: 0, leftBottom: 0, rightTop: 0, rightBottom: 0 }]
  private readonly _recycles: number[] = [] // 可回收的结点编号
  private readonly _idToValue: Map<S, number> = new Map() // 从1开始编号(0表示未分配)
  private readonly _valueToId: S[] = [null as S] // 从1开始编号

  constructor(rectangle: S[][]) {
    this._row = rectangle.length
    this._col = rectangle[0].length
    this._origin = rectangle
  }

  /**
   * 将子矩形左上角`[row1,col1]`到右下角`[row2,col2]`之间所有元素的值都修改为`newValue`.
   */
  updateSubrectangle(row1: number, col1: number, row2: number, col2: number, newValue: S): void {
    const id = this._getIdByValue(newValue)
    this._update(0, 0, 0, this._row, this._col, row1, col1, row2 + 1, col2 + 1, id)
  }

  getValue(row: number, col: number): S {
    const id = this._query(0, 0, 0, this._row, this._col, row, col)
    if (id) return this._getValueById(id)
    return this._origin[row][col]
  }

  private _allocNode(value = 0): number {
    if (this._recycles.length) {
      const id = this._recycles.pop()!
      const node = this._nodes[id]
      if (node.leftBottom) {
        this._recycles.push(node.leftTop)
        this._recycles.push(node.leftBottom)
        this._recycles.push(node.rightTop)
        this._recycles.push(node.rightBottom)
        node.leftBottom = 0
      }
      node.leftTop = value
      return id
    }
    const id = this._nodes.length
    this._nodes.push({ leftTop: value, leftBottom: 0, rightTop: 0, rightBottom: 0 })
    return id
  }

  private _freeNode(id: number): void {
    this._recycles.push(id)
  }

  private _update(
    nodeId: number,
    l: number,
    t: number,
    r: number,
    b: number,
    row1: number,
    col1: number,
    row2: number,
    col2: number,
    newValue: number
  ): void {
    if (l >= row1 && t >= col1 && r <= row2 && b <= col2) {
      const node = this._nodes[nodeId]
      if (node.leftBottom) {
        this._freeNode(node.leftTop)
        this._freeNode(node.leftBottom)
        this._freeNode(node.rightTop)
        this._freeNode(node.rightBottom)
        node.leftBottom = 0
      }
      node.leftTop = newValue
      return
    }

    const xmid = Math.floor((l + r) / 2)
    const ymid = Math.floor((t + b) / 2)

    // copy
    const o = this._nodes[nodeId]
    const node = {
      leftTop: o.leftTop,
      leftBottom: o.leftBottom,
      rightTop: o.rightTop,
      rightBottom: o.rightBottom
    }
    if (!node.leftBottom) {
      const oldValue = node.leftTop
      node.leftTop = this._allocNode(oldValue)
      node.leftBottom = this._allocNode(oldValue)
      node.rightTop = this._allocNode(oldValue)
      node.rightBottom = this._allocNode(oldValue)
      this._nodes[nodeId] = node
    }

    if (row1 < xmid) {
      if (col1 < ymid) {
        this._update(node.leftTop, l, t, xmid, ymid, row1, col1, row2, col2, newValue)
      }
      if (col2 > ymid) {
        this._update(node.leftBottom, l, ymid, xmid, b, row1, col1, row2, col2, newValue)
      }
    }
    if (row2 > xmid) {
      if (col1 < ymid) {
        this._update(node.rightTop, xmid, t, r, ymid, row1, col1, row2, col2, newValue)
      }
      if (col2 > ymid) {
        this._update(node.rightBottom, xmid, ymid, r, b, row1, col1, row2, col2, newValue)
      }
    }
  }

  private _query(
    nodeId: number,
    l: number,
    t: number,
    r: number,
    b: number,
    row: number,
    col: number
  ): number {
    const node = this._nodes[nodeId]
    if (!node.leftBottom) return node.leftTop
    const xmid = Math.floor((l + r) / 2)
    const ymid = Math.floor((t + b) / 2)
    if (row < xmid) {
      if (col < ymid) return this._query(node.leftTop, l, t, xmid, ymid, row, col)
      return this._query(node.leftBottom, l, ymid, xmid, b, row, col)
    }

    if (col < ymid) return this._query(node.rightTop, xmid, t, r, ymid, row, col)
    return this._query(node.rightBottom, xmid, ymid, r, b, row, col)
  }

  private _getIdByValue(v: S): number {
    const res = this._idToValue.get(v)
    if (res !== void 0) return res
    const newId = this._idToValue.size + 1 // 0 is reserved for null
    this._idToValue.set(v, newId)
    this._valueToId.push(v)
    return newId
  }

  private _getValueById(id: number): S {
    return this._valueToId[id]
  }
}

export {}

if (require.main === module) {
  const s = new SubrectangleQueries<number>([
    [1, 2, 1],
    [4, 3, 4],
    [3, 2, 1],
    [1, 1, 1]
  ])
  s.updateSubrectangle(0, 0, 2, 2, 0)
  console.log(s.getValue(0, 0))

  const ROW = 500
  const COL = 500
  const tree = new SubrectangleQueries<number>(
    Array(ROW)
      .fill(0)
      .map(() => Array(COL).fill(0))
  )
  console.time('update')
  for (let i = 0; i < 1e5; ++i) {
    tree.updateSubrectangle(0, 0, ROW - 1, COL - 1, i)
    tree.getValue(0, 0)
  }
  console.timeEnd('update') // update: 16.084ms
}
