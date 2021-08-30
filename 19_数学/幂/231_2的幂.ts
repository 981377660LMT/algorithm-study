/**
 * @param {number} n
 * @return {boolean}
 * 如果一个数字 n 是 2 的幂次方，那么 n & (n - 1) 一定等于 0
 */
let isPowerOfTwo = (n: number): boolean => n > 0 && (n & (n - 1)) === 0

console.log(isPowerOfTwo(4))
