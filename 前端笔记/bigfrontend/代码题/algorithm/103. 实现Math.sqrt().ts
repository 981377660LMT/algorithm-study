/**
 * @param {any} x
 * @return {number}
 */
function mySqrt(x: any): number {
  if (x === Infinity) return Infinity
  if (!Number.isFinite(x) || x < 0) return NaN

  let l = 0
  let r = ~~(x / 2) + 1
  while (l <= r) {
    const m = (l + r) >> 1
    if (m ** 2 <= x && x < (m + 1) ** 2) {
      return m
    } else if (x < m ** 2) {
      r = m - 1
    } else {
      l = m + 1
    }
  }
  return l
}

console.log(mySqrt(7))
console.log(mySqrt(9))
