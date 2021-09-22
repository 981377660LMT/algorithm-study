class Solution {
  constructor(public radius: number, public x_center: number, public y_center: number) {}

  /**
   * 写一个在圆中产生均匀随机点的函数 randPoint
   * @summary
   * 考虑一个半径为1的圆，点落在半径为x的圆内概率为x^2，
   * 所以为了使点均匀分布，x=sqrt(U),这样概率就是U
   */
  randPoint(): number[] {
    const r = this.radius * Math.sqrt(Math.random())
    const degree = Math.random() * Math.PI * 2
    const x = this.x_center + r * Math.cos(degree)
    const y = this.y_center + r * Math.sin(degree)
    return [x, y]
  }
}

const s = new Solution(1, 0, 0)

export {}
