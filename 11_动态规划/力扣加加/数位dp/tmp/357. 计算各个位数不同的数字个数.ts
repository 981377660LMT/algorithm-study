/**
 * @param {number} n 计算各位数字都不同的数字 x 的个数 其中 0 ≤ x < 10n
 * @return {number}
 * @description
 * n位数 n-1位数 ... 1位数的情况
 */
const countNumbersWithUniqueDigits = function (n: number): number {
  if (n === 0) return 1
  if (n === 1) return 10
  if (n >= 11) return 0
  // 最高位9中选法(1-9) 其余n-1位有 9 8 7 ...
  let res = 9
  for (let i = 0; i < n - 1; i++) {
    res *= 9 - i
  }
  return res + countNumbersWithUniqueDigits(n - 1)
}

// 91
console.log(countNumbersWithUniqueDigits(2))
// 答案应为除去 11,22,33,44,55,66,77,88,99 外，在 [0,100) 区间内的所有数字。
