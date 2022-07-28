// 如果一个数字是 4 的幂次方，那么只需要满足：
// 是二的倍数
// 减去 1 是三的倍数 或者 奇数位上有1 (用0b0101校验奇数位上的1)

// 判断4的幂
function isPowerOfFour(n: number): boolean {
  // return n > 0 && (n & (n - 1)) === 0 && (n - 1) % 3 === 0

  // !偶数二进制位都是 0，奇数二进制位都是 1
  return n > 0 && (n & (n - 1)) === 0 && (n & 0xaaaaaaaa) === 0
}

console.log(isPowerOfFour(16))
console.log(Number(64).toString(2))
