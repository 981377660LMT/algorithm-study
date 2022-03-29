/**
 * @param {number} n
 * @return {boolean}
 * @description 空间超限
 */
var canWinNim = function (n: number): boolean {
  // const dp = Array<boolean>(n + 1).fill(true)
  // dp[4] = false
  // for (let i = 5; i < n + 1; i++) {
  //   dp[i] = !dp[i - 1] || !dp[i - 2] || !dp[i - 3]
  // }

  // return dp[n]
  return n % 4 !== 0
}

console.log(canWinNim(4))
console.log(canWinNim(5))
console.log(canWinNim(6))
console.log(canWinNim(7))
console.log(canWinNim(8))
