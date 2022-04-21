/**
 *
 * @param matrix  1 <= matrix.length <= 100  1 <= matrix[0].length <= 100
 * @param target
 * 给出矩阵 matrix 和目标值 target，返回元素总和等于目标值的非空子矩阵的数量。
 * @summary
 * 定三移一，降级为两数之和问题
 */
function numSubmatrixSumTarget(matrix: number[][], target: number): number {
  const m = matrix.length
  const n = matrix[0].length
  const pre = Array.from<any, number[]>({ length: m + 1 }, () => Array(n + 1).fill(0))
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      pre[i][j] = matrix[i - 1][j - 1] + pre[i - 1][j] + pre[i][j - 1] - pre[i - 1][j - 1]
    }
  }

  let res = 0
  // 固定上下
  for (let top = 1; top <= m; top++) {
    for (let bottom = top; bottom <= m; bottom++) {
      const counter = new Map<number, number>([[0, 1]])
      // 遍历子矩阵的右边界，两数之和问题
      for (let right = 1; right <= n; right++) {
        const sum = pre[bottom][right] - pre[top - 1][right]
        if (counter.has(sum - target)) res += counter.get(sum - target)!
        counter.set(sum, (counter.get(sum) || 0) + 1)
      }
    }
  }

  return res
}

console.log(
  numSubmatrixSumTarget(
    [
      [0, 1, 0],
      [1, 1, 1],
      [0, 1, 0],
    ],
    0
  )
)
