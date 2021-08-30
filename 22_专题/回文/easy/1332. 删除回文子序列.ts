/**
 * @param {string} s
 * @return {number}
 * 每一次删除操作都可以从 s 中删除一个回文 子序列。
 * 返回删除给定字符串中所有字符（字符串为空）的最小删除次数。
 * 由于只有 a 和 b 两个字符。其实最多的消除次数就是 2。这是因为每次我们可以消除一个子序列
 */
const removePalindromeSub = function (s: string): number {
  if (!s) return 0
  for (let i = 0, j = s.length - 1; i < j; i++, j--)
    if (s.codePointAt(i) !== s.codePointAt(j)) return 2
  return 1
}

console.log(removePalindromeSub('ababa'))
