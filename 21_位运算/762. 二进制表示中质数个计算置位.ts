import { hammingWeight } from './191. 位 1 的个数'

/**
 * @param {number} left
 * @param {number} right
 * @return {number}
 * 191. 位 1 的个数
 * L, R 是 L <= R 且在 [1, 10^6] 中的整数。
   R - L 的最大值为 10000。
   可知10^6化二进制，质数位数最多19位
 */
const countPrimeSetBits = function (left: number, right: number): number {
  const primes = new Set([2, 3, 5, 7, 11, 13, 17, 19])
  let res = 0
  for (let value = left; value < right + 1; value++) {
    primes.has(hammingWeight(value)) && res++
  }
  return res
}
