// 果一个整数 n 在 b 进制下（b 为 2 到 n - 2 之间的所有整数）对应的字符串 全部 都是 回文的 ，那么我们称这个数 n 是 严格回文 的。

/**
 * @param {number} n
 * @return {boolean}
 */
let isStrictlyPalindromic = function (n) {
  for (let i = 2; i <= Math.min(36, n - 2); i++) {
    const str = n.toString(i)
    if (str !== str.split('').reverse().join('')) {
      return false
    }
  }

  return true
}

console.log(isStrictlyPalindromic(1000))
