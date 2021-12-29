import { randint } from '../randint'

// 尽量最少调用随机函数 Math.random()，并且优化时间和空间复杂度。
class Solution {
  private total: number
  private rows: number
  private cols: number
  private black: Map<number, number>

  constructor(m: number, n: number) {
    this.total = m * n
    this.rows = m
    this.cols = n
    this.black = new Map<number, number>()
  }

  /**
   * 均匀随机的将矩阵中的 0 变为 1，
   * 并返回该值的位置下标 [row_id,col_id]
   * @description
   * 一次调用random方法
   * 每次flip后都将flip的值记录到黑名单中，并将其映射到最后一个数上
   */
  flip(): number[] {
    if (this.total <= 0) return []
    this.total--
    const choose = randint(0, this.total)
    const white = this.getWhite(choose)
    // choose 映射到另一个白名单的数
    this.black.set(choose, this.getWhite(this.total))
    return [~~(white / this.cols), white % this.cols]
  }

  /**
   * 将所有的值都重新置为 0
   */
  reset(): void {
    this.total = this.rows * this.cols
    this.black.clear()
  }

  private getWhite(num: number): number {
    return this.black.get(num) || num
  }
}

const s = new Solution(2, 3)
console.log(s.flip())
s.reset()
console.log(s.flip())
export {}
