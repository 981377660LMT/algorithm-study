/**
 * @param {number} x
 * @return {number}
 */
const mySqrt = function (x: number): number {
  if (x === Infinity) return Infinity
  if (!Number.isFinite(x) || x < 0) return NaN

  let l = 0
  let r = x >> 1
  while (l <= r) {
    const mid = (l + r) >> 1

    if (mid ** 2 <= x && x < (mid + 1) ** 2) {
      return mid
    } else if (x < mid ** 2) {
      r = mid - 1
    } else {
      l = mid + 1
    }
  }
  return 0
}

console.log(mySqrt(8))

export {}
