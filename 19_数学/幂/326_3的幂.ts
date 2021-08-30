/**
 * @param {number} n
 * @return {boolean}
 */
var isPowerOfThree = function (n: number): boolean {
  const res = Math.log(n) / Math.log(3)
  return n > 0 && parseFloat(res.toFixed(10)) === Math.floor(res)
}

console.log(isPowerOfThree(9))
console.log(Number.EPSILON)
console.log(0.1 + 0.2 - 0.3 < Number.EPSILON)
