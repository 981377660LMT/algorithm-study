/**
 * @param {number} n - a positive integer
 * @return {number}
 * 也被称为汉明重量
 */
var hammingWeight = (n: number): number => {
  // let res = 0
  // while (n) {
  //   res += n & 1
  //   n = n >>> 1
  // }
  // return res
  let res = 0
  while (n) {
    // n & (n - 1) 可以消除 n 最后的一个 1 的原理。
    n &= n - 1
    res++
  }
  return res
}

console.log(hammingWeight(0b00000000000000000000000000001011))
