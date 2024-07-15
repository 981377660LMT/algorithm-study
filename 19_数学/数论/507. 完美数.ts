/**
 * @param {number} num  1 <= num <= 108
 * @return {boolean}
 * 对于一个 正整数，如果它和除了它自身以外的所有 正因子 之和相等，
 * 我们称它为 「完美数」
 */
var checkPerfectNumber = function (num: number): boolean {
  if (num <= 1) return false
  let res = 0
  for (let i = 1; i <= ~~Math.sqrt(num); i++) {
    if (num % i === 0) res += i + num / i
  }
  return res === num * 2
}

console.log(checkPerfectNumber(28))

export {}
