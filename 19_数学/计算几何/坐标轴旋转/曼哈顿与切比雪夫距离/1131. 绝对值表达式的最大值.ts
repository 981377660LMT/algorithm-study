const DIR4 = [
  [0, 1],
  [0, -1],
  [1, 0],
  [-1, 0]
]

/**
 *
 * @param arr1
 * @param arr2
 * 给你两个长度相等的整数数组，返回下面表达式的最大值：
   |arr1[i] - arr1[j]| + |arr2[i] - arr2[j]| + |i - j|
   !去绝对值+公式变形
 */
function maxAbsValExpr(arr1: number[], arr2: number[]): number {
  const n = arr1.length
  let res = 0

  // 四个方向
  for (const [dx, dy] of DIR4) {
    let min = Infinity
    let max = -Infinity

    for (let i = 0; i < n; i++) {
      min = Math.min(min, arr1[i] * dx + arr2[i] * dy - i)
      max = Math.max(max, arr1[i] * dx + arr2[i] * dy - i)
    }

    res = Math.max(res, max - min)
  }

  return res
}

console.log(maxAbsValExpr([1, -2, -5, 0, 10], [0, -2, -1, -7, -4]))
