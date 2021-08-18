/**
 * @param {number} x
 * @return {number}
 */
const mySqrt = function (x: number): number {
  let l = 0
  let r = ~~(x / 2) + 1
  while (l <= r) {
    const mid = ~~((l + r) / 2)

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
