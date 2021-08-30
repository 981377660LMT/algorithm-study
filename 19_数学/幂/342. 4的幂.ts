/**
 * @param {number} n
 * @return {boolean}
 * 你能不使用循环或者递归来完成本题吗？
 */
var isPowerOfFour = function (n: number): boolean {
  return /^1(00)*$/.test(n.toString(2))
}

// 如果一个数字是 4 的幂次方，那么只需要满足：
// 是二的倍数
// 减去 1 是三的倍数
var isPowerOfFour2 = function (n: number): boolean {
  return n > 0 && (n & (n - 1)) === 0 && (n - 1) % 3 === 0
}
console.log(isPowerOfFour2(16))
console.log(Number(64).toString(2))
