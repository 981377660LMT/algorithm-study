/**
 * @param {number} x
 * @return {number}
 */
const reverse = function (x: number): number {
  if (x < 0) return -1 * reverse(x * -1)

  let res = 0
  while (x > 0) {
    const [div, mod] = [~~(x / 10), x % 10]
    x = div
    res = res * 10 + mod
  }

  if (res < (-2) ** 31 || res > 2 ** 31 - 1) {
    return 0
  }

  return res
}

console.log(reverse(123))
