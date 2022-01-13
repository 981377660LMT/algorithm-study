/**
 * @param {number} n - a positive integer
 * @return {number}
 * 二进制位1的个数，汉明重量
 */
const hammingWeight = (n: number): number => {
  let res = 0
  while (n) {
    // n & (n - 1) 可以消除 n 最后的一个 1 的原理。
    n &= n - 1
    res++
  }
  return res
}

// console.log(hammingWeight(0b00000000000000000000000000001011))
export { hammingWeight }
