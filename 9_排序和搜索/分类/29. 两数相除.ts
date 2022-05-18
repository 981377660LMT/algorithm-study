/**
 * @param {number} dividend
 * @param {number} divisor
 * @return {number}
 * 将两数相除，要求不使用乘法、除法和 mod 运算符。
   返回被除数 dividend 除以除数 divisor 得到的商。
   整数除法的结果应当截去（truncate）其小数部分
   假设我们的环境只能存储 32 位有符号整数，其数值范围是 [−231,  231 − 1]。本题中，如果除法结果溢出，则返回 231 − 1。
   @summary
   从2^31试到2^0 
   直到被除数被减到比除数小，
   每个能满足除出来的最大的2的幂都加入答案
   也可以理解为每次计算出答案的32位中的某一位
 */
const divide = function (dividend: number, divisor: number): number {
  if (dividend === 0) return 0

  const [MAX, MIN] = [2 ** 31 - 1, -1 * 2 ** 31]
  const isNegative = (dividend ^ divisor) < 0
  let [posDividend, posDivisor] = [Math.abs(dividend), Math.abs(divisor)]
  let res = 0

  for (let i = 31; ~i; i--) {
    // 找出满足条件的最大的倍数,全部用二进制数考虑
    if (posDividend >>> i >= posDivisor) {
      // 累加上这个倍数
      res += 1 << i
      // 被除数减去这个倍数*b
      posDividend -= posDivisor << i
    }
  }

  if (res >= MAX || res <= MIN) {
    return MAX
  }

  return isNegative ? -1 * res : res
}

console.log(divide(-2147483648, 1))

export {}
