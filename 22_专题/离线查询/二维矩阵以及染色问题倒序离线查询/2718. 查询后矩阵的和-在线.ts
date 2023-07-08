// https://leetcode.cn/problems/sum-of-matrix-after-queries/solution/shu-zhuang-shu-zu-shi-jian-chuo-zai-xian-eqgo/
// 要修改一行时，用树状数组维护，上次修改这一行的时间，一直到现在，这段时间内哪些列进行了修改，就可以实现在线查询
// TODO

import { BITArray } from '../../../6_tree/树状数组/经典题/BIT'

class MatrixSumQueries {
  private readonly _sum: BITArray // 区间和
  private readonly _times: BITArray // 修改次数

  constructor(row: number, col: number, maxQueries: number) {}

  fillRow(row: number, val: number): void {}

  fillCol(col: number, val: number): void {}

  sum(): number {}
}

// https://leetcode.cn/problems/sum-of-matrix-after-queries/
// 2718. 查询后矩阵的和
function matrixSumQueries(n: number, queries: number[][]): number {
  const M = new MatrixSumQueries(n, n, queries.length)
  queries.forEach(([type, index, val]) => {
    if (!type) {
      M.fillRow(index, val)
    } else {
      M.fillCol(index, val)
    }
  })

  return M.sum()
}

export {}
