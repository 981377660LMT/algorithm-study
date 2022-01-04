/**
 * @param {number} n
 * @return {number}
 * 给定一个非负整数 N，找出小于或等于 N 的最大的整数，同时这个整数需要满足其各个位数上的数字是单调递增。
 * @summary
 * 思路是从左向右 如果strNum[i] < strNum[i+1] 可以的
 * 如果strNum[i] === strNum[i+1] 看后面
 * 如果strNum[i] > strNum[i+1] 跳出循环
 */
const monotoneIncreasingDigits = function (n: number): number {
  const digits = n.toString()

  let mid = 0
  let isIncreased = true
  for (let i = 0; i < digits.length - 1; i++) {
    if (Number(digits[i + 1]) > Number(digits[i])) mid = i + 1
    if (Number(digits[i + 1]) < Number(digits[i])) {
      isIncreased = false
      break
    }
  }
  console.log(digits, mid)

  if (isIncreased) return n
  return Number(
    digits.slice(0, mid) + String(Number(digits[mid]) - 1) + '9'.repeat(digits.length - 1 - mid)
  )
}

console.log(monotoneIncreasingDigits(343))
console.log(monotoneIncreasingDigits(332))
console.log(monotoneIncreasingDigits(345))
console.log(monotoneIncreasingDigits(321))

export default 1
