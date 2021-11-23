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
  if (n === 0) return 1
  else if (n < 0) return 1 / myPow2(x, -n)

  // 将n转换为2进制
  let res = 1
  while (n) {
    //如果n的当前末位为1,res乘上当前的x,x自乘,n往右移一位
    if (n & 1) res *= x
    x = x ** 2
    n >>>= 1
  }

  return res
}
console.log(myPow2(2.1, 3))

// 输出：9.26100
export default 1

// 最初ans为1，然后我们一位一位算：
// 1010的最后一位是0，所以a^1这一位不要。然后1010变为101，a变为a^2。
// 101的最后一位是1，所以a^2这一位是需要的，乘入ans。101变为10，a再自乘。
// 10的最后一位是0，跳过，右移，自乘。
// 然后1的最后一位是1，ans再乘上a^8。循环结束，返回结果。

// 矩阵快速幂的一个经典应用是求斐波那契数列
