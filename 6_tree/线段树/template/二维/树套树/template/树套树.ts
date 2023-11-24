/* eslint-disable no-inner-declarations */
/* eslint-disable eqeqeq */
/* eslint-disable max-len */

import { BITArray2 } from '../../../../../重链剖分/BIT'

/** 线段树套数据结构. */
class SegmentTreeDivideInterval<InnerTree> {
  private readonly _n: number
  private readonly _smallN: boolean
  /** 线段树中一共offset+n个节点,offset+i对应原来的第i个节点. */
  private readonly _offset: number
  private readonly _createInnerTree: () => InnerTree
  private readonly _innerTreeList: InnerTree[] = []
  private readonly _innerTreeMap: Map<number, InnerTree> = new Map()

  /**
   * @param n 第一个维度(一般是序列)的长度.
   * @param createTree 创建第二个维度(一般是线段树)的树.
   * @param smallN n较小时，会预先创建好所有的内层树; 否则会用map保存内层树，并在需要的时候创建.
   */
  constructor(n: number, createInnerTree: () => InnerTree, smallN: boolean) {
    this._n = n
    this._smallN = smallN
    this._offset = 1
    while (this._offset < n) this._offset *= 2
    this._createInnerTree = createInnerTree
    if (smallN) {
      const innerTreeList = Array<InnerTree>(this._offset + n)
      for (let i = 0; i < this._offset + n; i++) {
        innerTreeList[i] = createInnerTree()
      }
      this._innerTreeList = innerTreeList
    }
  }

  enumeratePoint(index: number, f: (tree: InnerTree) => void): void {
    if (index < 0 || index >= this._n) return
    index += this._offset
    while (index > 0) {
      f(this._getTree(index))
      index = Math.floor(index / 2)
    }
  }

  enumerateRange(start: number, end: number, f: (tree: InnerTree) => void): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    const leftSegments: InnerTree[] = []
    const rightSegments: InnerTree[] = []
    start += this._offset
    end += this._offset
    while (start < end) {
      if (start & 1) {
        leftSegments.push(this._getTree(start))
        start++
      }
      if (end & 1) {
        end--
        rightSegments.push(this._getTree(end))
      }
      start = Math.floor(start / 2)
      end = Math.floor(end / 2)
    }
    for (let i = 0; i < leftSegments.length; i++) {
      f(leftSegments[i])
    }
    for (let i = rightSegments.length - 1; i >= 0; i--) {
      f(rightSegments[i])
    }
  }

  private _getTree(segmentId: number): InnerTree {
    if (this._smallN) return this._innerTreeList[segmentId]
    const innerTree = this._innerTreeMap.get(segmentId)
    if (innerTree != undefined) return innerTree
    const newTree = this._createInnerTree()
    this._innerTreeMap.set(segmentId, newTree)
    return newTree
  }
}

/** 树状数组套数据结构. */
class FenwickTreeDivideInterval<InnerTree> {
  private readonly _n: number
  private readonly _smallN: boolean
  private readonly _createInnerTree: () => InnerTree
  private readonly _innerTreeList: InnerTree[] = []
  private readonly _innerTreeMap: Map<number, InnerTree> = new Map()

  /**
   * @param n 第一个维度(一般是序列)的长度.必须是int32范围内的整数.
   * @param createTree 创建第二个维度(一般是线段树)的树.
   * @param smallN n较小时，会预先创建好所有的内层树; 否则会用map保存内层树，并在需要的时候创建.
   */
  constructor(n: number, createInnerTree: () => InnerTree, smallN: boolean) {
    this._n = n
    this._smallN = smallN
    this._createInnerTree = createInnerTree
    if (smallN) {
      const innerTreeList = Array<InnerTree>(n)
      for (let i = 0; i < n; i++) {
        innerTreeList[i] = createInnerTree()
      }
      this._innerTreeList = innerTreeList
    }
  }

  update(index: number, f: (tree: InnerTree) => void): void {
    if (index < 0 || index >= this._n) return
    for (index++; index <= this._n; index += index & -index) {
      f(this._getTree(index - 1))
    }
  }

  queryPrefix(end: number, f: (tree: InnerTree) => void): void {
    if (end > this._n) end = this._n
    for (; end > 0; end &= end - 1) {
      f(this._getTree(end - 1))
    }
  }

  private _getTree(segmentId: number): InnerTree {
    if (this._smallN) return this._innerTreeList[segmentId]
    const innerTree = this._innerTreeMap.get(segmentId)
    if (innerTree != undefined) return innerTree
    const newTree = this._createInnerTree()
    this._innerTreeMap.set(segmentId, newTree)
    return newTree
  }
}

export { SegmentTreeDivideInterval, FenwickTreeDivideInterval, FenwickTreeDivideInterval as BITDivideInterval }

if (require.main === module) {
  testSegmentTreeDivideInterval()
  testFenwickTreeDivideInterval()

  function testSegmentTreeDivideInterval() {
    const seg = new SegmentTreeDivideInterval(10, () => 1, true)
    seg.enumerateRange(2, 9, tree => console.log(tree))
  }

  function testFenwickTreeDivideInterval() {
    const bit = new FenwickTreeDivideInterval(10, () => new BITArray2(10), true)
    bit.update(1, tree => tree.add(3, 9, 1))
    bit.update(2, tree => tree.add(3, 9, 1))
    bit.update(3, tree => tree.add(3, 10, 1))
    let res = 0
    bit.queryPrefix(4, tree => {
      res += tree.query(0, 5)
    })
    console.log(res, 'aaa')
  }

  // https://leetcode.cn/problems/increment-submatrices-by-one/description/
  function rangeAddQueries(n: number, queries: number[][]): number[][] {
    const seg = new SegmentTreeDivideInterval(n, () => new BITArray2(n), true)
    for (let qi = 0; qi < queries.length; qi++) {
      const { 0: row1, 1: col1, 2: row2, 3: col2 } = queries[qi]
      seg.enumerateRange(row1, row2 + 1, tree => {
        tree.add(col1, col2 + 1, 1)
      })
    }

    const res: number[][] = Array(n)
    for (let i = 0; i < n; i++) {
      const curRow = Array(n).fill(0)
      for (let j = 0; j < n; j++) {
        seg.enumeratePoint(i, tree => {
          curRow[j] += tree.query(j, j + 1)
        })
      }
      res[i] = curRow
    }

    return res
  }

  console.log(
    rangeAddQueries(3, [
      [1, 1, 2, 2],
      [0, 0, 1, 1]
    ])
  )
}
