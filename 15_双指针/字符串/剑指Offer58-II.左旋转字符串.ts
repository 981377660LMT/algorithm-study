/**
 * @param {string} s
 * @param {number} n
 * @return {string}
 * 不能申请额外空间，只能在本串上操作
 * @summary 相当于循环移动字符串
 */
const reverseLeftWords = function (s: string, n: number): string {
  const sb = s.split('')
  const reverse = (sb: string[], left: number, right: number) => {
    for (; left < right; left++, right--) {
      ;[sb[left], sb[right]] = [sb[right], sb[left]]
    }
  }

  reverse(sb, 0, n - 1)
  reverse(sb, n, sb.length - 1)
  reverse(sb, 0, sb.length - 1)
  return sb.join('')
}

console.log(reverseLeftWords('abcdefg', 2))

export default 1
