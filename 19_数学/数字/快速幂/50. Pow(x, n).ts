/**
 * @param {number} x
 * @param {number} n
 * @return {number}
 */
const myPow1 = function (x: number, n: number): number {
  if (n === 0) {
    return 1
  } else if (n < 0) {
    return 1 / myPow1(x, -n)
  } else {
    if (n % 2 === 1) {
      return x * myPow1(x, n - 1)
    } else {
      return myPow1(x * x, n / 2)
    }
  }
}

// 快速幂
// https://leetcode-cn.com/problems/powx-n/solution/50-powx-n-kuai-su-mi-qing-xi-tu-jie-by-jyd/
const myPow2 = function (x: number, n: number): number {
  if (n === 0) {
    return 1
  } else if (n < 0) {
    return 1 / myPow2(x, -n)
  } else {
    // 将n转换为2进制
    let res = 1
    while (n) {
      if (n & 1) res *= x
      x *= x
      n >>= 1
    }
    return res
  }
}
console.log(myPow2(2.1, 3))

// 输出：9.26100
export default 1
