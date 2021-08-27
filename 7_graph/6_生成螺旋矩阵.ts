/**
 * @param {number} n
 * @return {number[][]}
 * 生成一个包含 1 到 n2 所有元素，且元素按顺时针顺序螺旋排列的 n x n 正方形矩阵 matrix 。
 */
const generateMatrix = function (n: number): number[][] {
  const res = Array.from<number, number[]>({ length: n }, () => Array(n).fill(Infinity))
  let i = 0,
    j = 0,
    di = 0,
    dj = 1
  let cur = 1
  while (cur <= n ** 2) {
    res[i][j] = cur
    // 拐弯的时候判断一下是不是已经走过了
    // j + dj + col 还要加col是因为%计算的是商而不是mod 需要变成正数
    if (res[(i + di + n) % n][(j + dj + n) % n] < cur) {
      // 逆时针旋转90度
      ;[di, dj] = [dj, -di]
    }
    cur++
    i += di
    j += dj
  }

  return res
}

console.log(generateMatrix(3))
