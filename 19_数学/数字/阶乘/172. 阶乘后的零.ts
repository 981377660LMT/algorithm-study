/**
 * @param {number} n
 * @return {number}
 * @description 给定一个整数 n，返回 n! 结果尾数中零的数量。
 * 题目翻译： 数5
 * 想要结果末尾是 0，必须是分解质因数之后，2 和 5 相乘才行，同时因数分解之后发现 5 的个数远小于 2，
 * 因此我们只需要求解这 n 数字分解质因数之后一共有多少个 5 即可.
 * @summary 也可以勒让德定理做
 */
const trailingZeroes = function (n: number): number {
  let res = 0
  while (n !== 0) {
    const div = Math.floor(n / 5)
    res += div
    n = div
  }
  return res
  return n === 0 ? 0 : ~~(n / 5) + trailingZeroes(~~(n / 5))
}

console.log(trailingZeroes(5))
export { trailingZeroes }
