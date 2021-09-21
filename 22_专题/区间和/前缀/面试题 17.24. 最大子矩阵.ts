// 找出元素总和最大的子矩阵。
// 返回一个数组 [r1, c1, r2, c2]，
// 其中 r1, c1 分别代表子矩阵左上角的行号和列号，
// r2, c2 分别代表右下角的行号和列号。
// @summary
// 定三移一，降级为最大和的子数组问题
function getMaxMatrix(matrix: number[][]): number[] {
  const m = matrix.length
  const n = matrix[0].length
  const pre = Array.from<any, number[]>({ length: m + 1 }, () => Array(n + 1).fill(0))
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      pre[i][j] = matrix[i - 1][j - 1] + pre[i - 1][j] + pre[i][j - 1] - pre[i - 1][j - 1]
    }
  }

  let res = [0, 0, 0, 0]
  let globalMax = -Infinity
  // 固定上下
  for (let top = 1; top <= m; top++) {
    for (let bottom = top; bottom <= m; bottom++) {
      // 遍历子矩阵的右边界，降级为最大和的子数组问题
      let left = 0
      let sum = 0
      for (let right = 1; right <= n; right++) {
        // 注意这个表达式
        const num =
          pre[bottom][right] -
          pre[bottom][right - 1] -
          (pre[top - 1][right] - pre[top - 1][right - 1])

        if (sum < 0) {
          sum = num
          left = right - 1
        } else {
          sum += num
        }

        console.log(num, sum, 888, top, bottom, right)
        // if (sum < 0) {
        //   right + 1 <= n && (left = right)
        // }

        if (sum > globalMax) {
          globalMax = sum
          res[0] = top - 1
          res[1] = left
          res[2] = bottom - 1
          res[3] = right - 1
        }
      }
    }
  }

  return res
}

// console.log(
//   getMaxMatrix([
//     [1, 0],
//     [0, 11],
//   ])
// )
console.log(getMaxMatrix([[-1]]))
console.log(
  getMaxMatrix([
    [1, -3],
    [-4, 9],
  ])
)
// [1,1,1,1]
export {}
