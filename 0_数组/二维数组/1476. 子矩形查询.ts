class SubrectangleQueries {
  private matrix: number[][]
  private history: number[][]

  constructor(rectangle: number[][]) {
    this.matrix = rectangle
    this.history = []
  }

  /**
   *
   * @param row1
   * @param col1
   * @param row2
   * @param col2
   * @param newValue
   * 用 newValue 更新以 (row1,col1) 为左上角且以 (row2,col2) 为右下角的子矩形。
   */
  updateSubrectangle(
    row1: number,
    col1: number,
    row2: number,
    col2: number,
    newValue: number
  ): void {
    this.history.push([row1, col1, row2, col2, newValue])
  }

  /**
   *
   * @param row
   * @param col
   * 返回矩形中坐标 (row,col) 的当前值
   */
  getValue(row: number, col: number): number {
    for (let i = this.history.length - 1; i >= 0; i--) {
      const [row1, col1, row2, col2, value] = this.history[i]
      if (row >= row1 && row <= row2 && col >= col1 && col <= col2) return value
    }
    return this.matrix[row][col]
  }
}

export {}

// 读多写少 暴力更新即可
// 读少写多 history数组历史查询
// 本题假设读少写多
