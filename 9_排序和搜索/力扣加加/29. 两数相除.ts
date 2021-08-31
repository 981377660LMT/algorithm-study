/**
 * @param {number} dividend
 * @param {number} divisor
 * @return {number}
 * 将两数相除，要求不使用乘法、除法和 mod 运算符。
   返回被除数 dividend 除以除数 divisor 得到的商。
   整数除法的结果应当截去（truncate）其小数部分
   假设我们的环境只能存储 32 位有符号整数，其数值范围是 [−231,  231 − 1]。本题中，如果除法结果溢出，则返回 231 − 1。
   @summary
   10/3 可以减 1 次3 吗 可以减 2 次3 吗  可以减 4次3吗 ...
 */
const divide = function (dividend: number, divisor: number): number {
  const helper = (dividend: number, divisor: number): number => {
    if (dividend === 0) return 0
    if (dividend < divisor) return 0
    if (divisor === 1) return dividend
    let curDivisor = 2 * divisor
    let res = 1
    while (dividend - curDivisor > 0) {
      curDivisor += curDivisor
      res += res
    }
    const remain = dividend - Math.floor(curDivisor / 2)
    return res + helper(remain, divisor)
  }
  const MAX = 2 ** 31
  const isNegative = dividend !== 0 && dividend * divisor < 0
  const res = helper(Math.abs(dividend), Math.abs(divisor))

  if (res > MAX - 1 || res < -1 * MAX) {
    return MAX - 1
  }

  return isNegative ? -1 * res : res
}

console.log(divide(-2147483648, 1))

export {}
