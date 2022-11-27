// 1
// 1 1
// 1 2 1
// 1 3 3 1
// 1 4 6 4 1   C(4,0) C(4,1) C(4,2) C(4,3) C(4,4)
// C(n,k)=C(n,k-1)*(n-k+1)/k

// !给定一个非负整数 numRows，生成「杨辉三角」的前 numRows 行。
/**
 * @param {number} numRows
 * @return {number[][]}
 */
function generate(numRows: number): number[][] {
  const res = Array.from<number, number[]>({ length: numRows }, (_, i) => Array(i + 1).fill(1))
  for (let i = 0; i < numRows; i++) {
    for (let j = 1; j < i; j++) {
      res[i][j] = res[i - 1][j - 1] + res[i - 1][j]
    }
  }
  return res
}

// !给定一个非负索引 rowIndex，返回「杨辉三角」的第 rowIndex 行。
// 0 <= rowIndex <= 33
// 你可以优化你的算法到 O(rowIndex) 空间复杂度吗？
// 第i行的组合数
/**
 * @param {number} rowIndex
 * @return {number[]}
 */
function getRow(rowIndex: number): number[] {
  const res = Array(rowIndex + 1).fill(1)
  for (let i = 1; i <= rowIndex; i++) {
    res[i] = (res[i - 1] * (rowIndex - i + 1)) / i
  }
  return res
}

export {}
