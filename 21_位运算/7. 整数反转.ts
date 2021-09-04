/**
 * @param {number} x
 * @return {number}
 */
var reverse = function (x: number): number {
  if (x < 0) return -1 * reverse(x * -1)
  let res = 0
  while (x !== 0) {
    const cur = x % 10
    res = res * 10 + cur
    x = ~~(x / 10)
  }
  if (res < Math.pow(-2, 31) || res >= Math.pow(2, 31) - 1) {
    return 0
  }
  return res
}

console.log(reverse(123))
