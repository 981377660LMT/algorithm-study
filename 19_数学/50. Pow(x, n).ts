/**
 * @param {number} x
 * @param {number} n
 * @return {number}
 */
const myPow = function (x: number, n: number): number {
  if (n === 0) {
    return 1
  } else if (n < 0) {
    return 1 / myPow(x, -n)
  } else {
    if (n % 2 === 1) {
      return x * myPow(x, n - 1)
    } else {
      return myPow(x * x, n / 2)
    }
  }
}

console.log(myPow(2.1, 3))

// 输出：9.26100
export default 1
