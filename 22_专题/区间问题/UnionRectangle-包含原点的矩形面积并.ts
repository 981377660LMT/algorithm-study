/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { SortedList } from '../离线查询/根号分治/SortedList/_SortedList'

const INF = 2e15

/**
 * !`包含原点的`矩形合并/矩形面积并.
 */
class UnionRectangleRange {
  private readonly _ur = new UnionRectangle()
  private readonly _ul = new UnionRectangle()
  private readonly _dr = new UnionRectangle()
  private readonly _dl = new UnionRectangle()

  /**
   * Add [x1, x2] * [y1, y2].
   * x1 <= 0 <= x2.
   * y1 <= 0 <= y2.
   */
  add(x1: number, x2: number, y1: number, y2: number): void {
    this._ur.add(x2, y2)
    this._ul.add(-x1, y2)
    this._dr.add(x2, -y1)
    this._dl.add(-x1, -y1)
  }

  query(): number {
    return this._ur.sum + this._ul.sum + this._dr.sum + this._dl.sum
  }
}

/**
 * !`包含原点的`矩形合并/矩形面积并.
 */
class UnionRectangle {
  sum = 0
  private readonly _sl = new SortedList(
    [
      [0, INF],
      [INF, 0]
    ],
    (a, b) => a[0] - b[0]
  )

  /**
   * Add [0, x] * [0, y].
   * x >= 0.
   * y >= 0.
   */
  add(x: number, y: number): void {
    let pos = this._sl.bisectLeft([x, -INF])
    const item = this._sl.at(pos)!
    if (item[1] >= y) return
    const nextY = item[1]
    pos--
    let pre = this._sl.at(pos)!
    while (pre[1] <= y) {
      const x1 = pre[0]
      const y1 = pre[1]
      this._sl.pop(pos)
      pos--
      pre = this._sl.at(pos)!
      this.sum -= (x1 - pre[0]) * (y1 - nextY)
    }
    this.sum += (x - this._sl.at(pos)![0]) * (y - nextY)
    pos = this._sl.bisectLeft([x, -INF])
    if (this._sl.at(pos)![0] === x) this._sl.pop(pos)
    this._sl.add([x, y])
  }

  query(): number {
    return this.sum
  }
}

if (require.main === module) {
  const R = new UnionRectangleRange()
  R.add(-2, 1, -2, 1)
  console.log(R.query())

  // https://yukicoder.me/submissions/857898
  // import * as fs from 'fs'
  // import { resolve } from 'path'

  // function useInput(path?: string) {
  //   let data: string
  //   if (path) {
  //     data = fs.readFileSync(resolve(__dirname, path), 'utf8')
  //   } else {
  //     data = fs.readFileSync(process.stdin.fd, 'utf8')
  //   }

  //   const lines = data.split(/\r\n|\r|\n/)
  //   let lineId = 0
  //   const input = (): string => lines[lineId++]

  //   return {
  //     input
  //   }
  // }

  // const { input } = useInput()
  // const n = Number(input())
  // const R = new UnionRectangleRange()
  // let pre = 0
  // for (let i = 0; i < n; i++) {
  //   const [x1, x2, y1, y2] = input().split(' ').map(Number)
  //   R.add(x1, x2, y1, y2)
  //   const cur = R.query()
  //   console.log(cur - pre)
  //   pre = cur
  // }
}

export { UnionRectangleRange, UnionRectangle }
