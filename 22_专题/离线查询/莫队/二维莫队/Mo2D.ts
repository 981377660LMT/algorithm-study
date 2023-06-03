/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-inner-declarations */
// 二维莫队
// 时间复杂度O(row*col*(q**0.75))
// !一般 n,m<=200 q<=1e5

import assert from 'assert'

/**
 * 二维莫队.
 */
class Mo2D {
  private readonly _blockId: Uint32Array
  private readonly _queries: [x1: number, y1: number, x2: number, y2: number, qid: number][] = []

  constructor(row: number, col: number, q: number) {
    let chunkSize = ((row * col) ** 0.5 / Math.max(q ** 0.25, 1)) | 0
    if (chunkSize < 1) chunkSize = 1
    const blockId = new Uint32Array(Math.max(row, col) + 5)
    for (let i = 0; i < blockId.length; ++i) blockId[i] = (i / chunkSize) | 0
    this._blockId = blockId
  }

  /**
   * 添加查询矩形区域`[x1, x2) * [y1, y2)`.
   * 0 <= x1 < x2 <= row, 0 <= y1 < y2 <= col.
   */
  addQuery(x1: number, x2: number, y1: number, y2: number): void {
    this._queries.push([x1, y1, x2 - 1, y2 - 1, this._queries.length])
  }

  /**
   * 返回每个查询的结果.
   * @param addRow 将新的行添加到窗口. dir: 1 表示row变大，-1 表示row变小. 对应列的范围是[col1, col2).
   * @param addCol 将新的列添加到窗口. dir: 1 表示col变大，-1 表示col变小. 对应行的范围是[row1, row2).
   * @param removeRow 将行从窗口移除. dir: 1 表示row变大，-1 表示row变小. 对应列的范围是[col1, col2).
   * @param removeCol 将列从窗口移除. dir: 1 表示col变大，-1 表示col变小. 对应行的范围是[row1, row2).
   * @param query 查询窗口内的数据.
   */
  run(
    addRow: (row: number, dir: -1 | 1, col1: number, col2: number) => void,
    addCol: (col: number, dir: -1 | 1, row1: number, row2: number) => void,
    removeRow: (row: number, dir: -1 | 1, col1: number, col2: number) => void,
    removeCol: (col: number, dir: -1 | 1, row1: number, row2: number) => void,
    query: (qid: number) => void
  ): void {
    const bid = this._blockId
    const queries = this._queries
    queries.sort((q1, q2) => {
      const bid1 = bid[q1[0]]
      const bid2 = bid[q2[0]]
      if (bid1 !== bid2) return bid1 - bid2
      const bid3 = bid[q1[1]]
      const bid4 = bid[q2[1]]
      if (bid3 !== bid4) return bid1 & 1 ? q1[1] - q2[1] : q2[1] - q1[1]
      const bid5 = bid[q1[3]]
      const bid6 = bid[q2[3]]
      if (bid5 !== bid6) return bid3 & 1 ? q1[3] - q2[3] : q2[3] - q1[3]
      return bid5 & 1 ? q1[2] - q2[2] : q2[2] - q1[2]
    })

    let x1 = 0
    let y1 = 0
    let x2 = -1
    let y2 = -1
    for (let i = 0; i < queries.length; ++i) {
      const [qx1, qy1, qx2, qy2, qid] = queries[i]
      while (x1 > qx1) {
        x1--
        addRow(x1, -1, y1, y2 + 1)
      }
      while (x2 < qx2) {
        x2++
        addRow(x2, 1, y1, y2 + 1)
      }
      while (y1 > qy1) {
        y1--
        addCol(y1, -1, x1, x2 + 1)
      }
      while (y2 < qy2) {
        y2++
        addCol(y2, 1, x1, x2 + 1)
      }
      while (x1 < qx1) {
        removeRow(x1, 1, y1, y2 + 1)
        x1++
      }
      while (x2 > qx2) {
        removeRow(x2, -1, y1, y2 + 1)
        x2--
      }
      while (y1 < qy1) {
        removeCol(y1, 1, x1, x2 + 1)
        y1++
      }
      while (y2 > qy2) {
        removeCol(y2, -1, x1, x2 + 1)
        y2--
      }
      query(qid)
    }
  }
}

export { Mo2D }

if (require.main === module) {
  // https://hydro.ac/d/bzoj/p/2639
  const grid = [
    [1, 3, 2, 1],
    [1, 3, 2, 4],
    [1, 2, 3, 4]
  ]

  const queries = [
    [0, 2, 0, 2],
    [0, 2, 0, 1],
    [0, 3, 0, 4],
    [0, 1, 0, 1],
    [1, 3, 1, 3],
    [1, 3, 1, 4],
    [0, 3, 0, 3],
    [1, 3, 3, 4]
  ]

  assert.deepStrictEqual(solve(grid, queries), [8, 4, 38, 1, 8, 12, 27, 4])
  console.log('ok')

  function solve(
    grid: number[][],
    queries: [x1: number, x2: number, y1: number, y2: number][] | number[][]
  ): number[] {
    const ROW = grid.length
    const COL = grid[0].length
    const q = queries.length

    const pool = new Map<unknown, number>()
    const id = (o: unknown): number => {
      if (!pool.has(o)) pool.set(o, pool.size)
      return pool.get(o)!
    }

    const newGrid = new Uint32Array(ROW * COL)
    for (let i = 0; i < ROW; ++i) {
      for (let j = 0; j < COL; ++j) {
        newGrid[i * COL + j] = id(grid[i][j])
      }
    }

    const mo2d = new Mo2D(ROW, COL, q)
    queries.forEach(([x1, x2, y1, y2]) => {
      mo2d.addQuery(x1, x2, y1, y2)
    })

    let cur = 0
    const counter = new Uint32Array(pool.size)
    const res = Array(queries.length)
    mo2d.run(addRow, addCol, removeRow, removeCol, query)
    return res

    function addRow(row: number, dir: -1 | 1, col1: number, col2: number): void {
      for (let c = col1; c < col2; c++) {
        const v = newGrid[row * COL + c]
        const pre = counter[v]
        cur -= pre * pre
        counter[v]++
        cur += (pre + 1) * (pre + 1)
      }
    }

    function addCol(col: number, dir: -1 | 1, row1: number, row2: number): void {
      for (let r = row1; r < row2; r++) {
        const v = newGrid[r * COL + col]
        const pre = counter[v]
        cur -= pre * pre
        counter[v]++
        cur += (pre + 1) * (pre + 1)
      }
    }

    function removeRow(row: number, dir: -1 | 1, col1: number, col2: number): void {
      for (let c = col1; c < col2; c++) {
        const v = newGrid[row * COL + c]
        const pre = counter[v]
        cur -= pre * pre
        counter[v]--
        cur += (pre - 1) * (pre - 1)
      }
    }

    function removeCol(col: number, dir: -1 | 1, row1: number, row2: number): void {
      for (let r = row1; r < row2; r++) {
        const v = newGrid[r * COL + col]
        const pre = counter[v]
        cur -= pre * pre
        counter[v]--
        cur += (pre - 1) * (pre - 1)
      }
    }

    function query(qid: number): void {
      res[qid] = cur
    }
  }
}
