import { SegmentTree2DRangeTreeRangeUpdatePointGet } from '../template/SegmentTree2DRangeTreeRangeUpdatePointGet'

// https://leetcode.cn/problems/increment-submatrices-by-one/
function rangeAddQueries(n: number, queries: number[][]): number[][] {
  const seg2d = new SegmentTree2DRangeTreeRangeUpdatePointGet(
    n,
    n,
    n => new NaiveTree(n),
    (a, b) => a + b
  )

  queries.forEach(([x1, y1, x2, y2]) => {
    seg2d.update(x1, x2 + 1, y1, y2 + 1, 1)
  })

  const res = Array(n)
  for (let i = 0; i < n; i++) {
    res[i] = Array(n)
    for (let j = 0; j < n; j++) {
      res[i][j] = seg2d.get(i, j)
    }
  }

  return res
}

class NaiveTree {
  private readonly _nums: number[]

  constructor(n: number) {
    this._nums = Array(n).fill(0)
  }

  update(start: number, end: number, delta: number): void {
    for (let i = start; i < end; i++) {
      this._nums[i] += delta
    }
  }

  get(index: number): number {
    return this._nums[index]
  }

  set(index: number, value: number): void {
    this._nums[index] = value
  }
}

if (require.main === module) {
  const seg2d = new SegmentTree2DRangeTreeRangeUpdatePointGet(
    3,
    3,
    n => new NaiveTree(n),
    (a, b) => a + b
  )

  seg2d.update(0, 2, 0, 2, 1)
  console.log(seg2d.get(1, 1))
  console.log(seg2d.get(0, 1))
  console.log(seg2d.get(0, 2))
  console.log(seg2d.get(1, 2))
  seg2d.update(1, 2, 1, 2, 1)
  console.log(seg2d.get(1, 1))
  console.log(seg2d.get(0, 1))
  console.log(seg2d.get(0, 2))
  console.log(seg2d.get(1, 2))
}
