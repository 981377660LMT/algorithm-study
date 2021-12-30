/**
 * @param {number} n
 * @return {boolean}
 * 你能不使用循环或者递归来完成本题吗？
 */
function isPowerOfFour(n: number): boolean {
  return /^1(00)*$/.test(n.toString(2))
}

// 如果一个数字是 4 的幂次方，那么只需要满足：
// 是二的倍数
// 减去 1 是三的倍数 或者 奇数位上有1 (用0b0101校验奇数位上的1)
function isPowerOfFour2(n: number): boolean {
  // return n > 0 && (n & (n - 1)) === 0 && (n - 1) % 3 === 0
  return n > 0 && (n & (n - 1)) === 0 && (n & 0x55555555) !== 0
}

console.log(isPowerOfFour2(16))
console.log(Number(64).toString(2))
