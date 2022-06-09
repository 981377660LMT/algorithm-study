import { bisectRight } from '../../../9_排序和搜索/二分/bisect'

class Solution {
  private readonly rectangels: number[][]
  private readonly preSum: Uint32Array

  /**
   *
   * @param rects
   * 矩形无重复
   * 我们用 w[i] 表示第 i 个矩形 rect[i] 中整数点的数目，
   * 那么我们的随机算法需要使得每个矩形被选到的概率与 w[i] 成正比
   */
  constructor(rects: number[][]) {
    this.rectangels = rects
    this.preSum = new Uint32Array(rects.length + 1)
    for (let i = 0; i < rects.length; i++) {
      const [x1, y1, x2, y2] = rects[i]
      const area = (x2 - x1 + 1) * (y2 - y1 + 1)
      this.preSum[i + 1] = this.preSum[i] + area
    }
  }

  /**
   * 写一个函数 pick 随机均匀地选取矩形覆盖的空间中的整数点。
   */
  pick(): number[] {
    const rand = this.randint(0, this.preSum[this.preSum.length - 1] - 1)
    const pos = bisectRight(this.preSum, rand) - 1
    const [x1, y1, x2] = this.rectangels[pos]
    const offset = rand - this.preSum[pos]

    const row = ~~(offset / (x2 - x1 + 1))
    const col = offset % (x2 - x1 + 1)
    return [x1 + col, y1 + row]
  }

  private randint(start: number, end: number) {
    if (start > end) throw new Error('invalid interval')
    const diff = end - start
    return Math.floor((diff + 1) * Math.random()) + start
  }
}

// const solution = new Solution([
//   [-2, -2, -1, -1],
//   [1, 0, 3, 0],
// ])

if (require.main === module) {
  const solution = new Solution([
    [-2, -2, 1, 1],
    [2, 2, 4, 6],
  ])
  console.log(solution.pick())
  console.log(solution.pick())
  console.log(solution.pick())
  console.log(solution.pick())
  console.log(solution.pick())
  console.log(solution.pick())
  console.log(solution.pick())
  console.log(solution.pick())
}

export {}
