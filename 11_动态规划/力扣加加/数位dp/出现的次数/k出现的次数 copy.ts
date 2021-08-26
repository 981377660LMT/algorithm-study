/**
 * @param {number} n n <= 10^9
 * @return {number}
 */
const numberOfKsInRange = function (n: number, k: number): number {
  if (n <= k - 1) return 0
  if (n < 10) return 1

  const len = n.toString().length
  const base = 10 ** (len - 1) // n四位数的话基数就是1000
  const remainder = n % base
  const times = ~~(n / base)
  let baseCount // 最高位为2的次数
  if (times <= k - 1) baseCount = 0
  else if (times === k) baseCount = n - k * base + 1
  else baseCount = base

  return numberOfKsInRange(base - 1, k) * times + baseCount + numberOfKsInRange(remainder, k)
}

console.log(numberOfKsInRange(251, 2))
