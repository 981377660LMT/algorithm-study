/**
 * 
 * @param arr1 
 * @param arr2 
 * 给你两个长度相等的整数数组，返回下面表达式的最大值：
   |arr1[i] - arr1[j]| + |arr2[i] - arr2[j]| + |i - j|
   @summary
   我们将 arr1 中的所有值看做平面上点的横坐标，将 arr2 中的所有值看做平面上点的纵坐标，
   那么表达式 |arr1[i] - arr1[j]| + |arr2[i] - arr2[j]| + |i - j|，
   前两项就可以看成是平面上两点的曼哈顿距离，整个式子就是要求两个点的曼哈顿距离与两个点索引差的和。
   由于有取绝对值的操作存在，那么可能产生的计算有 4 种，分别为：右上减左下，右下减左上，左上减右下，左下减右上。
   @link https://leetcode-cn.com/problems/maximum-of-absolute-value-expression/solution/man-ha-dun-ju-chi-python3-by-smoon1989/
 */
function maxAbsValExpr(arr1: number[], arr2: number[]): number {
  let res = 0
  const directions = [
    [-1, -1],
    [-1, 1],
    [1, -1],
    [1, 1],
  ]

  // 四个方向
  for (const [dx, dy] of directions) {
    let min = Infinity
    let max = -Infinity

    for (let i = 0; i < arr1.length; i++) {
      min = Math.min(min, arr1[i] * dx + arr2[i] * dy + i)
      max = Math.max(max, arr1[i] * dx + arr2[i] * dy + i)
    }

    res = Math.max(res, max - min)
  }

  return res
}

console.log(maxAbsValExpr([1, -2, -5, 0, 10], [0, -2, -1, -7, -4]))
