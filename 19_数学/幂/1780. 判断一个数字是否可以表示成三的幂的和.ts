/**
 * @param {number} n
 * @return {boolean}
 * 给你一个整数 n ，如果你可以将 n 表示成若干个不同的三的幂之和
 */
const checkPowersOfThree = function (n: number): boolean {
  return !n.toString(3).includes('2')
}

console.log(checkPowersOfThree(91))
// 解释：91 = 3**0 + 3**2 + 3**4
