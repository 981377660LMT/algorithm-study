/**
 * @param {number} N
 * @param {number} M
 * @param {number} i
 * @param {number} j
 * @return {number}
 * 使 M 对应的二进制数字插入 N 对应的二进制数字的第 i ~ j 位区域，不足之处用 0 补齐
 */
function insertBits(N: number, M: number, i: number, j: number): number {
  // 1.把N的i到j位置为0
  // 2.把M的数值左移i位
  // 3.将N的i到j位加上M
  for (let bit = i; bit <= j; bit++) {
    const mask = 1 << bit
    if (N & mask) N ^= mask
  }

  return N + (M << i)
}

console.log(insertBits(1024, 19, 2, 6))
// 输入：N = 1024(10000000000), M = 19(10011), i = 2, j = 6
// 输出：N = 1100(10001001100)
