/**
 * @param {number} n
 * @param {number} start
 * @return {number}
 * 请返回 nums 中所有元素按位异或（XOR）后得到的结果。
 * @summary 别想多了 就是模拟
 */
const xorOperation = function (n: number, start: number): number {
  let res = 0
  for (let i = 0; i < n; i++) {
    res ^= start + 2 * i
  }

  return res
}

console.log(xorOperation(5, 0))
// 输出：8
// 解释：数组 nums 为 [0, 2, 4, 6, 8]，其中 (0 ^ 2 ^ 4 ^ 6 ^ 8) = 8 。
//      "^" 为按位异或 XOR 运算符。
