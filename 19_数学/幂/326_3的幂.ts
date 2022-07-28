/**
 * @param {number} n
 * @return {boolean}
 */
function isPowerOfThree(n: number): boolean {
  while (n && n % 3 === 0) {
    n = n / 3
  }

  return n === 1
}

console.log(isPowerOfThree(9))
console.log(Number.EPSILON)
console.log(0.1 + 0.2 - 0.3 < Number.EPSILON)
