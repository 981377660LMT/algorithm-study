/**
 *
 * @param n 0 <= n <= 109
 * 给定一个整数 n，计算所有小于等于 n 的非负整数中数字 1 出现的个数。
 * @link https://leetcode-cn.com/problems/number-of-digit-one/solution/xiang-xi-tong-su-de-si-lu-fen-xi-duo-jie-fa-by-50/
 * @description
 * 8192
 * countDigitOne(999)*8 + 1000 + countDigitOne(192)
   1192:
   countDigitOne(999)*1 + (1192 - 1000 + 1) + countDigitOne(192)
 */
function countDigitOne(n: number): number {
  // 对每一位讨论，从高位到地位算
  if (n <= 0) return 0
  if (n < 10) return 1

  const len = n.toString().length
  const base = 10 ** (len - 1) // n四位数的话基数就是1000
  const remainder = n % base
  const times = ~~(n / base)
  const oneInBase = times === 1 ? n - base + 1 : base // 最高位为1的len位数个数

  return countDigitOne(base - 1) * times + oneInBase + countDigitOne(remainder)
}

console.log(countDigitOne(13))
