/**
 * @param {number} x
 * @param {number} y
 * @return {number}
 * 两个整数之间的 汉明距离 指的是这两个数字对应二进制位不同的位置的数目。
 */
// var hammingDistance = function (x: number, y: number): number {
//   let res = 0
//   while (x || y) {
//     res += (x ^ y) & 1
//     x = x >> 1
//     y = y >> 1
//     console.log(x, y)
//   }

//   return res
// }
// 解法 2: 布赖恩·克尼根算法（推荐）
// 它是借助 num & (num - 1) 来直接去除 num 的二进制中最右边的 1。
// 例如6 & 5:0b110 变为 0b100
var hammingDistance = function (x: number, y: number): number {
  let res = 0
  let xor = x ^ y
  while (xor) {
    xor = xor & (xor - 1)
    res++
  }

  return res
}
console.log(hammingDistance(1, 4))
