/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { SortedListFast } from '../../../22_专题/离线查询/根号分治/SortedList/SortedListFast'

const INF = 2e9 // !超过int32使用2e15

/**
 * !`包含原点的`矩形合并/矩形面积并.
 */
class IncrementalRectangleUnionRange {
  private readonly _ur = new IncrementalRectangleUnion()
  private readonly _ul = new IncrementalRectangleUnion()
  private readonly _dr = new IncrementalRectangleUnion()
  private readonly _dl = new IncrementalRectangleUnion()

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
    return this._ur.query() + this._ul.query() + this._dr.query() + this._dl.query()
  }
}

/**
 * !`包含原点的`矩形合并/矩形面积并.
 */
class IncrementalRectangleUnion {
  private _sum = 0
  private readonly _sl = new SortedListFast(
    [
      { start: 0, end: INF },
      { start: INF, end: 0 }
    ],
    (a, b) => a.start - b.start
  )

  /**
   * Add [0, x] * [0, y].
   * x >= 0.
   * y >= 0.
   */
  add(x: number, y: number): void {
    let pos = this._sl.bisectLeft({ start: x, end: -INF })
    const item = this._sl.at(pos)!
    if (item.end >= y) return
    const nextY = item.end
    pos--
    let pre = this._sl.at(pos)!
    while (pre.end <= y) {
      const x1 = pre.start
      const y1 = pre.end
      this._sl.pop(pos)
      pos--
      pre = this._sl.at(pos)!
      this._sum -= (x1 - pre.start) * (y1 - nextY)
    }
    this._sum += (x - this._sl.at(pos)!.start) * (y - nextY)
    pos = this._sl.bisectLeft({ start: x, end: -INF })
    if (this._sl.at(pos)!.start === x) this._sl.pop(pos)
    this._sl.add({ start: x, end: y })
  }

  query(): number {
    return this._sum
  }
}

/**
 * 求出所有左下角为原点的立方体的体积并.
 * @param points [x, y, z] 每个点的坐标，表示一个[0, x] * [0, y] * [0, z]的立方体.
 */
function cuboidUnionVolumn(points: [x: number, y: number, z: number][]): number {
  points = points.slice().sort((a, b) => b[2] - a[2])
  let preZ = INF
  let res = 0
  let area = 0
  const I = new IncrementalRectangleUnion()
  for (let i = 0; i < points.length; i++) {
    const { 0: x, 1: y, 2: z } = points[i]
    res += (preZ - z) * area
    I.add(x, y)
    area = I.query()
    preZ = z
  }
  res += preZ * area
  return res
}

if (require.main === module) {
  const R = new IncrementalRectangleUnionRange()
  R.add(-2, 1, -2, 1)
  console.log(R.query())
  R.add(-1, 2, -1, 2)
  console.log(R.query())

  console.log(
    cuboidUnionVolumn([
      [1, 1, 1],
      [2, 2, 2],
      [3, 3, 3]
    ])
  )

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
  //   const [x1,y1, x2, y2] = input().split(' ').map(Number)
  //   R.add(x1, x2, y1, y2)
  //   const cur = R.query()
  //   console.log(cur - pre)
  //   pre = cur
  // }
}

export { IncrementalRectangleUnion, IncrementalRectangleUnionRange }
