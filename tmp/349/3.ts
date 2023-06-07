// 给你一个整数 n 和一个下标从 0 开始的 二维数组 queries ，其中 queries[i] = [typei, indexi, vali] 。

// 一开始，给你一个下标从 0 开始的 n x n 矩阵，所有元素均为 0 。每一个查询，你需要执行以下操作之一：

// 如果 typei == 0 ，将第 indexi 行的元素全部修改为 vali ，覆盖任何之前的值。
// 如果 typei == 1 ，将第 indexi 列的元素全部修改为 vali ，覆盖任何之前的值。
// 请你执行完所有查询以后，返回矩阵中所有整数的和。

// 维护每行每列剩余的个数
function matrixSumQueries(n: number, queries: number[][]): number {
  const visitedRow = new Uint8Array(n)
  const visitedCol = new Uint8Array(n)
  const rowRemain = new Uint16Array(n)
  const colRemain = new Uint16Array(n)
  for (let i = 0; i < n; i++) {
    rowRemain[i] = n
    colRemain[i] = n
  }

  let res = 0
  for (let i = queries.length - 1; i >= 0; i--) {
    const [type, rowOrCol, val] = queries[i]
    // 修改行
    if (type === 0) {
      if (!visitedRow[rowOrCol]) {
        visitedRow[rowOrCol] = 1
        res += rowRemain[rowOrCol] * val
        rowRemain[rowOrCol] = 0
        for (let j = 0; j < n; j++) colRemain[j]--
      }
    } else if (!visitedCol[rowOrCol]) {
      visitedCol[rowOrCol] = 1
      res += colRemain[rowOrCol] * val
      colRemain[rowOrCol] = 0
      for (let j = 0; j < n; j++) rowRemain[j]--
    }
  }

  return res
}

export {}

// console.log(
//   matrixSumQueries(8, [
//     [0, 6, 30094],
//     [0, 7, 99382],
//     [1, 2, 18599],
//     [1, 3, 49292],
//     [1, 0, 81549],
//     [1, 1, 38280],
//     [0, 0, 19405],
//     [0, 4, 30065],
//     [1, 4, 60826],
//     [1, 5, 9241],
//     [0, 5, 33729],
//     [0, 1, 41456],
//     [0, 2, 62692],
//     [0, 3, 30807],
//     [1, 7, 70613],
//     [1, 6, 9506],
//     [0, 5, 39344],
//     [1, 0, 44658],
//     [1, 1, 56485],
//     [1, 2, 48112],
//     [0, 6, 43384]
//   ])
// )

const n = 1e4
const qs = Array.from({ length: 5e4 }, () => [0, 0, 0])

console.time('matrixSumQueries')
matrixSumQueries(n, qs)
console.timeEnd('matrixSumQueries')

// !要开1e8很大的数组时,最好是在全局开一个pool，然后每次在函数开头memset一下需要使用的部分
const visited = new Uint8Array(1e8)
const visitedRow = new Uint8Array(1e4)
const visitedCol = new Uint8Array(1e4)
function matrixSumQueries2(n: number, queries: number[][]): number {
  for (let i = 0; i < n * n; i++) visited[i] = 0
  for (let i = 0; i < n; i++) visitedRow[i] = 0
  for (let i = 0; i < n; i++) visitedCol[i] = 0

  let res = 0
  for (let i = queries.length - 1; ~i; i--) {
    const [type, rowOrCol, val] = queries[i]
    // 修改行
    if (type === 0) {
      if (visitedRow[rowOrCol]) continue
      visitedRow[rowOrCol] = 1
      for (let j = 0; j < n; j++) {
        const pos = rowOrCol * n + j
        if (!visited[pos]) {
          visited[pos] = 1
          res += val
        }
      }
    } else {
      if (visitedCol[rowOrCol]) continue
      visitedCol[rowOrCol] = 1
      for (let j = 0; j < n; j++) {
        const pos = j * n + rowOrCol
        if (!visited[pos]) {
          visited[pos] = 1
          res += val
        }
      }
    }
  }

  return res
}
