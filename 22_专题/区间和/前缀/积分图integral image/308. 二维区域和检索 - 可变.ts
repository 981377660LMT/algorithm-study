import { BIT } from '../../../../6_tree/树状数组/BIT'

// 把二维矩阵转化为一维数组，然后使用树状数组
class NumMatrix {
  private m: number
  private n: number
  private matrix: number[][]
  private bit: BIT

  constructor(matrix: number[][]) {
    this.m = matrix.length
    this.n = matrix[0].length
    this.matrix = matrix
    this.bit = new BIT(this.m * this.n)
    for (let i = 0; i < this.m; i++) {
      for (let j = 0; j < this.n; j++) {
        const index = this.getIndex(i, j)
        this.bit.add(index, matrix[i][j])
      }
    }
  }

  update(row: number, col: number, val: number): void {
    const index = this.getIndex(row, col)
    const diff = val - this.matrix[row][col]
    this.bit.add(index, diff)
    this.matrix[row][col] = val // 不要忘了更新原来的矩阵
  }

  sumRegion(row1: number, col1: number, row2: number, col2: number): number {
    let res = 0
    for (let row = row1; row <= row2; row++) {
      const leftIndex = this.getIndex(row, col1)
      const rightIndex = this.getIndex(row, col2)
      res += this.bit.sumRange(leftIndex, rightIndex)
    }
    return res
  }

  private getIndex(row: number, col: number) {
    return this.n * row + col + 1 // 注意这里加一 让树状数组查询索引是正数
  }
}

const test = new NumMatrix([
  [2, 4],
  [-3, 5],
])
test.update(0, 1, 3)
test.update(1, 1, -3)
test.update(0, 1, 1)
console.log(test.sumRegion(0, 0, 1, 1))
// test.update
export {}
