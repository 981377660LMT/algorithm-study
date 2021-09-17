/**
 * @param {number[]} a
 * @param {number[]} b
 * @return {number} 计算具有最小差绝对值的一对数值（每个数组中取一个值），并返回该对数值的差
 */
var smallestDifference = function (a: number[], b: number[]): number {
  a.sort((a, b) => a - b)
  b.sort((a, b) => a - b)
  let res = Infinity
  let i = 0
  let j = 0
  while (i < a.length && j < b.length) {
    res = Math.min(res, Math.abs(a[i] - b[j]))
    a[i] < b[j] ? i++ : j++
  }

  return res
}

console.log(smallestDifference([1, 3, 15, 11, 2], [23, 127, 235, 19, 8]))
// 输出：3，即数值对(11, 8)

export {}
