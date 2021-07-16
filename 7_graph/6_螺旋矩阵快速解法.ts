// 给你一个 m 行 n 列的矩阵 matrix ，请按照 顺时针螺旋顺序 ，返回矩阵中的所有元素。

// 三步:push,pop,反转,push,pop,反转...
const spiralOrder = (matrix: number[][]) => {
  const res: number[] = []

  while (matrix.length) {
    res.push(...matrix.shift()!)
    for (const m of matrix) {
      res.push(m.pop()!)
      m.reverse()
    }
    matrix.reverse()
  }
  return res
}

console.log(
  spiralOrder([
    [1, 2, 3, 4],
    [5, 6, 7, 8],
    [9, 10, 11, 12],
  ])
)
// 输出：[1,2,3,4,8,12,11,10,9,5,6,7]
export {}
