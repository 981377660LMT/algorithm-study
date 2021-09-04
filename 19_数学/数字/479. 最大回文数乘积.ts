/**
 * @param {number} n  n 的取值范围为 [1,8]。
 * @return {number}
 * 你需要找到由两个 n 位数的乘积组成的最大回文数
 * @link https://leetcode-cn.com/problems/largest-palindrome-product/solution/9-line-in-java-by-wdw87/
 */
var largestPalindrome = function (n: number): number {
  const getPalindrome = (s: string) => s + s.split('').reverse().join('')
  if (n == 1) return 9
  //计算给定位数的最大值
  const max = 10 ** n - 1
  for (let i = max; i > max / 10; i--) {
    // 1. 构造回文数
    const pal = getPalindrome(i.toString())
    //2. 检验该回文数能否由给定的数相乘得到
    for (let j = max; j * j >= parseInt(pal); j--) {
      if (parseInt(pal) % j === 0) {
        console.log(pal)
        return parseInt(pal) % 1337
      }
    }
  }
  return -1
}

console.log(largestPalindrome(8))

// 输出: 987
// 解释: 99 x 91 = 9009, 9009 % 1337 = 987
