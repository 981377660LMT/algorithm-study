// 给你一个 m 行 n 列的矩阵 matrix ，请按照 顺时针螺旋顺序 ，返回矩阵中的所有元素。
// 将已经走过的地方置0，然后拐弯的时候判断一下是不是已经走过了，如果走过了就计算一下新的方向：
const spiralOrder = (matrix: number[][]): number[] => {
  const res: number[] = []
  const row = matrix.length
  const col = matrix[0].length
  console.log(row * col)
  let i = 0,
    j = 0,
    di = 0,
    dj = 1
  for (let step = 0; step < row * col; step++) {
    console.log(i, j)
    res.push(matrix[i][j])
    matrix[i][j] = 0
    // 拐弯的时候判断一下是不是已经走过了
    // j + dj + col 还要加col是因为%计算的是商而不是mod 需要变成正数
    if (matrix[(i + di + row) % row][(j + dj + col) % col] === 0) {
      // 逆时针旋转90度
      ;[di, dj] = [dj, -di]
    }
    i += di
    j += dj
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
