/* eslint-disable no-constant-condition */

/**
 * @param {number} n
 * @return {number}
 * @description 给定一个整数 n，返回 n! 结果尾数中零的数量。
 * 题目翻译： 数5
 * 想要结果末尾是 0，必须是分解质因数之后，2 和 5 相乘才行，同时因数分解之后发现 5 的个数远小于 2，
 * !因此我们只需要求解这 n! 分解质因数之后一共有多少个 5 即可.
 * !勒让德定理
 */
function trailingZeroes(n: number): number {
  let div = 5
  let res = 0
  while (true) {
    const count = Math.floor(n / div)
    if (count === 0) break
    res += count
    div *= 5
  }

  return res
}

console.log(trailingZeroes(5))

export { trailingZeroes }
