/**
 * @param {number} n n <= 10^9
 * @return {number}
 */
const numberOf2sInRange = function (n: number): number {
  if (n <= 1) return 0

  if (n < 10) return 1

  const len = n.toString().length
  const base = 10 ** (len - 1) // n四位数的话基数就是1000
  const remainder = n % base
  const times = ~~(n / base)
  let baseCount // 最高位为2的次数
  if (times <= 1) baseCount = 0
  else if (times === 2) baseCount = n - 2 * base + 1
  else baseCount = base

  return numberOf2sInRange(base - 1) * times + baseCount + numberOf2sInRange(remainder)
}

console.log(numberOf2sInRange(251))

export {}
