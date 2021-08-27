// 给定一个非负整数 numRows，生成「杨辉三角」的前 numRows 行。
/**
 * @param {number} numRows
 * @return {number[][]}
 */
var generate = function (numRows: number): number[][] {
  const res = Array.from<number, number[]>({ length: numRows }, (_, i) => Array(i + 1).fill(1))
  for (let i = 0; i < numRows; i++) {
    for (let j = 1; j < i; j++) {
      res[i][j] = res[i - 1][j - 1] + res[i - 1][j]
    }
  }
  return res
}
