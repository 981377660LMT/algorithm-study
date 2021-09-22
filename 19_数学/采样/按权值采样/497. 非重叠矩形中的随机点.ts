import { bisectLeft } from '../../../9_排序和搜索/二分api/7_二分搜索寻找最左插入位置'

class Solution {
  private rectangels: number[][]
  private pre: number[]
  /**
   *
   * @param rects
   * 矩形无重复
   * 我们用 w[i] 表示第 i 个矩形 rect[i] 中整数点的数目，
   * 那么我们的随机算法需要使得每个矩形被选到的概率与 w[i] 成正比
   */
  constructor(rects: number[][]) {
    this.rectangels = rects
    this.pre = Array(rects.length).fill(0)
    for (let i = 0; i < rects.length; i++) {
      const [x1, y1, x2, y2] = rects[i]
      const points = (x2 - x1 + 1) * (y2 - y1 + 1)
      this.pre[i] = (this.pre[i - 1] || 0) + points
    }
  }

  /**
   * 写一个函数 pick 随机均匀地选取矩形覆盖的空间中的整数点。
   */
  pick(): number[] {
    // 第几个点
    const rand = this.randint(1, this.pre[this.pre.length - 1])
    const idOfRectangle = bisectLeft(this.pre, rand)
    const [x1, y1, x2, y2] = this.rectangels[idOfRectangle]
    const offset = this.pre[idOfRectangle] - rand
    // console.log(offset, rand, idOfRectangle)
    // 矩形坐标公式
    const rowDiff = ~~(offset / (x2 - x1 + 1))
    const colDifff = offset % (x2 - x1 + 1)
    return [x1 + colDifff, y2 - rowDiff]
  }

  private randint(start: number, end: number) {
    if (start > end) throw new Error('invalid interval')
    const amplitude = end - start
    return Math.floor((amplitude + 1) * Math.random()) + start
  }
}

// const solution = new Solution([
//   [-2, -2, -1, -1],
//   [1, 0, 3, 0],
// ])
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
export {}
